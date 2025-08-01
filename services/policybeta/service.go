package policybeta

import (
	"context"

	"github.com/melvinodsa/go-iam/sdk"
)

type Service interface {
	GetAll(ctx context.Context) ([]sdk.PolicyBeta, error)
	Get(ctx context.Context, id string) (*sdk.PolicyBeta, error)
	Create(ctx context.Context, policy *sdk.PolicyBeta) error
	Update(ctx context.Context, policy *sdk.PolicyBeta) error
	Delete(ctx context.Context, id string) error
	GetPoliciesByRoleId(ctx context.Context, roleId string) ([]sdk.PolicyBeta, error)
	SyncResourcesbyPolicyId(ctx context.Context, policyId map[string]string, ResourceId string, name string) error
}
