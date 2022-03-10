package redis

import (
	"context"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/pkg/errors"
)

var (
	rdb = New()
	ctx = context.Background()
)

func New() *redis.Client {
	return redis.NewClient(&redis.Options{})
}

func WriteToRedis(key, value string) ([]byte, error) {
	expire := time.Second * 20
	_, err := rdb.SetEX(ctx, key, value, expire).Result()
	if err != nil {
		return nil, err
	}
	return nil, nil
}

func ReadFromRedis(key string) ([]byte, error) {
	result, err := rdb.Get(ctx, key).Result()
	if errors.Is(err, redis.Nil) {
		return nil, err
	}
	if err != nil {
		return nil, errors.Wrap(err, "failed to read cache from Redis")
	}
	return []byte(result), nil
}
