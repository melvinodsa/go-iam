package withpassword

import (
	"github.com/melvinodsa/go-iam/db/models"
	"github.com/melvinodsa/go-iam/sdk"
)

func fromModelToSdk(user *models.WithPasswordUser) *sdk.WithPasswordUser {
	if user == nil {
		return nil
	}
	return &sdk.WithPasswordUser{
		ProjectID: user.ProjectID,
		Email:     user.Email,
		Password:  user.Password,
		CreatedAt: user.CreatedAt,
		CreatedBy: user.CreatedBy,
		UpdatedAt: user.UpdatedAt,
		UpdatedBy: user.UpdatedBy,
	}
}
