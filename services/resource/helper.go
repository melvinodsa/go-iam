package resource

import (
	"github.com/melvinodsa/go-iam/db/models"
	"github.com/melvinodsa/go-iam/sdk"
)

func fromModelToSdk(m *models.Resource) *sdk.Resource {
	return &sdk.Resource{
		ID:          m.ID,
		Name:        m.Name,
		Description: m.Description,
		Key:         m.Key,
		Enabled:     m.Enabled,
		CreatedAt:   m.CreatedAt,
		CreatedBy:   m.CreatedBy,
		UpdatedAt:   m.UpdatedAt,
		UpdatedBy:   m.UpdatedBy,
		DeletedAt:   m.DeletedAt,
	}
}

func fromModelListToSdk(models []models.Resource) []sdk.Resource {
	resources := make([]sdk.Resource, len(models))
	for i, m := range models {
		r := fromModelToSdk(&m)
		resources[i] = *r
	}
	return resources
}

func fromSdkToModel(s sdk.Resource) *models.Resource {
	return &models.Resource{
		ID:          s.ID,
		Name:        s.Name,
		Description: s.Description,
		Key:         s.Key,
		Enabled:     s.Enabled,
		CreatedAt:   s.CreatedAt,
		CreatedBy:   s.CreatedBy,
		UpdatedAt:   s.UpdatedAt,
		UpdatedBy:   s.UpdatedBy,
		DeletedAt:   s.DeletedAt,
	}
}
