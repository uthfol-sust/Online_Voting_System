package database

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisStore struct {
	Client *redis.Client
}

var Ctx = context.Background()

func NewRedis(addr, password string) *RedisStore {
	rdb := redis.NewClient(&redis.Options{
		Addr:         addr,
		Password:     password,
		DB:           0,
		DialTimeout:  5 * time.Second,
		ReadTimeout:  2 * time.Second,
		WriteTimeout: 2 * time.Second,
	})

	return &RedisStore{Client: rdb}
}

func (r *RedisStore) Ping(ctx context.Context) error {
	return r.Client.Ping(ctx).Err()
}
