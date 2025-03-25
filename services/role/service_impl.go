package role

import (
	"context"

	"github.com/melvinodsa/go-iam/sdk"
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
	return s.store.Update(ctx, role)
}

func (s *service) GetById(ctx context.Context, id string) (*sdk.Role, error) {
	return s.store.GetById(ctx, id)
}

func (s *service) GetAll(ctx context.Context, query sdk.RoleQuery) ([]sdk.Role, error) {
	return s.store.GetAll(ctx, query)
}

func (s *service) AddRoleToUser(ctx context.Context, userId, roleId string) error {
	return s.store.AddRoleToUser(ctx, userId, roleId)
}

func (s *service) RemoveRoleFromUser(ctx context.Context, userId, roleId string) error {
	return s.store.RemoveRoleFromUser(ctx, userId, roleId)
}
