package user

import (
	"context"
	"errors"
	"fmt"

	"github.com/melvinodsa/go-iam/middlewares/projects"
	"github.com/melvinodsa/go-iam/sdk"
	"github.com/melvinodsa/go-iam/services/role"
)

type service struct {
	store   Store
	roleSvc role.Service
}

func NewService(store Store, roleSvc role.Service) Service {
	return &service{
		store:   store,
		roleSvc: roleSvc,
	}
}

func (s *service) Create(ctx context.Context, user *sdk.User) error {
	return s.store.Create(ctx, user)
}

func (s *service) Update(ctx context.Context, user *sdk.User) error {
	return s.store.Update(ctx, user)
}

func (s *service) GetByEmail(ctx context.Context, email string, projectId string) (*sdk.User, error) {
	return s.store.GetByEmail(ctx, email, projectId)
}

func (s *service) GetById(ctx context.Context, id string) (*sdk.User, error) {
	return s.store.GetById(ctx, id)
}

func (s *service) GetByPhone(ctx context.Context, phone string, projectId string) (*sdk.User, error) {
	return s.store.GetByPhone(ctx, phone, projectId)
}

func (s *service) GetAll(ctx context.Context, query sdk.UserQuery) (*sdk.UserList, error) {
	query.ProjectIds = projects.GetProjects(ctx)
	return s.store.GetAll(ctx, query)
}

func (s *service) AddRoleToUser(ctx context.Context, userId, roleId string) error {
	if userId == "" || roleId == "" {
		return errors.New("user ID and role ID are required")
	}

	user, err := s.GetById(ctx, userId)
	if err != nil {
		return err
	}

	// Fetch Role
	role, err := s.roleSvc.GetById(ctx, roleId)
	if err != nil {
		return err
	}

	// Skip if role already exists
	if _, exists := user.Roles[role.Id]; exists {
		return nil
	}

	addRoleToUserObj(user, *role)

	err = s.store.Update(ctx, user)
	if err != nil {
		return fmt.Errorf("failed to add role to user: %w", err)
	}

	return nil
}

// removing a role from user, handled all scenarios in it [hopefully T-T]
func (s *service) RemoveRoleFromUser(ctx context.Context, userId, roleId string) error {
	if userId == "" || roleId == "" {
		return errors.New("user ID and role ID are required")
	}

	user, err := s.GetById(ctx, userId)
	if err != nil {
		return err
	}

	// Skip if role does not exist
	if _, exists := user.Roles[roleId]; !exists {
		return nil
	}

	// Fetch Role
	role, err := s.roleSvc.GetById(ctx, roleId)
	if err != nil {
		return err
	}

	removeRoleFromUserObj(user, *role)

	// Update user in the database
	err = s.store.Update(ctx, user)
	if err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}

	return nil
}
