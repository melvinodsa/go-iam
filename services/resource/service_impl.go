package resource

import (
	"context"

	"github.com/melvinodsa/go-iam/middlewares"
	"github.com/melvinodsa/go-iam/sdk"
)

type service struct {
	s Store
}

func NewService(s Store) Service {
	return service{s: s}
}

func (s service) Search(ctx context.Context, query sdk.ResourceQuery) (*sdk.ResourceList, error) {
	query.ProjectIds = middlewares.GetProjects(ctx)
	return s.s.Search(ctx, query)
}

func (s service) Get(ctx context.Context, id string) (*sdk.Resource, error) {
	return s.s.Get(ctx, id)
}

func (s service) Create(ctx context.Context, resource *sdk.Resource) error {
	_, err := s.s.Create(ctx, resource)
	return err
}

func (s service) Update(ctx context.Context, resource *sdk.Resource) error {
	return s.s.Update(ctx, resource)
}

func (s service) Delete(ctx context.Context, id string) error {
	return s.s.Delete(ctx, id)
}
