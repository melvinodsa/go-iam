package user

import "github.com/gofiber/fiber/v2"

func RegisterRoutes(router fiber.Router) {
	user := router.Group("/v1")
	user.Post("/", Create)
	user.Get("/:id", GetById)
	user.Get("/", GetAll)
	user.Put("/:id", Update)
	user.Put("/:id/roles", UpdateRoles)
}
