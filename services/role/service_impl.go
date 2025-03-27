package role

import (
	"context"
	"errors"
	"fmt"

	"github.com/melvinodsa/go-iam/db/models"
	"github.com/melvinodsa/go-iam/sdk"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type service struct {
	store Store
}

func NewService(store Store) Service {
	return &service{
		store: store,
	}
}

func (s *service) Create(ctx context.Context, role *sdk.Role) error {
	return s.store.Create(ctx, role)
}

func (s *service) Update(ctx context.Context, role *sdk.Role) error {
	// Update the role in the database
	if err := s.store.Update(ctx, role); err != nil {
		return err
	}

	// Refresh users with this role
	return s.refreshUsersWithRole(ctx, role)
}

func (s *service) GetById(ctx context.Context, id string) (*sdk.Role, error) {
	return s.store.GetById(ctx, id)
}

func (s *service) GetAll(ctx context.Context, query sdk.RoleQuery) ([]sdk.Role, error) {
	return s.store.GetAll(ctx, query)
}

func (s *service) AddRoleToUser(ctx context.Context, userId, roleId string) error {
	if userId == "" || roleId == "" {
		return errors.New("user ID and role ID are required")
	}

	// Fetch user from the database
	userMd := models.GetUserModel()
	var user models.User
	err := s.store.(*store).db.FindOne(ctx, userMd, bson.M{userMd.IdKey: userId}).Decode(&user)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return fmt.Errorf("user with ID %s not found", userId)
		}
		return fmt.Errorf("failed to fetch user: %w", err)
	}

	// Initialize roles and resources if nil
	if user.Roles == nil {
		user.Roles = make(map[string]models.UserRoles)
	}
	if user.Resource == nil {
		user.Resource = make(map[string]models.UserResource)
	}

	// Check if role is already assigned
	if _, exists := user.Roles[roleId]; exists {
		return nil
	}

	// Get role details
	role, err := s.GetById(ctx, roleId)
	if err != nil {
		return err
	}

	// Add new role to user's roles
	user.Roles[roleId] = models.UserRoles{
		Id:   role.Id,
		Name: role.Name,
	}

	// Add unique resources from the role
	resourceSet := make(map[string]struct{})
	for _, res := range user.Resource {
		resourceSet[res.Key] = struct{}{}
	}
	for _, res := range role.Resources {
		if _, exists := resourceSet[res.Key]; !exists {
			user.Resource[res.Key] = models.UserResource{
				Id:   res.Id,
				Key:  res.Key,
				Name: res.Name,
			}
			resourceSet[res.Key] = struct{}{}
		}
	}

	// First add the role to the user's roles in the database
	if err := s.store.AddRoleToUser(ctx, userId, roleId); err != nil {
		return err
	}

	// Then update the roles and resources
	return s.updateUserRoleAndResources(ctx, &user)
}

func (s *service) RemoveRoleFromUser(ctx context.Context, userId, roleId string) error {
	if userId == "" || roleId == "" {
		return errors.New("user ID and role ID are required")
	}

	// Fetch user from the database
	userMd := models.GetUserModel()
	var user models.User
	err := s.store.(*store).db.FindOne(ctx, userMd, bson.M{userMd.IdKey: userId}).Decode(&user)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return fmt.Errorf("user with ID %s not found", userId)
		}
		return fmt.Errorf("failed to fetch user: %w", err)
	}

	// Track resource usage across remaining roles
	resourceUsage := s.calculateResourceUsage(user.Roles, roleId)

	// Remove resources that are no longer used
	user.Resource = s.cleanupUnusedResources(user.Resource, resourceUsage)

	// Remove role from user
	delete(user.Roles, roleId)

	// First remove the role from the user's roles in the database
	if err := s.store.RemoveRoleFromUser(ctx, userId, roleId); err != nil {
		return err
	}

	// Then update the roles and resources
	return s.updateUserRoleAndResources(ctx, &user)
}

// updateUserRoleAndResources updates a single user's roles and resources
func (s *service) updateUserRoleAndResources(ctx context.Context, user *models.User) error {
	updatedRoles := make(map[string]models.UserRoles)
	resourceSet := make(map[string]models.UserResource)

	for roleId := range user.Roles {
		updatedRole, err := s.GetById(ctx, roleId)
		if err != nil {
			return fmt.Errorf("failed to fetch role %s: %w", roleId, err)
		}

		updatedRoles[roleId] = models.UserRoles{
			Id:   updatedRole.Id,
			Name: updatedRole.Name,
		}

		for _, res := range updatedRole.Resources {
			resourceSet[res.Key] = models.UserResource{
				Id:   res.Id,
				Key:  res.Key,
				Name: res.Name,
			}
		}
	}

	// Persist changes to the database
	userMd := models.GetUserModel()
	update := bson.D{
		{Key: "$set", Value: bson.M{
			"roles":    updatedRoles,
			"resource": resourceSet,
		}},
	}

	_, err := s.store.(*store).db.UpdateOne(
		ctx,
		userMd,
		bson.D{{Key: userMd.IdKey, Value: user.Id}},
		update,
	)
	if err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}

	return nil
}

// refreshUsersWithRole updates all users with the modified role
func (s *service) refreshUsersWithRole(ctx context.Context, role *sdk.Role) error {
	userMd := models.GetUserModel()

	// Find users with this role
	users, err := s.store.(*store).db.Find(
		ctx,
		userMd,
		bson.D{{Key: "roles." + role.Id, Value: bson.D{{Key: "$exists", Value: true}}}},
	)
	if err != nil {
		return fmt.Errorf("failed to fetch users with role: %w", err)
	}
	defer users.Close(ctx)

	// Batch update users
	for users.Next(ctx) {
		var user models.User
		if err := users.Decode(&user); err != nil {
			return fmt.Errorf("failed to decode user: %w", err)
		}

		if err := s.updateUserRoleAndResources(ctx, &user); err != nil {
			return err
		}
	}

	return nil
}

// calculateResourceUsage counts resource usage across remaining roles
func (s *service) calculateResourceUsage(userRoles map[string]models.UserRoles, roleToSkip string) map[string]int {
	resourceUsage := make(map[string]int)

	for roleId := range userRoles {
		if roleId == roleToSkip {
			continue
		}

		role, err := s.GetById(context.Background(), roleId)
		if err != nil {
			// Log the error or handle it appropriately
			continue
		}

		for _, res := range role.Resources {
			resourceUsage[res.Key]++
		}
	}

	return resourceUsage
}

// cleanupUnusedResources removes resources that are no longer used by any role
func (s *service) cleanupUnusedResources(
	currentResources map[string]models.UserResource,
	resourceUsage map[string]int,
) map[string]models.UserResource {
	cleanedResources := make(map[string]models.UserResource)

	for key, resource := range currentResources {
		// Keep the resource if it's used in other roles
		if resourceUsage[key] > 0 {
			cleanedResources[key] = resource
		}
	}

	return cleanedResources
}
