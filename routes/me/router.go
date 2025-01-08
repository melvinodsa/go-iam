package me

import "github.com/gofiber/fiber/v2"

func RegisterRoutes(router fiber.Router) {
	v1 := router.Group("/v1")
	v1.Get("/", Me)
}
