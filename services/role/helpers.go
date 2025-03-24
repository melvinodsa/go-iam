package role

import (
	"github.com/melvinodsa/go-iam/db/models"
	"github.com/melvinodsa/go-iam/sdk"
)

func fromSdkToModel(role sdk.Role) models.Role {
	return models.Role{
		Id:        role.Id,
		ProjectId: role.ProjectId,
		Name:      role.Name,
		Resources: fromSdkResourceListToModel(role.Resources),
		CreatedAt: *role.CreatedAt,
		CreatedBy: role.CreatedBy,
		UpdatedAt: *role.UpdatedAt,
		UpdatedBy: role.UpdatedBy,
	}
}

func fromModelToSdk(role *models.Role) *sdk.Role {
	return &sdk.Role{
		Id:        role.Id,
		ProjectId: role.ProjectId,
		Name:      role.Name,
		Resources: fromModelResourceListToSdk(role.Resources),
		CreatedAt: &role.CreatedAt,
		CreatedBy: role.CreatedBy,
		UpdatedAt: &role.UpdatedAt,
		UpdatedBy: role.UpdatedBy,
	}
}

func fromModelListToSdk(roles []models.Role) []sdk.Role {
	result := []sdk.Role{}
	for i := range roles {
		result = append(result, *fromModelToSdk(&roles[i]))
	}
	return result
}

func fromSdkResourceListToModel(resources []sdk.Resources) []models.Resources {
	result := []models.Resources{}
	for _, res := range resources {
		result = append(result, models.Resources{
			Id:     res.Id,
			Name:   res.Name,
			Scopes: fromSdkScopeListToModel(res.Scopes),
		})
	}
	return result
}

func fromModelResourceListToSdk(resources []models.Resources) []sdk.Resources {
	result := []sdk.Resources{}
	for _, res := range resources {
		result = append(result, sdk.Resources{
			Id:     res.Id,
			Name:   res.Name,
			Scopes: fromModelScopeListToSdk(res.Scopes),
		})
	}
	return result
}

func fromSdkScopeListToModel(scopes []sdk.Scope) []models.Scope {
	result := []models.Scope{}
	for _, scope := range scopes {
		result = append(result, models.Scope{Name: scope.Name})
	}
	return result
}

func fromModelScopeListToSdk(scopes []models.Scope) []sdk.Scope {
	result := []sdk.Scope{}
	for _, scope := range scopes {
		result = append(result, sdk.Scope{Name: scope.Name})
	}
	return result
}
