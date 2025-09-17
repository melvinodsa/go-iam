package sdk

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/gofiber/fiber/v2"
)

// ErrAuthProviderNotFound is returned when a requested authentication provider cannot be found.
var ErrAuthProviderNotFound = errors.New("auth provider not found")

// AuthProviderType represents the type of external authentication provider.
type AuthProviderType string

const (
	// AuthProviderTypeGoogle represents Google OAuth2/OIDC authentication.
	AuthProviderTypeGoogle AuthProviderType = "GOOGLE"

	// AuthProviderTypeMicrosoft represents Microsoft Azure AD authentication.
	AuthProviderTypeMicrosoft AuthProviderType = "MICROSOFT"

	// AuthProviderTypeGitHub represents GitHub OAuth2 authentication.
	AuthProviderTypeGitHub AuthProviderType = "GITHUB"

	// AuthProviderTypeOIDC represents generic OpenID Connect authentication.
	AuthProviderTypeOIDC AuthProviderType = "OIDC"
)

// AuthProvider represents an external authentication provider configuration.
// This contains the settings and credentials needed to integrate with
// external identity providers like Google, Microsoft, or GitHub.
type AuthProvider struct {
	Id        string              `json:"id"`         // Unique identifier for the auth provider
	Name      string              `json:"name"`       // Display name of the auth provider
	Icon      string              `json:"icon"`       // Icon URL or identifier for UI display
	Provider  AuthProviderType    `json:"provider"`   // Type of the authentication provider
	Params    []AuthProviderParam `json:"params"`     // Configuration parameters for the provider
	ProjectId string              `json:"project_id"` // ID of the project this provider belongs to
	Enabled   bool                `json:"enabled"`    // Whether this provider is active
	CreatedAt *time.Time          `json:"created_at"` // Timestamp when provider was created
	UpdatedAt *time.Time          `json:"updated_at"` // Timestamp when provider was last updated
	CreatedBy string              `json:"created_by"` // ID of the user who created this provider
	UpdatedBy string              `json:"updated_by"` // ID of the user who last updated this provider
}

// GetParam retrieves the value of a configuration parameter by key.
// Returns an empty string if the parameter is not found.
func (a AuthProvider) GetParam(key string) string {
	for _, p := range a.Params {
		if p.Key == key {
			return p.Value
		}
	}
	return ""
}

// AuthProviderParam represents a configuration parameter for an authentication provider.
// These parameters contain provider-specific settings like client IDs, secrets, and URLs.
type AuthProviderParam struct {
	Label    string `json:"label"`     // Human-readable label for the parameter
	Value    string `json:"value"`     // The parameter value
	Key      string `json:"key"`       // Unique key identifying the parameter
	IsSecret bool   `json:"is_secret"` // Whether this parameter contains sensitive information
}

// AuthProviderResponse represents an API response containing a single authentication provider.
type AuthProviderResponse struct {
	Success bool          `json:"success"` // Indicates if the operation was successful
	Message string        `json:"message"` // Human-readable message about the operation
	Data    *AuthProvider `json:"data"`    // The auth provider data (present only on success)
}

// NewErrorAuthProviderResponse creates a new error response for auth provider operations.
func NewErrorAuthProviderResponse(msg string, status int, c *fiber.Ctx) error {
	return c.Status(status).JSON(ClientResponse{
		Success: false,
		Message: msg,
	})
}

// AuthProviderBadRequest returns a 400 Bad Request error response for auth provider operations.
func AuthProviderBadRequest(msg string, c *fiber.Ctx) error {
	return NewErrorAuthProviderResponse(msg, http.StatusBadRequest, c)
}

// AuthProviderNotFound returns a 404 Not Found error response for auth provider operations.
func AuthProviderNotFound(msg string, c *fiber.Ctx) error {
	return NewErrorAuthProviderResponse(msg, http.StatusNotFound, c)
}

// AuthProviderInternalServerError returns a 500 Internal Server Error response for auth provider operations.
func AuthProviderInternalServerError(msg string, c *fiber.Ctx) error {
	return NewErrorAuthProviderResponse(msg, http.StatusInternalServerError, c)
}

// AuthProvidersResponse represents an API response containing a list of authentication providers.
type AuthProvidersResponse struct {
	Success bool           `json:"success"` // Indicates if the operation was successful
	Message string         `json:"message"` // Human-readable message about the operation
	Data    []AuthProvider `json:"data"`    // Array of auth provider data
}

// NewErrorAuthProvidersResponse creates a new error response for auth providers list operations.
func NewErrorAuthProvidersResponse(msg string, status int, c *fiber.Ctx) error {
	return c.Status(status).JSON(AuthProvidersResponse{
		Success: false,
		Message: msg,
	})
}

// AuthProvidersInternalServerError returns a 500 Internal Server Error response for auth providers list operations.
func AuthProvidersInternalServerError(msg string, c *fiber.Ctx) error {
	return NewErrorAuthProvidersResponse(msg, http.StatusInternalServerError, c)
}

// ServiceProvider interface defines the contract for external authentication service implementations.
// This interface must be implemented by each supported authentication provider (Google, Microsoft, GitHub, etc.)
// to provide OAuth2/OIDC functionality.
type ServiceProvider interface {
	// GetAuthCodeUrl returns the authorization URL where users should be redirected for authentication.
	GetAuthCodeUrl(state string) string

	// VerifyCode exchanges an authorization code for access and refresh tokens.
	VerifyCode(ctx context.Context, code string) (*AuthToken, error)

	// RefreshToken uses a refresh token to obtain a new access token.
	RefreshToken(refreshToken string) (*AuthToken, error)

	// GetIdentity retrieves user identity information using an access token.
	GetIdentity(token string) ([]AuthIdentity, error)

	// HasRefreshTokenFlow indicates whether this provider supports refresh token flow.
	HasRefreshTokenFlow() bool
}

// AuthProviderQueryParams represents query parameters for filtering authentication providers.
type AuthProviderQueryParams struct {
	ProjectIds []string `json:"project_id"` // Filter providers by project IDs
}
