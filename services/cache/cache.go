package cache

import "github.com/redis/go-redis/v9"

type Service struct {
	Redis *redis.Client
}
