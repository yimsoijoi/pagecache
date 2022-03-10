package redis

import (
	"context"

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

func WriteToRedis(key, value string) error {
	// expire := time.Second * 20
	_, err := rdb.Set(ctx, key, value, -1).Result()
	if err != nil {
		return err
	}
	return nil
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
