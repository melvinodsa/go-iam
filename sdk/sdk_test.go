package sdk

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
)

func TestUser(t *testing.T) {
	t.Run("User struct initialization", func(t *testing.T) {
		now := time.Now()
		user := User{
			Id:        "user-123",
			ProjectId: "project-456",
			Name:      "Test User",
			Email:     "test@example.com",
			Enabled:   true,
			CreatedAt: &now,
		}

		assert.Equal(t, "user-123", user.Id)
		assert.Equal(t, "project-456", user.ProjectId)
		assert.Equal(t, "Test User", user.Name)
		assert.Equal(t, "test@example.com", user.Email)
		assert.True(t, user.Enabled)
		assert.Equal(t, &now, user.CreatedAt)
	})
}

func TestClient(t *testing.T) {
	t.Run("IsServiceAccount returns true for Go-IAM client with linked user", func(t *testing.T) {
		client := Client{
			GoIamClient:  true,
			LinkedUserId: "user-123",
		}
		assert.True(t, client.IsServiceAccount())
	})

	t.Run("IsServiceAccount returns false for client without linked user", func(t *testing.T) {
		client := Client{
			GoIamClient: true,
		}
		assert.False(t, client.IsServiceAccount())
	})

	t.Run("HasGoIamAuthProvider returns true when DefaultAuthProviderId is empty", func(t *testing.T) {
		client := Client{}
		assert.True(t, client.HasGoIamAuthProvider())
	})

	t.Run("HasGoIamAuthProvider returns false when DefaultAuthProviderId is set", func(t *testing.T) {
		client := Client{
			DefaultAuthProviderId: "provider-123",
		}
		assert.False(t, client.HasGoIamAuthProvider())
	})
}

func TestClientErrorResponses(t *testing.T) {
	app := fiber.New()

	t.Run("NewErrorClientResponse", func(t *testing.T) {
		app.Get("/test", func(c *fiber.Ctx) error {
			return NewErrorClientResponse("test error", http.StatusBadRequest, c)
		})

		req := httptest.NewRequest("GET", "/test", nil)
		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("ClientBadRequest", func(t *testing.T) {
		app.Get("/bad", func(c *fiber.Ctx) error {
			return ClientBadRequest("bad request", c)
		})

		req := httptest.NewRequest("GET", "/bad", nil)
		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("ClientNotFound", func(t *testing.T) {
		app.Get("/notfound", func(c *fiber.Ctx) error {
			return ClientNotFound("not found", c)
		})

		req := httptest.NewRequest("GET", "/notfound", nil)
		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusNotFound, resp.StatusCode)
	})

	t.Run("ClientInternalServerError", func(t *testing.T) {
		app.Get("/error", func(c *fiber.Ctx) error {
			return ClientInternalServerError("server error", c)
		})

		req := httptest.NewRequest("GET", "/error", nil)
		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
	})

	t.Run("ClientsInternalServerError", func(t *testing.T) {
		app.Get("/clients-error", func(c *fiber.Ctx) error {
			return ClientsInternalServerError("clients error", c)
		})

		req := httptest.NewRequest("GET", "/clients-error", nil)
		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
	})
}

func TestErrors(t *testing.T) {
	t.Run("ErrUserNotFound", func(t *testing.T) {
		assert.Equal(t, "user not found", ErrUserNotFound.Error())
	})

	t.Run("ErrClientNotFound", func(t *testing.T) {
		assert.Equal(t, "client not found", ErrClientNotFound.Error())
	})
}

func TestUserTypeVal(t *testing.T) {
	t.Run("UserTypeVal is initialized", func(t *testing.T) {
		assert.NotNil(t, UserTypeVal)
		assert.IsType(t, UserType{}, UserTypeVal)
	})
}
func TestMaskedBytes(t *testing.T) {
	t.Run("MarshalJSON returns masked bytes", func(t *testing.T) {
		mb := MaskedBytes([]byte("secret"))
		data, err := mb.MarshalJSON()
		assert.NoError(t, err)
		assert.Equal(t, []byte(`"*****"`), data)
	})

	t.Run("String returns masked string", func(t *testing.T) {
		mb := MaskedBytes([]byte("password"))
		assert.Equal(t, "*****", mb.String())
	})
}
func TestAuthProviderErrorResponses(t *testing.T) {
	app := fiber.New()

	t.Run("NewErrorAuthProviderResponse", func(t *testing.T) {
		app.Get("/test", func(c *fiber.Ctx) error {
			return NewErrorAuthProviderResponse("test error", http.StatusBadRequest, c)
		})

		req := httptest.NewRequest("GET", "/test", nil)
		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("AuthProviderBadRequest", func(t *testing.T) {
		app.Get("/bad", func(c *fiber.Ctx) error {
			return AuthProviderBadRequest("bad request", c)
		})

		req := httptest.NewRequest("GET", "/bad", nil)
		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("AuthProviderNotFound", func(t *testing.T) {
		app.Get("/notfound", func(c *fiber.Ctx) error {
			return AuthProviderNotFound("not found", c)
		})

		req := httptest.NewRequest("GET", "/notfound", nil)
		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusNotFound, resp.StatusCode)
	})

	t.Run("AuthProviderInternalServerError", func(t *testing.T) {
		app.Get("/error", func(c *fiber.Ctx) error {
			return AuthProviderInternalServerError("server error", c)
		})

		req := httptest.NewRequest("GET", "/error", nil)
		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
	})

	t.Run("AuthProvidersInternalServerError", func(t *testing.T) {
		app.Get("/providers-error", func(c *fiber.Ctx) error {
			return AuthProvidersInternalServerError("providers error", c)
		})

		req := httptest.NewRequest("GET", "/providers-error", nil)
		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
	})
}
// TestAuthMetadataType is a test implementation of AuthMetadataType
type TestAuthMetadataType struct{}

func (t TestAuthMetadataType) UpdateUserDetails(user *User) {
	user.Name = "Updated Name"
}

func TestAuthIdentity(t *testing.T) {
	t.Run("UpdateUserDetails calls metadata UpdateUserDetails", func(t *testing.T) {
		user := &User{
			Id:   "user-123",
			Name: "Original Name",
		}

		identity := AuthIdentity{
			Type:     AuthIdentityTypeEmail,
			Metadata: TestAuthMetadataType{},
		}

		identity.UpdateUserDetails(user)

		assert.Equal(t, "Updated Name", user.Name)
	})
}