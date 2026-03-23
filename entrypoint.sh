#!/bin/sh
set -e

# Wait for Postgres to be ready before running migrations.
# pg_isready returns 0 when the server is accepting connections.
# This loop is especially important in local dev where the Postgres container
# and the app container start at the same time.
echo "Waiting for Postgres at ${DATABASE_HOST}:${DATABASE_PORT}..."
until pg_isready -h "$DATABASE_HOST" -p "$DATABASE_PORT" -U "$DATABASE_USERNAME" -q; do
  sleep 2
done
echo "Postgres is ready."

# Build the database URL for golang-migrate from individual env vars.
DATABASE_URL="postgres://${DATABASE_USERNAME}:${DATABASE_PASSWORD}@${DATABASE_HOST}:${DATABASE_PORT}/${DATABASE_NAME}?sslmode=disable"

# Run all pending migrations.
# golang-migrate is idempotent: if all migrations are already applied, this is a no-op.
# If a migration fails, the script exits (set -e) and the container does not start.
echo "Running database migrations..."
migrate -database "$DATABASE_URL" -path /app/db/migrations up
echo "Migrations complete."

# Start the server. exec replaces this shell with the server process,
# making the server PID 1 so it receives Docker signals directly.
exec /app/bin/server
