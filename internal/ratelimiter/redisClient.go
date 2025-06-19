// Package ratelimiter provides Redis client abstraction for rate limiting state.
package ratelimiter

import (
	"context"
	"time"

	"github.com/go-redis/redis/v8"
)

// RedisClientInterface defines the methods required for a Redis client in rate limiting.
type RedisClientInterface interface {
	IncrWithTTL(key string, window time.Duration) (int, time.Duration, error)
}

// RedisClient wraps a go-redis client for use in rate limiting.
type RedisClient struct {
	client *redis.Client
}

// NewRedisClient creates a new RedisClient from a Redis URL.
func NewRedisClient(url string) *RedisClient {
	opt, err := redis.ParseURL(url)
	if err != nil {
		panic(err)
	}
	client := redis.NewClient(opt)
	return &RedisClient{client: client}
}

// IncrWithTTL atomically increments a key and sets expiry if new.
// Returns (count, ttl, error). Used for token bucket logic.
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

// Ensure RedisClient implements RedisClientInterface
var _ RedisClientInterface = (*RedisClient)(nil) 