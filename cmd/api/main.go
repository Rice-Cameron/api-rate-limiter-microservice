package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
	"api-rate-limiter-microservice/internal/ratelimiter"
	"github.com/joho/godotenv"
)

// Config holds environment-based configuration
type Config struct {
	Port         string
	RedisURL     string
	GlobalLimit  int
	GlobalWindow time.Duration
}

var config Config

func loadConfig() Config {
	_ = godotenv.Load()
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
		fmt.Sscanf(v, "%d", &globalLimit)
	}
	globalWindow := time.Minute
	if v := os.Getenv("WINDOW_SIZE"); v != "" {
		d, err := time.ParseDuration(v)
		if err == nil {
			globalWindow = d
		}
	}
	return Config{
		Port:         port,
		RedisURL:     redisURL,
		GlobalLimit:  globalLimit,
		GlobalWindow: globalWindow,
	}
}

func main() {
	config = loadConfig()
	redisClient := ratelimiter.NewRedisClient(config.RedisURL)
	rateLimiter := ratelimiter.NewRateLimiter(redisClient, config.GlobalLimit, config.GlobalWindow)

	http.HandleFunc("/check", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		var req struct {
			ClientID string `json:"client_id"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.ClientID == "" {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		allowed, retryAfter := rateLimiter.Allow(req.ClientID)
		resp := map[string]interface{}{"allowed": allowed}
		if !allowed {
			resp["retry_after"] = retryAfter.Seconds()
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	})

	log.Printf("Listening on :%s", config.Port)
	log.Fatal(http.ListenAndServe(":"+config.Port, nil))
} 