# =============================================================================
# Stage 1: Builder
# Uses the full Go toolchain to compile the binary and download the migrate CLI.
# This stage is discarded after build — none of its size ends up in the final image.
# =============================================================================
FROM golang:1.25.8-alpine AS builder

WORKDIR /app

# Download Go module dependencies first.
# Docker caches each instruction as a layer. By copying go.mod/go.sum before
# the source code, this layer is only invalidated when dependencies change —
# not on every source code edit. This makes rebuilds much faster.
COPY go.mod go.sum ./
RUN go mod download

# Install the golang-migrate CLI with the postgres driver tag.
# The binary will be at /go/bin/migrate inside this stage.
RUN go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@v4.19.1

# Copy source code and build the server binary.
# CGO_ENABLED=0: produces a fully static binary (no C library dependencies).
# GOOS=linux: cross-compile for Linux if building from macOS/Windows.
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o bin/server cmd/web/main.go

# =============================================================================
# Stage 2: Runtime
# Minimal Alpine image containing only what is needed to run the server.
# =============================================================================
FROM alpine:3.21

WORKDIR /app

# ca-certificates: required for outbound HTTPS connections.
# curl: used by the Docker HEALTHCHECK command.
# postgresql-client: provides pg_isready, used in entrypoint.sh to wait for
#   Postgres to be ready before running migrations.
RUN apk add --no-cache ca-certificates curl postgresql-client

# Create a non-root user and group.
# Running as non-root limits damage if the container is ever compromised.
RUN addgroup -S appgroup && adduser -S appuser -G appgroup

# Copy compiled artifacts from the builder stage.
COPY --from=builder /app/bin/server bin/server
COPY --from=builder /go/bin/migrate /usr/local/bin/migrate

# Copy runtime configuration and migration files.
# config.json contains non-secret app settings (port, log level, DB pool).
# db/migrations/ contains the SQL migration files the migrate CLI will run.
COPY config.json .
COPY db/migrations db/migrations
COPY entrypoint.sh .
RUN chmod +x entrypoint.sh bin/server

# Switch to non-root user for all subsequent commands and at runtime.
USER appuser

EXPOSE 3000

# HEALTHCHECK tells Docker how to verify the container is healthy.
# --interval=30s: check every 30 seconds.
# --timeout=5s: fail the check if it takes longer than 5 seconds.
# --start-period=10s: give the app 10 seconds to start before health checks begin.
# --retries=3: mark unhealthy after 3 consecutive failures.
HEALTHCHECK --interval=30s --timeout=5s --start-period=10s --retries=3 \
  CMD curl -f http://localhost:3000/health || exit 1

ENTRYPOINT ["./entrypoint.sh"]
