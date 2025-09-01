package cache

import (
	"context"
	"errors"
	"sync"
	"time"
)

// RedisService represents a simple in-memory Redis-like service.
type RedisService struct {
	data map[string]string
	mu   sync.RWMutex
	ttl  map[string]time.Time
}

// NewRedisService creates a new instance of RedisService.
func NewMockService() *RedisService {
	return &RedisService{
		data: make(map[string]string),
		ttl:  make(map[string]time.Time),
	}
}

// Set stores a key-value pair in the Redis service with an optional TTL.
func (r *RedisService) Set(ctx context.Context, key string, value string, ttl time.Duration) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.data[key] = value
	if ttl > 0 {
		r.ttl[key] = time.Now().Add(ttl)
	} else {
		delete(r.ttl, key) // Remove TTL if no duration is provided
	}
	return nil
}

// Get retrieves the value for a given key from the Redis service.
// It returns an error if the key does not exist or has expired.
func (r *RedisService) Get(ctx context.Context, key string) (string, error) {
	if r == nil {
		return "", errors.New("redis service is nil")
	}

	r.mu.RLock()
	defer r.mu.RUnlock()

	// Check if the key exists
	value, exists := r.data[key]
	if !exists {
		return "", errors.New("key not found")
	}

	// Check if the key has expired
	if expiry, hasTTL := r.ttl[key]; hasTTL && time.Now().After(expiry) {
		// Key has expired, delete it
		r.mu.RUnlock() // Unlock read lock before acquiring write lock
		r.mu.Lock()
		delete(r.data, key)
		delete(r.ttl, key)
		r.mu.Unlock()
		r.mu.RLock() // Reacquire read lock
		return "", errors.New("key not found (expired)")
	}

	return value, nil
}

// Delete removes a key-value pair from the Redis service.
func (r *RedisService) Delete(ctx context.Context, key string) error {
	if r == nil {
		return errors.New("redis service is nil")
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.data[key]; !exists {
		return errors.New("key not found")
	}
	delete(r.data, key)
	delete(r.ttl, key)
	return nil
}

// Expire sets a TTL (time-to-live) for a key.
func (r *RedisService) Expire(ctx context.Context, key string, ttl time.Duration) error {
	if r == nil {
		return errors.New("redis service is nil")
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.data[key]; !exists {
		return errors.New("key not found")
	}
	r.ttl[key] = time.Now().Add(ttl)
	return nil
}
