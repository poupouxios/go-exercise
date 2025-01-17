package main

import (
	"os"
	"github.com/redis/go-redis/v9"
)

var RedisClient *redis.Client

func initRedis(){
	redisHost := os.Getenv("REDIS_HOST") + ":" + os.Getenv("REDIS_PORT")
	RedisClient = redis.NewClient(&redis.Options{
		Addr: redisHost,
		Password: "",
		DB: 0,
	})
}

func getRedisClient() *redis.Client {
	if RedisClient == nil {
		initRedis()
	}
	return RedisClient
}	