package hashing

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHashSecret(t *testing.T) {
	t.Run("hash_simple_secret", func(t *testing.T) {
		secret := "mysecret123"

		hashedSecret, err := HashSecret(secret)

		assert.NoError(t, err)
		assert.NotEmpty(t, hashedSecret)
		assert.NotEqual(t, secret, hashedSecret)
		// SHA256 hash encoded in base64 should be 44 characters long
		assert.Equal(t, 44, len(hashedSecret))
	})

	t.Run("hash_empty_secret", func(t *testing.T) {
		secret := ""

		hashedSecret, err := HashSecret(secret)

		assert.Error(t, err)
		assert.Empty(t, hashedSecret)
	})

	t.Run("hash_same_secret_produces_same_hash", func(t *testing.T) {
		secret := "consistent_secret"

		hash1, err1 := HashSecret(secret)
		hash2, err2 := HashSecret(secret)

		assert.NoError(t, err1)
		assert.NoError(t, err2)
		assert.Equal(t, hash1, hash2)
	})

	t.Run("hash_different_secrets_produce_different_hashes", func(t *testing.T) {
		secret1 := "secret1"
		secret2 := "secret2"

		hash1, err1 := HashSecret(secret1)
		hash2, err2 := HashSecret(secret2)

		assert.NoError(t, err1)
		assert.NoError(t, err2)
		assert.NotEqual(t, hash1, hash2)
	})

	t.Run("hash_special_characters", func(t *testing.T) {
		secret := "special!@#$%^&*()_+-={}[]|\\:;\"'<>,.?/~`"

		hashedSecret, err := HashSecret(secret)

		assert.NoError(t, err)
		assert.NotEmpty(t, hashedSecret)
		assert.Equal(t, 44, len(hashedSecret))
	})
}
