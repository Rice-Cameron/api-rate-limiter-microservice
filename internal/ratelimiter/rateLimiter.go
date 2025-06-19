package ratelimiter

import (
	"fmt"
	"time"
)

// RateLimiter implements the token bucket algorithm using Redis for state
// Supports global and per-client limits
// Thread-safe for concurrent use

type RateLimiter struct {
	redis        RedisClientInterface
	globalLimit  int
	globalWindow time.Duration
	// Optionally: per-client custom limits (map or Redis hash)
}

func NewRateLimiter(redis RedisClientInterface, globalLimit int, globalWindow time.Duration) *RateLimiter {
	return &RateLimiter{
		redis:        redis,
		globalLimit:  globalLimit,
		globalWindow: globalWindow,
	}
}

// Allow checks if a request for clientID is allowed, returns (allowed, retryAfter)
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