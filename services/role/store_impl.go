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
	s.addToResourceMap(ctx, role.Id, role.Resources)
	return nil
}

func (s *store) Update(ctx context.Context, role *sdk.Role) error {
	now := time.Now()
	role.UpdatedAt = &now
	if role.Id == "" {
		return errors.New("role not found")
	}

	// Get roleMap from context
	roleMap, ok := ctx.Value("roleMap").(map[string][]string)
	if !ok {
		return errors.New("roleMap missing in context")
	}

	// Remove role from all assigned users
	for _, user := range roleMap[role.Id] {
		if err := s.RemoveRoleFromUser(ctx, user, role.Id); err != nil {
			return fmt.Errorf("error removing role from user %s: %w", user, err)
		}
	}

	// Fetch existing role data
	existingRole, err := s.GetById(ctx, role.Id)
	if err != nil {
		return fmt.Errorf("error finding role: %w", err)
	}
	role.CreatedAt = existingRole.CreatedAt

	// Convert role to DB model and update in database
	d := fromSdkToModel(*role)
	md := models.GetRoleModel()
	_, err = s.db.UpdateOne(ctx, md, bson.M{md.IdKey: role.Id}, bson.M{"$set": d})
	if err != nil {
		return fmt.Errorf("error updating role: %w", err)
	}

	// Reassign role to users
	for _, userId := range roleMap[role.Id] {
		if err := s.AddRoleToUser(ctx, userId, role.Id); err != nil {
			return fmt.Errorf("error adding role to user %s: %w", userId, err)
		}
	}

	// Update the resource map
	if err := s.addToResourceMap(ctx, role.Id, role.Resources); err != nil {
		return fmt.Errorf("error updating resource map: %w", err)
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

func (s *store) AddRoleToUser(ctx context.Context, userId string, roleId string) error {

	// get the user
	userMd := models.GetUserModel()
	var user models.User
	if err := s.db.FindOne(ctx, userMd, bson.D{{Key: userMd.IdKey, Value: userId}}).Decode(&user); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return errors.New("user not found")
		}
		return fmt.Errorf("error finding user: %w", err)
	}

	// get the role
	roleMd := models.GetRoleModel()
	var role models.Role
	if err := s.db.FindOne(ctx, roleMd, bson.D{{Key: roleMd.IdKey, Value: roleId}}).Decode(&role); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return errors.New("role not found")
		}
		return fmt.Errorf("error finding role: %w", err)
	}

	// check if the user already has the role
	for _, r := range user.Roles {
		if r.Id == roleId {
			return errors.New("role already assigned to user")
		}
	}

	// add the role to the user
	user.Roles = append(user.Roles, models.UserRoles{Id: role.Id, Name: role.Name})

	// ensure unique resources when adding to the user
	resourceSet := make(map[string]struct{}, len(user.Resource))
	for _, res := range user.Resource {
		resourceSet[res.Key] = struct{}{}
	}

	for _, res := range role.Resources {
		if _, exists := resourceSet[res.Key]; !exists {
			user.Resource = append(user.Resource, models.UserResource{Key: res.Key, Name: res.Name})
			resourceSet[res.Key] = struct{}{}
		}
	}

	if _, err := s.db.UpdateOne(ctx, userMd, bson.D{{Key: userMd.IdKey, Value: userId}}, bson.D{{Key: "$set", Value: user}}); err != nil {
		return fmt.Errorf("error adding resource to user: %w", err)
	}

	s.addToRoleMap(ctx, roleId, userId)
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

	// Remove role from user roles
	s.removeFromRoleMap(ctx, roleId, userId)

	// Remove the role from the user's role list
	filteredRoles := []models.UserRoles{}
	for _, r := range user.Roles {
		if r.Id != roleId { // Keep only roles that are NOT being removed
			filteredRoles = append(filteredRoles, r)
		}
	}
	user.Roles = filteredRoles

	// Remove the associated resources with the roleId
	resourceMap, ok := ctx.Value("resourceMap").(map[string][]string)
	if !ok {
		return errors.New("resourceMap missing in context")
	}

	var updatedResources []models.UserResource
	for _, resource := range user.Resource {
		// Filter roles that are NOT the one being removed
		remainingRoles := []string{}
		for _, r := range resourceMap[resource.Key] {
			if r != roleId {
				remainingRoles = append(remainingRoles, r)
			}
		}

		// If at least one role remains for the resource, keep it
		if len(remainingRoles) > 0 {
			updatedResources = append(updatedResources, resource)
		}
	}
	user.Resource = updatedResources

	// Update the user document in the database
	_, err = s.db.UpdateOne(ctx, md, bson.D{{Key: md.IdKey, Value: userId}}, bson.D{{Key: "$set", Value: user}})
	if err != nil {
		return fmt.Errorf("error removing role from user: %w", err)
	}

	return nil
}
