package database

import (
	"context"
	"os"

	"github.com/go-redis/redis/v9"
)

var Ctx = context.Background()

func CreateClient(DatabaseNo int) *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr:     os.Getenv("DB_ADDRESS"),
		Password: os.Getenv("DB_PASSWORD"),
		DB:       DatabaseNo,
	})

	return rdb
}
