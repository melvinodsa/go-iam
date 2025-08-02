package auth

import "github.com/gofiber/fiber/v2"

func RegisterRoutes(router fiber.Router, path string) {
	v1Path := path + "/v1"
	v1 := router.Group(v1Path)
	LoginRoute(v1, v1Path)
	RedirectRoute(v1, v1Path)
	VerifyRoute(v1, v1Path)
}

var routeTags = []string{"Auth"}
