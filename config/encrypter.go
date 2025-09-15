package config

import "github.com/melvinodsa/go-iam/sdk"

// Encrypter holds encryption configuration settings.
type Encrypter struct {
	key sdk.MaskedBytes // Encryption key (private field, use Key() method to access)
}

// Key returns the encryption key used for encrypting and decrypting sensitive data.
// The key is stored as MaskedBytes for security purposes.
//
// Returns the encryption key configured via ENCRYPTER_KEY environment variable.
func (e Encrypter) Key() sdk.MaskedBytes {
	return e.key
}
