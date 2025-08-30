package auth

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/melvinodsa/go-iam/sdk"
	"github.com/melvinodsa/go-iam/utils"
	"github.com/melvinodsa/go-iam/utils/goiamuniverse"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockAuthService is a mock implementation of auth.Service
type MockAuthService struct {
	mock.Mock
}

func (m *MockAuthService) GetLoginUrl(ctx context.Context, clientId, authProviderId, state, redirectUrl, codeChallengeMethod, codeChallenge string) (string, error) {
	args := m.Called(ctx, clientId, authProviderId, state, redirectUrl, codeChallengeMethod, codeChallenge)
	return args.String(0), args.Error(1)
}

func (m *MockAuthService) Redirect(ctx context.Context, code, state string) (*sdk.AuthRedirectResponse, error) {
	args := m.Called(ctx, code, state)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*sdk.AuthRedirectResponse), args.Error(1)
}

func (m *MockAuthService) ClientCallback(ctx context.Context, code, codeChallenge, clientId, clietSecret string) (*sdk.AuthVerifyCodeResponse, error) {
	args := m.Called(ctx, code, codeChallenge, clientId, clietSecret)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*sdk.AuthVerifyCodeResponse), args.Error(1)
}

func (m *MockAuthService) GetIdentity(ctx context.Context, token string) (*sdk.User, error) {
	args := m.Called(ctx, token)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*sdk.User), args.Error(1)
}

func (m *MockAuthService) HandleEvent(event utils.Event[sdk.Client]) {
	m.Called(event)
}

// MockClientService is a mock implementation of client.Service
type MockClientService struct {
	mock.Mock
}

func (m *MockClientService) GetAll(ctx context.Context, queryParams sdk.ClientQueryParams) ([]sdk.Client, error) {
	args := m.Called(ctx, queryParams)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]sdk.Client), args.Error(1)
}

func (m *MockClientService) GetGoIamClients(ctx context.Context, params sdk.ClientQueryParams) ([]sdk.Client, error) {
	args := m.Called(ctx, params)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]sdk.Client), args.Error(1)
}

func (m *MockClientService) Get(ctx context.Context, id string, dontCheckProjects bool) (*sdk.Client, error) {
	args := m.Called(ctx, id, dontCheckProjects)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*sdk.Client), args.Error(1)
}

func (m *MockClientService) Create(ctx context.Context, client *sdk.Client) error {
	args := m.Called(ctx, client)
	return args.Error(0)
}

func (m *MockClientService) Update(ctx context.Context, client *sdk.Client) error {
	args := m.Called(ctx, client)
	return args.Error(0)
}

func (m *MockClientService) Emit(event utils.Event[sdk.Client]) {
	m.Called(event)
}

func (m *MockClientService) Subscribe(eventName goiamuniverse.Event, subscriber utils.Subscriber[utils.Event[sdk.Client], sdk.Client]) {
	m.Called(eventName, subscriber)
}

func setupTestApp() (*fiber.App, *MockAuthService, *MockClientService) {
	app := fiber.New()
	mockAuthSvc := new(MockAuthService)
	mockClientSvc := new(MockClientService)
	return app, mockAuthSvc, mockClientSvc
}

func setupTestAppWithAuthClient() (*fiber.App, *MockAuthService, *MockClientService, *Middlewares) {
	app, mockAuthSvc, mockClientSvc := setupTestApp()

	// Mock the GetGoIamClients call for middleware creation
	testClient := createTestClient()
	mockClientSvc.On("GetGoIamClients", mock.Anything, mock.MatchedBy(func(params sdk.ClientQueryParams) bool {
		return params.GoIamClient == true
	})).Return([]sdk.Client{*testClient}, nil)

	middlewares, err := NewMiddlewares(mockAuthSvc, mockClientSvc)
	if err != nil {
		return app, mockAuthSvc, mockClientSvc, nil
	}
	return app, mockAuthSvc, mockClientSvc, middlewares
}

func setupTestAppWithoutAuthClient() (*fiber.App, *MockAuthService, *MockClientService, *Middlewares) {
	app, mockAuthSvc, mockClientSvc := setupTestApp()

	// Mock the GetGoIamClients call to return empty slice (no auth client)
	mockClientSvc.On("GetGoIamClients", mock.Anything, mock.MatchedBy(func(params sdk.ClientQueryParams) bool {
		return params.GoIamClient == true
	})).Return([]sdk.Client{}, nil)

	middlewares, err := NewMiddlewares(mockAuthSvc, mockClientSvc)
	if err != nil {
		return app, mockAuthSvc, mockClientSvc, nil
	}
	return app, mockAuthSvc, mockClientSvc, middlewares
}

func createTestUser() *sdk.User {
	return &sdk.User{
		Id:    "test-user-id",
		Name:  "Test User",
		Email: "test@example.com",
	}
}

func createTestClient() *sdk.Client {
	return &sdk.Client{
		Id:   "test-client-id",
		Name: "Test Client",
	}
}

func TestNewMiddlewares(t *testing.T) {
	_, mockAuthSvc, mockClientSvc := setupTestApp()

	// Mock the GetGoIamClients call that happens during NewMiddlewares
	testClient := createTestClient()
	mockClientSvc.On("GetGoIamClients", mock.Anything, mock.MatchedBy(func(params sdk.ClientQueryParams) bool {
		return params.GoIamClient == true
	})).Return([]sdk.Client{*testClient}, nil)

	middlewares, err := NewMiddlewares(mockAuthSvc, mockClientSvc)
	assert.NoError(t, err)

	assert.NotNil(t, middlewares)
	assert.Equal(t, mockAuthSvc, middlewares.authSvc)
	assert.Equal(t, mockClientSvc, middlewares.clientSvc)
	assert.NotNil(t, middlewares.AuthClient)
	assert.Equal(t, testClient.Id, middlewares.AuthClient.Id)

	mockClientSvc.AssertExpectations(t)
}

func TestNewMiddlewares_NoAuthClient(t *testing.T) {
	_, mockAuthSvc, mockClientSvc := setupTestApp()

	// Mock the GetGoIamClients call to return empty slice (no auth client)
	mockClientSvc.On("GetGoIamClients", mock.Anything, mock.MatchedBy(func(params sdk.ClientQueryParams) bool {
		return params.GoIamClient == true
	})).Return([]sdk.Client{}, nil)

	middlewares, err := NewMiddlewares(mockAuthSvc, mockClientSvc)
	assert.NoError(t, err)

	assert.NotNil(t, middlewares)
	assert.Equal(t, mockAuthSvc, middlewares.authSvc)
	assert.Equal(t, mockClientSvc, middlewares.clientSvc)
	assert.Nil(t, middlewares.AuthClient) // No auth client when empty result

	mockClientSvc.AssertExpectations(t)
}

func TestNewMiddlewares_GetGoIamClientsError(t *testing.T) {
	_, mockAuthSvc, mockClientSvc := setupTestApp()

	// Mock the GetGoIamClients call to return an error
	mockClientSvc.On("GetGoIamClients", mock.Anything, mock.MatchedBy(func(params sdk.ClientQueryParams) bool {
		return params.GoIamClient == true
	})).Return(nil, errors.New("database error"))

	middlewares, err := NewMiddlewares(mockAuthSvc, mockClientSvc)
	assert.Error(t, err)

	assert.Nil(t, middlewares)

	mockClientSvc.AssertExpectations(t)
}

func TestMiddlewares_User_Success(t *testing.T) {
	app, mockAuthSvc, mockClientSvc, middlewares := setupTestAppWithAuthClient()

	testUser := createTestUser()
	mockAuthSvc.On("GetIdentity", mock.Anything, "valid-token").Return(testUser, nil)

	// Create a test route
	app.Get("/test", middlewares.User, func(c *fiber.Ctx) error {
		user := c.Context().UserValue("user").(*sdk.User)
		return c.JSON(user)
	})

	// Create request with valid authorization header
	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	req.Header.Set("Authorization", "Bearer valid-token")

	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	mockAuthSvc.AssertExpectations(t)
	mockClientSvc.AssertExpectations(t)
}

func TestMiddlewares_User_NoAuthClient(t *testing.T) {
	app, mockAuthSvc, mockClientSvc, middlewares := setupTestAppWithoutAuthClient()

	// Create a test route
	app.Get("/test", middlewares.User, func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"message": "success"})
	})

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	// Should not call auth service when no auth client
	mockAuthSvc.AssertNotCalled(t, "GetIdentity")
	mockClientSvc.AssertExpectations(t)
}

func TestMiddlewares_User_AuthError(t *testing.T) {
	app, mockAuthSvc, mockClientSvc, middlewares := setupTestAppWithAuthClient()

	mockAuthSvc.On("GetIdentity", mock.Anything, "invalid-token").Return(nil, errors.New("invalid token"))

	app.Get("/test", middlewares.User, func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"message": "success"})
	})

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	req.Header.Set("Authorization", "Bearer invalid-token")

	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)

	mockAuthSvc.AssertExpectations(t)
	mockClientSvc.AssertExpectations(t)
}

func TestMiddlewares_DashboardUser_Success(t *testing.T) {
	app, mockAuthSvc, mockClientSvc, middlewares := setupTestAppWithAuthClient()

	testUser := createTestUser()
	mockAuthSvc.On("GetIdentity", mock.Anything, "valid-token").Return(testUser, nil)

	app.Get("/dashboard", middlewares.DashboardUser, func(c *fiber.Ctx) error {
		user := c.Context().UserValue("user").(*sdk.User)
		return c.JSON(user)
	})

	req := httptest.NewRequest(http.MethodGet, "/dashboard", nil)
	req.Header.Set("Authorization", "Bearer valid-token")

	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	mockAuthSvc.AssertExpectations(t)
	mockClientSvc.AssertExpectations(t)
}

func TestMiddlewares_DashboardUser_NoAuthClient(t *testing.T) {
	app, mockAuthSvc, mockClientSvc, middlewares := setupTestAppWithoutAuthClient()

	app.Get("/dashboard", middlewares.DashboardUser, func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"message": "success"})
	})

	req := httptest.NewRequest(http.MethodGet, "/dashboard", nil)
	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	mockAuthSvc.AssertNotCalled(t, "GetIdentity")
	mockClientSvc.AssertExpectations(t)
}

func TestMiddlewares_DashboardUser_AuthError(t *testing.T) {
	app, mockAuthSvc, mockClientSvc, middlewares := setupTestAppWithAuthClient()

	mockAuthSvc.On("GetIdentity", mock.Anything, "invalid-token").Return(nil, errors.New("invalid token"))

	app.Get("/dashboard", middlewares.DashboardUser, func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"message": "success"})
	})

	req := httptest.NewRequest(http.MethodGet, "/dashboard", nil)
	req.Header.Set("Authorization", "Bearer invalid-token")

	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)

	mockAuthSvc.AssertExpectations(t)
	mockClientSvc.AssertExpectations(t)
}

func TestMiddlewares_GetUser_Success(t *testing.T) {
	_, mockAuthSvc, mockClientSvc, middlewares := setupTestAppWithAuthClient()

	testUser := createTestUser()
	mockAuthSvc.On("GetIdentity", mock.Anything, "valid-token").Return(testUser, nil)

	// Test GetUser indirectly through middleware since direct testing requires internal context
	app := fiber.New()
	app.Get("/test", func(c *fiber.Ctx) error {
		user, err := middlewares.GetUser(c)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": err.Error()})
		}
		return c.JSON(user)
	})

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	req.Header.Set("Authorization", "Bearer valid-token")

	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	mockAuthSvc.AssertExpectations(t)
	mockClientSvc.AssertExpectations(t)
}

func TestMiddlewares_GetUser_NoAuthHeader(t *testing.T) {
	_, mockAuthSvc, mockClientSvc, middlewares := setupTestAppWithAuthClient()

	app := fiber.New()
	app.Get("/test", func(c *fiber.Ctx) error {
		user, err := middlewares.GetUser(c)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": err.Error()})
		}
		return c.JSON(user)
	})

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	// No Authorization header

	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)

	mockAuthSvc.AssertNotCalled(t, "GetIdentity")
	mockClientSvc.AssertExpectations(t)
}

func TestMiddlewares_GetUser_InvalidAuthHeader(t *testing.T) {
	_, mockAuthSvc, mockClientSvc, middlewares := setupTestAppWithAuthClient()

	// Test cases for invalid headers
	testCases := []struct {
		name       string
		authHeader string
	}{
		{"Too short", "Invalid"},
		{"No token", "Bearer"},
		{"Empty token", "Bearer "},
		{"Wrong scheme", "Basic dGVzdA=="},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			app := fiber.New()
			app.Get("/test", func(c *fiber.Ctx) error {
				user, err := middlewares.GetUser(c)
				if err != nil {
					return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": err.Error()})
				}
				return c.JSON(user)
			})

			req := httptest.NewRequest(http.MethodGet, "/test", nil)
			req.Header.Set("Authorization", tc.authHeader)

			resp, err := app.Test(req)
			assert.NoError(t, err)
			assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
		})
	}

	mockAuthSvc.AssertNotCalled(t, "GetIdentity")
	mockClientSvc.AssertExpectations(t)
}

func TestMiddlewares_GetUser_AuthServiceError(t *testing.T) {
	_, mockAuthSvc, mockClientSvc, middlewares := setupTestAppWithAuthClient()

	mockAuthSvc.On("GetIdentity", mock.Anything, "invalid-token").Return(nil, errors.New("token expired"))

	app := fiber.New()
	app.Get("/test", func(c *fiber.Ctx) error {
		user, err := middlewares.GetUser(c)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": err.Error()})
		}
		return c.JSON(user)
	})

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	req.Header.Set("Authorization", "Bearer invalid-token")

	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)

	mockAuthSvc.AssertExpectations(t)
	mockClientSvc.AssertExpectations(t)
}

func TestMiddlewares_User_MissingAuthHeader(t *testing.T) {
	app, mockAuthSvc, mockClientSvc, middlewares := setupTestAppWithAuthClient()

	app.Get("/test", middlewares.User, func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"message": "success"})
	})

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	// No Authorization header

	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)

	mockAuthSvc.AssertNotCalled(t, "GetIdentity")
	mockClientSvc.AssertExpectations(t)
}

func TestMiddlewares_DashboardUser_MissingAuthHeader(t *testing.T) {
	app, mockAuthSvc, mockClientSvc, middlewares := setupTestAppWithAuthClient()

	app.Get("/dashboard", middlewares.DashboardUser, func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"message": "success"})
	})

	req := httptest.NewRequest(http.MethodGet, "/dashboard", nil)
	// No Authorization header

	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)

	mockAuthSvc.AssertNotCalled(t, "GetIdentity")
	mockClientSvc.AssertExpectations(t)
}

// Benchmark tests
func BenchmarkMiddlewares_GetUser(b *testing.B) {
	_, mockAuthSvc, _, middlewares := setupTestAppWithAuthClient()

	testUser := createTestUser()
	mockAuthSvc.On("GetIdentity", mock.Anything, "valid-token").Return(testUser, nil)

	app := fiber.New()
	app.Get("/test", func(c *fiber.Ctx) error {
		user, err := middlewares.GetUser(c)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": err.Error()})
		}
		return c.JSON(user)
	})

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	req.Header.Set("Authorization", "Bearer valid-token")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := app.Test(req)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkMiddlewares_User(b *testing.B) {
	app, mockAuthSvc, _, middlewares := setupTestAppWithAuthClient()

	testUser := createTestUser()
	mockAuthSvc.On("GetIdentity", mock.Anything, "valid-token").Return(testUser, nil)

	app.Get("/test", middlewares.User, func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"message": "success"})
	})

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	req.Header.Set("Authorization", "Bearer valid-token")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := app.Test(req)
		if err != nil {
			b.Fatal(err)
		}
	}
}
