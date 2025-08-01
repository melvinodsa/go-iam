package authprovider

import (
	"context"

	"github.com/melvinodsa/go-iam/sdk"
)

type Service interface {
	GetAll(ctx context.Context, params sdk.AuthProviderQueryParams) ([]sdk.AuthProvider, error)
	Get(ctx context.Context, id string, dontCheckProjects bool) (*sdk.AuthProvider, error)
	Create(ctx context.Context, provider *sdk.AuthProvider) error
	Update(ctx context.Context, provider *sdk.AuthProvider) error
	GetProvider(ctx context.Context, v sdk.AuthProvider) (sdk.ServiceProvider, error)
}
