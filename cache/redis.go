package cache

import (
"context"
"fmt"
"github.com/go-redis/redis/v8"
"crab/config"
"time"
)

var Client *redis.Client

func Init(conf *config.Redis) error {
	Client = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%v", conf.Host, conf.Port),
		Password: conf.Password,
		DB:       0,
	})
	ctx, cancel := context.WithTimeout(context.Background(), time.Second * 5)
	defer cancel()
	return Client.Ping(ctx).Err()
}