# Docker Operational Guide

This document explains how to build, run, and maintain the application using Docker.
For the design rationale and architecture decisions, see `docs/superpowers/specs/2026-03-23-docker-design.md`.

---

## Local Development

### Prerequisites

- Docker Desktop (macOS/Windows) or Docker Engine + Docker Compose plugin (Linux)
- A filled-in `.env` file (copy `.env.example` to `.env` and fill in values)

### Start the full stack

```bash
docker compose up --build
```

This starts three containers: the Go app, Postgres, and Redis.
The `--build` flag rebuilds the app image on every run. Omit it after the first run for faster startup.

### Start in the background

```bash
docker compose up -d --build
```

View logs:
```bash
docker compose logs -f app
```

### Stop everything

```bash
docker compose down
```

This stops containers but keeps the Postgres data volume. To also delete all data:
```bash
docker compose down -v
```

### Rebuild after code changes

```bash
docker compose up --build
```

### Running tests

**Unit tests** have no external dependencies and can always be run directly:

```bash
go test ./test/unit/... -v
```

**Integration tests** run outside Docker — on the host machine, against native PostgreSQL
and Redis. They are not run inside a container.

Why outside Docker? The integration tests use Fiber's `app.Test()` helper, which bootstraps
the entire application in-memory (no HTTP server started). They just need a real database
and Redis to connect to. Running them natively is faster and doesn't require a Docker rebuild
on every code change.

Prerequisites for integration tests:
- PostgreSQL running natively with a `slido_clone_test` database created
- Redis running natively
- `.env.test` file present at project root with `DATABASE_NAME=slido_clone_test`,
  `DATABASE_HOST=localhost`, and `REDIS_DB=1`

```bash
# Create the test database (one-time setup)
psql -U your_user -c "CREATE DATABASE slido_clone_test OWNER your_user;"

# Run integration tests
go test ./test/integration/... -v
```

The Postgres and Redis containers started by `docker compose up` are **not accessible from
the host** (no ports are exposed) — they exist solely for the running app container. Keep
your native PostgreSQL and Redis running alongside Docker for local development.

---

## Building the Production Image

Build and tag the image:
```bash
docker build -t slido-clone-backend:latest .
```

Tag for a registry (e.g., GitHub Container Registry):
```bash
docker tag slido-clone-backend:latest ghcr.io/youruser/slido-clone-backend:latest
docker push ghcr.io/youruser/slido-clone-backend:latest
```

---

## Production Deployment (Single VPS)

### First-time VPS setup

**1. Install Docker Engine:**
Follow the official guide: https://docs.docker.com/engine/install/ubuntu/

**2. Configure native Postgres to accept Docker connections:**

The app container connects to Postgres on the host machine. Postgres must be configured
to listen on the Docker bridge interface and accept connections from Docker's subnet.

Edit `postgresql.conf` (usually `/etc/postgresql/<version>/main/postgresql.conf`):
```
listen_addresses = '*'
```

Edit `pg_hba.conf` (usually `/etc/postgresql/<version>/main/pg_hba.conf`), add:
```
host  all  all  172.17.0.0/16  scram-sha-256
```
`172.17.0.0/16` is Docker's default bridge subnet on Linux.

Restart Postgres:
```bash
sudo systemctl restart postgresql
```

**3. Create a `.env` file on the VPS** (never commit this to git):
```bash
# /home/deploy/.env  (or wherever you run docker compose)
DATABASE_USERNAME=your_db_user
DATABASE_PASSWORD=your_db_password
DATABASE_PORT=5432
DATABASE_NAME=slido_clone
JWT_SECRET=your_jwt_secret
REDIS_DB=0
REDIS_PORT=6379
APP_IMAGE=ghcr.io/youruser/slido-clone-backend:latest
```

### Deploy

```bash
# Pull the latest image
docker compose -f docker-compose.prod.yml pull

# Start (or restart) containers
docker compose -f docker-compose.prod.yml up -d
```

### View production logs

```bash
docker compose -f docker-compose.prod.yml logs -f app
```

### Stop production containers

```bash
docker compose -f docker-compose.prod.yml down
```

---

## Troubleshooting

### Migrations fail with "dirty database version"

This means a previous migration failed partway through. Fix:
1. Check which version is dirty: `docker compose logs app | grep "dirty"`
2. SSH into the VPS and run the migrate CLI manually:
   ```bash
   migrate -database "postgres://user:pass@host:5432/dbname?sslmode=disable" \
           -path /path/to/db/migrations \
           force <last-good-version-number>
   ```
3. Fix the migration SQL if needed, then redeploy.

### App cannot connect to Postgres in production

Checklist:
- Is Postgres running? `sudo systemctl status postgresql`
- Does `pg_hba.conf` allow `172.17.0.0/16`? Check with `sudo cat /etc/postgresql/*/main/pg_hba.conf`
- Is `listen_addresses = '*'` set in `postgresql.conf`?
- Did you restart Postgres after the config change? `sudo systemctl restart postgresql`
- Test from inside the container: `docker compose -f docker-compose.prod.yml exec app pg_isready -h host.docker.internal -U $DATABASE_USERNAME`

### Port 3000 is already in use

Another process is using port 3000. Find it:
```bash
sudo lsof -i :3000
```
Kill it or change the port mapping in the compose file (`"3001:3000"` to expose on a different host port).

### Container health check is failing

```bash
docker inspect $(docker compose ps -q app) | grep -A 10 '"Health"'
```

This shows the last health check result and exit code.
