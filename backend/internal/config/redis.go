package config

import (
	"context"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
)

const (
	// Redis key patterns
	RedisInvalidatedUserKey = "invalidated:user:%s"
	RedisActiveUsersKey     = "channel:%s:active_users:%s"
	RedisMembersCountKey    = "community:%s:member_count"
)

// NewRedisClient creates and returns a new Redis client using the global AppConfig.
func NewRedisClient() *redis.Client {
	// Create the client with configuration from the global Cfg variable.
	client := redis.NewClient(&redis.Options{
		Addr:     Cfg.Redis.Addr,
		Password: Cfg.Redis.Password,
		DB:       Cfg.Redis.DB,
	})

	// default timeout
	timeoutSec := 5
	timeoutStr := os.Getenv("REDIS_TIMEOUT")
	if timeoutStr != "" {
		if v, err := strconv.Atoi(timeoutStr); err == nil {
			timeoutSec = v
		} else {
			log.Printf("Invalid REDIS_TIMEOUT value: %s (using default 5s)", timeoutStr)
		}
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(timeoutSec)*time.Second)
	defer cancel()

	// Ping the Redis server to ensure the connection is alive.
	if err := client.Ping(ctx).Err(); err != nil {
		// Log a warning instead of a fatal error.
		// This allows the application to continue running even if Redis is unavailable.
		// Features that depend on Redis (like token invalidation) will be gracefully disabled.
		log.Printf("WARNING: Could not connect to Redis at %s. Features depending on Redis may be disabled. Error: %v", Cfg.Redis.Addr, err)
	} else {
		log.Println("Successfully connected to Redis.")
	}

	return client
}

func ResetRedisAppKeys(rdb *redis.Client) {
	ctx := context.Background()

	patterns := []string{
		"invalidated:user:*",
		"channel:*:active_users",
		"community:*:member_count",
	}

	for _, pattern := range patterns {
		iter := rdb.Scan(ctx, 0, pattern, 0).Iterator()

		for iter.Next(ctx) {
			key := iter.Val()
			if err := rdb.Del(ctx, key).Err(); err != nil {
				log.Printf("Failed to delete redis key %s: %v", key, err)
			}
		}

		if err := iter.Err(); err != nil {
			log.Printf("Error during scan for pattern %s: %v", pattern, err)
		}
	}
}
