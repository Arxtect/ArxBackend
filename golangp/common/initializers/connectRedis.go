package initializers

import (
	"context"
	"log"

	"github.com/Arxtect/ArxBackend/golangp/config"

	"github.com/go-redis/redis/v8"
)

var Rdb *redis.Client

func InitRedisClient(config *config.Config) {
	Rdb = redis.NewClient(&redis.Options{
		Addr:     config.RedisHost,
		Password: config.RedisPassword,
		DB:       config.RedisDB,
	})
	_, err := Rdb.Ping(context.Background()).Result()

	if err != nil {
		log.Println("Redis connection failed", err)
		panic(err)
	}
}
