package repository

import (
	"context"
	"encoding/json"
	"github.com/romang000/winter-service/internal/redisClient"
	"time"
)

type Reading struct {
	Timestamp   time.Time `bson:"timestamp"`
	Temperature float64   `bson:"temperature"`
}

type CacheRepository struct {
	rdb *redisClient.Client
}

func NewCacheRepository(rdb *redisClient.Client) *CacheRepository {
	return &CacheRepository{rdb: rdb}
}

func (cr *CacheRepository) Set(ctx context.Context, city string, reading Reading) error {
	data, err := json.Marshal(&reading)
	if err != nil {
		return err
	}
	
	return cr.rdb.Rdb.RPush(ctx, city, string(data)).Err()
}
