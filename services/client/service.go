package client

import (
	"context"

	"github.com/melvinodsa/go-iam/sdk"
)

type Service interface {
	GetAll(ctx context.Context, queryParams sdk.ClientQueryParams) ([]sdk.Client, error)
	Get(ctx context.Context, id string, dontCheckProjects bool) (*sdk.Client, error)
	Create(ctx context.Context, client *sdk.Client) error
	Update(ctx context.Context, client *sdk.Client) error
}
