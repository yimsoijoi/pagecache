package main

import (
	"context"
	"fmt"

	"github.com/go-redis/redis/v8"
)

func main() {
	rdb := redis.NewClient(&redis.Options{})
	if _, err := rdb.Set(context.Background(), "kuy", "heeee", -1).Result(); err != nil {
		fmt.Println(err.Error())
	}
}
