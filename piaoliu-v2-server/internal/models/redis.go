package models

import (
	"context"
	"drift-bottle-server/internal/config"
	"log"

	"github.com/redis/go-redis/v9"
)

var RDB *redis.Client

func InitRedis(cfg *config.Config) {
	RDB = redis.NewClient(&redis.Options{
		Addr:     cfg.RedisAddr,
		Password: cfg.RedisPassword,
		DB:       cfg.RedisDB,
	})

	ctx := context.Background()
	if err := RDB.Ping(ctx).Err(); err != nil {
		log.Printf("Warning: Redis connection failed: %v, running without Redis", err)
		return
	}

	log.Println("Redis connected successfully")
}
