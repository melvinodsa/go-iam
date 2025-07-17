package user

import (
	"context"

	"github.com/melvinodsa/go-iam/middlewares/projects"
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
