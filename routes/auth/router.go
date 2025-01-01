package auth

import "github.com/gofiber/fiber/v2"

func RegisterRoutes(router fiber.Router) {
	v1 := router.Group("/v1")
	v1.Get("/login", Login)
	v1.Get("/authp-callback", Redirect)
	v1.Get("/verify", Verify)
}
