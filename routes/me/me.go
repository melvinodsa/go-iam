package me

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/melvinodsa/go-iam/middlewares/auth"
	"github.com/melvinodsa/go-iam/providers"
	"github.com/melvinodsa/go-iam/sdk"
)

func Me(c *fiber.Ctx) error {
	// get access token from auth bearer token
	user := auth.GetUser(c.Context())
	log.Debug("user fetched successfully")
	return c.Status(http.StatusOK).JSON(sdk.UserResponse{
		Success: true,
		Message: "User fetched successfully",
		Data:    user,
	})
}

func DashboardMe(c *fiber.Ctx) error {

	pr := providers.GetProviders(c)
	if pr.AuthClient == nil {
		res := sdk.DashboardUserResponse{
			Success: true,
			Message: "auth is not setup yet.",
		}
		res.Data.Setup.ClientAdded = false
		return c.Status(http.StatusOK).JSON(res)
	}
	res := sdk.DashboardUserResponse{
		Success: false,
	}
	res.Data.Setup.ClientAdded = true
	res.Data.Setup.ClientId = pr.AuthClient.Id
	// get access token from auth bearer token
	user := auth.GetUser(c.Context())
	log.Debug("user fetched successfully")
	res = sdk.DashboardUserResponse{
		Success: true,
		Message: "User fetched successfully",
	}
	res.Data.Setup.ClientAdded = true
	res.Data.User = user
	return c.Status(http.StatusOK).JSON(res)
}
