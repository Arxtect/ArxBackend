package initializers

import (
	"github.com/arxtect/ArxBackend/golangp/config"
	"context"
	"log"

	"github.com/go-redis/redis/v8"
	"github.com/toheart/functrace"
)

var Rdb *redis.Client

func InitRedisClient(config *config.Config) {
	defer functrace.Trace([]interface {
	}{config})()
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
