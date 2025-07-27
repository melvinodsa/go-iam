package policybeta

import (
	"github.com/melvinodsa/go-iam/db/models"
	"github.com/melvinodsa/go-iam/sdk"
)

func fromSdkToModel(policy *sdk.PolicyBeta) *models.Policy {
	return &models.Policy{
		Id:          policy.Id,
		Name:        policy.Name,
		Roles:       policy.Roles,
		Description: policy.Description,
		CreatedBy:   policy.CreatedBy,
	}
}

func fromModelToSdk(policy *models.Policy) *sdk.PolicyBeta {
	return &sdk.PolicyBeta{
		Id:          policy.Id,
		Name:        policy.Name,
		Roles:       policy.Roles,
		Description: policy.Description,
		CreatedBy:   policy.CreatedBy,
	}
}

func fromModelListToSdk(policies []models.Policy) []sdk.PolicyBeta {
	var sdkPolicies []sdk.PolicyBeta
	for i := range policies {
		sdkPolicies = append(sdkPolicies, *fromModelToSdk(&policies[i]))
	}
	return sdkPolicies
}
