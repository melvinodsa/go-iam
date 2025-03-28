package resource

import (
	"context"

	"github.com/melvinodsa/go-iam/sdk"
)

type Service interface {
	Search(ctx context.Context, query sdk.ResourceQuery) (*sdk.ResourceList, error)
	Get(ctx context.Context, id string) (*sdk.Resource, error)
	Create(ctx context.Context, resource *sdk.Resource) error
	Update(ctx context.Context, resource *sdk.Resource) error
	Delete(ctx context.Context, id string) error
}
