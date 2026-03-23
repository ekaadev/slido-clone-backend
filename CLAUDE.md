# CLAUDE.md

Interactive QnA backend. Core differentiator: **XP Ranking gamification system** — rewards quality participation, making the platform usable as a formative assessment tool.

## Commands

```bash
go run cmd/web/main.go                                            # Run server
go build -o bin/server cmd/web/main.go                            # Build
go mod tidy

# Docker (local dev — starts app + Postgres + Redis)
docker compose up --build                                         # Start full stack
docker compose up -d --build                                      # Start in background
docker compose down                                               # Stop (keep data)
docker compose down -v                                            # Stop + delete data volume
docker compose logs -f app                                        # Stream app logs

# Unit tests (use go-sqlmock, no DB required)
go test ./test/unit/... -v
go test ./test/unit/... -run TestFunctionName -v
go test ./test/unit/user_usecase_test.go ./test/mocks/*.go -v     # Single file

# Integration tests — run on HOST (outside Docker), against native PostgreSQL + Redis
# Requires: slido_clone_test DB, .env.test with DATABASE_HOST=localhost, REDIS_DB=1
go test ./test/integration/... -v
go test ./test/integration/... -run TestRegister_Success -v

# Run migrations (adjust DSN to match .env values)
migrate -database "postgres://user:pass@localhost:5432/dbname?sslmode=disable" -path db/migrations up
migrate -database "postgres://user:pass@localhost:5432/dbname?sslmode=disable" -path db/migrations down
```

## Environment

Copy `.env.example` to `.env`. Required: `DATABASE_USERNAME`, `DATABASE_PASSWORD`, `DATABASE_HOST`, `DATABASE_PORT`, `DATABASE_NAME`, `JWT_SECRET`, `REDIS_HOST`, `REDIS_PORT`, `REDIS_DB`

## Tech Stack

Go Fiber (HTTP + WebSocket) · GORM + PostgreSQL · Redis (JWT blacklist) · JWT auth · Pion WebRTC SFU (`internal/sfu/`)

## Rules

Load only the files relevant to your task:

| Rule file | When to load |
|-----------|--------------|
| [coding-standards.md](.claude/rules/coding-standards.md) | Always |
| [architecture.md](.claude/rules/architecture.md) | Adding/modifying any domain, controller, use case, or repository |
| [auth.md](.claude/rules/auth.md) | Auth, JWT, middleware, login/logout |
| [database.md](.claude/rules/database.md) | Schema changes, migrations, GORM queries |
| [websocket.md](.claude/rules/websocket.md) | WebSocket events, hub, broadcasting, WebRTC/conference |
| [xp-gamification.md](.claude/rules/xp-gamification.md) | XP awards, leaderboard, timeline |
| [domain-rules.md](.claude/rules/domain-rules.md) | Rooms, participants, questions, polls |
| [testing.md](.claude/rules/testing.md) | Writing or running tests |
