package me

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/melvinodsa/go-iam/config"
	"github.com/melvinodsa/go-iam/providers"
	"github.com/melvinodsa/go-iam/sdk"
	"github.com/melvinodsa/go-iam/services/cache"
	"github.com/melvinodsa/go-iam/utils/test"
	"github.com/melvinodsa/go-iam/utils/test/server"
	"github.com/melvinodsa/go-iam/utils/test/services"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestMe(t *testing.T) {
	err := os.Setenv("JWT_SECRET", "abcd")
	require.NoError(t, err)
	cnf := config.NewAppConfig()
	log.Infow("Loaded Configurations",
		"host", cnf.Server.Host,
		"port", cnf.Server.Port,
		"env", cnf.Deployment.Environment,
		"app_name", cnf.Deployment.Name,
	)

	t.Run("success - returns user information", func(t *testing.T) {
		app := fiber.New(fiber.Config{
			ReadBufferSize: 8192,
		})

		d := test.SetupMockDB()
		cs := cache.NewMockService()
		svcs, err := server.GetServices(*cnf, cs, d)
		if err != nil {
			t.Errorf("error getting services: %s", err)
			return
		}

		prv := server.SetupTestServer(app, cnf, svcs, cs, d)

		app.Use(providers.Handle(prv))

		// Create a test user and set it in context
		testUser := &sdk.User{
			Id:    "user-123",
			Email: "test@example.com",
			Name:  "Test User",
		}

		// Add middleware to set user in context
		app.Use(func(c *fiber.Ctx) error {
			c.Context().SetUserValue(sdk.UserTypeVal, testUser)
			return c.Next()
		})

		// Register the route directly
		app.Get("/me/v1/", Me)

		req, _ := http.NewRequest("GET", "/me/v1/", nil)
		res, err := app.Test(req, -1)
		assert.Nil(t, err)
		assert.Equal(t, http.StatusOK, res.StatusCode)

		var resp sdk.UserResponse
		err = json.NewDecoder(res.Body).Decode(&resp)
		assert.Nil(t, err)
		assert.True(t, resp.Success)
		assert.Equal(t, "User fetched successfully", resp.Message)
		assert.Equal(t, testUser, resp.Data)
	})

	t.Run("unauthorized when user missing", func(t *testing.T) {
		app := fiber.New(fiber.Config{
			ReadBufferSize: 8192,
		})

		d := test.SetupMockDB()
		cs := cache.NewMockService()
		svcs, err := server.GetServices(*cnf, cs, d)
		if err != nil {
			t.Errorf("error getting services: %s", err)
			return
		}

		prv := server.SetupTestServer(app, cnf, svcs, cs, d)
		app.Use(providers.Handle(prv))
		app.Get("/me/v1/", Me)

		req, _ := http.NewRequest("GET", "/me/v1/", nil)
		res, err := app.Test(req, -1)
		assert.Nil(t, err)
		assert.Equal(t, http.StatusUnauthorized, res.StatusCode)
	})
}

func TestResetPassword(t *testing.T) {
	err := os.Setenv("JWT_SECRET", "abcd")
	require.NoError(t, err)
	cnf := config.NewAppConfig()

	t.Run("success - returns password reset url in postback mode", func(t *testing.T) {
		app := fiber.New(fiber.Config{
			ReadBufferSize: 8192,
		})

		d := test.SetupMockDB()
		cs := cache.NewMockService()
		svcs, err := server.GetServices(*cnf, cs, d)
		if err != nil {
			t.Errorf("error getting services: %s", err)
			return
		}

		mockAuthSvc := services.MockAuthService{}
		mockAuthSvc.On("GetResetPasswordUrl", mock.Anything, "provider-1").Return("https://idp.example.com/reset", nil).Once()
		svcs.Auth = &mockAuthSvc

		prv := server.SetupTestServer(app, cnf, svcs, cs, d)
		prv.AuthClient = &sdk.Client{
			Id:                    "go-iam-client",
			DefaultAuthProviderId: "provider-1",
		}

		app.Use(providers.Handle(prv))
		app.Use(func(c *fiber.Ctx) error {
			c.Context().SetUserValue(sdk.UserTypeVal, &sdk.User{Id: "user-123"})
			return c.Next()
		})
		app.Get("/me/v1/password-reset", ResetPassword)

		req, _ := http.NewRequest("GET", "/me/v1/password-reset?postback=true", nil)
		res, err := app.Test(req, -1)
		assert.Nil(t, err)
		assert.Equal(t, http.StatusOK, res.StatusCode)

		var resp sdk.AuthRedirectResponse
		err = json.NewDecoder(res.Body).Decode(&resp)
		assert.Nil(t, err)
		assert.Equal(t, "https://idp.example.com/reset", resp.RedirectUrl)
		mockAuthSvc.AssertExpectations(t)
	})

	t.Run("success - redirects to password reset url", func(t *testing.T) {
		app := fiber.New(fiber.Config{
			ReadBufferSize: 8192,
		})

		d := test.SetupMockDB()
		cs := cache.NewMockService()
		svcs, err := server.GetServices(*cnf, cs, d)
		if err != nil {
			t.Errorf("error getting services: %s", err)
			return
		}

		mockAuthSvc := services.MockAuthService{}
		mockAuthSvc.On("GetResetPasswordUrl", mock.Anything, "provider-1").Return("https://idp.example.com/reset", nil).Once()
		svcs.Auth = &mockAuthSvc

		prv := server.SetupTestServer(app, cnf, svcs, cs, d)
		prv.AuthClient = &sdk.Client{
			Id:                    "go-iam-client",
			DefaultAuthProviderId: "provider-1",
		}

		app.Use(providers.Handle(prv))
		app.Use(func(c *fiber.Ctx) error {
			c.Context().SetUserValue(sdk.UserTypeVal, &sdk.User{Id: "user-123"})
			return c.Next()
		})
		app.Get("/me/v1/password-reset", ResetPassword)

		req, _ := http.NewRequest("GET", "/me/v1/password-reset", nil)
		res, err := app.Test(req, -1)
		assert.Nil(t, err)
		assert.Equal(t, http.StatusTemporaryRedirect, res.StatusCode)
		assert.Equal(t, "https://idp.example.com/reset", res.Header.Get("Location"))
	})

	t.Run("unauthorized when user missing", func(t *testing.T) {
		app := fiber.New()

		d := test.SetupMockDB()
		cs := cache.NewMockService()
		svcs, err := server.GetServices(*cnf, cs, d)
		require.NoError(t, err)

		prv := server.SetupTestServer(app, cnf, svcs, cs, d)
		prv.AuthClient = &sdk.Client{
			Id:                    "go-iam-client",
			DefaultAuthProviderId: "provider-1",
		}
		app.Use(providers.Handle(prv))
		app.Get("/me/v1/password-reset", ResetPassword)

		req := httptest.NewRequest(http.MethodGet, "/me/v1/password-reset", nil)
		res, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusUnauthorized, res.StatusCode)
	})

	t.Run("bad request when auth client is not configured", func(t *testing.T) {
		app := fiber.New()

		d := test.SetupMockDB()
		cs := cache.NewMockService()
		svcs, err := server.GetServices(*cnf, cs, d)
		require.NoError(t, err)

		prv := server.SetupTestServer(app, cnf, svcs, cs, d)
		prv.AuthClient = nil
		app.Use(providers.Handle(prv))
		app.Use(func(c *fiber.Ctx) error {
			c.Context().SetUserValue(sdk.UserTypeVal, &sdk.User{Id: "user-123"})
			return c.Next()
		})
		app.Get("/me/v1/password-reset", ResetPassword)

		req := httptest.NewRequest(http.MethodGet, "/me/v1/password-reset", nil)
		res, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, res.StatusCode)
	})
}

func TestAuthClientCheck(t *testing.T) {
	err := os.Setenv("JWT_SECRET", "abcd")
	require.NoError(t, err)
	cnf := config.NewAppConfig()

	t.Run("auth client is set up - continues to next", func(t *testing.T) {
		app := fiber.New()

		d := test.SetupMockDB()
		cs := cache.NewMockService()
		svcs, err := server.GetServices(*cnf, cs, d)
		if err != nil {
			t.Errorf("error getting services: %s", err)
			return
		}

		// Set up auth client
		authClient := &sdk.Client{
			Id: "test-client",
		}

		prv := &providers.Provider{
			S:          svcs,
			D:          d,
			C:          cs,
			AuthClient: authClient,
		}

		app.Use(providers.Handle(prv))

		called := false
		app.Get("/test", AuthClientCheck, func(c *fiber.Ctx) error {
			called = true
			return c.SendString("success")
		})

		req := httptest.NewRequest(http.MethodGet, "/test", nil)
		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
		assert.True(t, called)
	})

	t.Run("auth client not set up - returns setup not complete", func(t *testing.T) {
		app := fiber.New()

		d := test.SetupMockDB()
		cs := cache.NewMockService()
		svcs, err := server.GetServices(*cnf, cs, d)
		if err != nil {
			t.Errorf("error getting services: %s", err)
			return
		}

		// No auth client set
		prv := &providers.Provider{
			S: svcs,
			D: d,
			C: cs,
		}

		app.Use(providers.Handle(prv))

		app.Get("/test", AuthClientCheck, func(c *fiber.Ctx) error {
			return c.SendString("should not reach here")
		})

		req := httptest.NewRequest(http.MethodGet, "/test", nil)
		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var response sdk.DashboardUserResponse
		err = json.NewDecoder(resp.Body).Decode(&response)
		assert.NoError(t, err)
		assert.True(t, response.Success)
		assert.Equal(t, "auth is not setup yet.", response.Message)
		assert.False(t, response.Data.Setup.ClientAdded)
	})
}

func TestRegisterRoutes(t *testing.T) {
	app := fiber.New()

	RegisterRoutes(app, "/api")

	// Check if routes are registered
	routes := app.GetRoutes()
	meRouteFound := false
	resetRouteFound := false
	for _, route := range routes {
		if route.Path == "/api/v1/" && route.Method == "GET" {
			meRouteFound = true
		}
		if route.Path == "/api/v1/password-reset" && route.Method == "GET" {
			resetRouteFound = true
		}
	}
	assert.True(t, meRouteFound, "Me route should be registered")
	assert.True(t, resetRouteFound, "Reset password route should be registered")
}

func TestRegisterOpenRoutes(t *testing.T) {
	app := fiber.New()

	d := test.SetupMockDB()
	cs := cache.NewMockService()
	cnf := config.NewAppConfig()
	svcs, err := server.GetServices(*cnf, cs, d)
	require.NoError(t, err)

	authClient := &sdk.Client{Id: "test-client"}
	prv := &providers.Provider{
		S:          svcs,
		D:          d,
		C:          cs,
		AuthClient: authClient,
	}

	RegisterOpenRoutes(app, "/api", prv)

	// Check if routes are registered
	routes := app.GetRoutes()
	dashboardRouteFound := false
	for _, route := range routes {
		if route.Path == "/api/v1/dashboard" && route.Method == "GET" {
			dashboardRouteFound = true
			break
		}
	}
	assert.True(t, dashboardRouteFound, "Dashboard route should be registered")
}
func TestDashboardMe(t *testing.T) {
	err := os.Setenv("JWT_SECRET", "abcd")
	require.NoError(t, err)
	cnf := config.NewAppConfig()

	t.Run("success - returns dashboard user information", func(t *testing.T) {
		app := fiber.New(fiber.Config{
			ReadBufferSize: 8192,
		})

		d := test.SetupMockDB()
		cs := cache.NewMockService()
		svcs, err := server.GetServices(*cnf, cs, d)
		if err != nil {
			t.Errorf("error getting services: %s", err)
			return
		}

		authClient := &sdk.Client{
			Id: "test-client",
		}

		prv := server.SetupTestServer(app, cnf, svcs, cs, d)
		prv.AuthClient = authClient

		app.Use(providers.Handle(prv))

		// Create a test user and set it in context
		testUser := &sdk.User{
			Id:    "user-123",
			Email: "test@example.com",
			Name:  "Test User",
		}

		// Add middleware to set user in context
		app.Use(func(c *fiber.Ctx) error {
			c.Context().SetUserValue(sdk.UserTypeVal, testUser)
			return c.Next()
		})

		// Register the route directly
		app.Get("/dashboard", DashboardMe)

		req, _ := http.NewRequest("GET", "/dashboard", nil)
		res, err := app.Test(req, -1)
		assert.Nil(t, err)
		assert.Equal(t, http.StatusOK, res.StatusCode)

		var resp sdk.DashboardUserResponse
		err = json.NewDecoder(res.Body).Decode(&resp)
		assert.Nil(t, err)
		assert.True(t, resp.Success)
		assert.Equal(t, "User fetched successfully", resp.Message)
		assert.Equal(t, testUser, resp.Data.User)
		assert.True(t, resp.Data.Setup.ClientAdded)
		assert.Equal(t, "test-client", resp.Data.Setup.ClientId)
	})
}
