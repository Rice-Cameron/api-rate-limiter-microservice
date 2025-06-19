package main

import (
	"testing"
	"time"
)

type mockRedis struct {
	count int
	window time.Duration
}

func (m *mockRedis) IncrWithTTL(key string, window time.Duration) (int, time.Duration, error) {
	m.count++
	return m.count, window, nil
}

func TestRateLimiter_Allow(t *testing.T) {
	redis := &mockRedis{}
	rl := &RateLimiter{
		redis:        redis,
		globalLimit:  2,
		globalWindow: time.Second,
	}
	allowed, _ := rl.Allow("client1")
	if !allowed {
		t.Fatal("first request should be allowed")
	}
	allowed, _ = rl.Allow("client1")
	if !allowed {
		t.Fatal("second request should be allowed")
	}
	allowed, _ = rl.Allow("client1")
	if allowed {
		t.Fatal("third request should be rate limited")
	}
} 