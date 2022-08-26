package helper

import (
	"log"
	"context"
	"github.com/go-redis/redis/v8"
)

var (
	rdb *redis.Client
)

func redis_init() error {
	rdb = redis.NewClient(&redis.Options{
		Addr:     Settings.Server.REDIS_SERVER,
		Password: Settings.Server.REDIS_PASSWD,
		DB:       0,  // use default DB
	})

	if _, err := rdb.Ping(context.Background()).Result(); err!=nil {
		return err
	}

	log.Println("Redis connected.")

	return nil
}
