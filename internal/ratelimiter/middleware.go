// Package ratelimiter provides HTTP middleware for rate limiting.
package ratelimiter

import (
	"net/http"
)

// RateLimitMiddleware returns a middleware that checks rate limits before passing to the next handler.
// Usage: wrap your handler with this middleware, providing a function to extract the client ID from the request.
func RateLimitMiddleware(rl *RateLimiter, extractID func(*http.Request) string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			clientID := extractID(r)
			if clientID == "" {
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			allowed, retryAfter := rl.Allow(clientID)
			if !allowed {
				w.Header().Set("Retry-After", retryAfter.String())
				w.WriteHeader(http.StatusTooManyRequests)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
} 