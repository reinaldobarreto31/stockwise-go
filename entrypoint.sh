#!/bin/sh
set -e

# Wait for PostgreSQL to accept connections
echo "⏳  Waiting for PostgreSQL at ${DB_HOST}:${DB_PORT}..."
until pg_isready -h "${DB_HOST}" -p "${DB_PORT}" -U "${DB_USER}" -q; do
  sleep 1
done
echo "✅  PostgreSQL is ready"

# Apply migrations (idempotent — all use IF NOT EXISTS)
echo "🗄️   Running migrations..."
psql "${DATABASE_URL}" -f db/migrations/001_create_users.sql
psql "${DATABASE_URL}" -f db/migrations/002_create_products.sql
psql "${DATABASE_URL}" -f db/migrations/003_create_movements.sql
echo "✅  Migrations applied"

# Hand off to the API binary
exec ./stockwise
