---
description: Test file location, mock patterns, unit and integration test patterns, and how to run tests.
---

# Testing Rules

## Location

```
test/
  unit/         # unit tests — no DB required
  integration/  # integration tests — require real PostgreSQL + Redis
  mocks/        # hand-written testify mocks shared by unit tests
```

## Unit Tests (`test/unit/`)

**Tooling:** `go-sqlmock` + `testify/mock` with hand-written mocks.

- Instantiate use case structs directly using exported fields — do not use constructors.
- New mocks must implement the full interface that the use case depends on.
- Mock files live in `test/mocks/`.

```bash
go test ./test/unit/... -v
go test ./test/unit/user_usecase_test.go ./test/mocks/*.go -v  # single file
go test ./test/unit/... -run TestFunctionName -v
```

## Integration Tests (`test/integration/`)

**Tooling:** Fiber's `app.Test()` + `testify/assert` against a real PostgreSQL database.

**Prerequisites:**
1. `docker-compose.test.yml` running: `docker compose -f docker-compose.test.yml up -d`
   - postgres-test on port 5434, redis-test on port 6381
2. `.env.test` present at project root with `DATABASE_NAME=reisify_test`, `DATABASE_PORT=5434`, `REDIS_PORT=6381`, `REDIS_DB=1`

**Important — Docker:** Integration tests run on the host against `docker-compose.test.yml` (postgres on port 5434, redis on port 6381). Start it with `docker compose -f docker-compose.test.yml up -d` before running tests.

**How it works:**
- `TestMain` in `setup_test.go` loads `.env.test`, connects to the test DB, runs all migrations via `golang-migrate`, bootstraps the full Fiber app, and runs tests.
- `cleanDB(t)` truncates all tables at the start of each test for isolation.
- `makeRequest()` uses `app.Test()` to call handlers in-process (no network needed).
- After all tests, `teardown()` rolls back all migrations (`migrate down`).

```bash
go test ./test/integration/... -v
go test ./test/integration/... -run TestRegister_Success -v
```

**Key helpers in `setup_test.go`:**
- `makeRequest(t, method, path, body, token)` — builds and sends an HTTP request
- `readBody(t, resp)` — decodes the JSON response body
- `cleanDB(t)` — truncates all tables
- `registerUser`, `loginUser`, `createRoom`, `joinRoom` — common setup shortcuts
- `formatID(id float64) string` — converts JSON float64 ID to URL path string
