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
	if authHeader == "" {
		return c.Status(http.StatusUnauthorized).JSON(sdk.UserResponse{
			Success: false,
			Message: "Authorizationnot found in header",
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
