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
		CreatedAt: user.CreatedAt,
		CreatedBy: user.CreatedBy,
		UpdatedAt: user.UpdatedAt,
		UpdatedBy: user.UpdatedBy,
	}
}

func fromModelListToSdk(users []models.User) []sdk.User {
	result := []sdk.User{}
	for i := range users {
		result = append(result, *fromModelToSdk(&users[i]))
	}
	return result
}
