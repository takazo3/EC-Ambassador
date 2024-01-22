package database

import (
	"context"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
)

var Cache *redis.Client
var CacheChannel chan string

func SetupRedis() {
	Cache = redis.NewClient(&redis.Options{
		Addr:     "redis:6379",
		Password: "",
		DB:       0,
	})
}

func SetUpDeleteCache() {
	CacheChannel = make(chan string)

	go func(ch chan string) {
		for {
			key := <-ch
			time.Sleep(5 * time.Second)
			Cache.Del(context.Background(), key)
			fmt.Println("Cache deleted %s", key)
		}
	}(CacheChannel)
}

func ClearCache(keys ...string) {
	for _, key := range keys {
		CacheChannel <- key
	}
}
