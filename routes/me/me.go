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
		Parameters: []docs.ApiParameter{
			{
				Name:        "force_fetch",
				In:          "query",
				Description: "Force fetch user information",
				Required:    false,
			},
		},
		Tags:                 routeTags,
		ProjectIDNotRequired: true,
	})
}

func Me(c *fiber.Ctx) error {
	// get access token from auth bearer token
	user := middlewares.GetUser(c.Context())
	if user == nil {
		return c.Status(http.StatusUnauthorized).JSON(sdk.UserResponse{
			Success: false,
			Message: "Unauthorized",
		})
	}
	log.Debug("user fetched successfully")
	return c.Status(http.StatusOK).JSON(sdk.UserResponse{
		Success: true,
		Message: "User fetched successfully",
		Data:    user,
	})
}

func ResetPasswordRoute(router fiber.Router, basePath string) {
	routePath := "/password-reset"
	path := basePath + routePath
	router.Get(routePath, ResetPassword)
	docs.RegisterApi(docs.ApiWrapper{
		Path:        path,
		Method:      http.MethodGet,
		Name:        "Reset My Password",
		Description: "Redirect the authenticated user to the configured password reset page",
		Response: &docs.ApiResponse{
			Description: "Password reset URL generated successfully",
			Content:     new(sdk.AuthRedirectResponse),
		},
		Parameters: []docs.ApiParameter{
			{
				Name:        "postback",
				In:          "query",
				Description: "Whether to return the reset URL in the response instead of redirecting",
				Required:    false,
			},
		},
		Tags:                 routeTags,
		ProjectIDNotRequired: true,
	})
}

func ResetPassword(c *fiber.Ctx) error {
	user := middlewares.GetUser(c.Context())
	if user == nil {
		return c.Status(http.StatusUnauthorized).JSON(sdk.UserResponse{
			Success: false,
			Message: "Unauthorized",
		})
	}

	pr := providers.GetProviders(c)
	if pr.AuthClient == nil {
		return sdk.AuthProviderBadRequest("auth client is not setup", c)
	}
	if pr.AuthClient.DefaultAuthProviderId == "" {
		return sdk.AuthProviderBadRequest("password reset is not available for the configured auth provider", c)
	}

	resetURL, err := pr.S.Auth.GetResetPasswordUrl(c.Context(), pr.AuthClient.DefaultAuthProviderId)
	if err != nil {
		log.Errorw("failed to get password reset url", "error", err)
		return sdk.AuthProviderInternalServerError(err.Error(), c)
	}

	if c.Query("postback", "false") == "true" {
		return c.Status(http.StatusOK).JSON(sdk.AuthRedirectResponse{
			RedirectUrl: resetURL,
		})
	}

	return c.Redirect(resetURL, http.StatusTemporaryRedirect)
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
		Parameters: []docs.ApiParameter{
			{
				Name:        "force_fetch",
				In:          "query",
				Description: "Force fetch user information",
				Required:    false,
			},
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
	res.Success = true
	res.Message = "User fetched successfully"
	res.Data.Setup.ClientAdded = true
	res.Data.User = user
	return c.Status(http.StatusOK).JSON(res)
}
