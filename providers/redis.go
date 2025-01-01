package providers

import (
	"github.com/melvinodsa/go-iam/config"
	"github.com/melvinodsa/go-iam/services/cache"
	"github.com/redis/go-redis/v9"
)

func NewCache(cnf config.AppConfig) *cache.Service {
	client := redis.NewClient(&redis.Options{
		Addr:     cnf.Redis.Host,
		Password: string(cnf.Redis.Password), // no password set
		DB:       0,                          // use default DB
	})
	return &cache.Service{
		Redis: client,
	}
}
