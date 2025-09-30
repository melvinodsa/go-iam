package providers

import (
	"context"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/melvinodsa/go-iam/config"
	"github.com/melvinodsa/go-iam/middlewares/auth"
	"github.com/stretchr/testify/mock"
	"github.com/melvinodsa/go-iam/sdk"
	testservices "github.com/melvinodsa/go-iam/utils/test/services"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHandle(t *testing.T) {
	p := &Provider{} // empty provider for test
	handler := Handle(p)

	app := fiber.New()
	app.Use(handler)
	app.Get("/test", func(c *fiber.Ctx) error {
		retrieved := GetProviders(c)
		assert.Equal(t, p, retrieved, "GetProviders should return the set provider")
		return c.SendString("ok")
	})

	req := httptest.NewRequest("GET", "/test", nil)
	resp, err := app.Test(req, -1)
	assert.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)
}

func TestInjectDefaultProviders(t *testing.T) {
	err := os.Setenv("JWT_SECRET", "abcd")
	require.NoError(t, err)
	cnf := config.NewAppConfig()

	t.Run("successful provider injection", func(t *testing.T) {
		provider, err := InjectDefaultProviders(*cnf)
		assert.NoError(t, err)
		assert.NotNil(t, provider)
		assert.NotNil(t, provider.S)
		assert.NotNil(t, provider.D)
		assert.NotNil(t, provider.C)
		assert.NotNil(t, provider.PM)
		assert.NotNil(t, provider.AM)

		// Clean up
		if provider.D != nil {
			provider.D.Disconnect(context.Background())
		}
	})
}

func TestProviderStruct(t *testing.T) {
	t.Run("Provider struct initialization", func(t *testing.T) {
		provider := &Provider{
			S:          &Service{},
			D:          nil,
			C:          nil,
			PM:         nil,
			AM:         &auth.Middlewares{},
			AuthClient: &sdk.Client{Id: "test-client"},
		}

		assert.NotNil(t, provider.S)
		assert.NotNil(t, provider.AM)
		assert.NotNil(t, provider.AuthClient)
		assert.Equal(t, "test-client", provider.AuthClient.Id)
	})
}

func TestCheckAndAddDefaultProject(t *testing.T) {
	t.Run("default project already exists", func(t *testing.T) {
		mockSvc := &testservices.MockProjectService{}
		mockSvc.On("GetByName", context.Background(), "Default Project").Return(&sdk.Project{Name: "Default Project"}, nil)

		err := CheckAndAddDefaultProject(mockSvc)
		assert.NoError(t, err)
		mockSvc.AssertExpectations(t)
	})

	t.Run("default project does not exist, create it", func(t *testing.T) {
		mockSvc := &testservices.MockProjectService{}
		mockSvc.On("GetByName", context.Background(), "Default Project").Return((*sdk.Project)(nil), sdk.ErrProjectNotFound)
		mockSvc.On("Create", context.Background(), mock.AnythingOfType("*sdk.Project")).Return(nil)

		err := CheckAndAddDefaultProject(mockSvc)
		assert.NoError(t, err)
		mockSvc.AssertExpectations(t)
	})

	t.Run("error fetching project", func(t *testing.T) {
		mockSvc := &testservices.MockProjectService{}
		mockSvc.On("GetByName", context.Background(), "Default Project").Return((*sdk.Project)(nil), assert.AnError)

		err := CheckAndAddDefaultProject(mockSvc)
		assert.Error(t, err)
		mockSvc.AssertExpectations(t)
	})
}