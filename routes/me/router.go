package me

import (
	"github.com/gofiber/fiber/v2"
	"github.com/melvinodsa/go-iam/providers"
)

func RegisterRoutes(router fiber.Router) {
	v1 := router.Group("/v1")
	v1.Get("/", Me)
}

func RegisterOpenRoutes(router fiber.Router, prv *providers.Provider) {
	v1 := router.Group("/v1")
	v1.Get("/dashboard", AuthClientCheck, prv.AM.User, DashboardMe)
}
