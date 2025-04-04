package role

import (
	"time"

	"github.com/melvinodsa/go-iam/db/models"
	"github.com/melvinodsa/go-iam/sdk"
)

// Conversion functions with added input validation
func fromSdkToModel(role sdk.Role) models.Role {
	// Use zero time if pointers are nil to prevent nil pointer dereference
	createdAt := time.Time{}
	updatedAt := time.Time{}

	if role.CreatedAt != nil {
		createdAt = *role.CreatedAt
	}
	if role.UpdatedAt != nil {
		updatedAt = *role.UpdatedAt
	}

	return models.Role{
		Id:        role.Id,
		ProjectId: role.ProjectId,
		Name:      role.Name,
		Resources: fromSdkResourceMapToModel(role.Resources),
		CreatedAt: createdAt,
		CreatedBy: role.CreatedBy,
		UpdatedAt: updatedAt,
		UpdatedBy: role.UpdatedBy,
	}
}

func fromModelToSdk(role *models.Role) *sdk.Role {
	if role == nil {
		return nil
	}

	return &sdk.Role{
		Id:        role.Id,
		ProjectId: role.ProjectId,
		Name:      role.Name,
		Resources: fromModelResourceMapToSdk(role.Resources),
		CreatedAt: &role.CreatedAt,
		CreatedBy: role.CreatedBy,
		UpdatedAt: &role.UpdatedAt,
		UpdatedBy: role.UpdatedBy,
	}
}

func fromModelListToSdk(roles []models.Role) []sdk.Role {
	result := make([]sdk.Role, 0, len(roles))
	for i := range roles {
		sdkRole := fromModelToSdk(&roles[i])
		if sdkRole != nil {
			result = append(result, *sdkRole)
		}
	}
	return result
}

// More robust resource map conversion functions
func fromSdkResourceMapToModel(resources map[string]sdk.Resources) map[string]models.Resources {
	result := make(map[string]models.Resources)
	for _, res := range resources {
		// Only add resources with non-empty keys
		if res.Key != "" {
			result[res.Key] = models.Resources{
				Id:   res.Id,
				Key:  res.Key,
				Name: res.Name,
			}
		}
	}
	return result
}

func fromModelResourceMapToSdk(resources map[string]models.Resources) map[string]sdk.Resources {
	result := make(map[string]sdk.Resources)
	for key, res := range resources {
		// Only add resources with non-empty keys
		if key != "" {
			result[key] = sdk.Resources{
				Id:   res.Id,
				Key:  res.Key,
				Name: res.Name,
			}
		}
	}
	return result
}
