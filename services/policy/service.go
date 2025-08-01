package policy

import (
	"context"

	"github.com/melvinodsa/go-iam/sdk"
)

type Service interface {
	GetAll(ctx context.Context) ([]sdk.Project, error)
}
