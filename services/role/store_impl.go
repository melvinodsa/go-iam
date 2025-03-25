package role

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2/log"
	"github.com/google/uuid"
	"github.com/melvinodsa/go-iam/db"
	"github.com/melvinodsa/go-iam/db/models"
	"github.com/melvinodsa/go-iam/sdk"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type store struct {
	db db.DB
}

func NewStore(db db.DB) Store {
	return &store{
		db: db,
	}
}

func (s *store) Create(ctx context.Context, role *sdk.Role) error {
	id := uuid.New().String()
	role.Id = id
	t := time.Now()
	role.CreatedAt = &t
	d := fromSdkToModel(*role)
	md := models.GetRoleModel()
	_, err := s.db.InsertOne(ctx, md, d)
	if err != nil {
		return fmt.Errorf("error creating role: %w", err)
	}

	for _, res := range role.Resources {
		resMd := models.GetResourceMap()
		var resourceMap models.ResourceMap
		err = s.db.FindOne(ctx, resMd, bson.M{"resource_id": res.Id}).Decode(&resourceMap)
		if err != nil {
			if errors.Is(err, mongo.ErrNoDocuments) {
				// resource does not exist in resourceMap collection, create it
				_, err = s.db.InsertOne(ctx, resMd, models.ResourceMap{
					ResourceId: res.Id,
					RoleId:     []string{role.Id},
				})
				if err != nil {
					return fmt.Errorf("error creating resource map: %w", err)
				}
			} else {
				return fmt.Errorf("error finding resource map: %w", err)
			}
		} else {
			// resource exists in resourceMap collection, update it
			roleIds := append(resourceMap.RoleId, role.Id)
			_, err = s.db.UpdateOne(ctx, resMd, bson.M{"resource_id": res.Id}, bson.M{"$set": bson.M{"role_id": roleIds}})
			if err != nil {
				return fmt.Errorf("error updating resource map: %w", err)
			}
		}
	}

	return nil
}

func (s *store) Update(ctx context.Context, role *sdk.Role) error {
	now := time.Now()
	role.UpdatedAt = &now
	if role.Id == "" {
		return errors.New("role not found")
	}

	var users []models.User
	Md := models.GetUserModel()
	cursor, err := s.db.Find(ctx, Md, bson.D{})
	if err != nil {
		return fmt.Errorf("error finding users with role: %w", err)
	}
	defer cursor.Close(ctx)

	if err := cursor.All(ctx, &users); err != nil {
		return fmt.Errorf("error decoding users with role: %w", err)
	}

	// filter the users with the role id in the roles array
	var usersWithRole []models.User
	for _, user := range users {
		for _, r := range user.Roles {
			if r.Id == role.Id {
				usersWithRole = append(usersWithRole, user)
				s.RemoveRoleFromUser(ctx, user.Id, role.Id)
				break
			}
		}
	}

	o, err := s.GetById(ctx, role.Id)
	if err != nil {
		return fmt.Errorf("error finding role: %w", err)
	}
	role.CreatedAt = o.CreatedAt
	d := fromSdkToModel(*role)
	md := models.GetRoleModel()
	_, err = s.db.UpdateOne(ctx, md, bson.D{{Key: md.IdKey, Value: role.Id}}, bson.D{{Key: "$set", Value: d}})
	if err != nil {
		return fmt.Errorf("error updating role: %w", err)
	}

	// add the role to the users
	for _, user := range usersWithRole {
		err := s.AddRoleToUser(ctx, user.Id, role.Id)
		if err != nil {
			return fmt.Errorf("error adding role to user %s: %w", user.Id, err)
		}
	}

	return nil
}

func (s *store) GetById(ctx context.Context, id string) (*sdk.Role, error) {
	md := models.GetRoleModel()
	var role models.Role
	err := s.db.FindOne(ctx, md, bson.D{{Key: md.IdKey, Value: id}}).Decode(&role)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, errors.New("role not found")
		}
		return nil, fmt.Errorf("error finding role: %w", err)
	}
	return fromModelToSdk(&role), nil
}

func (s *store) GetAll(ctx context.Context, query sdk.RoleQuery) ([]sdk.Role, error) {
	md := models.GetRoleModel()
	var roles []models.Role
	filter := bson.D{}

	if query.ProjectId != "" {
		filter = append(filter, bson.E{Key: md.ProjectIdKey, Value: query.ProjectId})
	}
	if query.SearchQuery != "" {
		filter = append(filter, bson.E{
			Key: "$or", Value: bson.A{
				bson.D{{Key: md.NameKey, Value: bson.D{{Key: "$regex", Value: query.SearchQuery}, {Key: "$options", Value: "i"}}}},
			},
		})
	}

	cursor, err := s.db.Find(ctx, md, filter)
	if err != nil {
		return nil, fmt.Errorf("error finding roles: %w", err)
	}
	defer func() {
		if err := cursor.Close(ctx); err != nil {
			log.Errorw("error closing cursor after reading roles", "error", err)
		}
	}()

	err = cursor.All(ctx, &roles)
	if err != nil {
		return nil, fmt.Errorf("error reading roles: %w", err)
	}
	return fromModelListToSdk(roles), nil
}

// function to add a role to a user

func (s *store) AddRoleToUser(ctx context.Context, userId string, roleId string) error {
	md := models.GetUserModel()
	var user models.User
	err := s.db.FindOne(ctx, md, bson.D{{Key: md.IdKey, Value: userId}}).Decode(&user)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return errors.New("user not found")
		}
		return fmt.Errorf("error finding user: %w", err)
	}
	roleMd := models.GetRoleModel()
	var role models.Role
	err = s.db.FindOne(ctx, roleMd, bson.D{{Key: roleMd.IdKey, Value: roleId}}).Decode(&role)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return errors.New("role not found")
		}
		return fmt.Errorf("error finding role: %w", err)
	}
	if user.Roles == nil {
		user.Roles = []models.UserRoles{}
	}

	// check if the role is already assigned to the user
	for _, r := range user.Roles {
		if r.Id == roleId {
			return errors.New("role already assigned to user")
		}
	}

	user.Roles = append(user.Roles, models.UserRoles{
		Id:   role.Id,
		Name: role.Name,
	})
	_, err = s.db.UpdateOne(ctx, md, bson.D{{Key: md.IdKey, Value: userId}}, bson.D{{Key: "$set", Value: user}})
	if err != nil {
		return fmt.Errorf("error adding role to user: %w", err)
	}

	// add resources corresponding to the role to the user
	if user.Resource == nil {
		user.Resource = []models.UserResource{}
	}
	for _, res := range role.Resources {
		// get the resource model using the role.resource id
		resMd := models.GetResourceModel()
		var resource models.Resource
		err = s.db.FindOne(ctx, resMd, bson.D{{Key: resMd.IdKey, Value: res.Id}}).Decode(&resource)
		if err != nil {
			if errors.Is(err, mongo.ErrNoDocuments) {
				return errors.New("resource not found")
			}
			return fmt.Errorf("error finding resource: %w", err)
		}
		user.Resource = append(user.Resource, models.UserResource{
			RoleID: role.Id,
			Key:    resource.Key,
			Name:   resource.Name,
			Scope:  res.Scopes,
		})
	}
	_, err = s.db.UpdateOne(ctx, md, bson.D{{Key: md.IdKey, Value: userId}}, bson.D{{Key: "$set", Value: user}})
	if err != nil {
		return fmt.Errorf("error adding resource to user: %w", err)
	}

	// add the role_id to user_id in the roleMap collection
	roleMapMd := models.GetRoleMap()
	fmt.Println("roleid", roleId)
	// check if the roleMap already exists
	var roleMap models.RoleMap
	err = s.db.FindOne(ctx, roleMapMd, bson.M{"role_id": roleId}).Decode(&roleMap)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			// roleMap does not exist, create it
			roleMap = models.RoleMap{
				RoleId: roleId,
				UserId: []string{userId},
			}
			_, err = s.db.InsertOne(ctx, roleMapMd, roleMap)
			if err != nil {
				return fmt.Errorf("error creating roleMap: %w", err)
			}
		} else {
			return fmt.Errorf("error finding roleMap: %w", err)
		}
	} else {
		// roleMap exists, update it
		roleMap.UserId = append(roleMap.UserId, userId)
		_, err = s.db.UpdateOne(ctx, roleMapMd, bson.M{"role_id": roleId}, bson.M{"$set": bson.M{"user_id": roleMap.UserId}})
		if err != nil {
			return fmt.Errorf("error updating roleMap: %w", err)
		}
	}
	// return the user with the updated roles and resources
	return nil
}

func (s *store) RemoveRoleFromUser(ctx context.Context, userId string, roleId string) error {
	md := models.GetUserModel()
	var user models.User
	err := s.db.FindOne(ctx, md, bson.D{{Key: md.IdKey, Value: userId}}).Decode(&user)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return errors.New("user not found")
		}
		return fmt.Errorf("error finding user: %w", err)
	}

	// Find the index of the role to remove
	roleIndex := -1
	for i, r := range user.Roles {
		if r.Id == roleId {
			roleIndex = i
			break
		}
	}

	if roleIndex == -1 {
		return errors.New("role not assigned to user")
	}

	// Remove the role from the user's Roles
	user.Roles = append(user.Roles[:roleIndex], user.Roles[roleIndex+1:]...)

	// Remove the resources associated with the role
	var updatedResources []models.UserResource
	for _, resource := range user.Resource {
		if resource.RoleID != roleId {
			updatedResources = append(updatedResources, resource)
		}
	}
	user.Resource = updatedResources

	// Update the user document in the database
	_, err = s.db.UpdateOne(ctx, md, bson.D{{Key: md.IdKey, Value: userId}}, bson.D{{Key: "$set", Value: user}})
	if err != nil {
		return fmt.Errorf("error removing role from user: %w", err)
	}

	// Return nil to indicate success
	return nil
}
