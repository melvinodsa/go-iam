package resource

import "github.com/gofiber/fiber/v2"

// RegisterRoutes registers all resource related routes
func RegisterRoutes(router fiber.Router) {
	resource := router.Group("/v1")
	resource.Post("/", Create)      // Create a new resource
	resource.Get("/:id", Get)       // Get a specific resource
	resource.Get("/search", Search) // Search resources
	resource.Put("/:id", Update)    // Update a resource
}
