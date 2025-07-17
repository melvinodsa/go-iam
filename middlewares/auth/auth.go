package auth

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/melvinodsa/go-iam/sdk"
	"github.com/melvinodsa/go-iam/services/auth"
)

type Middlewares struct {
	authSvc auth.Service
}

func NewMiddlewares(authSvc auth.Service) *Middlewares {
	return &Middlewares{
		authSvc: authSvc,
	}
}

func (m Middlewares) User(c *fiber.Ctx) error {
	// This middleware can be used to check if the user is authenticated
	// For now, we just pass the request to the next handler
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

	user, err := m.authSvc.GetIdentity(c.Context(), token)
	if err != nil {
		message := fmt.Errorf("failed to fetch user. %w", err).Error()
		log.Errorw("failed to fetch user", "error", err)
		return c.Status(http.StatusInternalServerError).JSON(sdk.UserResponse{
			Success: false,
			Message: message,
		})
	}
	c.Context().SetUserValue("user", user)
	return c.Next()
}

func GetUser(ctx context.Context) *sdk.User {
	user := ctx.Value("user")
	if user == nil {
		return nil
	}
	authUser, ok := user.(*sdk.User)
	if !ok {
		return nil
	}
	return authUser
}
