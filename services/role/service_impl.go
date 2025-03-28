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

func NewService(store Store) *service {
	return &service{store: store}
}

func (s *service) Create(ctx context.Context, role *sdk.Role) error {
	return s.store.Create(ctx, role)
}

func (s *service) Update(ctx context.Context, role *sdk.Role) error {
	if err := s.store.Update(ctx, role); err != nil {
		return err
	}
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

	userMd := models.GetUserModel()
	var user models.User
	err := s.store.(*store).db.FindOne(ctx, userMd, bson.M{userMd.IdKey: userId}).Decode(&user)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return fmt.Errorf("user with ID %s not found", userId)
		}
		return fmt.Errorf("failed to fetch user: %w", err)
	}

	if user.Roles == nil {
		user.Roles = make(map[string]models.UserRoles)
	}
	if user.Resource == nil {
		user.Resource = make(map[string]models.UserResource)
	}

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
