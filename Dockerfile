# syntax=docker/dockerfile:1

# Build stage
FROM golang:1.22-alpine AS builder
WORKDIR /app
COPY . .
RUN go mod init api-rate-limiter-microservice || true
RUN go mod tidy
# Build the main binary from the cmd/api directory
RUN go build -o rate-limiter ./cmd/api

# Final image
FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/rate-limiter .
EXPOSE 8080
CMD ["./rate-limiter"] 