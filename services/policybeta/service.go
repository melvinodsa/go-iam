package policybeta

import (
	"context"

	"github.com/melvinodsa/go-iam/sdk"
)

type Service interface {
	GetAll(ctx context.Context) ([]sdk.Policy, error)
	Get(ctx context.Context, id string) (*sdk.Policy, error)
	Create(ctx context.Context, policy *sdk.Policy) error
	Update(ctx context.Context, policy *sdk.Policy) error
	Delete(ctx context.Context, id string) error
	GetPoliciesByRoleId(ctx context.Context, roleId string) ([]sdk.Policy, error)
	SyncResourcesbyPolicyId(ctx context.Context, policyId map[string]string, ResourceId string, name string) error
}
