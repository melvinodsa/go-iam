package providers

import (
	"github.com/melvinodsa/go-iam/config"
	"github.com/redis/go-redis/v9"
)

type Cache struct {
	Redis *redis.Client
}

func NewCache(cnf config.AppConfig) *Cache {
	client := redis.NewClient(&redis.Options{
		Addr:     cnf.Redis.Host,
		Password: string(cnf.Redis.Password), // no password set
		DB:       0,                          // use default DB
	})
	return &Cache{
		Redis: client,
	}
}
