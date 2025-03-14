package config

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"log"
	"os"
)

func NewRedis(ctx context.Context) *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr: os.Getenv("REDIS_CLIENT"),
	})

	result, err := client.Ping(ctx).Result()
	if err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}
	fmt.Printf("Connected to Redis successfully. %s \n", result)

	return client
}
