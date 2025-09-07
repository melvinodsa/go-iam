package auth

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/melvinodsa/go-iam/sdk"
	"github.com/melvinodsa/go-iam/services/auth"
	"github.com/melvinodsa/go-iam/services/client"
	goiamclient "github.com/melvinodsa/go-iam/utils/goiamclient"
)

type Middlewares struct {
	authSvc    auth.Service
	clientSvc  client.Service
	AuthClient *sdk.Client
}

func NewMiddlewares(authSvc auth.Service, clientSvc client.Service) (*Middlewares, error) {
	authClient, err := goiamclient.GetGoIamClient(clientSvc)
	if err != nil {
		return nil, err
	}
	return &Middlewares{
		authSvc:    authSvc,
		clientSvc:  clientSvc,
		AuthClient: authClient,
	}, nil
}

func (m *Middlewares) User(c *fiber.Ctx) error {
	if m.AuthClient == nil {
		// If the auth client is not set, we cannot authenticate the user
		return c.Next()
	}
	// This middleware can be used to check if the user is authenticated
	user, err := m.GetUser(c)
	if err != nil {
		log.Warnw("failed to fetch user", "error", err)
		return c.Status(http.StatusUnauthorized).JSON(sdk.UserResponse{
			Success: false,
			Message: err.Error(),
		})
	}
	c.Context().SetUserValue(sdk.UserTypeVal, user)
	return c.Next()
}

func (m *Middlewares) DashboardUser(c *fiber.Ctx) error {
	if m.AuthClient == nil {
		// If the auth client is not set, we cannot authenticate the user
		return c.Next()
	}
	res := &sdk.DashboardUserResponse{
		Success: false,
	}
	res.Data.Setup.ClientAdded = m.AuthClient != nil
	res.Data.Setup.ClientId = m.AuthClient.Id
	// This middleware can be used to check if the user is authenticated
	user, err := m.GetUser(c)
	if err != nil {
		log.Warnw("failed to fetch user", "error", err)
		res.Message = err.Error()
		return c.Status(http.StatusUnauthorized).JSON(res)
	}

	c.Context().SetUserValue(sdk.UserTypeVal, user)
	return c.Next()
}

func (m *Middlewares) GetUser(c *fiber.Ctx) (*sdk.User, error) {
	authHeader := c.Get("Authorization")
	if authHeader == "" || len(authHeader) < 7 || authHeader[:7] != "Bearer " {
		return nil, errors.New("`Authorization` not found in header")
	}

	// Extract the token from the header
	token := authHeader[len("Bearer "):]
	if token == "" {
		return nil, errors.New("`Bearer` token not found in header")
	}

	user, err := m.authSvc.GetIdentity(c.Context(), token)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch identity: %w", err)
	}
	return user, nil
}
