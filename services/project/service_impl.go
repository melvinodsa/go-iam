package project

import (
	"context"

	"github.com/melvinodsa/go-iam/sdk"
)

type service struct {
	s Store
}

func NewService(s Store) Service {
	return service{s: s}
}

func (s service) GetAll(ctx context.Context) ([]sdk.Project, error) {
	return s.s.GetAll(ctx)
}

func (s service) GetByName(ctx context.Context, name string) (*sdk.Project, error) {
	if len(name) == 0 {
		return nil, sdk.ErrProjectNotFound
	}
	return s.s.GetByName(ctx, name)
}

func (s service) Get(ctx context.Context, id string) (*sdk.Project, error) {
	if len(id) == 0 {
		return nil, sdk.ErrProjectNotFound
	}
	return s.s.Get(ctx, id)
}

func (s service) Create(ctx context.Context, project *sdk.Project) error {
	return s.s.Create(ctx, project)
}

func (s service) Update(ctx context.Context, project *sdk.Project) error {
	return s.s.Update(ctx, project)
}
