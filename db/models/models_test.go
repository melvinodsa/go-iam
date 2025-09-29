package models

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUserModel(t *testing.T) {
	t.Run("Name returns correct collection name", func(t *testing.T) {
		m := GetUserModel()
		assert.Equal(t, "users", m.Name())
	})

	t.Run("GetUserModel returns correct field keys", func(t *testing.T) {
		m := GetUserModel()
		assert.Equal(t, "id", m.IdKey)
		assert.Equal(t, "name", m.NameKey)
		assert.Equal(t, "email", m.EmailKey)
		assert.Equal(t, "project_id", m.ProjectIDKey)
	})
}

func TestRoleModel(t *testing.T) {
	t.Run("Name returns correct collection name", func(t *testing.T) {
		m := GetRoleModel()
		assert.Equal(t, "roles", m.Name())
	})

	t.Run("GetRoleModel returns correct field keys", func(t *testing.T) {
		m := GetRoleModel()
		assert.Equal(t, "id", m.IdKey)
		assert.Equal(t, "name", m.NameKey)
		assert.Equal(t, "project_id", m.ProjectIdKey)
		assert.Equal(t, "resources", m.ResourcesKey)
	})
}

func TestPolicyModel(t *testing.T) {
	t.Run("Name returns correct collection name", func(t *testing.T) {
		m := GetPolicyModel()
		assert.Equal(t, "policies", m.Name())
	})

	t.Run("GetPolicyModel returns correct field keys", func(t *testing.T) {
		m := GetPolicyModel()
		assert.Equal(t, "id", m.IdKey)
		assert.Equal(t, "name", m.NameKey)
		assert.Equal(t, "roles", m.RolesKey)
		assert.Equal(t, "description", m.DescriptionKey)
	})
}

func TestProjectModel(t *testing.T) {
	t.Run("Name returns correct collection name", func(t *testing.T) {
		m := GetProjectModel()
		assert.Equal(t, "projects", m.Name())
	})

	t.Run("GetProjectModel returns correct field keys", func(t *testing.T) {
		m := GetProjectModel()
		assert.Equal(t, "id", m.IdKey)
		assert.Equal(t, "name", m.NameKey)
		assert.Equal(t, "tags", m.TagsKey)
		assert.Equal(t, "description", m.DescriptionKey)
	})
}

func TestResourceModel(t *testing.T) {
	t.Run("Name returns correct collection name", func(t *testing.T) {
		m := GetResourceModel()
		assert.Equal(t, "resources", m.Name())
	})

	t.Run("GetResourceModel returns correct field keys", func(t *testing.T) {
		m := GetResourceModel()
		assert.Equal(t, "id", m.IdKey)
		assert.Equal(t, "name", m.NameKey)
		assert.Equal(t, "key", m.KeyKey)
		assert.Equal(t, "project_id", m.ProjectIdKey)
	})
}

func TestClientModel(t *testing.T) {
	t.Run("Name returns correct collection name", func(t *testing.T) {
		m := GetClientModel()
		assert.Equal(t, "clients", m.Name())
	})

	t.Run("GetClientModel returns correct field keys", func(t *testing.T) {
		m := GetClientModel()
		assert.Equal(t, "id", m.IdKey)
		assert.Equal(t, "name", m.NameKey)
		assert.Equal(t, "project_id", m.ProjectIdKey)
		assert.Equal(t, "go_iam_client", m.GoIamClientKey)
	})
}

func TestAuthProviderModel(t *testing.T) {
	t.Run("Name returns correct collection name", func(t *testing.T) {
		m := GetAuthProviderModel()
		assert.Equal(t, "auth_providers", m.Name())
	})

	t.Run("GetAuthProviderModel returns correct field keys", func(t *testing.T) {
		m := GetAuthProviderModel()
		assert.Equal(t, "id", m.IdKey)
		assert.Equal(t, "name", m.NameKey)
		assert.Equal(t, "provider", m.ProviderKey)
		assert.Equal(t, "project_id", m.ProjectIdKey)
	})
}

func TestMigrationModel(t *testing.T) {
	t.Run("Name returns correct collection name", func(t *testing.T) {
		m := GetMigrationModel()
		assert.Equal(t, "migrations", m.Name())
	})

	t.Run("GetMigrationModel returns correct field keys", func(t *testing.T) {
		m := GetMigrationModel()
		assert.Equal(t, "version", m.VersionKey)
		assert.Equal(t, "name", m.NameKey)
		assert.Equal(t, "applied_at", m.AppliedAtKey)
	})
}

func TestAllModelsDbName(t *testing.T) {
	t.Run("All models return correct database name", func(t *testing.T) {
		models := []interface{ DbName() string }{
			GetUserModel(),
			GetRoleModel(),
			GetPolicyModel(),
			GetProjectModel(),
			GetResourceModel(),
			GetClientModel(),
			GetAuthProviderModel(),
			GetMigrationModel(),
		}

		for _, model := range models {
			assert.Equal(t, "iam", model.DbName())
		}
	})
}

func TestModelStructs(t *testing.T) {
	t.Run("User struct initialization", func(t *testing.T) {
		user := User{
			Id:        "user-123",
			ProjectId: "project-456",
			Name:      "Test User",
			Email:     "test@example.com",
			Enabled:   true,
		}

		assert.Equal(t, "user-123", user.Id)
		assert.Equal(t, "project-456", user.ProjectId)
		assert.Equal(t, "Test User", user.Name)
		assert.Equal(t, "test@example.com", user.Email)
		assert.True(t, user.Enabled)
	})

	t.Run("UserRoles struct initialization", func(t *testing.T) {
		userRole := UserRoles{
			Id:   "role-123",
			Name: "Test Role",
		}

		assert.Equal(t, "role-123", userRole.Id)
		assert.Equal(t, "Test Role", userRole.Name)
	})

	t.Run("UserResource struct initialization", func(t *testing.T) {
		userResource := UserResource{
			RoleIds:   map[string]bool{"role1": true},
			PolicyIds: map[string]bool{"policy1": true},
			Key:       "resource-key",
			Name:      "Test Resource",
		}

		assert.Equal(t, map[string]bool{"role1": true}, userResource.RoleIds)
		assert.Equal(t, map[string]bool{"policy1": true}, userResource.PolicyIds)
		assert.Equal(t, "resource-key", userResource.Key)
		assert.Equal(t, "Test Resource", userResource.Name)
	})

	t.Run("UserPolicy struct initialization", func(t *testing.T) {
		userPolicy := UserPolicy{
			Name: "Test Policy",
			Mapping: UserPolicyMapping{
				Arguments: map[string]UserPolicyMappingValue{
					"arg1": {Static: "value1"},
				},
			},
		}

		assert.Equal(t, "Test Policy", userPolicy.Name)
		assert.Equal(t, "value1", userPolicy.Mapping.Arguments["arg1"].Static)
	})
}