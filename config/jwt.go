package config

import "github.com/melvinodsa/go-iam/sdk"

// Jwt holds JWT token configuration settings.
type Jwt struct {
	secret sdk.MaskedBytes // JWT secret key (private field, use Secret() method to access)
}

// Secret returns the JWT secret key used for signing and verifying JWT tokens.
// The secret is stored as MaskedBytes for security purposes.
//
// Returns the JWT secret configured via JWT_SECRET environment variable.
func (j Jwt) Secret() sdk.MaskedBytes {
	return j.secret
}
