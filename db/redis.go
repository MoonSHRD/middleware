package db

import (
    "github.com/go-redis/redis"
)

type redi struct {
    instance *redis.Client
}

var redisInstance *redi

func GetRedis() *redis.Client {
    if redisInstance == nil {
        client := redis.NewClient(&redis.Options{
            Addr:     "localhost:6379",
            Password: "",
            DB:       0,
        })
        redisInstance = &redi{instance: client}
    }
    return redisInstance.instance
}