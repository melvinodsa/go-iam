package config

// ServiceAccount holds service account token configuration settings.
// These settings control the time-to-live (TTL) for service account tokens.
// All fields are public and can be accessed directly.
type ServiceAccount struct {
	AccessTokenTTLInMinutes int64 // Access token time-to-live in minutes
	RefreshTokenTTLInDays   int64 // Refresh token time-to-live in days
}
