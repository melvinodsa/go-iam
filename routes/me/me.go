package me

import (
	"fmt"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/melvinodsa/go-iam/providers"
	"github.com/melvinodsa/go-iam/sdk"
)

func Me(c *fiber.Ctx) error {
	// get access token from auth bearer token
	authHeader := c.Get("Authorization")
	if authHeader == "" || len(authHeader) < 7 || authHeader[:7] != "Bearer " {
		return c.Status(http.StatusUnauthorized).JSON(sdk.UserResponse{
			Success: false,
			Message: "Authorization not found in header",
		})
	}

	// Extract the token from the header
	token := authHeader[len("Bearer "):]
	if token == "" {
		return c.Status(http.StatusUnauthorized).JSON(sdk.UserResponse{
			Success: false,
			Message: "Bearer token not found in header",
		})
	}

	pr := providers.GetProviders(c)
	user, err := pr.S.Auth.GetIdentity(c.Context(), token)
	if err != nil {
		message := fmt.Errorf("failed to fetch user. %w", err).Error()
		log.Errorw("failed to fetch user", "error", err)
		return c.Status(http.StatusInternalServerError).JSON(sdk.UserResponse{
			Success: false,
			Message: message,
		})
	}
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
	authHeader := c.Get("Authorization")
	if authHeader == "" || len(authHeader) < 7 || authHeader[:7] != "Bearer " {
		res.Message = "Authorization not found in header"
		return c.Status(http.StatusUnauthorized).JSON(res)
	}

	// Extract the token from the header
	token := authHeader[len("Bearer "):]
	if token == "" {
		res.Message = "Bearer token not found in header"
		return c.Status(http.StatusUnauthorized).JSON(res)
	}

	user, err := pr.S.Auth.GetIdentity(c.Context(), token)
	if err != nil {
		message := fmt.Errorf("failed to fetch user. %w", err).Error()
		log.Errorw("failed to fetch user", "error", err)
		res.Message = message
		// if error is not found, return 500
		return c.Status(http.StatusInternalServerError).JSON(res)
	}
	log.Debug("user fetched successfully")
	res = sdk.DashboardUserResponse{
		Success: true,
		Message: "User fetched successfully",
	}
	res.Data.Setup.ClientAdded = true
	res.Data.User = user
	return c.Status(http.StatusOK).JSON(res)
}
