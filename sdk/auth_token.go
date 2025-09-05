package sdk

import "time"

type AuthToken struct {
	AccessToken          string    `json:"access_token"`
	RefreshToken         string    `json:"refresh_token"`
	ExpiresAt            time.Time `json:"expires_at"`
	AuthProviderID       string    `json:"auth_provider_id"`
	CodeChallengeMethod  string    `json:"code_challenge_method"`
	CodeChallenge        string    `json:"code_challenge"`
	ClientId             string    `json:"client_id"`
	ServiceAccountUserId string    `json:"service_account_user_id"`
}
