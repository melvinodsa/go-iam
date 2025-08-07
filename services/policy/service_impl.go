package policy

import (
	"context"

	"github.com/melvinodsa/go-iam/sdk"
)

type service struct {
	store Store
}

func NewService(store Store) Service {
	return &service{
		store: store,
	}
}

func (s *service) GetAll(ctx context.Context, query sdk.PolicyQuery) (*sdk.PolicyList, error) {
	return s.store.GetAll(ctx, query)
}
