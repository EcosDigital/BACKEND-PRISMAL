package core

import (
	"context"
	"fmt"

	"github.com/ecosistema/core/src/shared/logging"
	"github.com/redis/go-redis/v9"
)

var RedisClient *redis.Client

func InitRedis() error {
	RedisClient = redis.NewClient(&redis.Options{
		Addr:     Cfg.RedisAddr,
		Password: Cfg.RedisPassword,
		DB:       0,
	})

	ctx := context.Background()
	if err := RedisClient.Ping(ctx).Err(); err != nil {
		logging.Error.Printf("redis connection failed: %v", err)
	}

	fmt.Println("✓ Redis conectado exitosamente")
	return nil

}
