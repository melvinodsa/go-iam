package client

import (
	"context"

	"github.com/melvinodsa/go-iam/sdk"
)

type Store interface {
	GetAll(ctx context.Context) ([]sdk.Client, error)
	Get(ctx context.Context, id string) (*sdk.Client, error)
	Create(ctx context.Context, client *sdk.Client) error
	Update(ctx context.Context, client *sdk.Client) error
}
