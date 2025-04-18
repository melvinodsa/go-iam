package policy

import (
	"context"

	"github.com/melvinodsa/go-iam/sdk"
)

type Store interface {
	GetAll(ctx context.Context) ([]sdk.Policy, error)
	Get(ctx context.Context, id string) (*sdk.Policy, error)
	Create(ctx context.Context, policy *sdk.Policy) error
	Update(ctx context.Context, policy *sdk.Policy) error
	Delete(ctx context.Context, id string) error
	GetPoliciesByRoleId(ctx context.Context, roleId string) ([]sdk.Policy, error)
	GetRolesByPolicyId(ctx context.Context, policies []string) ([]string, error)
	AddResourceToRole(ctx context.Context, roleId string, resourceId string, name string) error
	SyncResourcesByPolicyId(ctx context.Context, policies map[string]string, resourceId string, name string) error
}
