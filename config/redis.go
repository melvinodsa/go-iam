package config

import "github.com/melvinodsa/go-iam/sdk"

// Redis holds Redis cache configuration settings.
// All fields are public and can be accessed directly.
type Redis struct {
	Host     string          `json:"host"`     // Redis server address (host:port)
	Password sdk.MaskedBytes `json:"password"` // Redis password (optional, stored as MaskedBytes for security)
	DB       int             `json:"db"`       // Redis database number to use
}
