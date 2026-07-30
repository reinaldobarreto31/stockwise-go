# ── Stage 1: build ───────────────────────────────────────────────────────────
FROM golang:1.25-alpine AS builder

WORKDIR /app

# Download dependencies first (layer cache)
COPY go.mod go.sum ./
RUN go mod download

# Build binary
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o stockwise ./cmd/api/main.go

# ── Stage 2: run ─────────────────────────────────────────────────────────────
FROM alpine:3.20

# postgresql-client for pg_isready + psql (migrations); ca-certificates for TLS
RUN apk add --no-cache ca-certificates postgresql-client

WORKDIR /app

COPY --from=builder /app/stockwise .
COPY db/migrations ./db/migrations
COPY entrypoint.sh .
RUN chmod +x entrypoint.sh

EXPOSE 8080

ENTRYPOINT ["./entrypoint.sh"]
