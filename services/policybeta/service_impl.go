package policybeta

import (
	"context"

	"github.com/melvinodsa/go-iam/sdk"
)

type service struct {
	s Store
}

func NewService(s Store) Service {
	return &service{
		s: s,
	}
}

func (s service) GetAll(ctx context.Context) ([]sdk.Policy, error) {
	return s.s.GetAll(ctx)
}

func (s service) Get(ctx context.Context, id string) (*sdk.Policy, error) {
	if len(id) == 0 {
		return nil, ErrPolicyNotFound
	}
	return s.s.Get(ctx, id)
}

func (s service) Create(ctx context.Context, policy *sdk.Policy) error {
	return s.s.Create(ctx, policy)
}

func (s service) Update(ctx context.Context, policy *sdk.Policy) error {
	return s.s.Update(ctx, policy)
}

func (s service) Delete(ctx context.Context, id string) error {
	return s.s.Delete(ctx, id)
}

func (s service) GetPoliciesByRoleId(ctx context.Context, roleId string) ([]sdk.Policy, error) {
	return s.s.GetPoliciesByRoleId(ctx, roleId)
}

func (s service) SyncResourcesbyPolicyId(ctx context.Context, policies map[string]string, resourceId string, name string) error {
	return s.s.SyncResourcesByPolicyId(ctx, policies, resourceId, name)
}
