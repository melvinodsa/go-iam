package role

import "github.com/gofiber/fiber/v2"

func RegisterRoutes(router fiber.Router) {
	role := router.Group("/v1")
	role.Post("/", Create)
	role.Get("/:id", Get)
	role.Put("/:id", Update)
	role.Post("/:userid/:roleid", AddRoleToUser)
	role.Get("/:userid/:roleid", RemoveRoleFromUser)
}
