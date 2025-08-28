package sdk

import "time"

type AuthToken struct {
	AccessToken    string    `json:"access_token"`
	RefreshToken   string    `json:"refresh_token"`
	ExpiresAt      time.Time `json:"expires_at"`
	AuthProviderID string    `json:"auth_provider_id"`
	CodeChallenge  string    `json:"code_challenge"`
	CodeVerifier   string    `json:"code_verifier"`
	ClientId       string    `json:"client_id"`
}
