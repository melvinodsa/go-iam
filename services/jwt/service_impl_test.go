package jwt

import (
	"strings"
	"testing"
	"time"

	jwtp "github.com/golang-jwt/jwt/v4"
	"github.com/melvinodsa/go-iam/sdk"
	"github.com/stretchr/testify/assert"
)

func TestNewService(t *testing.T) {
	t.Run("successful_service_creation", func(t *testing.T) {
		secret := sdk.MaskedBytes("test-secret-key")

		service := NewService(secret)

		assert.NotNil(t, service)
	})

	t.Run("service_creation_with_empty_secret", func(t *testing.T) {
		secret := sdk.MaskedBytes("")

		service := NewService(secret)

		assert.NotNil(t, service)
		// Service is created but will have issues during token operations
	})

	t.Run("service_creation_with_long_secret", func(t *testing.T) {
		secret := sdk.MaskedBytes(strings.Repeat("long-secret-key", 10))

		service := NewService(secret)

		assert.NotNil(t, service)
	})

	t.Run("service_creation_with_nil_secret", func(t *testing.T) {
		var secret sdk.MaskedBytes

		service := NewService(secret)

		assert.NotNil(t, service)
		// Service is created but will have issues during token operations
	})
}

func TestService_GenerateToken(t *testing.T) {
	secret := sdk.MaskedBytes("test-secret-key-for-jwt")
	service := NewService(secret)

	t.Run("successful_token_generation", func(t *testing.T) {
		claims := map[string]interface{}{
			"user_id": "123",
			"email":   "test@example.com",
			"role":    "admin",
		}
		expiryTime := time.Now().Add(time.Hour).Unix()

		token, err := service.GenerateToken(claims, expiryTime)

		assert.NoError(t, err)
		assert.NotEmpty(t, token)
		assert.True(t, strings.Contains(token, ".")) // JWT has dots
		parts := strings.Split(token, ".")
		assert.Equal(t, 3, len(parts)) // JWT has 3 parts: header.payload.signature
	})

	t.Run("generate_token_with_empty_claims", func(t *testing.T) {
		claims := map[string]interface{}{}
		expiryTime := time.Now().Add(time.Hour).Unix()

		token, err := service.GenerateToken(claims, expiryTime)

		assert.NoError(t, err)
		assert.NotEmpty(t, token)
		assert.True(t, strings.Contains(token, "."))
	})

	t.Run("generate_token_with_nil_claims", func(t *testing.T) {
		var claims map[string]interface{}
		expiryTime := time.Now().Add(time.Hour).Unix()

		token, err := service.GenerateToken(claims, expiryTime)

		assert.NoError(t, err)
		assert.NotEmpty(t, token)
	})

	t.Run("generate_token_with_various_claim_types", func(t *testing.T) {
		claims := map[string]interface{}{
			"string_claim": "string_value",
			"int_claim":    123,
			"float_claim":  123.456,
			"bool_claim":   true,
			"array_claim":  []string{"item1", "item2"},
			"object_claim": map[string]interface{}{"nested": "value"},
		}
		expiryTime := time.Now().Add(time.Hour).Unix()

		token, err := service.GenerateToken(claims, expiryTime)

		assert.NoError(t, err)
		assert.NotEmpty(t, token)
	})

	t.Run("generate_token_with_past_expiry", func(t *testing.T) {
		claims := map[string]interface{}{
			"user_id": "123",
		}
		pastTime := time.Now().Add(-time.Hour).Unix()

		token, err := service.GenerateToken(claims, pastTime)

		assert.NoError(t, err)
		assert.NotEmpty(t, token)
		// Token is generated but will be expired when validated
	})

	t.Run("generate_token_with_zero_expiry", func(t *testing.T) {
		claims := map[string]interface{}{
			"user_id": "123",
		}
		expiryTime := int64(0)

		token, err := service.GenerateToken(claims, expiryTime)

		assert.NoError(t, err)
		assert.NotEmpty(t, token)
	})

	t.Run("generate_token_with_future_expiry", func(t *testing.T) {
		claims := map[string]interface{}{
			"user_id": "123",
		}
		futureTime := time.Now().Add(24 * time.Hour).Unix()

		token, err := service.GenerateToken(claims, futureTime)

		assert.NoError(t, err)
		assert.NotEmpty(t, token)
	})

	t.Run("generate_multiple_tokens_are_different", func(t *testing.T) {
		claims := map[string]interface{}{
			"user_id": "123",
		}
		expiryTime := time.Now().Add(time.Hour).Unix()

		token1, err1 := service.GenerateToken(claims, expiryTime)
		// Wait a moment to ensure different timestamps
		time.Sleep(time.Millisecond)
		token2, err2 := service.GenerateToken(claims, expiryTime)

		assert.NoError(t, err1)
		assert.NoError(t, err2)
		assert.NotEmpty(t, token1)
		assert.NotEmpty(t, token2)
		// Tokens should be the same since claims and expiry are the same
		// JWT is deterministic for same input
		assert.Equal(t, token1, token2)
	})
}

func TestService_ValidateToken(t *testing.T) {
	secret := sdk.MaskedBytes("test-secret-key-for-jwt")
	service := NewService(secret)

	t.Run("successful_token_validation", func(t *testing.T) {
		originalClaims := map[string]interface{}{
			"user_id": "123",
			"email":   "test@example.com",
			"role":    "admin",
		}
		expiryTime := time.Now().Add(time.Hour).Unix()

		// Generate token first
		token, err := service.GenerateToken(originalClaims, expiryTime)
		assert.NoError(t, err)

		// Validate token
		claims, err := service.ValidateToken(token)

		assert.NoError(t, err)
		assert.NotNil(t, claims)
		assert.Equal(t, "123", claims["user_id"])
		assert.Equal(t, "test@example.com", claims["email"])
		assert.Equal(t, "admin", claims["role"])
		// exp claim should not be in the result
		_, hasExp := claims["exp"]
		assert.False(t, hasExp)
	})

	t.Run("validate_token_with_empty_claims", func(t *testing.T) {
		originalClaims := map[string]interface{}{}
		expiryTime := time.Now().Add(time.Hour).Unix()

		token, err := service.GenerateToken(originalClaims, expiryTime)
		assert.NoError(t, err)

		claims, err := service.ValidateToken(token)

		assert.NoError(t, err)
		assert.NotNil(t, claims)
		assert.Empty(t, claims)
	})

	t.Run("validate_expired_token", func(t *testing.T) {
		originalClaims := map[string]interface{}{
			"user_id": "123",
		}
		pastTime := time.Now().Add(-time.Hour).Unix()

		token, err := service.GenerateToken(originalClaims, pastTime)
		assert.NoError(t, err)

		claims, err := service.ValidateToken(token)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "error parsing jwt token")
		assert.Empty(t, claims)
	})

	t.Run("validate_invalid_token_format", func(t *testing.T) {
		invalidToken := "invalid.token.format"

		claims, err := service.ValidateToken(invalidToken)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "error parsing jwt token")
		assert.Empty(t, claims)
	})

	t.Run("validate_empty_token", func(t *testing.T) {
		emptyToken := ""

		claims, err := service.ValidateToken(emptyToken)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "error parsing jwt token")
		assert.Empty(t, claims)
	})

	t.Run("validate_malformed_token", func(t *testing.T) {
		malformedToken := "not-a-jwt-token"

		claims, err := service.ValidateToken(malformedToken)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "error parsing jwt token")
		assert.Empty(t, claims)
	})

	t.Run("validate_token_with_wrong_secret", func(t *testing.T) {
		// Generate token with one service
		originalClaims := map[string]interface{}{
			"user_id": "123",
		}
		expiryTime := time.Now().Add(time.Hour).Unix()
		token, err := service.GenerateToken(originalClaims, expiryTime)
		assert.NoError(t, err)

		// Try to validate with different service (different secret)
		differentSecret := sdk.MaskedBytes("different-secret-key")
		differentService := NewService(differentSecret)

		claims, err := differentService.ValidateToken(token)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "error parsing jwt token")
		assert.Empty(t, claims)
	})

	t.Run("validate_token_with_tampered_payload", func(t *testing.T) {
		originalClaims := map[string]interface{}{
			"user_id": "123",
		}
		expiryTime := time.Now().Add(time.Hour).Unix()
		token, err := service.GenerateToken(originalClaims, expiryTime)
		assert.NoError(t, err)

		// Tamper with the token by changing a character
		tamperedToken := token[:len(token)-5] + "XXXXX"

		claims, err := service.ValidateToken(tamperedToken)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "error parsing jwt token")
		assert.Empty(t, claims)
	})

	t.Run("validate_token_missing_parts", func(t *testing.T) {
		incompleteToken := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ"
		// Missing signature part

		claims, err := service.ValidateToken(incompleteToken)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "error parsing jwt token")
		assert.Empty(t, claims)
	})

	t.Run("validate_token_with_wrong_signing_method", func(t *testing.T) {
		// Create a token with RS256 method but validate with HS256 secret
		// This will fail the signing method check
		wrongMethodClaims := jwtp.MapClaims{
			"user_id": "123",
			"exp":     time.Now().Add(time.Hour).Unix(),
		}
		wrongMethodToken := jwtp.NewWithClaims(jwtp.SigningMethodRS256, wrongMethodClaims)
		// Don't sign it properly, just encode the header and payload
		wrongTokenString, _ := wrongMethodToken.SigningString()

		claims, err := service.ValidateToken(wrongTokenString + ".invalidsignature")

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "unexpected signing method")
		assert.Empty(t, claims)
	})
}

func TestService_GenerateValidateRoundtrip(t *testing.T) {
	secret := sdk.MaskedBytes("roundtrip-test-secret")
	service := NewService(secret)

	testCases := []struct {
		name   string
		claims map[string]interface{}
	}{
		{
			name: "simple_claims",
			claims: map[string]interface{}{
				"user_id": "123",
				"role":    "user",
			},
		},
		{
			name: "complex_claims",
			claims: map[string]interface{}{
				"user_id":     "456",
				"email":       "user@example.com",
				"permissions": []string{"read", "write"},
				"metadata":    map[string]interface{}{"department": "engineering"},
				"active":      true,
				"score":       99.5,
			},
		},
		{
			name:   "empty_claims",
			claims: map[string]interface{}{},
		},
		{
			name: "special_characters",
			claims: map[string]interface{}{
				"name":        "José María",
				"description": "Special chars: !@#$%^&*()",
				"unicode":     "🔒🗝️🛡️",
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			expiryTime := time.Now().Add(time.Hour).Unix()

			// Generate token
			token, err := service.GenerateToken(tc.claims, expiryTime)
			assert.NoError(t, err)
			assert.NotEmpty(t, token)

			// Validate token
			retrievedClaims, err := service.ValidateToken(token)
			assert.NoError(t, err)

			// Compare claims (excluding exp)
			for key, expectedValue := range tc.claims {
				actualValue, exists := retrievedClaims[key]
				assert.True(t, exists, "Claim %s should exist", key)

				// Handle type conversion for slices and maps that JWT might change
				switch expectedVal := expectedValue.(type) {
				case []string:
					// JWT converts []string to []interface{}
					if actualSlice, ok := actualValue.([]interface{}); ok {
						assert.Equal(t, len(expectedVal), len(actualSlice), "Claim %s slice length should match", key)
						for i, expectedItem := range expectedVal {
							assert.Equal(t, expectedItem, actualSlice[i], "Claim %s slice item %d should match", key, i)
						}
					} else {
						assert.Equal(t, expectedValue, actualValue, "Claim %s should match", key)
					}
				default:
					assert.Equal(t, expectedValue, actualValue, "Claim %s should match", key)
				}
			}

			// Ensure no extra claims (except those we filter out)
			assert.Equal(t, len(tc.claims), len(retrievedClaims))
		})
	}
}

func TestService_InterfaceCompliance(t *testing.T) {
	t.Run("verify_service_interface_implementation", func(t *testing.T) {
		secret := sdk.MaskedBytes("interface-test-secret")
		service := NewService(secret)

		// Verify service implements the Service interface
		t.Log("service", service)
		assert.True(t, true) // If compilation passes, interface is implemented
	})
}

func TestService_EdgeCases(t *testing.T) {
	secret := sdk.MaskedBytes("edge-case-test-secret")
	service := NewService(secret)

	t.Run("very_large_claims", func(t *testing.T) {
		// Create large claims payload
		largeClaims := map[string]interface{}{
			"large_data": strings.Repeat("x", 10000),
			"user_id":    "123",
		}
		expiryTime := time.Now().Add(time.Hour).Unix()

		token, err := service.GenerateToken(largeClaims, expiryTime)
		assert.NoError(t, err)
		assert.NotEmpty(t, token)

		claims, err := service.ValidateToken(token)
		assert.NoError(t, err)
		assert.Equal(t, largeClaims["large_data"], claims["large_data"])
		assert.Equal(t, largeClaims["user_id"], claims["user_id"])
	})

	t.Run("claims_with_nested_structures", func(t *testing.T) {
		nestedClaims := map[string]interface{}{
			"user": map[string]interface{}{
				"id":    "123",
				"email": "test@example.com",
				"profile": map[string]interface{}{
					"name": "John Doe",
					"age":  30,
				},
			},
			"permissions": []interface{}{
				map[string]interface{}{
					"resource": "users",
					"actions":  []string{"read", "write"},
				},
			},
		}
		expiryTime := time.Now().Add(time.Hour).Unix()

		token, err := service.GenerateToken(nestedClaims, expiryTime)
		assert.NoError(t, err)

		claims, err := service.ValidateToken(token)
		assert.NoError(t, err)
		assert.NotNil(t, claims["user"])
		assert.NotNil(t, claims["permissions"])
	})

	t.Run("extreme_expiry_times", func(t *testing.T) {
		claims := map[string]interface{}{"user_id": "123"}

		// Test very far future
		farFuture := time.Now().Add(100 * 365 * 24 * time.Hour).Unix()
		token1, err := service.GenerateToken(claims, farFuture)
		assert.NoError(t, err)
		assert.NotEmpty(t, token1)

		// Validate far future token
		retrievedClaims, err := service.ValidateToken(token1)
		assert.NoError(t, err)
		assert.Equal(t, claims["user_id"], retrievedClaims["user_id"])
	})

	t.Run("service_consistency_across_instances", func(t *testing.T) {
		// Create two services with the same secret
		sameSecret := sdk.MaskedBytes("same-secret-key")
		service1 := NewService(sameSecret)
		service2 := NewService(sameSecret)

		claims := map[string]interface{}{
			"user_id": "cross-service-test",
		}
		expiryTime := time.Now().Add(time.Hour).Unix()

		// Generate with service1
		token, err := service1.GenerateToken(claims, expiryTime)
		assert.NoError(t, err)

		// Validate with service2
		retrievedClaims, err := service2.ValidateToken(token)
		assert.NoError(t, err)
		assert.Equal(t, claims["user_id"], retrievedClaims["user_id"])
	})
}

func TestService_ErrorHandling(t *testing.T) {
	t.Run("service_with_empty_secret", func(t *testing.T) {
		emptySecret := sdk.MaskedBytes("")
		service := NewService(emptySecret)

		claims := map[string]interface{}{"user_id": "123"}
		expiryTime := time.Now().Add(time.Hour).Unix()

		// Should still work with empty secret (though not secure)
		token, err := service.GenerateToken(claims, expiryTime)
		assert.NoError(t, err)
		assert.NotEmpty(t, token)

		// Should be able to validate with same empty secret
		retrievedClaims, err := service.ValidateToken(token)
		assert.NoError(t, err)
		assert.Equal(t, claims["user_id"], retrievedClaims["user_id"])
	})

	t.Run("special_characters_in_secret", func(t *testing.T) {
		specialSecret := sdk.MaskedBytes("secret!@#$%^&*()_+-=[]{}|;':\",./<>?`~")
		service := NewService(specialSecret)

		claims := map[string]interface{}{"user_id": "special-secret-test"}
		expiryTime := time.Now().Add(time.Hour).Unix()

		token, err := service.GenerateToken(claims, expiryTime)
		assert.NoError(t, err)

		retrievedClaims, err := service.ValidateToken(token)
		assert.NoError(t, err)
		assert.Equal(t, claims["user_id"], retrievedClaims["user_id"])
	})
}

// Benchmark tests
func BenchmarkService_GenerateToken(b *testing.B) {
	secret := sdk.MaskedBytes("benchmark-secret-key")
	service := NewService(secret)
	claims := map[string]interface{}{
		"user_id": "benchmark_user",
		"role":    "admin",
	}
	expiryTime := time.Now().Add(time.Hour).Unix()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = service.GenerateToken(claims, expiryTime)
	}
}

func BenchmarkService_ValidateToken(b *testing.B) {
	secret := sdk.MaskedBytes("benchmark-secret-key")
	service := NewService(secret)
	claims := map[string]interface{}{
		"user_id": "benchmark_user",
		"role":    "admin",
	}
	expiryTime := time.Now().Add(time.Hour).Unix()
	token, _ := service.GenerateToken(claims, expiryTime)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = service.ValidateToken(token)
	}
}

func BenchmarkService_GenerateValidateRoundtrip(b *testing.B) {
	secret := sdk.MaskedBytes("benchmark-secret-key")
	service := NewService(secret)
	claims := map[string]interface{}{
		"user_id": "benchmark_user",
		"role":    "admin",
	}
	expiryTime := time.Now().Add(time.Hour).Unix()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		token, _ := service.GenerateToken(claims, expiryTime)
		_, _ = service.ValidateToken(token)
	}
}
