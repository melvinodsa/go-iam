package resource

import "github.com/gofiber/fiber/v2"

// RegisterRoutes registers all resource related routes
func RegisterRoutes(router fiber.Router) {
	router.Post("/", Create)            // Create a new resource
	router.Get("/:id", Get)             // Get a specific resource
	router.Get("/search", Search)         // Search resources
	router.Put("/:id", Update)          // Update a resource
}