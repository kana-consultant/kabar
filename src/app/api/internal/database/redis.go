// internal/database/redis.go
package database

import (
	"context"
	"log"
	"os"

	"github.com/go-redis/redis/v8"
)

var RedisClient *redis.Client

func InitRedis() error {
	redisAddr := os.Getenv("REDIS_ADDR")
	if redisAddr == "" {
		redisAddr = "localhost:6379"
	}

	redisPassword := os.Getenv("REDIS_PASSWORD")
	redisDB := 0

	RedisClient = redis.NewClient(&redis.Options{
		Addr:     redisAddr,
		Password: redisPassword,
		DB:       redisDB,
	})

	// Test connection
	ctx := context.Background()
	if err := RedisClient.Ping(ctx).Err(); err != nil {
		return err
	}

	log.Println("Redis connected successfully")
	return nil
}

func CloseRedis() {
	if RedisClient != nil {
		RedisClient.Close()
		log.Println("Redis connection closed")
	}
}
