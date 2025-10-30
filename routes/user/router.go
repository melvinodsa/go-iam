package user

import "github.com/gofiber/fiber/v2"

func RegisterRoutes(router fiber.Router, path string) {
	v1Path := path + "/v1"
	v1 := router.Group(v1Path)
	CreateRoute(v1, v1Path)
	GetByIdRoute(v1, v1Path)
	GetAllRoute(v1, v1Path)
	UpdateRoute(v1, v1Path)
	UpdateRolesRoute(v1, v1Path)
	UpdatePoliciesRoute(v1, v1Path)
	TransferOwnershipRoute(v1, v1Path)
	CopyResourcesRoute(v1, v1Path)
}

var routeTags = []string{"User"}
