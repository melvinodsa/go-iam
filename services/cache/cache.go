package cache

import (
	"github.com/melvinodsa/go-iam/services/mockredis"
	"github.com/redis/go-redis/v9"
)

type Service struct {
	// actual redis
	Redis *redis.Client

	MockRedisSvc *mockredis.RedisService // Updated to use the correct type
}
