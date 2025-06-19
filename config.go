package main

import (
	"os"
	"strconv"
	"time"
)

type AppConfig struct {
	Port         string
	RedisURL     string
	GlobalLimit  int
	GlobalWindow time.Duration
}

func LoadAppConfig() AppConfig {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	redisURL := os.Getenv("REDIS_URL")
	if redisURL == "" {
		redisURL = "redis://localhost:6379"
	}
	globalLimit := 100
	if v := os.Getenv("RATE_LIMIT"); v != "" {
		if n, err := strconv.Atoi(v); err == nil {
			globalLimit = n
		}
	}
	globalWindow := time.Minute
	if v := os.Getenv("WINDOW_SIZE"); v != "" {
		if d, err := time.ParseDuration(v); err == nil {
			globalWindow = d
		}
	}
	return AppConfig{
		Port:         port,
		RedisURL:     redisURL,
		GlobalLimit:  globalLimit,
		GlobalWindow: globalWindow,
	}
} 