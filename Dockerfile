# ── Builder stage ──────────────────────────────────────────────────────────────
FROM golang:1.22-alpine AS builder

WORKDIR /app

# Copy dependency manifests BEFORE the rest of the source.
# Docker caches this layer independently — go mod download only re-runs when
# go.mod or go.sum actually change, not on every source file edit.
COPY go.mod go.sum ./
RUN go mod download

# Now copy the source and compile a static binary.
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o stockwise ./cmd/api

# ── Runtime stage ──────────────────────────────────────────────────────────────
FROM alpine:3.19

RUN apk --no-cache add ca-certificates

WORKDIR /app

COPY --from=builder /app/stockwise .

EXPOSE 8080

CMD ["./stockwise"]
