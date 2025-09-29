package routes

import (
	"os"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/melvinodsa/go-iam/config"
	"github.com/melvinodsa/go-iam/providers"
	"github.com/melvinodsa/go-iam/sdk"
	"github.com/melvinodsa/go-iam/services/cache"
	"github.com/melvinodsa/go-iam/utils/test"
	"github.com/melvinodsa/go-iam/utils/test/server"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRegisterOpenRoutes(t *testing.T) {
	err := os.Setenv("JWT_SECRET", "abcd")
	require.NoError(t, err)
	cnf := config.NewAppConfig()

	app := fiber.New()

	d := test.SetupMockDB()
	cs := cache.NewMockService()
	svcs, err := server.GetServices(*cnf, cs, d)
	require.NoError(t, err)

	authClient := &sdk.Client{Id: "test-client"}
	prv := &providers.Provider{
		S:          svcs,
		D:          d,
		C:          cs,
		AuthClient: authClient,
	}

	RegisterOpenRoutes(app, prv)

	// Check if routes are registered
	routes := app.GetRoutes()
	dashboardRouteFound := false
	healthRouteFound := false
	authRouteFound := false

	for _, route := range routes {
		if route.Path == "/me/v1/dashboard" && route.Method == "GET" {
			dashboardRouteFound = true
		}
		if route.Path == "/health/v1/" && route.Method == "GET" {
			healthRouteFound = true
		}
		if route.Path == "/auth/v1/login" && route.Method == "GET" {
			authRouteFound = true
		}
	}

	assert.True(t, dashboardRouteFound, "Dashboard route should be registered")
	assert.True(t, healthRouteFound, "Health route should be registered")
	assert.True(t, authRouteFound, "Auth route should be registered")
}

func TestRegisterAuthRoutes(t *testing.T) {
	err := os.Setenv("JWT_SECRET", "abcd")
	require.NoError(t, err)
	cnf := config.NewAppConfig()

	app := fiber.New()

	d := test.SetupMockDB()
	cs := cache.NewMockService()
	svcs, err := server.GetServices(*cnf, cs, d)
	require.NoError(t, err)

	authClient := &sdk.Client{Id: "test-client"}
	prv := &providers.Provider{
		S:          svcs,
		D:          d,
		C:          cs,
		AuthClient: authClient,
	}

	RegisterAuthRoutes(app, prv)

	// Check if authenticated routes are registered
	routes := app.GetRoutes()
	projectRouteFound := false
	clientRouteFound := false
	userRouteFound := false
	roleRouteFound := false
	policyRouteFound := false

	for _, route := range routes {
		if route.Path == "/project/v1/" && route.Method == "GET" {
			projectRouteFound = true
		}
		if route.Path == "/client/v1/" && route.Method == "GET" {
			clientRouteFound = true
		}
		if route.Path == "/user/v1/" && route.Method == "GET" {
			userRouteFound = true
		}
		if route.Path == "/role/v1/" && route.Method == "GET" {
			roleRouteFound = true
		}
		if route.Path == "/policy/v1/" && route.Method == "GET" {
			policyRouteFound = true
		}
	}

	assert.True(t, projectRouteFound, "Project route should be registered")
	assert.True(t, clientRouteFound, "Client route should be registered")
	assert.True(t, userRouteFound, "User route should be registered")
	assert.True(t, roleRouteFound, "Role route should be registered")
	assert.True(t, policyRouteFound, "Policy route should be registered")
}

func TestRegisterRoutes(t *testing.T) {
	err := os.Setenv("JWT_SECRET", "abcd")
	require.NoError(t, err)
	cnf := config.NewAppConfig()

	app := fiber.New()

	d := test.SetupMockDB()
	cs := cache.NewMockService()
	svcs, err := server.GetServices(*cnf, cs, d)
	require.NoError(t, err)

	authClient := &sdk.Client{Id: "test-client"}
	prv := &providers.Provider{
		S:          svcs,
		D:          d,
		C:          cs,
		AuthClient: authClient,
	}

	RegisterRoutes(app, prv)

	// Check if both open and auth routes are registered
	routes := app.GetRoutes()
	assert.NotEmpty(t, routes, "Routes should be registered")

	// Should have both open routes (health, auth) and protected routes (user, project, etc.)
	routeCount := len(routes)
	assert.Greater(t, routeCount, 10, "Should have multiple routes registered")
}