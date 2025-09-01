package cache

import (
	"context"
	"testing"
	"time"

	"github.com/melvinodsa/go-iam/sdk"
	"github.com/stretchr/testify/assert"
)

func TestNewRedisService(t *testing.T) {
	host := "localhost:6379"
	password := sdk.MaskedBytes("testpassword")

	service := NewRedisService(host, password)

	assert.NotNil(t, service)
	redisService, ok := service.(*redisService)
	assert.True(t, ok)
	assert.NotNil(t, redisService.client)
}

func TestRedisService_WithMockService(t *testing.T) {
	// Use the existing mock service for comprehensive testing
	mockService := NewMockService()

	t.Run("set_and_get_operations", func(t *testing.T) {
		ctx := context.Background()
		key := "test-key"
		value := "test-value"
		ttl := time.Hour

		// Test Set
		err := mockService.Set(ctx, key, value, ttl)
		assert.NoError(t, err)

		// Test Get
		result, err := mockService.Get(ctx, key)
		assert.NoError(t, err)
		assert.Equal(t, value, result)
	})

	t.Run("get_non_existent_key", func(t *testing.T) {
		ctx := context.Background()
		key := "non-existent-key"

		result, err := mockService.Get(ctx, key)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "key not found")
		assert.Equal(t, "", result)
	})

	t.Run("delete_existing_key", func(t *testing.T) {
		ctx := context.Background()
		key := "delete-key"
		value := "delete-value"

		// Set the key first
		err := mockService.Set(ctx, key, value, 0)
		assert.NoError(t, err)

		// Delete the key
		err = mockService.Delete(ctx, key)
		assert.NoError(t, err)

		// Verify it's deleted
		result, err := mockService.Get(ctx, key)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "key not found")
		assert.Equal(t, "", result)
	})

	t.Run("delete_non_existent_key", func(t *testing.T) {
		ctx := context.Background()
		key := "non-existent-delete-key"

		err := mockService.Delete(ctx, key)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "key not found")
	})

	t.Run("expire_existing_key", func(t *testing.T) {
		ctx := context.Background()
		key := "expire-key"
		value := "expire-value"
		ttl := time.Millisecond * 100

		// Set the key first
		err := mockService.Set(ctx, key, value, 0)
		assert.NoError(t, err)

		// Set expiration
		err = mockService.Expire(ctx, key, ttl)
		assert.NoError(t, err)

		// Wait for expiration
		time.Sleep(ttl + time.Millisecond*50)

		// Verify it's expired
		result, err := mockService.Get(ctx, key)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "key not found")
		assert.Equal(t, "", result)
	})

	t.Run("expire_non_existent_key", func(t *testing.T) {
		ctx := context.Background()
		key := "non-existent-expire-key"
		ttl := time.Hour

		err := mockService.Expire(ctx, key, ttl)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "key not found")
	})

	t.Run("set_with_ttl", func(t *testing.T) {
		ctx := context.Background()
		key := "ttl-key"
		value := "ttl-value"
		ttl := time.Millisecond * 100

		// Set with TTL
		err := mockService.Set(ctx, key, value, ttl)
		assert.NoError(t, err)

		// Get immediately should work
		result, err := mockService.Get(ctx, key)
		assert.NoError(t, err)
		assert.Equal(t, value, result)

		// Wait for expiration
		time.Sleep(ttl + time.Millisecond*50)

		// Should be expired now
		result, err = mockService.Get(ctx, key)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "expired")
		assert.Equal(t, "", result)
	})

	t.Run("set_with_zero_ttl", func(t *testing.T) {
		ctx := context.Background()
		key := "no-ttl-key"
		value := "no-ttl-value"

		// Set without TTL (ttl = 0)
		err := mockService.Set(ctx, key, value, 0)
		assert.NoError(t, err)

		// Should persist
		result, err := mockService.Get(ctx, key)
		assert.NoError(t, err)
		assert.Equal(t, value, result)
	})
}

func TestRedisService_ParameterValidation(t *testing.T) {
	mockService := NewMockService()

	t.Run("empty_key_operations", func(t *testing.T) {
		ctx := context.Background()
		emptyKey := ""
		value := "test-value"
		ttl := time.Hour

		// Set with empty key should work (Redis allows it)
		err := mockService.Set(ctx, emptyKey, value, ttl)
		assert.NoError(t, err)

		// Get with empty key should work
		result, err := mockService.Get(ctx, emptyKey)
		assert.NoError(t, err)
		assert.Equal(t, value, result)

		// Delete with empty key should work
		err = mockService.Delete(ctx, emptyKey)
		assert.NoError(t, err)
	})

	t.Run("empty_value_set", func(t *testing.T) {
		ctx := context.Background()
		key := "empty-value-key"
		emptyValue := ""

		err := mockService.Set(ctx, key, emptyValue, 0)
		assert.NoError(t, err)

		result, err := mockService.Get(ctx, key)
		assert.NoError(t, err)
		assert.Equal(t, emptyValue, result)
	})

	t.Run("get_by_name_empty_name_handling", func(t *testing.T) {
		ctx := context.Background()

		// Test that the mock service handles empty names correctly
		result, err := mockService.Get(ctx, "")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "key not found")
		assert.Equal(t, "", result)
	})
}

func TestRedisService_ErrorScenarios(t *testing.T) {
	t.Run("service_error_handling", func(t *testing.T) {
		mockService := NewMockService()
		ctx := context.Background()

		// Test getting a key that was never set
		result, err := mockService.Get(ctx, "never-set-key")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "key not found")
		assert.Equal(t, "", result)

		// Test deleting a key that doesn't exist
		err = mockService.Delete(ctx, "never-set-key")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "key not found")

		// Test expiring a key that doesn't exist
		err = mockService.Expire(ctx, "never-set-key", time.Hour)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "key not found")
	})

	t.Run("nil_service_error_handling", func(t *testing.T) {
		// Test that the Get method specifically handles nil service
		var nilService *RedisService
		ctx := context.Background()
		key := "test-key"

		// Only test the methods that have explicit nil checks in the mock
		result, err := nilService.Get(ctx, key)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "redis service is nil")
		assert.Equal(t, "", result)

		// Test Delete nil check
		err = nilService.Delete(ctx, key)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "redis service is nil")

		// Test Expire nil check
		err = nilService.Expire(ctx, key, time.Hour)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "redis service is nil")

		// Note: Set method doesn't have nil check in mock, so we skip testing it
		// to avoid panic. In a real Redis client, this would be handled by the client library
	})
}

func TestRedisService_IntegrationScenarios(t *testing.T) {
	mockService := NewMockService()

	t.Run("complete_lifecycle_test", func(t *testing.T) {
		ctx := context.Background()
		key := "lifecycle-key"
		value := "lifecycle-value"
		newValue := "updated-value"
		ttl := time.Hour

		// Create
		err := mockService.Set(ctx, key, value, ttl)
		assert.NoError(t, err)

		// Read
		result, err := mockService.Get(ctx, key)
		assert.NoError(t, err)
		assert.Equal(t, value, result)

		// Update
		err = mockService.Set(ctx, key, newValue, ttl)
		assert.NoError(t, err)

		// Verify update
		result, err = mockService.Get(ctx, key)
		assert.NoError(t, err)
		assert.Equal(t, newValue, result)

		// Delete
		err = mockService.Delete(ctx, key)
		assert.NoError(t, err)

		// Verify deletion
		result, err = mockService.Get(ctx, key)
		assert.Error(t, err)
		assert.Equal(t, "", result)
	})

	t.Run("multiple_keys_operations", func(t *testing.T) {
		ctx := context.Background()
		keys := []string{"key1", "key2", "key3"}
		values := []string{"value1", "value2", "value3"}

		// Set multiple keys
		for i, key := range keys {
			err := mockService.Set(ctx, key, values[i], 0)
			assert.NoError(t, err)
		}

		// Get all keys
		for i, key := range keys {
			result, err := mockService.Get(ctx, key)
			assert.NoError(t, err)
			assert.Equal(t, values[i], result)
		}

		// Delete all keys
		for _, key := range keys {
			err := mockService.Delete(ctx, key)
			assert.NoError(t, err)
		}

		// Verify all deleted
		for _, key := range keys {
			result, err := mockService.Get(ctx, key)
			assert.Error(t, err)
			assert.Equal(t, "", result)
		}
	})

	t.Run("expire_workflow", func(t *testing.T) {
		ctx := context.Background()
		key := "expire-workflow-key"
		value := "expire-workflow-value"
		initialTTL := time.Hour
		newTTL := time.Millisecond * 100

		// Set with initial TTL
		err := mockService.Set(ctx, key, value, initialTTL)
		assert.NoError(t, err)

		// Update TTL to shorter duration
		err = mockService.Expire(ctx, key, newTTL)
		assert.NoError(t, err)

		// Wait for expiration
		time.Sleep(newTTL + time.Millisecond*50)

		// Verify expired
		result, err := mockService.Get(ctx, key)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "expired")
		assert.Equal(t, "", result)
	})
}

func TestRedisService_ContextHandling(t *testing.T) {
	mockService := NewMockService()

	t.Run("operations_with_cancelled_context", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		key := "context-key"
		value := "context-value"

		// Set operation before cancellation
		err := mockService.Set(ctx, key, value, 0)
		assert.NoError(t, err)

		// Cancel context
		cancel()

		// Operations should still work with mock service (it doesn't check context cancellation)
		result, err := mockService.Get(ctx, key)
		assert.NoError(t, err)
		assert.Equal(t, value, result)
	})

	t.Run("operations_with_timeout_context", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*10)
		defer cancel()

		key := "timeout-key"
		value := "timeout-value"

		// Operations should complete quickly with mock
		err := mockService.Set(ctx, key, value, 0)
		assert.NoError(t, err)

		result, err := mockService.Get(ctx, key)
		assert.NoError(t, err)
		assert.Equal(t, value, result)
	})
}

func TestRedisService_ServiceInterface(t *testing.T) {
	t.Run("verify_service_interface_implementation", func(t *testing.T) {
		// Verify that redisService implements Service interface
		var _ Service = &redisService{}

		// Verify that mock RedisService also implements Service interface
		var _ Service = &RedisService{}

		assert.True(t, true) // If compilation passes, interfaces are implemented correctly
	})

	t.Run("service_creation_with_different_parameters", func(t *testing.T) {
		testCases := []struct {
			name     string
			host     string
			password sdk.MaskedBytes
		}{
			{"localhost_with_password", "localhost:6379", sdk.MaskedBytes("password123")},
			{"remote_host_with_password", "redis.example.com:6379", sdk.MaskedBytes("secretpass")},
			{"localhost_no_password", "localhost:6379", sdk.MaskedBytes("")},
			{"custom_port", "localhost:6380", sdk.MaskedBytes("testpass")},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				service := NewRedisService(tc.host, tc.password)
				assert.NotNil(t, service)

				redisService, ok := service.(*redisService)
				assert.True(t, ok)
				assert.NotNil(t, redisService.client)
			})
		}
	})
}

func TestRedisService_EdgeCases(t *testing.T) {
	mockService := NewMockService()

	t.Run("very_long_key_and_value", func(t *testing.T) {
		ctx := context.Background()
		longKey := string(make([]byte, 1000))    // Very long key
		longValue := string(make([]byte, 10000)) // Very long value

		err := mockService.Set(ctx, longKey, longValue, 0)
		assert.NoError(t, err)

		result, err := mockService.Get(ctx, longKey)
		assert.NoError(t, err)
		assert.Equal(t, longValue, result)
	})

	t.Run("special_characters_in_key_and_value", func(t *testing.T) {
		ctx := context.Background()
		specialKey := "key:with:colons:and@symbols#and$percent%and^caret&and*asterisk"
		specialValue := "value with spaces\nand\tnewlines\rand\x00null\x01bytes"

		err := mockService.Set(ctx, specialKey, specialValue, 0)
		assert.NoError(t, err)

		result, err := mockService.Get(ctx, specialKey)
		assert.NoError(t, err)
		assert.Equal(t, specialValue, result)
	})

	t.Run("unicode_characters", func(t *testing.T) {
		ctx := context.Background()
		unicodeKey := "키_🔑_ключ_キー"
		unicodeValue := "값_💎_значение_バリュー"

		err := mockService.Set(ctx, unicodeKey, unicodeValue, 0)
		assert.NoError(t, err)

		result, err := mockService.Get(ctx, unicodeKey)
		assert.NoError(t, err)
		assert.Equal(t, unicodeValue, result)
	})
}

// Benchmark tests
func BenchmarkRedisService_Set(b *testing.B) {
	mockService := NewMockService()
	ctx := context.Background()
	key := "benchmark-key"
	value := "benchmark-value"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = mockService.Set(ctx, key, value, 0)
	}
}

func BenchmarkRedisService_Get(b *testing.B) {
	mockService := NewMockService()
	ctx := context.Background()
	key := "benchmark-key"
	value := "benchmark-value"

	// Setup
	_ = mockService.Set(ctx, key, value, 0)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = mockService.Get(ctx, key)
	}
}

func BenchmarkRedisService_Delete(b *testing.B) {
	mockService := NewMockService()
	ctx := context.Background()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		key := "benchmark-key-" + string(rune(i))
		value := "benchmark-value"

		// Setup key for deletion
		_ = mockService.Set(ctx, key, value, 0)

		// Benchmark the delete operation
		_ = mockService.Delete(ctx, key)
	}
}

func BenchmarkRedisService_Expire(b *testing.B) {
	mockService := NewMockService()
	ctx := context.Background()
	ttl := time.Hour

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		key := "benchmark-key-" + string(rune(i))
		value := "benchmark-value"

		// Setup key
		_ = mockService.Set(ctx, key, value, 0)

		// Benchmark the expire operation
		_ = mockService.Expire(ctx, key, ttl)
	}
}
