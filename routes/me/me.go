package me

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/melvinodsa/go-iam/middlewares"
	"github.com/melvinodsa/go-iam/providers"
	"github.com/melvinodsa/go-iam/sdk"
	"github.com/melvinodsa/go-iam/utils/docs"
)

func MeRoute(router fiber.Router, basePath string) {
	routePath := "/"
	path := basePath + routePath
	router.Get(routePath, Me)
	docs.RegisterApi(docs.ApiWrapper{
		Path:        path,
		Method:      http.MethodGet,
		Name:        "Get Me",
		Description: "Get current user information",
		Response: &docs.ApiResponse{
			Description: "User fetched successfully",
			Content:     new(sdk.UserResponse),
		},
		Tags:                 routeTags,
		ProjectIDNotRequired: true,
	})
}

func Me(c *fiber.Ctx) error {
	// get access token from auth bearer token
	user := middlewares.GetUser(c.Context())
	log.Debug("user fetched successfully")
	return c.Status(http.StatusOK).JSON(sdk.UserResponse{
		Success: true,
		Message: "User fetched successfully",
		Data:    user,
	})
}

func AuthClientCheck(c *fiber.Ctx) error {
	pr := providers.GetProviders(c)
	if pr.AuthClient == nil {
		res := sdk.DashboardUserResponse{
			Success: true,
			Message: "auth is not setup yet.",
		}
		res.Data.Setup.ClientAdded = false
		return c.Status(http.StatusOK).JSON(res)
	}
	return c.Next()
}

func DashboardMeRoute(router fiber.Router, basePath string, prv *providers.Provider) {
	routePath := "/dashboard"
	path := basePath + routePath
	router.Get(routePath, AuthClientCheck, prv.AM.DashboardUser, DashboardMe)
	docs.RegisterApi(docs.ApiWrapper{
		Path:        path,
		Method:      http.MethodGet,
		Name:        "Get Dashboard Me",
		Description: "Get current user information for dashboard",
		Response: &docs.ApiResponse{
			Description: "User fetched successfully",
			Content:     new(sdk.DashboardUserResponse),
		},
		Tags:                 routeTags,
		ProjectIDNotRequired: true,
	})
}

func DashboardMe(c *fiber.Ctx) error {
	pr := providers.GetProviders(c)
	res := sdk.DashboardUserResponse{
		Success: false,
	}
	res.Data.Setup.ClientAdded = true
	res.Data.Setup.ClientId = pr.AuthClient.Id
	// get access token from auth bearer token
	user := middlewares.GetUser(c.Context())
	log.Debug("user fetched successfully")
	res.Message = "User fetched successfully"
	res.Data.Setup.ClientAdded = true
	res.Data.User = user
	return c.Status(http.StatusOK).JSON(res)
}
