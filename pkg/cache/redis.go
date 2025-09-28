package cache

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisCache struct {
	client *redis.Client
	ctx    context.Context
}

type RedisConnectionSettings struct{
	Host string
	Port int
	Password string
	Db int
}
func NewRedisCache(options RedisConnectionSettings) *RedisCache {
	if options.Port == 0 {
		options.Port = 6379
	}
	fmt.Println(options)
	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", options.Host, options.Port),     // e.g. "localhost:6379"
		Password: options.Password, // "" if no password
		DB:       options.Db,       // 0 is default DB
	})
	ctx := context.Background()
	_, err := rdb.Ping(ctx).Result()
    if err != nil {
        panic(err)
    }
	return &RedisCache{
		client: rdb,
		ctx:    ctx,
	}
}

func (r *RedisCache) Get(key string) ([]byte, error) {
	val, err := r.client.Get(r.ctx, key).Bytes()
	if err == redis.Nil {
		// Key not found
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return val, nil
}

func (r *RedisCache) Set(key string, data []byte, expiresIn int64) error {
	return r.client.Set(r.ctx, key, data, time.Duration(expiresIn)*time.Second).Err()
}
