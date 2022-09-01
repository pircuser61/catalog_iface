package main

import (
	"context"
	"fmt"

	"github.com/go-redis/redis"
	"gitlab.ozon.dev/pircuser61/catalog_iface/config"
)

var redisClient *redis.Client

func cacheListener(ctx context.Context) {
	redisClient = redis.NewClient(&redis.Options{
		Addr:     config.RedisAddr,
		DB:       config.RedisResponseDb,
		Password: config.RedisPassword})
	subs := redisClient.Subscribe("response")
	for {
		msg, err := subs.ReceiveMessage()
		if err == nil {
			fmt.Println(msg)
		} else {
			fmt.Println("Redis ReceiveMessage error:", err.Error())
		}
		fmt.Print("\nCatalog_iface>")
	}
}
