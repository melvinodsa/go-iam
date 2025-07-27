package policybeta

import (
	"context"

	"github.com/melvinodsa/go-iam/sdk"
)

type Store interface {
	GetAll(ctx context.Context) ([]sdk.PolicyBeta, error)
	Get(ctx context.Context, id string) (*sdk.PolicyBeta, error)
	Create(ctx context.Context, policy *sdk.PolicyBeta) error
	Update(ctx context.Context, policy *sdk.PolicyBeta) error
	Delete(ctx context.Context, id string) error
	GetPoliciesByRoleId(ctx context.Context, roleId string) ([]sdk.PolicyBeta, error)
	GetRolesByPolicyId(ctx context.Context, policies []string) ([]string, error)
	AddResourceToRole(ctx context.Context, roleId string, resourceId string, name string) error
	SyncResourcesByPolicyId(ctx context.Context, policies map[string]string, resourceId string, name string) error
}
