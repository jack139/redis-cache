package helper

import (
	"log"
	"time"
	"context"
	"github.com/bsm/redislock"
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

// 查找缓存， 并更新
func Redis_shoot(key string) (string, error) {
	ctx := context.Background()

	// 查找 redis
	value, err := rdb.Get(ctx, key).Result()
	if err != redis.Nil { // 找到key 或 出错
		if err != nil {
			return "", err
		}

		log.Println("redis shot:", key, value)
		return value, nil
	}

	// 锁
	locker := redislock.New(rdb)

	// Retry every 100ms, for up-to 3x
	backoff := redislock.LimitRetry(redislock.LinearBackoff(100*time.Millisecond), 3)

	// Obtain lock with retry
	lock, err := locker.Obtain(ctx, key, 1*time.Second, &redislock.Options{
		RetryStrategy: backoff,
	})
	if err == redislock.ErrNotObtained {
		// 未取得 lock，用db数据返回
		log.Println("ErrNotObtained shot:", key)

		return Mssql_shoot(key) 
	} else if err != nil {
		return "", err
	}
	defer lock.Release(ctx)


	// 读取数据库
	value, err = Mssql_shoot(key)
	if err != nil {
		return "", err
	}

	// 判断锁失效
	if ttl, err := lock.TTL(ctx); err != nil {
		return "", err
	} else if ttl == 0 {
		// 锁已失效，直接返回数据
		log.Println("Lock expired:", key, value)

		return value, nil
	}

	// 写redis, 释放锁
	err = rdb.Set(ctx, key, value, 24*time.Hour).Err()
	if err != nil {
		return "", err
	}

	log.Println("Update redis:", key, value)

	return value, nil
}
