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
	"github.com/stretchr/testify/assert"
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
	routeFound := false
	for _, route := range routes {
		if route.Path == "/api/v1/" && route.Method == "GET" {
			routeFound = true
			break
		}
	}
	assert.True(t, routeFound, "Me route should be registered")
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