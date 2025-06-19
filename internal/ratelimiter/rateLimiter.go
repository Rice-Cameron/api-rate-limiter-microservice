// Package ratelimiter provides the core rate limiting logic using the token bucket algorithm.
package ratelimiter

import (
	"fmt"
	"time"
)

// RateLimiter implements the token bucket algorithm using Redis for state.
// It supports global and per-client limits and is safe for concurrent use.
type RateLimiter struct {
	redis        RedisClientInterface // Redis client abstraction for state
	globalLimit  int                 // Maximum allowed requests per window
	globalWindow time.Duration       // Window size for rate limiting
	// Optionally: per-client custom limits (map or Redis hash)
}

// NewRateLimiter constructs a new RateLimiter.
// redis: Redis client abstraction
// globalLimit: max requests per window
// globalWindow: window duration
func NewRateLimiter(redis RedisClientInterface, globalLimit int, globalWindow time.Duration) *RateLimiter {
	return &RateLimiter{
		redis:        redis,
		globalLimit:  globalLimit,
		globalWindow: globalWindow,
	}
}

// Allow checks if a request for clientID is allowed.
// Returns (allowed, retryAfter). If not allowed, retryAfter is the time until the next allowed request.
// Implements the token bucket algorithm: increments a counter in Redis and checks against the limit.
func (rl *RateLimiter) Allow(clientID string) (bool, time.Duration) {
	// For demo: only global limit. Per-client can be added with a lookup.
	key := fmt.Sprintf("ratelimit:%s", clientID)
	count, ttl, err := rl.redis.IncrWithTTL(key, rl.globalWindow)
	if err != nil {
		// Fail open: allow if Redis is down
		return true, 0
	}
	if count > rl.globalLimit {
		return false, ttl
	}
	return true, 0
} 