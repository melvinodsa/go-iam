package authprovider

import (
	"github.com/melvinodsa/go-iam/db/models"
	"github.com/melvinodsa/go-iam/sdk"
)

func fromModelToSdk(provider *models.AuthProvider) *sdk.AuthProvider {
	return &sdk.AuthProvider{
		Id:        provider.Id,
		Name:      provider.Name,
		Icon:      provider.Icon,
		Provider:  sdk.AuthProviderType(provider.Provider),
		Params:    fromProviderParamsModelToSdk(provider.Params),
		ProjectId: provider.ProjectId,
		Enabled:   provider.Enabled,
		CreatedAt: provider.CreatedAt,
		UpdatedAt: provider.UpdatedAt,
		CreatedBy: provider.CreatedBy,
		UpdatedBy: provider.UpdatedBy,
	}
}

func fromModelListToSdk(providers []models.AuthProvider) []sdk.AuthProvider {
	res := []sdk.AuthProvider{}
	for _, p := range providers {
		res = append(res, *fromModelToSdk(&p))
	}
	return res
}

func fromProviderParamsModelToSdk(params []models.AuthProviderParam) []sdk.AuthProviderParam {
	var res []sdk.AuthProviderParam
	for _, p := range params {
		res = append(res, sdk.AuthProviderParam{
			Label:    p.Label,
			Value:    p.Value,
			Key:      p.Key,
			IsSecret: p.IsSecret,
		})
	}
	return res
}

func fromSdkToModel(provider sdk.AuthProvider) models.AuthProvider {
	return models.AuthProvider{
		Id:        provider.Id,
		Name:      provider.Name,
		Icon:      provider.Icon,
		Provider:  models.AuthProviderType(provider.Provider),
		Params:    fromProviderParamsSdkToModel(provider.Params),
		ProjectId: provider.ProjectId,
		Enabled:   provider.Enabled,
		CreatedAt: provider.CreatedAt,
		UpdatedAt: provider.UpdatedAt,
		CreatedBy: provider.CreatedBy,
		UpdatedBy: provider.UpdatedBy,
	}
}

func fromProviderParamsSdkToModel(params []sdk.AuthProviderParam) []models.AuthProviderParam {
	var res []models.AuthProviderParam
	for _, p := range params {
		res = append(res, models.AuthProviderParam{
			Label:    p.Label,
			Value:    p.Value,
			Key:      p.Key,
			IsSecret: p.IsSecret,
		})
	}
	return res
}
