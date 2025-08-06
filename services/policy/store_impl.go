package policy

import (
	"context"

	"github.com/melvinodsa/go-iam/sdk"
	"github.com/melvinodsa/go-iam/services/policy/system"
)

type storeImpl struct{}

func NewStore() Store {
	return &storeImpl{}
}

func (s *storeImpl) GetAll(ctx context.Context, query sdk.PolicyQuery) (*sdk.PolicyList, error) {
	// This is a stub implementation. Replace with actual logic to fetch policies.
	policies := []sdk.Policy{
		system.NewAccessToCreatedResource(nil).PolicyDef(), // Assuming user service is nil for stub
		// Add other policies as needed
		system.NewAddResourcesToRole(nil, nil).PolicyDef(),
		system.NewAddResourcesToUser(nil).PolicyDef(),
	}

	limit := query.Limit
	skip := query.Skip
	if limit <= 0 {
		limit = 10 // Default limit
	}
	if skip < 0 {
		skip = 0 // Ensure skip is not negative
	}
	if skip+limit > int64(len(policies)) {
		limit = int64(len(policies)) - skip // Adjust limit if it exceeds available policies
	}
	if limit < 0 {
		limit = 0 // Ensure limit is not negative
	}

	result := policies[skip : limit+skip]

	return &sdk.PolicyList{
		Policies: result,
		Total:    len(policies),
		Skip:     skip,
		Limit:    limit,
	}, nil
}
