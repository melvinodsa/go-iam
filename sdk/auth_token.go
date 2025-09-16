package sdk

import "time"

// AuthToken represents a complete authentication token with OAuth2/OIDC information.
// This structure contains both access and refresh tokens along with metadata
// about the OAuth2 flow and associated entities.
type AuthToken struct {
	AccessToken          string    `json:"access_token"`            // JWT access token for API authentication
	RefreshToken         string    `json:"refresh_token"`           // Token used to refresh the access token
	ExpiresAt            time.Time `json:"expires_at"`              // Timestamp when the access token expires
	AuthProviderID       string    `json:"auth_provider_id"`        // ID of the authentication provider used
	CodeChallengeMethod  string    `json:"code_challenge_method"`   // PKCE code challenge method (e.g., "S256")
	CodeChallenge        string    `json:"code_challenge"`          // PKCE code challenge value
	ClientId             string    `json:"client_id"`               // OAuth2 client identifier
	ServiceAccountUserId string    `json:"service_account_user_id"` // Associated service account user ID (if applicable)
}
