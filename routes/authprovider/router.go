package authprovider

import "github.com/gofiber/fiber/v2"

func RegisterRoutes(router fiber.Router) {
	v1 := router.Group("/v1")
	v1.Post("/", Create)
	v1.Get("/:id", Get)
	v1.Get("/", FetchAll)
	v1.Put("/:id", Update)
}
