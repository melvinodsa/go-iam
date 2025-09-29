package migrations

import (
	"testing"
	"time"

	"github.com/melvinodsa/go-iam/db/models"
	"github.com/stretchr/testify/assert"
)

func TestOldUserStruct(t *testing.T) {
	t.Run("OldUser struct initialization", func(t *testing.T) {
		now := time.Now()
		oldUser := OldUser{
			Id:        "user-123",
			ProjectId: "project-456",
			Name:      "Test User",
			Email:     "test@example.com",
			Enabled:   true,
			CreatedAt: &now,
			Policies:  map[string]string{"policy1": "value1"},
			Roles:     map[string]models.UserRoles{"role1": {Id: "role1", Name: "Test Role"}},
			Resources: map[string]models.UserResource{"res1": {Key: "res1", Name: "Test Resource"}},
		}

		assert.Equal(t, "user-123", oldUser.Id)
		assert.Equal(t, "project-456", oldUser.ProjectId)
		assert.Equal(t, "Test User", oldUser.Name)
		assert.Equal(t, "test@example.com", oldUser.Email)
		assert.True(t, oldUser.Enabled)
		assert.Equal(t, &now, oldUser.CreatedAt)
		assert.Equal(t, map[string]string{"policy1": "value1"}, oldUser.Policies)
		assert.Len(t, oldUser.Roles, 1)
		assert.Len(t, oldUser.Resources, 1)
	})

	t.Run("OldUser with nil time fields", func(t *testing.T) {
		oldUser := OldUser{
			Id:        "user-456",
			ProjectId: "project-789",
			Name:      "Another User",
			Email:     "another@example.com",
			Enabled:   false,
			Expiry:    nil,
			CreatedAt: nil,
			UpdatedAt: nil,
		}

		assert.Equal(t, "user-456", oldUser.Id)
		assert.False(t, oldUser.Enabled)
		assert.Nil(t, oldUser.Expiry)
		assert.Nil(t, oldUser.CreatedAt)
		assert.Nil(t, oldUser.UpdatedAt)
	})
}