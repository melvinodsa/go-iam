package policy

import "github.com/gofiber/fiber/v2"

func RegisterRoutes(router fiber.Router) {
	policy := router.Group("/v1")
	policy.Post("/", Create)
	policy.Get("/:id", Get)
	policy.Get("/", FetchAll)
	policy.Put("/:id", Update)
	policy.Delete("/:id", Delete)
	policy.Get("/role/:id", GetPoliciesByRoleId)
}
