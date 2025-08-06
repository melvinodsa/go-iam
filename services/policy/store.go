package policy

import (
	"context"

	"github.com/melvinodsa/go-iam/sdk"
)

type Store interface {
	GetAll(ctx context.Context, query sdk.PolicyQuery) (*sdk.PolicyList, error)
}
