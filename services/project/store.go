package project

import (
	"context"

	"github.com/melvinodsa/go-iam/sdk"
)

type Store interface {
	GetAll(ctx context.Context) ([]sdk.Project, error)
	Get(ctx context.Context, id string) (*sdk.Project, error)
	Create(ctx context.Context, project *sdk.Project) error
	Update(ctx context.Context, project *sdk.Project) error
}
