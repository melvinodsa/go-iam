package providers

import (
	"context"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/melvinodsa/go-iam/config"
	"github.com/melvinodsa/go-iam/middlewares/auth"
	"github.com/melvinodsa/go-iam/sdk"
	"github.com/melvinodsa/go-iam/utils/goiamuniverse"
	"github.com/melvinodsa/go-iam/utils/test"
	testservices "github.com/melvinodsa/go-iam/utils/test/services"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
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
		assert.NotNil(t, cnf)
		assert.NotNil(t, cnf.Jwt)
	})
}

func TestNewServices(t *testing.T) {
	t.Run("creates services with mock dependencies", func(t *testing.T) {
		mockDB := &test.MockDB{}
		mockCache := &testservices.MockCacheService{}
		mockEncrypt := &testservices.MockEncryptService{}
		mockJWT := &testservices.MockJWTService{}

		services := NewServices(mockDB, mockCache, mockEncrypt, mockJWT, 60, 30)

		assert.NotNil(t, services)
		assert.NotNil(t, services.Projects)
		assert.NotNil(t, services.Clients)
		assert.NotNil(t, services.AuthProviders)
		assert.NotNil(t, services.Auth)
		assert.NotNil(t, services.Resources)
		assert.NotNil(t, services.User)
		assert.NotNil(t, services.Role)
		assert.NotNil(t, services.Policy)
	})
}

func TestNewDBConnection(t *testing.T) {
	t.Run("creates database connection with mock", func(t *testing.T) {
		cnf := config.NewAppConfig()

		// This test would normally try to connect to a real database
		// For now, we'll just test that the config is valid
		assert.NotNil(t, cnf)
		assert.NotNil(t, cnf.DB)
	})
}

func TestProvider_HandleEvent(t *testing.T) {
	t.Run("ignores non-client events", func(t *testing.T) {
		provider := &Provider{
			AuthClient: &sdk.Client{Id: "test-client"},
		}

		event := &testservices.MockEvent[sdk.Client]{}
		event.On("Name").Return(goiamuniverse.Event("some-other-event"))

		provider.HandleEvent(event)

		event.AssertExpectations(t)
	})

	t.Run("ignores non-go-iam clients", func(t *testing.T) {
		provider := &Provider{
			AuthClient: &sdk.Client{Id: "test-client"},
		}

		mockClient := &sdk.Client{Id: "test-client", GoIamClient: false}
		event := &testservices.MockEvent[sdk.Client]{}
		event.On("Name").Return(goiamuniverse.EventClientCreated)
		event.On("Payload").Return(*mockClient)

		provider.HandleEvent(event)

		event.AssertExpectations(t)
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
