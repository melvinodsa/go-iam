package authprovider

import (
	"context"
	"fmt"

	"github.com/melvinodsa/go-iam/sdk"
	"github.com/melvinodsa/go-iam/services/authprovider/google"
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

func (s service) GetProvider(v sdk.AuthProvider) (sdk.ServiceProvider, error) {
	switch v.Provider {
	case sdk.AuthProviderTypeGoogle:
		return google.NewAuthProvider(v), nil
	default:
		return nil, fmt.Errorf("unknown auth provider: %s", v.Provider)
	}
}
