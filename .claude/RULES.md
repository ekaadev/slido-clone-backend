# Project Rules — Slido Clone Backend

## Coding Agent Rules

1. Read existing code before making any changes.
2. Make the minimum necessary changes. Do not refactor or clean up code outside the scope of the task.
3. Do not introduce breaking changes unless explicitly requested.
4. Follow the **Golang Clean Architecture** pattern used throughout this project.
5. Never hardcode environment variables or secret values — use `.env` and `config.json`.
6. Add comments to explain non-obvious logic (e.g., after creating a function with complex behavior).
7. Do not use emoji or emoticons in code, comments, or markdown output.
8. Keep chat responses concise; avoid generating markdown summaries unless explicitly requested.
9. Prefer simple, maintainable solutions over clever ones.

## Application Context

An interactive QnA backend. The core USP is the **XP Ranking gamification system**, designed to identify and reward quality participation. This transforms the platform from a simple interaction tool into a formative assessment tool usable by educators, presenters, or event hosts.

## Tech Stack

- **Framework:** Go Fiber (HTTP + WebSocket)
- **ORM:** GORM with MySQL 8.0
- **Cache:** Redis (JWT blacklist storage)
- **Auth:** JWT
- **Real-time:** Fiber WebSocket (hub-based, room-scoped)
- **Conference:** Pion WebRTC SFU (`internal/sfu/`)

## Architecture Rules

### Layer Responsibilities
- **Delivery** (`internal/delivery/`) — HTTP controllers and WebSocket handler only. No business logic here.
- **Use Case** (`internal/usecase/`) — All business logic. One file per domain.
- **Repository** (`internal/repository/`) — Data access only via GORM. No business logic.
- **Entity** (`internal/entity/`) — Domain structs only. No methods with business logic.

### Dependency Direction
Controllers depend on use cases. Use cases depend on repositories. Repositories depend on GORM. Never reverse this.

### DI Root
All dependencies are wired in `internal/config/app.go` (`Bootstrap` function). When adding a new domain, register everything here.

### Adding a New Domain (follow this order)
1. `internal/entity/` — entity struct
2. `internal/repository/` — repo embedding `Repository[YourEntity]`
3. `internal/usecase/` — use case with DB, Log, Validate, and repo deps
4. `internal/delivery/http/` — controller
5. `internal/model/` — request/response DTOs
6. `internal/model/converter/` — conversion functions
7. `internal/config/app.go` — wire into `Bootstrap()`
8. `internal/delivery/http/route/route.go` — register routes

### HTTP Response Format
Always use `model.WebResponse`:
```go
// Success
c.JSON(model.WebResponse{Data: response})
// Error
c.Status(fiber.StatusBadRequest).JSON(model.WebResponse{Error: "message"})
// Paginated
c.JSON(model.WebResponse{Data: items, Paging: &model.PaginationResponse{...}})
```

### Auth
- Routes after `c.App.Use(c.AuthMiddleware)` require Bearer token.
- Get auth claims via `middleware.GetUser(c)` — returns `*model.Auth`.
- Room-scoped tokens (with `RoomID`, `ParticipantID`, `IsRoomOwner`) are issued on `Join` and `Create Room`.
- WebSocket connections authenticate via `?token=` query param; token must be room-scoped.

## Domain Rules

### Auth & Users
- Passwords must be hashed with bcrypt before storing.
- Email and username must be validated for uniqueness before insert.
- Logout must blacklist the token in Redis (not just discard it client-side).
- Anonymous users get a participant record but no `users` row. Their JWT has `IsAnonymous: true`.

### Rooms
- Room codes are generated with `crypto/rand` (cryptographically secure) — do not use `math/rand`.
- Presenter is always auto-enrolled as a participant on room creation (inside the same transaction).
- Only the room presenter can close, delete, or send announcements.
- A room must be `closed` before it can be deleted.
- Room status is one-way: `active` → `closed`. No re-opening.

### Participants
- `Join` is idempotent: if the user already has a participant in this room, return the existing one.
- `IsRoomOwner` is determined by comparing `participant.UserID` with `room.PresenterID`.
- `XPScore` is a denormalized field on `participants` — updated via `XPTransactionRepository.AddXP`, never manually.
- Leaderboard returns top 10; caller rank is computed separately via a count query.

### Q&A (Questions)
- A participant cannot upvote their own question.
- Duplicate votes are prevented at both application level and DB unique constraint.
- Upvote removal must also reverse the XP grant (negative XP transaction).
- A question can only be validated once (`IsValidatedByPresenter` is a one-way flag).
- `upvote_count` is managed by DB triggers — do NOT manually update it in application code.
- Only the room presenter can call the validate endpoint.

### Polling
- Only the room presenter can create or close polls.
- A participant can only vote once per poll (DB unique constraint + application check).
- Validate that the chosen option actually belongs to the poll before recording the vote.
- `vote_count` on poll options is managed by DB triggers — do NOT manually update it.
- Poll status is one-way: `active` → `closed`. No re-opening.

### XP Gamification
- All XP changes must go through `XPTransactionRepository.AddXP` — never update `xp_score` directly.
- Every XP grant must create an `xp_transaction` record with the correct `source_type` and `source_id`.
- XP is room-scoped. A participant's score is per-room.
- After every XP-awarding action, broadcast `leaderboard:updated` via `broadcastLeaderboardUpdate(hub, roomID)`.

| Action | XP | Recipient | Source Type |
|--------|-----|-----------|-------------|
| Submit question | +10 | Author | `question_created` |
| Receive upvote | +3 | Question author | `upvote_received` |
| Upvote removed | -3 | Question author | `upvote_received` |
| Presenter validates | +25 | Question author | `presenter_validated` |
| Vote on poll | +5 | Voter | `poll` |
| Send message | +1 | Sender | `message_created` |

### WebSocket
- All room broadcasts use `hub.BroadcastToRoom(roomID, message)`.
- Controllers that need to broadcast must receive `*hub.Hub` as a dependency.
- Never call hub methods directly from use cases — only from controllers/event handlers.
- Conference control actions (`start`, `stop`, `promote`, `demote`) are restricted to `client.isRoomOwner`.
- Conference state is in-memory only; it is not persisted to the database.

### Timeline
- The timeline is read-only. Never modify data through the timeline endpoint.
- Timeline uses a UNION ALL query across `messages`, `questions`, and `polls` — do not add new types without updating the UNION.
- Cursor pagination uses RFC3339 timestamps (`before`/`after`). Do not mix with page-based pagination.

## Testing Rules

- Tests are in `tests/` (not colocated with source).
- Use `go-sqlmock` to mock the GORM DB connection.
- Use `testify/mock` with hand-written mocks in `tests/mocks/`.
- Instantiate use case structs directly (exported fields) in tests — do not use constructors.
- New mocks must implement the full interface that the use case depends on.

## Database Rules

- All schema changes require a new migration file in `db/migrations/` using the numbered up/down naming convention.
- Do not modify existing migration files — add a new one.
- DB triggers manage `upvote_count` and `vote_count` denormalization. Do not replicate this logic in application code.
- Use GORM transactions (`db.Transaction(func(tx *gorm.DB) error {...})`) for any multi-step write operation.

## File Reference

| Concern | File |
|---------|------|
| DI root | `internal/config/app.go` |
| Route registration | `internal/delivery/http/route/route.go` |
| Auth middleware | `internal/delivery/http/middleware/auth_middleware.go` |
| WebSocket hub | `internal/delivery/websocket/hub.go` |
| WebSocket events | `internal/delivery/websocket/message.go` |
| Base repository | `internal/repository/repository.go` |
| Common response model | `internal/model/common.go` |
| JWT claims model | `internal/model/auth.go` |
| DB migrations | `db/migrations/` |
| Feature docs | `docs/` |
