package me

import (
	"github.com/gofiber/fiber/v2"
	"github.com/melvinodsa/go-iam/providers"
)

func RegisterRoutes(router fiber.Router, path string) {
	v1Path := path + "/v1"
	v1 := router.Group(v1Path)
	MeRoute(v1, v1Path)
}

func RegisterOpenRoutes(router fiber.Router, path string, prv *providers.Provider) {
	v1Path := path + "/v1"
	v1 := router.Group(v1Path)
	DashboardMeRoute(v1, v1Path, prv)
	v1.Get("/dashboard", DashboardMe)
}

var routeTags = []string{"Me"}
