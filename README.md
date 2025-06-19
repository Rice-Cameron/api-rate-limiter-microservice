# API Rate Limiter Microservice (Go)

A production-ready backend microservice for API rate limiting using the **Token Bucket** algorithm and Redis. Designed for high performance, clarity, and easy integration.

## Features

- **Token Bucket** algorithm (allows bursts, smooths traffic)
- **Global and per-client** rate limits
- **Redis** for distributed, atomic state
- **REST API**: `POST /check` to verify if a request is allowed
- **Configurable** via environment variables
- **Dockerized** for local development
- **Middleware** for easy integration
- **Unit tests** for core logic and Redis

## Why Token Bucket?

- Allows short bursts while enforcing a long-term average rate
- More flexible than fixed/sliding window for real-world APIs
- Efficient with Redis atomic operations

## Architecture

```
Client → [POST /check] → Go API → RateLimiter (Token Bucket) → Redis
```

- Each client (by IP, API key, etc) has a token bucket in Redis
- On each request, the bucket is checked/updated atomically
- If tokens remain, request is allowed; else, denied with `retry_after`

## Setup

### Prerequisites

- Docker & Docker Compose

### Run Locally

```
git clone <repo>
cd api-rate-limiter-microservice
docker-compose up --build
```

### Environment Variables

- `REDIS_URL` (default: redis://localhost:6379)
- `RATE_LIMIT` (default: 100)
- `WINDOW_SIZE` (default: 1m)
- `PORT` (default: 8080)

## API Usage

### Check Rate Limit

```
POST /check
Content-Type: application/json
{
  "client_id": "user-123"
}
```

**Response:**

- Allowed: `{ "allowed": true }`
- Rate limited: `{ "allowed": false, "retry_after": 42 }`

## Integration

- Use the provided Go middleware for easy plug-in to other services.

## Testing

```
docker-compose exec api-rate-limiter go test ./...
```

## Extending

- Add per-client custom limits in `rateLimiter.go`
- Add dashboard/CLI for usage stats (bonus)

## License

MIT
