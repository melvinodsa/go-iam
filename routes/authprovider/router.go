package authprovider

import "github.com/gofiber/fiber/v2"

func RegisterRoutes(router fiber.Router, path string) {
	v1Path := path + "/v1"
	v1 := router.Group(v1Path)
	CreateRoute(v1, v1Path)
	GetRoute(v1, v1Path)
	FetchAllRoute(v1, v1Path)
	UpdateRoute(v1, v1Path)
	EnableServiceAccountRoute(v1, v1Path)
}

var routeTags = []string{"Auth Provider"}
