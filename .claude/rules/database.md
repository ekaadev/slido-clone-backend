---
description: PostgreSQL/GORM patterns, migrations, DB triggers for denormalization, and transaction usage.
---

# Database Rules

## Stack

- PostgreSQL with GORM; connection pool configured in `config.json`.
- The generic `Repository[T]` base provides `Create`, `Update`, `Delete`, `FindById`, `CountById`; domain repos add their own query methods.

## Migrations

- All schema changes require a new migration file in `db/migrations/` using the numbered up/down naming convention.
- Do not modify existing migration files — add a new one.
- Migration tooling: [golang-migrate](https://github.com/golang-migrate/migrate).

## DB Triggers (Denormalization)

DB triggers manage `upvote_count` (questions) and `vote_count` (poll options) automatically:

| Trigger                      | Effect                                     |
| ---------------------------- | ------------------------------------------ |
| `after_vote_insert`          | Increments `upvote_count` on the question  |
| `after_vote_delete`          | Decrements `upvote_count` on the question  |
| `after_poll_response_insert` | Increments `vote_count` on the poll option |
| `after_poll_response_delete` | Decrements `vote_count` on the poll option |

Do not replicate this logic in application code. Never manually update these denormalized counts.

## Transactions

Use GORM transactions for any multi-step write operation:

```go
db.Transaction(func(tx *gorm.DB) error {
    // all steps here
    return nil
})
```
