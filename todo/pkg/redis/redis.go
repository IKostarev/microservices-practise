package redis

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
)

type Config struct {
	Address  string `envconfig:"REDIS_ADDRES" required:"true" default:"localhost:6379"`
	Password string `envconfig:"REDIS_PASSWORD" required:"true" default:"redisPassword"`
	Username string `envconfig:"REDIS_USERNAME" required:"true" default:"redisUser"`
	DB       int    `envconfig:"REDIS_DB" required:"true" default:"0"`
}

// RedisManager represents the Redis client manager.
type RedisManager struct {
	Client *redis.Client
}

// NewRedisManager creates a new RedisManager instance.
func NewRedisManager(config Config) (*RedisManager, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     config.Address,
		Password: config.Password,
		DB:       config.DB,
	})

	_, err := client.Ping(context.Background()).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to connect to Redis %s, %s, %d: %v", config.Address, config.Username, config.DB, err)
	}

	return &RedisManager{Client: client}, nil
}
