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
	for _, route := range routes {
		if route.Path == "/me/v1/dashboard" && route.Method == "GET" {
			dashboardRouteFound = true
			break
		}
	}
	assert.True(t, dashboardRouteFound, "Dashboard route should be registered")
}