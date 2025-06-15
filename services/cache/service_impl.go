package cache

import (
	"context"
	"time"

	"github.com/melvinodsa/go-iam/sdk"
	"github.com/redis/go-redis/v9"
)

type redisService struct {
	client *redis.Client
}

func NewRedisService(host string, password sdk.MaskedBytes) Service {
	client := redis.NewClient(&redis.Options{
		Addr:     host,
		Password: string(password), // no password set
		DB:       0,                // use default DB
	})
	return &redisService{client}
}

func (s redisService) Set(ctx context.Context, key string, value string, ttl time.Duration) error {
	res := s.client.Set(ctx, key, value, ttl)
	return res.Err()
}
func (s redisService) Get(ctx context.Context, key string) (string, error) {
	res := s.client.Get(ctx, key)
	return res.Val(), res.Err()
}
func (s redisService) Delete(ctx context.Context, key string) error {
	res := s.client.Del(ctx, key)
	return res.Err()
}
func (s redisService) Expire(ctx context.Context, key string, ttl time.Duration) error {
	res := s.client.Expire(ctx, key, ttl)
	return res.Err()
}
