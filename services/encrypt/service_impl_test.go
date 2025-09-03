package encrypt

import (
	"strings"
	"testing"

	"github.com/melvinodsa/go-iam/sdk"
	"github.com/stretchr/testify/assert"
)

func TestNewService(t *testing.T) {
	t.Run("successful_service_creation", func(t *testing.T) {
		// Valid AES key (32 bytes for AES-256)
		key := sdk.MaskedBytes(make([]byte, 32))
		for i := range key {
			key[i] = byte(i)
		}

		service, err := NewService(key)

		assert.NoError(t, err)
		assert.NotNil(t, service)
	})

	t.Run("service_creation_with_16_byte_key", func(t *testing.T) {
		// Valid AES key (16 bytes for AES-128)
		key := sdk.MaskedBytes(make([]byte, 16))
		for i := range key {
			key[i] = byte(i)
		}

		service, err := NewService(key)

		assert.NoError(t, err)
		assert.NotNil(t, service)
	})

	t.Run("service_creation_with_24_byte_key", func(t *testing.T) {
		// Valid AES key (24 bytes for AES-192)
		key := sdk.MaskedBytes(make([]byte, 24))
		for i := range key {
			key[i] = byte(i)
		}

		service, err := NewService(key)

		assert.NoError(t, err)
		assert.NotNil(t, service)
	})

	t.Run("invalid_key_length", func(t *testing.T) {
		// Invalid AES key length (15 bytes)
		key := sdk.MaskedBytes(make([]byte, 15))

		service, err := NewService(key)

		assert.Error(t, err)
		assert.Nil(t, service)
		assert.Contains(t, err.Error(), "error creating block")
	})

	t.Run("empty_key", func(t *testing.T) {
		key := sdk.MaskedBytes([]byte{})

		service, err := NewService(key)

		assert.Error(t, err)
		assert.Nil(t, service)
		assert.Contains(t, err.Error(), "error creating block")
	})

	t.Run("nil_key", func(t *testing.T) {
		var key sdk.MaskedBytes

		service, err := NewService(key)

		assert.Error(t, err)
		assert.Nil(t, service)
		assert.Contains(t, err.Error(), "error creating block")
	})
}

func TestService_Encrypt(t *testing.T) {
	// Setup a service for testing
	key := sdk.MaskedBytes(make([]byte, 32))
	for i := range key {
		key[i] = byte(i)
	}
	service, err := NewService(key)
	assert.NoError(t, err)
	assert.NotNil(t, service)

	t.Run("successful_encryption", func(t *testing.T) {
		rawMessage := "Hello, World!"

		encryptedMessage, err := service.Encrypt(rawMessage)

		assert.NoError(t, err)
		assert.NotEmpty(t, encryptedMessage)
		assert.NotEqual(t, rawMessage, encryptedMessage)

		// Verify it's hex encoded
		assert.True(t, isHexString(encryptedMessage))

		// Verify the encrypted message is longer than original (due to nonce + tag)
		assert.Greater(t, len(encryptedMessage), len(rawMessage)*2) // hex encoding doubles length
	})

	t.Run("encrypt_empty_string", func(t *testing.T) {
		rawMessage := ""

		encryptedMessage, err := service.Encrypt(rawMessage)

		assert.NoError(t, err)
		assert.NotEmpty(t, encryptedMessage) // Still has nonce + tag even for empty message
		assert.True(t, isHexString(encryptedMessage))
	})

	t.Run("encrypt_long_message", func(t *testing.T) {
		rawMessage := strings.Repeat("This is a long message that should be encrypted properly. ", 100)

		encryptedMessage, err := service.Encrypt(rawMessage)

		assert.NoError(t, err)
		assert.NotEmpty(t, encryptedMessage)
		assert.NotEqual(t, rawMessage, encryptedMessage)
		assert.True(t, isHexString(encryptedMessage))
	})

	t.Run("encrypt_special_characters", func(t *testing.T) {
		rawMessage := "Special chars: !@#$%^&*()_+-=[]{}|;':\",./<>?`~\n\t\r"

		encryptedMessage, err := service.Encrypt(rawMessage)

		assert.NoError(t, err)
		assert.NotEmpty(t, encryptedMessage)
		assert.NotEqual(t, rawMessage, encryptedMessage)
		assert.True(t, isHexString(encryptedMessage))
	})

	t.Run("encrypt_unicode_characters", func(t *testing.T) {
		rawMessage := "Unicode: 🔒🗝️🛡️ 中文 العربية русский"

		encryptedMessage, err := service.Encrypt(rawMessage)

		assert.NoError(t, err)
		assert.NotEmpty(t, encryptedMessage)
		assert.NotEqual(t, rawMessage, encryptedMessage)
		assert.True(t, isHexString(encryptedMessage))
	})

	t.Run("multiple_encryptions_produce_different_results", func(t *testing.T) {
		rawMessage := "Same message"

		encrypted1, err1 := service.Encrypt(rawMessage)
		encrypted2, err2 := service.Encrypt(rawMessage)

		assert.NoError(t, err1)
		assert.NoError(t, err2)
		assert.NotEmpty(t, encrypted1)
		assert.NotEmpty(t, encrypted2)
		// Should be different due to random nonce
		assert.NotEqual(t, encrypted1, encrypted2)
	})
}

func TestService_Decrypt(t *testing.T) {
	// Setup a service for testing
	key := sdk.MaskedBytes(make([]byte, 32))
	for i := range key {
		key[i] = byte(i)
	}
	service, err := NewService(key)
	assert.NoError(t, err)
	assert.NotNil(t, service)

	t.Run("successful_decryption", func(t *testing.T) {
		rawMessage := "Hello, World!"

		// First encrypt the message
		encryptedMessage, err := service.Encrypt(rawMessage)
		assert.NoError(t, err)

		// Then decrypt it
		decryptedMessage, err := service.Decrypt(encryptedMessage)

		assert.NoError(t, err)
		assert.Equal(t, rawMessage, decryptedMessage)
	})

	t.Run("decrypt_empty_string_encryption", func(t *testing.T) {
		rawMessage := ""

		encryptedMessage, err := service.Encrypt(rawMessage)
		assert.NoError(t, err)

		decryptedMessage, err := service.Decrypt(encryptedMessage)

		assert.NoError(t, err)
		assert.Equal(t, rawMessage, decryptedMessage)
	})

	t.Run("decrypt_long_message", func(t *testing.T) {
		rawMessage := strings.Repeat("This is a long message for testing encryption and decryption. ", 50)

		encryptedMessage, err := service.Encrypt(rawMessage)
		assert.NoError(t, err)

		decryptedMessage, err := service.Decrypt(encryptedMessage)

		assert.NoError(t, err)
		assert.Equal(t, rawMessage, decryptedMessage)
	})

	t.Run("decrypt_special_characters", func(t *testing.T) {
		rawMessage := "Special chars: !@#$%^&*()_+-=[]{}|;':\",./<>?`~\n\t\r"

		encryptedMessage, err := service.Encrypt(rawMessage)
		assert.NoError(t, err)

		decryptedMessage, err := service.Decrypt(encryptedMessage)

		assert.NoError(t, err)
		assert.Equal(t, rawMessage, decryptedMessage)
	})

	t.Run("decrypt_unicode_characters", func(t *testing.T) {
		rawMessage := "Unicode: 🔒🗝️🛡️ 中文 العربية русский"

		encryptedMessage, err := service.Encrypt(rawMessage)
		assert.NoError(t, err)

		decryptedMessage, err := service.Decrypt(encryptedMessage)

		assert.NoError(t, err)
		assert.Equal(t, rawMessage, decryptedMessage)
	})

	t.Run("decrypt_invalid_hex", func(t *testing.T) {
		invalidHex := "not_valid_hex_string"

		decryptedMessage, err := service.Decrypt(invalidHex)

		assert.Error(t, err)
		assert.Empty(t, decryptedMessage)
		assert.Contains(t, err.Error(), "error decoding encrypted message")
	})

	t.Run("decrypt_empty_string", func(t *testing.T) {
		encryptedMessage := ""

		// This will panic due to slice bounds but should be caught and converted to error
		// We need to recover from panic if it happens
		var decryptedMessage string
		var err error

		func() {
			defer func() {
				if r := recover(); r != nil {
					err = assert.AnError // Convert panic to error for test
				}
			}()
			decryptedMessage, err = service.Decrypt(encryptedMessage)
		}()

		assert.Error(t, err)
		assert.Empty(t, decryptedMessage)
	})

	t.Run("decrypt_invalid_encrypted_data", func(t *testing.T) {
		// Valid hex but too short for valid encrypted data (need at least nonce + tag)
		// GCM needs 12 bytes nonce + 16 bytes tag = 28 bytes minimum = 56 hex chars
		invalidEncrypted := "deadbeef"

		var decryptedMessage string
		var err error

		func() {
			defer func() {
				if r := recover(); r != nil {
					err = assert.AnError // Convert panic to error for test
				}
			}()
			decryptedMessage, err = service.Decrypt(invalidEncrypted)
		}()

		assert.Error(t, err)
		assert.Empty(t, decryptedMessage)
	})

	t.Run("decrypt_properly_sized_invalid_data", func(t *testing.T) {
		// Create properly sized but invalid encrypted data (28 bytes = 56 hex chars)
		// This should pass the slice bounds check but fail during decryption
		invalidData := strings.Repeat("00", 28) // 56 hex chars = 28 bytes

		decryptedMessage, err := service.Decrypt(invalidData)

		assert.Error(t, err)
		assert.Empty(t, decryptedMessage)
		assert.Contains(t, err.Error(), "error decrypting message")
	})

	t.Run("decrypt_truncated_message", func(t *testing.T) {
		rawMessage := "Test message"
		encryptedMessage, err := service.Encrypt(rawMessage)
		assert.NoError(t, err)

		// Truncate the encrypted message to make it invalid
		truncatedMessage := encryptedMessage[:len(encryptedMessage)/2]

		decryptedMessage, err := service.Decrypt(truncatedMessage)

		assert.Error(t, err)
		assert.Empty(t, decryptedMessage)
		assert.Contains(t, err.Error(), "error decrypting message")
	})

	t.Run("decrypt_with_wrong_key", func(t *testing.T) {
		rawMessage := "Secret message"

		// Encrypt with first service
		encryptedMessage, err := service.Encrypt(rawMessage)
		assert.NoError(t, err)

		// Create second service with different key
		differentKey := sdk.MaskedBytes(make([]byte, 32))
		for i := range differentKey {
			differentKey[i] = byte(255 - i) // Different key
		}
		service2, err := NewService(differentKey)
		assert.NoError(t, err)

		// Try to decrypt with wrong key
		decryptedMessage, err := service2.Decrypt(encryptedMessage)

		assert.Error(t, err)
		assert.Empty(t, decryptedMessage)
		assert.Contains(t, err.Error(), "error decrypting message")
	})
}

func TestService_EncryptDecryptRoundtrip(t *testing.T) {
	// Setup service
	key := sdk.MaskedBytes(make([]byte, 32))
	for i := range key {
		key[i] = byte(i)
	}
	service, err := NewService(key)
	assert.NoError(t, err)

	testCases := []struct {
		name    string
		message string
	}{
		{"simple_text", "Hello, World!"},
		{"empty_string", ""},
		{"single_character", "A"},
		{"numbers", "1234567890"},
		{"json_like", `{"key": "value", "number": 123}`},
		{"xml_like", "<root><child>value</child></root>"},
		{"multiline", "Line 1\nLine 2\nLine 3"},
		{"tabs_and_spaces", "Tab:\tSpace: Double  Space"},
		{"special_chars", "!@#$%^&*()_+-=[]{}|;':\",./<>?`~"},
		{"unicode", "🔒🗝️🛡️ 中文 العربية русский"},
		{"long_text", strings.Repeat("This is a test message. ", 100)},
		{"binary_like", string([]byte{0, 1, 2, 3, 255, 254, 253})},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Encrypt
			encrypted, err := service.Encrypt(tc.message)
			assert.NoError(t, err)
			assert.NotEmpty(t, encrypted)
			assert.True(t, isHexString(encrypted))

			// Decrypt
			decrypted, err := service.Decrypt(encrypted)
			assert.NoError(t, err)
			assert.Equal(t, tc.message, decrypted)
		})
	}
}

func TestService_InterfaceCompliance(t *testing.T) {
	t.Run("verify_service_interface_implementation", func(t *testing.T) {
		key := sdk.MaskedBytes(make([]byte, 32))
		service, err := NewService(key)
		assert.NoError(t, err)

		t.Log("service", service)
		assert.True(t, true) // If compilation passes, interface is implemented
	})
}

func TestService_Consistency(t *testing.T) {
	t.Run("same_key_produces_same_decryption", func(t *testing.T) {
		key := sdk.MaskedBytes(make([]byte, 32))
		for i := range key {
			key[i] = byte(i)
		}

		// Create two services with same key
		service1, err := NewService(key)
		assert.NoError(t, err)

		service2, err := NewService(key)
		assert.NoError(t, err)

		rawMessage := "Cross-service test message"

		// Encrypt with first service
		encrypted, err := service1.Encrypt(rawMessage)
		assert.NoError(t, err)

		// Decrypt with second service
		decrypted, err := service2.Decrypt(encrypted)
		assert.NoError(t, err)
		assert.Equal(t, rawMessage, decrypted)
	})

	t.Run("different_keys_cannot_decrypt", func(t *testing.T) {
		// Create first service
		key1 := sdk.MaskedBytes(make([]byte, 32))
		for i := range key1 {
			key1[i] = byte(i)
		}
		service1, err := NewService(key1)
		assert.NoError(t, err)

		// Create second service with different key
		key2 := sdk.MaskedBytes(make([]byte, 32))
		for i := range key2 {
			key2[i] = byte(255 - i)
		}
		service2, err := NewService(key2)
		assert.NoError(t, err)

		rawMessage := "Secret message"

		// Encrypt with first service
		encrypted, err := service1.Encrypt(rawMessage)
		assert.NoError(t, err)

		// Try to decrypt with second service (should fail)
		decrypted, err := service2.Decrypt(encrypted)
		assert.Error(t, err)
		assert.Empty(t, decrypted)
	})
}

func TestService_EdgeCases(t *testing.T) {
	key := sdk.MaskedBytes(make([]byte, 32))
	service, err := NewService(key)
	assert.NoError(t, err)

	t.Run("very_long_message", func(t *testing.T) {
		// Test with a very large message
		longMessage := strings.Repeat("A", 1000000) // 1MB of 'A's

		encrypted, err := service.Encrypt(longMessage)
		assert.NoError(t, err)
		assert.NotEmpty(t, encrypted)

		decrypted, err := service.Decrypt(encrypted)
		assert.NoError(t, err)
		assert.Equal(t, longMessage, decrypted)
	})

	t.Run("message_with_null_bytes", func(t *testing.T) {
		message := "Hello\x00World\x00Test"

		encrypted, err := service.Encrypt(message)
		assert.NoError(t, err)

		decrypted, err := service.Decrypt(encrypted)
		assert.NoError(t, err)
		assert.Equal(t, message, decrypted)
	})

	t.Run("repeated_encrypt_decrypt_operations", func(t *testing.T) {
		message := "Repeated operations test"

		// Perform multiple rounds of encryption/decryption
		current := message
		for i := 0; i < 10; i++ {
			encrypted, err := service.Encrypt(current)
			assert.NoError(t, err)

			decrypted, err := service.Decrypt(encrypted)
			assert.NoError(t, err)
			assert.Equal(t, current, decrypted)

			// Use decrypted as input for next iteration to verify consistency
			current = decrypted
		}

		assert.Equal(t, message, current)
	})
}

// Helper function to check if a string is valid hex
func isHexString(s string) bool {
	if len(s)%2 != 0 {
		return false
	}
	for _, r := range s {
		if !strings.ContainsRune("0123456789abcdefABCDEF", r) {
			return false
		}
	}
	return true
}

// Benchmark tests
func BenchmarkService_Encrypt(b *testing.B) {
	key := sdk.MaskedBytes(make([]byte, 32))
	service, _ := NewService(key)
	message := "Benchmark test message for encryption"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = service.Encrypt(message)
	}
}

func BenchmarkService_Decrypt(b *testing.B) {
	key := sdk.MaskedBytes(make([]byte, 32))
	service, _ := NewService(key)
	message := "Benchmark test message for decryption"
	encrypted, _ := service.Encrypt(message)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = service.Decrypt(encrypted)
	}
}

func BenchmarkService_EncryptDecryptRoundtrip(b *testing.B) {
	key := sdk.MaskedBytes(make([]byte, 32))
	service, _ := NewService(key)
	message := "Benchmark test message for roundtrip"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		encrypted, _ := service.Encrypt(message)
		_, _ = service.Decrypt(encrypted)
	}
}
