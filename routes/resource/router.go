package resource

import "github.com/gofiber/fiber/v2"

func RegisterRoutes(router fiber.Router) {
	resource := router.Group("/v1")
	resource.Post("/", Create)
	resource.Get("/:id", Get)
	resource.Get("/search", Search)
	resource.Put("/:id", Update)
	resource.Delete("/:id", Delete)
}
