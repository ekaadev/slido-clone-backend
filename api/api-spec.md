# WebSocket API Specification

## Connection

**Endpoint:** `ws://localhost:3000/ws?token={jwt_token}`

Token harus didapatkan dari login atau join room terlebih dahulu.

---

## Client -> Server Events

Format message:
```json
{
  "event": "event_name",
  "data": { ... }
}
```

| Event | Payload | Description |
|-------|---------|-------------|
| `message:send` | `{content: string}` | Kirim chat message ke room |
| `chat:typing` | `{is_typing: boolean}` | Typing indicator |
| `question:submit` | `{content: string}` | Submit pertanyaan Q&A |
| `question:upvote` | `{question_id: number}` | Upvote pertanyaan |
| `question:remove_upvote` | `{question_id: number}` | Hapus upvote |
| `poll:vote` | `{poll_id: number, option_id: number}` | Submit vote pada poll |
| `leaderboard:request` | `{}` | Request data leaderboard |

---

## Server -> Client Events

### Room Events

#### `room:user_joined`
Broadcast ketika participant baru join room.
```json
{
  "event": "room:user_joined",
  "data": {
    "participant_id": 123,
    "display_name": "John Doe",
    "is_anonymous": false,
    "joined_at": "2026-01-26T08:00:00+07:00"
  }
}
```

#### `room:user_left`
Broadcast ketika participant disconnect dari room.
```json
{
  "event": "room:user_left",
  "data": {
    "participant_id": 123,
    "display_name": "John Doe",
    "left_at": "2026-01-26T08:30:00+07:00"
  }
}
```

#### `room:announce`
Broadcast announcement dari presenter (trigger via HTTP POST).
```json
{
  "event": "room:announce",
  "data": {
    "message": "Announcement text from presenter",
    "announced_at": "2026-01-26T08:15:00+07:00"
  }
}
```

#### `room:closed`
Broadcast ketika room di-close oleh presenter.
```json
{
  "event": "room:closed",
  "data": {
    "room_id": 1
  }
}
```

---

### Message Events

#### `message:send`
Broadcast chat message baru.
```json
{
  "event": "message:send",
  "data": {
    "id": 456,
    "content": "Hello everyone!",
    "participant": {
      "id": 123,
      "display_name": "John"
    },
    "created_at": "2026-01-26T08:00:00+07:00"
  }
}
```

#### `chat:typing`
Broadcast typing indicator.
```json
{
  "event": "chat:typing",
  "data": {
    "participant_id": 123,
    "is_typing": true
  }
}
```

---

### Question (Q&A) Events

#### `question:created`
Broadcast pertanyaan baru.
```json
{
  "event": "question:created",
  "data": {
    "id": 789,
    "content": "What is the deadline?",
    "participant": { "id": 123, "display_name": "John" },
    "upvote_count": 0,
    "is_validated": false,
    "created_at": "2026-01-26T08:00:00+07:00"
  }
}
```

#### `question:upvoted`
Broadcast update upvote count.
```json
{
  "event": "question:upvoted",
  "data": {
    "question_id": 789,
    "upvote_count": 5
  }
}
```

#### `question:validated`
Broadcast ketika presenter validate pertanyaan.
```json
{
  "event": "question:validated",
  "data": {
    "question_id": 789,
    "is_validated": true
  }
}
```

---

### Poll Events

#### `poll:created`
Broadcast poll baru dari presenter.
```json
{
  "event": "poll:created",
  "data": {
    "id": 101,
    "question": "What topic should we cover next?",
    "options": [
      { "id": 1, "text": "Topic A", "vote_count": 0 },
      { "id": 2, "text": "Topic B", "vote_count": 0 }
    ],
    "status": "active"
  }
}
```

#### `poll:results_updated`
Broadcast update hasil poll setelah ada vote baru.
```json
{
  "event": "poll:results_updated",
  "data": {
    "poll_id": 101,
    "options": [
      { "id": 1, "text": "Topic A", "vote_count": 5 },
      { "id": 2, "text": "Topic B", "vote_count": 3 }
    ]
  }
}
```

#### `poll:closed`
Broadcast ketika presenter close poll.
```json
{
  "event": "poll:closed",
  "data": {
    "poll_id": 101
  }
}
```

---

### Leaderboard Events

#### `leaderboard:updated`
Broadcast update leaderboard (setelah ada perubahan XP).
```json
{
  "event": "leaderboard:updated",
  "data": {
    "leaderboard": [
      { "rank": 1, "participant": { "id": 123, "display_name": "John" }, "xp_score": 150 },
      { "rank": 2, "participant": { "id": 456, "display_name": "Jane" }, "xp_score": 120 }
    ],
    "total_participants": 25
  }
}
```

#### `xp:awarded`
Notifikasi ketika participant mendapatkan XP.
```json
{
  "event": "xp:awarded",
  "data": {
    "participant_id": 123,
    "points": 10,
    "reason": "question_validated"
  }
}
```

---

## HTTP Endpoint untuk WebSocket Features

| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/api/v1/rooms/:room_id/announcement` | Kirim announcement ke room (presenter only) |

---

## Timeline API (Unified Activity Feed)

### GET /api/v1/rooms/:room_id/timeline

Mendapatkan gabungan messages, questions, dan polls dalam satu timeline seperti Discord.

**Query Parameters:**
| Param | Type | Description |
|-------|------|-------------|
| `before` | RFC3339 | Cursor untuk load older items |
| `after` | RFC3339 | Cursor untuk load newer items |
| `limit` | number | Jumlah items (default: 50, max: 100) |

**Response:**
```json
{
  "data": {
    "items": [
      {
        "type": "message",
        "id": 1,
        "created_at": "2026-01-26T09:03:00Z",
        "data": {"content": "Hello!", "participant": {"id": 1, "display_name": "John"}}
      },
      {
        "type": "poll",
        "id": 2,
        "created_at": "2026-01-26T09:02:00Z",
        "data": {"question": "Choose topic", "options": [...], "status": "active"}
      },
      {
        "type": "question",
        "id": 5,
        "created_at": "2026-01-26T09:01:00Z",
        "data": {"content": "What is...?", "upvote_count": 3, "is_validated": false}
      }
    ],
    "has_more": true,
    "oldest_at": "2026-01-26T09:01:00Z",
    "newest_at": "2026-01-26T09:03:00Z"
  }
}
```

**Item Types:**
- `message` - Chat message
- `question` - Q&A question
- `poll` - Polling
- `announcement` - Room announcement

---

## Timeline WebSocket Event

#### `activity:new`
Broadcast ketika ada item baru (message, question, poll, announcement).
```json
{
  "event": "activity:new",
  "data": {
    "type": "message",
    "id": 123,
    "created_at": "2026-01-26T09:05:00Z",
    "data": {"content": "New message", "participant": {...}}
  }
}
```