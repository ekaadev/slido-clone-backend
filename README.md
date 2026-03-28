# Reisify

Backend service untuk aplikasi Reisify — platform interaktif Q&A, Polling, Live Chat, dan Leaderboard berbasis XP Gamification.

## Getting Started

### Quick Start (Docker — Postgres & Redis)

Run Postgres and Redis in Docker, then run the app on the host.

```bash
cp .env.example .env   # fill in DATABASE_USERNAME, DATABASE_PASSWORD, JWT_SECRET
docker compose up -d
go run cmd/web/main.go
```

Server: `http://localhost:3000`

See [docs/docker.md](docs/docker.md) for the full Docker guide (production deploy, troubleshooting, etc.).

---

### Manual Setup

## Prerequisites

- Go 1.25+
- PostgreSQL 14+
- Redis 7.0+

### 1. Install Dependencies

#### Install Golang Migrate

**Linux:**
```bash
curl -L https://github.com/golang-migrate/migrate/releases/download/v4.17.0/migrate.linux-amd64.tar.gz | tar xvz
sudo mv migrate /usr/local/bin/migrate
```

**macOS:**
```bash
brew install golang-migrate
```

**Windows:**
```bash
scoop install migrate
```

#### Install Redis

**Linux (Ubuntu/Debian):**
```bash
sudo apt update
sudo apt install redis-server
sudo systemctl start redis-server
sudo systemctl enable redis-server
```

**macOS:**
```bash
brew install redis
brew services start redis
```

**Windows:**
Download dari [Redis for Windows](https://github.com/microsoftarchive/redis/releases) atau gunakan WSL.

### 2. Setup Database

Access PostgreSQL shell:
```bash
psql -U postgres
```

Create user and database:
```sql
CREATE USER reisify_user WITH PASSWORD 'password';
CREATE DATABASE reisify OWNER reisify_user;
\q
```

Run migrations:
```bash
migrate -database "postgres://reisify_user:password@localhost:5432/reisify?sslmode=disable" -path db/migrations up
```

> **Note:** Ganti `password` dengan password yang aman untuk production.

### 3. Install Go Dependencies

```bash
go mod tidy
go build ./...
```

### 4. Setup Configuration

Copy `.env.example` to `.env` dan isi nilainya:

```bash
DATABASE_USERNAME=reisify_user
DATABASE_PASSWORD=password
DATABASE_HOST=localhost
DATABASE_PORT=5432
DATABASE_NAME=reisify
```

> **Docker dev ports:** If using Docker Compose for Postgres/Redis, use `DATABASE_PORT=5433` and `REDIS_PORT=6380`.

### 5. Run Application

```bash
go run cmd/web/main.go
# or
make run
```

Server akan berjalan di `http://localhost:3000`.

---

## Testing

### Unit Tests

No external dependencies required. Run anytime:

```bash
go test ./test/unit/... -v
# or
make test-unit
```

### Integration Tests

Integration tests run **on the host machine** against PostgreSQL and Redis provided by `docker-compose.test.yml`. They use Fiber's `app.Test()` in-memory — no running server needed.

**One-time setup:**

```bash
# Start the test containers
docker compose -f docker-compose.test.yml up -d

# Create .env.test
cp .env.example .env.test
```

Edit `.env.test`:
```
DATABASE_USERNAME=testuser
DATABASE_PASSWORD=testpass
DATABASE_HOST=localhost
DATABASE_PORT=5434
DATABASE_NAME=reisify_test
JWT_SECRET=test_secret
REDIS_HOST=localhost
REDIS_PORT=6381
REDIS_DB=1
```

**Run:**

```bash
go test ./test/integration/... -v
go test ./test/integration/... -run TestRegister_Success -v
# or
make test-integration
```

**Stop test containers when done:**
```bash
docker compose -f docker-compose.test.yml down -v
```
