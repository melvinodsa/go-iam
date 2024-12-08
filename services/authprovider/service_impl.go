package authprovider

import (
	"context"

	"github.com/melvinodsa/go-iam/sdk"
)

type service struct {
	s Store
}

func NewService(s Store) Service {
	return &service{
		s: s,
	}
}

func (s service) GetAll(ctx context.Context) ([]sdk.AuthProvider, error) {
	return s.s.GetAll(ctx)
}
func (s service) Get(ctx context.Context, id string) (*sdk.AuthProvider, error) {
	return s.s.Get(ctx, id)
}
func (s service) Create(ctx context.Context, provider *sdk.AuthProvider) error {
	return s.s.Create(ctx, provider)
}
func (s service) Update(ctx context.Context, provider *sdk.AuthProvider) error {
	return s.s.Update(ctx, provider)
}
