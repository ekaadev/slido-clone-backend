# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Commands

```bash
# Run the server
go run cmd/web/main.go

# Run all tests
go test ./tests/...

# Run a single test file
go test ./tests/user_usecase_test.go ./tests/mocks/*.go -v

# Run a single test function
go test ./tests/... -run TestFunctionName -v

# Build
go build -o bin/server cmd/web/main.go

# Tidy dependencies
go mod tidy
```

## Environment Setup

Copy `.env.example` to `.env` and fill in values. Required variables:
- `DATABASE_USERNAME`, `DATABASE_PASSWORD`, `DATABASE_HOST`, `DATABASE_PORT`, `DATABASE_NAME`
- `JWT_SECRET`
- `REDIS_HOST`, `REDIS_PORT`, `REDIS_DB`

Database migrations are in `db/migrations/` using [golang-migrate](https://github.com/golang-migrate/migrate) with numbered up/down SQL files.

## Architecture

This is a **Clean Architecture** Go backend (Fiber framework) with three layers:

1. **Delivery** (`internal/delivery/`) — HTTP controllers + WebSocket handler
2. **Use Case** (`internal/usecase/`) — Business logic, one file per domain
3. **Repository** (`internal/repository/`) — Data access via GORM

All dependencies are wired together in `internal/config/app.go` (`Bootstrap` function) — this is the DI root. When adding a new domain, follow the chain: entity → repository → use case → controller → route registration in `Bootstrap` and `route.go`.

### Adding a New Domain

1. Create entity in `internal/entity/`
2. Create repository in `internal/repository/` embedding `Repository[YourEntity]`
3. Create use case in `internal/usecase/` accepting DB, Log, Validate, and repo deps
4. Create controller in `internal/delivery/http/`
5. Add request/response DTOs in `internal/model/` and converters in `internal/model/converter/`
6. Wire everything in `internal/config/app.go` and register routes in `internal/delivery/http/route/route.go`

### HTTP Response Format

All HTTP responses use `model.WebResponse`:
```go
// Success
c.JSON(model.WebResponse{Data: response})
// Error
c.Status(fiber.StatusBadRequest).JSON(model.WebResponse{Error: "message"})
// Paginated
c.JSON(model.WebResponse{Data: items, Paging: &model.PaginationResponse{...}})
```

### WebSocket

- `GET /ws?token={jwt_token}` — token must contain `RoomID` and `ParticipantID` claims (set on `/join`)
- The `Hub` manages room-scoped connections (`internal/delivery/websocket/hub.go`)
- Incoming events are routed by `EventHandler.HandleMessage` based on `event` field in the JSON message
- All event type constants are defined in `internal/delivery/websocket/message.go`
- Controllers that need to broadcast (e.g., `PollController`, `QuestionController`) receive the `Hub` as a dependency and call `hub.BroadcastToRoom`
- After most actions, `broadcastLeaderboardUpdate` is called to push updated XP leaderboard to the room

**Key WebSocket events (Client → Server):** `message:send`, `chat:typing`, `question:submit`, `question:upvote`, `question:remove_upvote`, `poll:vote`, `leaderboard:request`, `webrtc:offer/answer/candidate`, `conference:start/stop/join/leave/raise_hand/lower_hand/promote/demote`

### JWT Auth

- All routes after `c.App.Use(c.AuthMiddleware)` in `route.go` require Bearer token
- Tokens store `UserID`, `Username`, and optionally `RoomID` + `ParticipantID` (populated after joining a room)
- Logout blacklists tokens in Redis

### Database

- MySQL 8.0 with GORM; connection pool configured in `config.json`
- DB triggers handle vote count denormalization automatically (see migrations for `after_vote_insert`, `after_vote_delete`, `after_poll_response_insert/delete`)
- The generic `Repository[T]` base provides `Create`, `Update`, `Delete`, `FindById`, `CountById`; domain repos add their own query methods

### XP / Gamification

Questions and messages award XP to participants, tracked in `xp_transaction` table. Constants in `question_usecase.go`:
- Submit question: 10 XP
- Receive upvote: 3 XP
- Presenter validates question: 25 XP

The leaderboard (`GET /api/v1/rooms/:room_id/leaderboard`) and timeline (`GET /api/v1/rooms/:room_id/timeline`) aggregate this activity.

### Conference (WebRTC SFU)

The `internal/sfu/` package is a Pion-based Selective Forwarding Unit for video conferencing. It is wired into the WebSocket `EventHandler` but is independent of the HTTP/domain flow. Conference controls (start/stop/promote speaker) are restricted to the room owner (`client.isRoomOwner`).

### Testing

Tests are in `tests/` (not colocated with source). They use:
- `go-sqlmock` to mock the GORM DB connection
- `testify/mock` with hand-written mocks in `tests/mocks/`
- New mocks must implement the interface methods the use case depends on
- Use case structs are instantiated directly in tests (fields are exported), not via constructors, to inject mocks
