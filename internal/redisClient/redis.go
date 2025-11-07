package redisClient

import (
	"context"
	"github.com/redis/go-redis/v9"
	"time"
)

const ctxTimeout = 2

type Client struct {
	Rdb *redis.Client
}

func New(addr, password string, db int, readTimeout, writeTimeout time.Duration) (*Client, error) {
	rbd := redis.NewClient(&redis.Options{
		Addr:         addr,
		Password:     password,
		DB:           db,
		ReadTimeout:  readTimeout * time.Second,
		WriteTimeout: writeTimeout * time.Second,
	})
	
	ctx, cancel := context.WithTimeout(context.Background(), ctxTimeout*time.Second)
	defer cancel()
	
	if err := rbd.Ping(ctx).Err(); err != nil {
		return nil, err
	}
	return &Client{rbd}, nil
}
