package policy

import (
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
)

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
	assert.True(t, routeFound, "Policy route should be registered")
}