// Package auth provides authentication middleware for the Go IAM system.
// It handles user authentication via Bearer tokens and integrates with
// the Go IAM authentication service to validate and extract user information.
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

// Middlewares provides authentication middleware functionality.
// It encapsulates the authentication service, client service, and Go IAM client
// needed for user authentication and authorization.
type Middlewares struct {
	authSvc    auth.Service   // Authentication service for token validation
	clientSvc  client.Service // Client service for client operations
	AuthClient *sdk.Client    // Go IAM client configuration
}

// NewMiddlewares creates a new authentication middleware instance.
// It initializes the middleware with the required services and retrieves
// the Go IAM client configuration for authentication operations.
//
// Parameters:
//   - authSvc: Authentication service for token validation
//   - clientSvc: Client service for client operations
//
// Returns:
//   - *Middlewares: Configured middleware instance
//   - error: Error if Go IAM client retrieval fails
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

// User is a Fiber middleware that authenticates users via Bearer token.
// It extracts the Authorization header, validates the token, and stores
// the authenticated user in the request context. If authentication fails,
// it returns a 401 Unauthorized response.
//
// Usage:
//
//	app.Use(authMiddleware.User)
//
// Parameters:
//   - c: Fiber context containing the HTTP request
//
// Returns:
//   - error: nil on success, HTTP error response on authentication failure
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

// DashboardUser is a specialized Fiber middleware for dashboard authentication.
// It performs user authentication similar to User middleware but returns
// dashboard-specific response format including setup information about
// the Go IAM client configuration.
//
// Usage:
//
//	app.Use(authMiddleware.DashboardUser)
//
// Parameters:
//   - c: Fiber context containing the HTTP request
//
// Returns:
//   - error: nil on success, HTTP error response with dashboard format on failure
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

// GetUser extracts and validates the user from the Authorization header.
// This is a helper method used by the middleware functions to perform
// the actual token validation and user retrieval.
//
// Parameters:
//   - c: Fiber context containing the HTTP request with Authorization header
//
// Returns:
//   - *sdk.User: Authenticated user information
//   - error: Error if token is missing, invalid, or authentication fails
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
