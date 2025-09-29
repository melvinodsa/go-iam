package hashing

import (
	"crypto/sha256"
	"encoding/base64"
	"fmt"
)

func HashSecret(secret string) (string, error) {
	// hash the secret then convert it to base64
	if secret == "" {
		return "", fmt.Errorf("secret cannot be empty")
	}
	hashedSecret := sha256.Sum256([]byte(secret))
	return base64.StdEncoding.EncodeToString(hashedSecret[:]), nil
}
