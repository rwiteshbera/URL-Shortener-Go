package database

import (
	"github.com/redis/go-redis/v9"
	"os"
)

func CreateRedisClient(DatabaseNumber int) *redis.Client {
	db := redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_ADDRESS"),
		Password: os.Getenv("REDIS_PASSWORD"),
		DB:       DatabaseNumber,
	})
	return db
}
