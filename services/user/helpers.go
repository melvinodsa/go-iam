package user

import (
	"github.com/melvinodsa/go-iam/db/models"
	"github.com/melvinodsa/go-iam/sdk"
)

func fromSdkToModel(user sdk.User) models.User {
	return models.User{
		Id:        user.Id,
		Email:     user.Email,
		Phone:     user.Phone,
		Name:      user.Name,
		ProjectId: user.ProjectId,
		Enabled:   user.Enabled,
		Expiry:    user.Expiry,
		Roles:     fromSdkUserRoleMapToModel(user.Roles),
		Resources: fromSdkUserResourceMapToModel(user.Resources),
		Policies:  user.Policies,
		CreatedAt: user.CreatedAt,
		CreatedBy: user.CreatedBy,
		UpdatedAt: user.UpdatedAt,
		UpdatedBy: user.UpdatedBy,
	}
}

func fromModelToSdk(user *models.User) *sdk.User {
	return &sdk.User{
		Id:        user.Id,
		Email:     user.Email,
		Phone:     user.Phone,
		Name:      user.Name,
		ProjectId: user.ProjectId,
		Expiry:    user.Expiry,
		Enabled:   user.Enabled,
		Roles:     fromModelUserRoleMapToSdk(user.Roles),
		Resources: fromModelUserResourceMapToSdk(user.Resources),
		Policies:  user.Policies,
		CreatedAt: user.CreatedAt,
		CreatedBy: user.CreatedBy,
		UpdatedAt: user.UpdatedAt,
		UpdatedBy: user.UpdatedBy,
	}
}

// Convert SDK UserRole map to Model UserRoles map (Key: Name)
func fromSdkUserRoleMapToModel(roles map[string]sdk.UserRole) map[string]models.UserRoles {
	userRoles := make(map[string]models.UserRoles)
	for key, role := range roles {
		userRoles[key] = models.UserRoles{
			Name: role.Name,
			Id:   role.Id,
		}
	}
	return userRoles
}

// Convert Model UserRoles map to SDK UserRole map (Key: Name)
func fromModelUserRoleMapToSdk(roles map[string]models.UserRoles) map[string]sdk.UserRole {
	userRoles := make(map[string]sdk.UserRole)
	for key, role := range roles {
		userRoles[key] = sdk.UserRole{
			Name: role.Name,
			Id:   role.Id,
		}
	}
	return userRoles
}

// Convert SDK UserResource map to Model UserResource map (Key: Key)
func fromSdkUserResourceMapToModel(resources map[string]sdk.UserResource) map[string]models.UserResource {
	userResources := make(map[string]models.UserResource)
	for key, res := range resources {
		userResources[key] = models.UserResource{
			Id:   res.Id,
			Key:  res.Key,
			Name: res.Name,
		}
	}
	return userResources
}

// Convert Model UserResource map to SDK UserResource map (Key: Key)
func fromModelUserResourceMapToSdk(resources map[string]models.UserResource) map[string]sdk.UserResource {
	userResources := make(map[string]sdk.UserResource)
	for key, res := range resources {
		userResources[key] = sdk.UserResource{
			Id:   res.Id,
			Key:  res.Key,
			Name: res.Name,
		}
	}
	return userResources
}

// Convert list of Model Users to list of SDK Users
func fromModelListToSdk(users []models.User) []sdk.User {
	result := []sdk.User{}
	for i := range users {
		result = append(result, *fromModelToSdk(&users[i]))
	}
	return result
}
