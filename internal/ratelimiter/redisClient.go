package main

import (
	"context"
	"time"

	"github.com/go-redis/redis/v8"
)

type RedisClient struct {
	client *redis.Client
}

func NewRedisClient(url string) *RedisClient {
	opt, err := redis.ParseURL(url)
	if err != nil {
		panic(err)
	}
	client := redis.NewClient(opt)
	return &RedisClient{client: client}
}

// IncrWithTTL atomically increments a key and sets expiry if new
// Returns (count, ttl, error)
func (r *RedisClient) IncrWithTTL(key string, window time.Duration) (int, time.Duration, error) {
	ctx := context.Background()
	pipe := r.client.TxPipeline()
	incr := pipe.Incr(ctx, key)
	pipe.Expire(ctx, key, window)
	_, err := pipe.Exec(ctx)
	if err != nil {
		return 0, 0, err
	}
	count := int(incr.Val())
	ttl, err := r.client.TTL(ctx, key).Result()
	if err != nil {
		return count, 0, err
	}
	return count, ttl, nil
} 