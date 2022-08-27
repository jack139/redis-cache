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
		//PoolSize: 500,
	})

	if _, err := rdb.Ping(context.Background()).Result(); err!=nil {
		return err
	}

	log.Println("Redis connected.")

	return nil
}



// 查找缓存， 并更新
func Redis_shoot(key string) (string, error) {
	var value string
	var err error

	lock_key := "__lock__" + key
	ctx := context.Background()

	// 查找 redis
	value, err = redis_get(key)
	if err != nil {
		return "", err
	}

	if value!="" { // 命中
		return value, nil
	}

	// 锁
	locker := redislock.New(rdb)

	// 取锁失败时 每50ms 重试2次
	backoff := redislock.LimitRetry(redislock.LinearBackoff(50*time.Millisecond), 2)
	// 取锁， 锁超时 2秒
	lock, err := locker.Obtain(ctx, lock_key, 2*time.Second, &redislock.Options{
		RetryStrategy: backoff,
	})
	if err == redislock.ErrNotObtained { // 未取得 lock，
		log.Println("ErrNotObtained lock:", lock_key)

		// 再次GET，如果还没命中，就直接从DB返回
		value, err = redis_get(key)
		if err != nil {
			return "", err
		}

		if value!="" {
			return value, nil
		} else {
			return Mssql_shoot(key)  // 没命中，直接从DB返回
		}
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
		log.Println("Lock expired:", key)

		return value, nil
	}

	// 写redis, 释放锁
	err = rdb.Set(ctx, key, value, 24*time.Hour).Err()
	if err != nil {
		return "", err
	}

	log.Println("Update redis:", key)

	return value, nil
}


// GET
func redis_get(key string) (string, error) {
	value, err := rdb.Get(context.Background(), key).Result()
	if err != redis.Nil { // 找到key 或 出错
		if err != nil {
			return "", err
		}

		log.Println("redis shot:", key)
		return value, nil
	} else {
		return "", nil		
	}
}