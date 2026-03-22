# Architecture Overview

## Overview

Slido Clone Backend is a Go-based backend for an interactive QnA platform. The primary USP is a gamification system ("XP Ranking") designed to measure and reward quality participation, turning the platform into a formative assessment tool for presenters and educators.

## Tech Stack

- **Framework:** Go Fiber (HTTP + WebSocket)
- **ORM:** GORM
- **Database:** MySQL 8.0
- **Cache/Token Store:** Redis (JWT blacklist)
- **Auth:** JWT (golang-jwt)
- **Real-time:** Fiber WebSocket (gorilla/websocket under the hood)
- **Conference:** Pion WebRTC SFU (`internal/sfu/`)

## Clean Architecture Layers

```
HTTP Request
     |
[Delivery Layer]       internal/delivery/http/         (controllers, middleware)
     |                 internal/delivery/websocket/     (hub, client, event handler)
     v
[Use Case Layer]       internal/usecase/               (business logic)
     |
[Repository Layer]     internal/repository/            (data access via GORM)
     |
[Entity Layer]         internal/entity/                (domain models)
     |
[Database]             MySQL via GORM
```

## Dependency Injection

All dependencies are wired in `internal/config/app.go` — `Bootstrap()` function. Order of wiring:
1. Repositories (data access, injected with `*gorm.DB`)
2. Utilities (`TokenUtil` — JWT + Redis)
3. Use Cases (business logic, injected with repos + utils)
4. WebSocket Hub (room-scoped connection manager)
5. HTTP Controllers (injected with use cases + hub)
6. Middleware (auth, injected with TokenUtil)
7. WebSocket Handler (injected with use cases + hub)
8. Routes registered via `route.SetupRoutes()`

## HTTP Response Format

All HTTP responses use `model.WebResponse`:

```go
// Success
c.JSON(model.WebResponse{Data: response})

// Error
c.Status(fiber.StatusBadRequest).JSON(model.WebResponse{Error: "message"})

// Paginated
c.JSON(model.WebResponse{Data: items, Paging: &model.PaginationResponse{
    Page: page, Size: size, TotalItem: total, TotalPage: totalPages,
}})
```

## Auth Flow

1. Register/Login → returns JWT with `UserID`, `Username`, `Email`, `Role`
2. Join Room → returns new JWT with added `RoomID`, `ParticipantID`, `IsRoomOwner` claims
3. WebSocket connection → uses room-scoped JWT via `?token=` query param
4. Logout → blacklists token in Redis; all subsequent requests with that token return 401

JWT claims are defined in `model.Auth` (`internal/model/auth.go`).

## Database Triggers (Denormalization)

MySQL triggers handle vote count denormalization to avoid N+1 update queries:

| Trigger | Action |
|---------|--------|
| `after_vote_insert` | Increments `questions.upvote_count` |
| `after_vote_delete` | Decrements `questions.upvote_count` |
| `after_poll_response_insert` | Increments `poll_options.vote_count` |
| `after_poll_response_delete` | Decrements `poll_options.vote_count` |

## Adding a New Domain

Follow this chain:
1. `internal/entity/` — create entity struct
2. `internal/repository/` — create repo embedding `Repository[YourEntity]`
3. `internal/usecase/` — create use case with DB, Log, Validate, repo deps
4. `internal/delivery/http/` — create controller
5. `internal/model/` — add request/response DTOs
6. `internal/model/converter/` — add conversion functions
7. `internal/config/app.go` — wire dependencies in `Bootstrap()`
8. `internal/delivery/http/route/route.go` — register routes

## Key File Paths

| Concern | File |
|---------|------|
| DI root | `internal/config/app.go` |
| Route registration | `internal/delivery/http/route/route.go` |
| Auth middleware | `internal/delivery/http/middleware/auth_middleware.go` |
| WebSocket hub | `internal/delivery/websocket/hub.go` |
| WebSocket events | `internal/delivery/websocket/message.go` |
| Base repository | `internal/repository/repository.go` |
| Common models | `internal/model/common.go` |
| Auth model/claims | `internal/model/auth.go` |
| DB migrations | `db/migrations/` |
