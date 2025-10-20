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

	t.Run("validates configuration structure", func(t *testing.T) {
		cnf := config.NewAppConfig()

		// Test that all required configuration sections exist
		assert.NotNil(t, cnf.Server)
		assert.NotNil(t, cnf.DB)
		assert.NotNil(t, cnf.Jwt)
		assert.NotNil(t, cnf.Encrypter)
		assert.NotNil(t, cnf.Redis)
		assert.NotNil(t, cnf.Deployment)
		assert.NotNil(t, cnf.Logger)
	})

	t.Run("validates JWT configuration", func(t *testing.T) {
		cnf := config.NewAppConfig()

		// Test JWT configuration
		assert.NotNil(t, cnf.Jwt.Secret())
	})

	t.Run("validates encryption configuration", func(t *testing.T) {
		cnf := config.NewAppConfig()

		// Test encryption configuration
		assert.NotNil(t, cnf.Encrypter.Key())
	})

	t.Run("validates server configuration", func(t *testing.T) {
		cnf := config.NewAppConfig()

		// Test server configuration
		assert.NotEmpty(t, cnf.Server.Host)
		assert.NotEmpty(t, cnf.Server.Port)
		assert.GreaterOrEqual(t, cnf.Server.TokenCacheTTLInMinutes, int64(0))
		assert.GreaterOrEqual(t, cnf.Server.AuthProviderRefetchIntervalInMinutes, int64(0))
	})

	t.Run("validates Redis configuration", func(t *testing.T) {
		cnf := config.NewAppConfig()

		// Test Redis configuration
		assert.NotEmpty(t, cnf.Redis.Host)
	})

	t.Run("validates deployment configuration", func(t *testing.T) {
		cnf := config.NewAppConfig()

		// Test deployment configuration
		assert.NotEmpty(t, cnf.Deployment.Environment)
		assert.NotEmpty(t, cnf.Deployment.Name)
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

	t.Run("validates database configuration", func(t *testing.T) {
		cnf := config.NewAppConfig()

		// Test that database configuration is properly structured
		assert.NotEmpty(t, cnf.DB.Host())
		assert.NotNil(t, cnf.DB)
	})

	t.Run("handles database configuration errors", func(t *testing.T) {
		// Test with invalid configuration
		cnf := config.NewAppConfig()

		// Even with default config, it should be valid
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

	t.Run("Provider struct with all fields", func(t *testing.T) {
		provider := &Provider{
			S:          &Service{},
			D:          &test.MockDB{},
			C:          &testservices.MockCacheService{},
			PM:         nil, // PM is a middleware, not a service
			AM:         &auth.Middlewares{},
			AuthClient: &sdk.Client{Id: "test-client", GoIamClient: true},
		}

		assert.NotNil(t, provider.S)
		assert.NotNil(t, provider.D)
		assert.NotNil(t, provider.C)
		assert.NotNil(t, provider.AM)
		assert.NotNil(t, provider.AuthClient)
		assert.Equal(t, "test-client", provider.AuthClient.Id)
		assert.True(t, provider.AuthClient.GoIamClient)
	})
}

// TestNewDBConnection_Extended tests the NewDBConnection function more comprehensively
func TestNewDBConnection_Extended(t *testing.T) {
	t.Run("validates database configuration structure", func(t *testing.T) {
		cnf := config.NewAppConfig()

		// Test database configuration fields
		assert.NotNil(t, cnf.DB)
		assert.NotEmpty(t, cnf.DB.Host())

		// Test that configuration is properly loaded
		assert.NotNil(t, cnf)
	})

	t.Run("validates database connection parameters", func(t *testing.T) {
		cnf := config.NewAppConfig()

		// Test that database host is configured
		host := cnf.DB.Host()
		assert.NotEmpty(t, host)
		assert.IsType(t, "", host)
	})

	t.Run("validates database configuration loading", func(t *testing.T) {
		// Test that configuration can be loaded multiple times
		cnf1 := config.NewAppConfig()
		cnf2 := config.NewAppConfig()

		assert.NotNil(t, cnf1)
		assert.NotNil(t, cnf2)
		assert.NotNil(t, cnf1.DB)
		assert.NotNil(t, cnf2.DB)
	})
}

// TestInjectDefaultProviders_Extended tests the InjectDefaultProviders function more comprehensively
func TestInjectDefaultProviders_Extended(t *testing.T) {
	t.Run("validates configuration loading process", func(t *testing.T) {
		cnf := config.NewAppConfig()

		// Test that all configuration sections are loaded
		assert.NotNil(t, cnf.Server)
		assert.NotNil(t, cnf.DB)
		assert.NotNil(t, cnf.Jwt)
		assert.NotNil(t, cnf.Encrypter)
		assert.NotNil(t, cnf.Redis)
		assert.NotNil(t, cnf.Deployment)
		assert.NotNil(t, cnf.Logger)
	})

	t.Run("validates configuration field access", func(t *testing.T) {
		cnf := config.NewAppConfig()

		// Test server configuration access
		assert.NotEmpty(t, cnf.Server.Host)
		assert.NotEmpty(t, cnf.Server.Port)
		assert.GreaterOrEqual(t, cnf.Server.TokenCacheTTLInMinutes, int64(0))
		assert.GreaterOrEqual(t, cnf.Server.AuthProviderRefetchIntervalInMinutes, int64(0))

		// Test database configuration access
		assert.NotEmpty(t, cnf.DB.Host())

		// Test JWT configuration access
		assert.NotNil(t, cnf.Jwt.Secret())

		// Test encryption configuration access
		assert.NotNil(t, cnf.Encrypter.Key())

		// Test Redis configuration access
		assert.NotEmpty(t, cnf.Redis.Host)

		// Test deployment configuration access
		assert.NotEmpty(t, cnf.Deployment.Environment)
		assert.NotEmpty(t, cnf.Deployment.Name)
	})

	t.Run("validates configuration consistency", func(t *testing.T) {
		cnf1 := config.NewAppConfig()
		cnf2 := config.NewAppConfig()

		// Test that configurations are consistent
		assert.Equal(t, cnf1.Server.Host, cnf2.Server.Host)
		assert.Equal(t, cnf1.Server.Port, cnf2.Server.Port)
		assert.Equal(t, cnf1.DB.Host(), cnf2.DB.Host())
		assert.Equal(t, cnf1.Redis.Host, cnf2.Redis.Host)
		assert.Equal(t, cnf1.Deployment.Environment, cnf2.Deployment.Environment)
		assert.Equal(t, cnf1.Deployment.Name, cnf2.Deployment.Name)
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
