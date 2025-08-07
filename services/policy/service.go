package policy

import (
	"context"

	"github.com/melvinodsa/go-iam/sdk"
)

type Service interface {
	GetAll(ctx context.Context, query sdk.PolicyQuery) (*sdk.PolicyList, error)
}
