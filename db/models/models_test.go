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