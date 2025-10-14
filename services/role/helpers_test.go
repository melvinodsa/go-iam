package role

import (
	"testing"
	"time"

	"github.com/melvinodsa/go-iam/db/models"
	"github.com/melvinodsa/go-iam/sdk"
	"github.com/stretchr/testify/assert"
)

func TestFromSdkToModel_Role(t *testing.T) {
	t.Run("converts fully populated role", func(t *testing.T) {
		createdAt := time.Now().Add(-time.Hour)
		updatedAt := time.Now()

		sdkRole := sdk.Role{
			Id:          "role-1",
			ProjectId:   "project-1",
			Name:        "Admin",
			Description: "Admin role",
			Enabled:     true,
			Resources: map[string]sdk.Resources{
				"users":    {Id: "res-1", Key: "users", Name: "Users"},
				"projects": {Id: "res-2", Key: "projects", Name: "Projects"},
			},
			CreatedAt: &createdAt,
			CreatedBy: "creator",
			UpdatedAt: &updatedAt,
			UpdatedBy: "updater",
		}

		modelRole := fromSdkToModel(sdkRole)

		assert.Equal(t, sdkRole.Id, modelRole.Id)
		assert.Equal(t, sdkRole.ProjectId, modelRole.ProjectId)
		assert.Equal(t, sdkRole.Name, modelRole.Name)
		assert.Equal(t, sdkRole.Description, modelRole.Description)
		assert.Equal(t, sdkRole.Enabled, modelRole.Enabled)
		// created/updated timestamps should not be nil (zero value allowed)
		assert.Equal(t, createdAt, modelRole.CreatedAt)
		assert.Equal(t, updatedAt, modelRole.UpdatedAt)

		// resource map conversion
		assert.Equal(t, 2, len(modelRole.Resources))
		assert.Equal(t, models.Resources{Id: "res-1", Key: "users", Name: "Users"}, modelRole.Resources["users"])
		assert.Equal(t, models.Resources{Id: "res-2", Key: "projects", Name: "Projects"}, modelRole.Resources["projects"])
	})
}

func TestFromModelToSdk_Role(t *testing.T) {
	t.Run("converts fully populated role model", func(t *testing.T) {
		createdAt := time.Now().Add(-2 * time.Hour)
		updatedAt := time.Now()

		modelRole := &models.Role{
			Id:          "role-2",
			ProjectId:   "project-2",
			Name:        "Viewer",
			Description: "Viewer role",
			Enabled:     false,
			Resources: map[string]models.Resources{
				"resources": {Id: "res-3", Key: "resources", Name: "Resources"},
			},
			CreatedAt: createdAt,
			CreatedBy: "creator-2",
			UpdatedAt: updatedAt,
			UpdatedBy: "updater-2",
		}

		sdkRole := fromModelToSdk(modelRole)

		assert.NotNil(t, sdkRole)
		assert.Equal(t, modelRole.Id, sdkRole.Id)
		assert.Equal(t, modelRole.ProjectId, sdkRole.ProjectId)
		assert.Equal(t, modelRole.Name, sdkRole.Name)
		assert.Equal(t, modelRole.Description, sdkRole.Description)
		assert.Equal(t, modelRole.Enabled, sdkRole.Enabled)
		// timestamps converted to pointers
		assert.NotNil(t, sdkRole.CreatedAt)
		assert.Equal(t, modelRole.CreatedAt, *sdkRole.CreatedAt)
		assert.NotNil(t, sdkRole.UpdatedAt)
		assert.Equal(t, modelRole.UpdatedAt, *sdkRole.UpdatedAt)

		// resource map conversion
		assert.Equal(t, 1, len(sdkRole.Resources))
		assert.Equal(t, sdk.Resources{Id: "res-3", Key: "resources", Name: "Resources"}, sdkRole.Resources["resources"])
	})

	t.Run("handles nil model pointer", func(t *testing.T) {
		var modelRole *models.Role
		result := fromModelToSdk(modelRole)
		assert.Nil(t, result)
	})
}

func TestFromModelListToSdk_Roles(t *testing.T) {
	t.Run("converts list with multiple items", func(t *testing.T) {
		now := time.Now()
		modelsList := []models.Role{
			{Id: "r1", Name: "A", CreatedAt: now, UpdatedAt: now},
			{Id: "r2", Name: "B", CreatedAt: now, UpdatedAt: now},
		}

		result := fromModelListToSdk(modelsList)

		assert.Equal(t, 2, len(result))
		assert.Equal(t, "r1", result[0].Id)
		assert.Equal(t, "A", result[0].Name)
		assert.Equal(t, "r2", result[1].Id)
		assert.Equal(t, "B", result[1].Name)
	})
}

func TestResourceMapConversions_Role(t *testing.T) {
	t.Run("fromSdkResourceMapToModel filters empty keys", func(t *testing.T) {
		src := map[string]sdk.Resources{
			"":      {Id: "x", Key: "", Name: "Ignore"},
			"users": {Id: "1", Key: "users", Name: "Users"},
		}
		dst := fromSdkResourceMapToModel(src)
		assert.Equal(t, 1, len(dst))
		_, ok := dst[""]
		assert.False(t, ok)
		assert.Equal(t, models.Resources{Id: "1", Key: "users", Name: "Users"}, dst["users"])
	})

	t.Run("fromModelResourceMapToSdk preserves keyed entries and filters empty keys", func(t *testing.T) {
		src := map[string]models.Resources{
			"":         {Id: "x", Key: "", Name: "Ignore"},
			"projects": {Id: "2", Key: "projects", Name: "Projects"},
		}
		dst := fromModelResourceMapToSdk(src)
		assert.Equal(t, 1, len(dst))
		_, ok := dst[""]
		assert.False(t, ok)
		assert.Equal(t, sdk.Resources{Id: "2", Key: "projects", Name: "Projects"}, dst["projects"])
	})
}
