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
		Roles:     fromSdkUserRoleListToModel(user.Roles),
		Resource:  fromSdkUserResourceListToModel(user.Resource),
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
		Roles:     fromModelUserRoleListToSdk(user.Roles),
		Resource:  fromModelUserResourceListToSdk(user.Resource),
		CreatedAt: user.CreatedAt,
		CreatedBy: user.CreatedBy,
		UpdatedAt: user.UpdatedAt,
		UpdatedBy: user.UpdatedBy,
	}
}

func fromSdkUserRoleListToModel(roles []sdk.UserRole) []models.UserRoles {
	var userRoles []models.UserRoles
	for _, role := range roles {
		userRoles = append(userRoles, models.UserRoles{
			Name: role.Name,
			Id:   role.Id,
		})
	}
	return userRoles
}

func fromModelUserRoleListToSdk(roles []models.UserRoles) []sdk.UserRole {
	var userRoles []sdk.UserRole
	for _, role := range roles {
		userRoles = append(userRoles, sdk.UserRole{
			Name: role.Name,
			Id:   role.Id,
		})
	}
	return userRoles
}

func fromSdkUserResourceListToModel(resources []sdk.UserResource) []models.UserResource {
	var userResources []models.UserResource
	for _, res := range resources {
		userResources = append(userResources, models.UserResource{
			Key:   res.Key,
			Name:  res.Name,
			Scope: res.Scope,
		})
	}
	return userResources
}

func fromModelUserResourceListToSdk(resources []models.UserResource) []sdk.UserResource {
	var userResources []sdk.UserResource
	for _, res := range resources {
		userResources = append(userResources, sdk.UserResource{
			Key:   res.Key,
			Name:  res.Name,
			Scope: res.Scope,
		})
	}
	return userResources
}

func fromModelListToSdk(users []models.User) []sdk.User {
	var result []sdk.User
	for i := range users {
		result = append(result, *fromModelToSdk(&users[i]))
	}
	return result
}
