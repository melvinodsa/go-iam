package authprovider

import (
	"context"

	"github.com/melvinodsa/go-iam/sdk"
)

type Service interface {
	GetAll(ctx context.Context) ([]sdk.AuthProvider, error)
	Get(ctx context.Context, id string) (*sdk.AuthProvider, error)
	Create(ctx context.Context, provider *sdk.AuthProvider) error
	Update(ctx context.Context, provider *sdk.AuthProvider) error
	GetProvider(v sdk.AuthProvider) (sdk.ServiceProvider, error)
}
