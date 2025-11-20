# Api Documentation

## Base URL

```
Development: http://localhost:3000/api/v1
```

## User

### Register

```http
POST /users/register
```

Request Body:

```json
{
  "username": "string",
  "email": "string",
  "password": "string",
  "role": "string"
}
```

Response 201:

```json
{
  "data": {
    "user": {
      "id": 1,
      "username": "string",
      "email": "string",
      "role": "string",
      "created_at": "2024-11-11T10:00:00Z"
    },
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
  }
}
```

### Login

```http
POST /users/login
```

Request Body:

```json
{
  "username": "string",
  "password": "string"
}
```

Response 200:

```json
{
  "data": {
    "user": {
      "id": 1,
      "username": "string",
      "email": "string",
      "role": "string",
      "created_at": "2024-11-11T10:00:00Z"
    },
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
  }
}
```

### Anon

```http
POST /users/anonymous
```

Request Body:

```json
{
  "room_code": "ABC123"
}
```

Response 200:

```json
{
  "data": {
    "participant": {
      "id": 67,
      "display_name": "string",
      "is_anonymous": true
    },
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
  }
}
```

## Room

### Create a Room

```http
POST /rooms
Authorization: Bearer {token}
```

Request Body

```json
{
  "title": "string",
  "presenter_id": 1
}
```

Response: 201

```json
{
  "data": {
    "room": {
      "id": 10,
      "room_code": "XY7K9M",
      "title": "string",
      "presenter_id": 1,
      "status": "active",
      "created_at": "2024-11-11T10:00:00Z",
      "closed_at": null
    }
  }
}
```

### Get Room by ID (Detail)

```http
GET /rooms/:room_code
```

Response: 200

```json
{
  "data": {
    "room": {
      "id": 10,
      "room_code": "XY7K9M",
      "title": "string",
      "status": "string",
      "presenter": {
        "id": 1,
        "username": "string"
      },
      "stats": {
        "total_participants": 45,
        "total_questions": 12,
        "total_polls": 3,
        "active_poll_id": 5
      },
      "created_at": "2024-11-11T10:00:00Z"
    }
  }
}
```

### Close Room

```http
PATCH /rooms/:room_id/close
Authorization: Bearer {token}
```

Request Body

```json
{
  "status": "closed"
}
```

Response: 200

```json
{
  "data": {
    "room": {
      "id": 10,
      "status": "closed",
      "closed_at": "2024-11-11T12:00:00Z"
    }
  }
}
```

### Get Room Presenter / List Room Presenter

```http
GET /rooms/my-rooms
Authorization: Bearer {token}
```

```json
{
  "data": {
    "rooms": [
      {
        "id": 10,
        "room_code": "XY7K9M",
        "title": "string",
        "status": "active",
        "participants_count": 45,
        "created_at": "2024-11-11T10:00:00Z"
      },
      {
        "id": 9,
        "room_code": "AB3C5D",
        "title": "string",
        "status": "closed",
        "participants_count": 38,
        "created_at": "2024-11-08T10:00:00Z"
      }
    ]
  }
}
```

## Participant

### Join Room

```http
POST /rooms/:room_code/join
Authorization: Bearer {token}
```

Request Body:

```json
{
  "display_name": "string"
  // Optional
}
```

Response: 200

```json
{
  "data": {
    "participant": {
      "id": 42,
      "room_id": 10,
      "display_name": "string",
      "xp_score": 0,
      "is_anonymous": false,
      "joined_at": "2024-11-11T10:05:00Z"
    },
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
  }
}
```

### Get Participants List

```http
GET /rooms/:room_id/participants
Authorization: Bearer {token}
```

Query Parameters:

- `page_size` (default: 10)
- `page` (default: 1)

Response: 200

```json
{
  "data": {
    "partipants": [
      {
        "id": 42,
        "display_name": "string",
        "xp_score": 75,
        "is_anonymous": false
      },
      {
        "id": 43,
        "display_name": "string",
        "xp_score": 50,
        "is_anonymous": true
      }
    ]
  },
  "paging": {
    "page": 1,
    "size": 10,
    "total_item": 10,
    "total_page": 1
  }
}
```

## Question (Q&A)

### Submit Question / Create

```http
POST /rooms/:room_id/questions
Authorization: Bearer {token}
```

Request Body:

```json
{
  "content": "string"
}
```

Response: 201

```json
{
  "data": {
    "question": {
      "id": 123,
      "room_id": 10,
      "participant_id": 42,
      "content": "string",
      "upvote_count": 0,
      "status": "pending",
      "is_validated_by_presenter": false,
      "xp_awarded": 10,
      "created_at": "2024-11-11T10:15:00Z"
    },
    "xp_earned": {
      "points": 10,
      // points yang didapat
      "new_total": 10
      // hasil kalkulasi (0 + 10)
    }
  }
}
```

### Get Question List

```http
GET /rooms/:room_id/questions
```

Query Parameters:

- `status` (pending/answered/highlighted)
- `sort_by` (recent/upvotes/validated) - default: upvotes
- `limit` (default: 20)
- `offset` (default: 0)

Response: 200

```json
{
  "data": {
    "questions": [
      {
        "id": 123,
        "participant": {
          "id": 42,
          "display_name": "string"
        },
        "content": "string",
        "upvote_count": 15,
        "status": "highlighted",
        "is_validated_by_presenter": true,
        "has_voted": false,
        "created_at": "2024-11-11T10:15:00Z"
      }
    ],
    "paging": {
      "total": 12,
      "limit": 20,
      "offset": 0
    }
  }
}
```

### Add Upvote

```http
POST /questions/:question_id/upvote
Authorization: Bearer {token}
```

Response: 200

```json
{
  "data": {
    "vote": {
      "id": 456,
      "question_id": 123,
      "participant_id": 43,
      "created_at": "2024-11-11T10:20:00Z"
    },
    "question": {
      "id": 123,
      "upvote_count": 16
    },
    "xp_earned": {
      "recipient_participant_id": 42,
      // id penerima
      "points": 3,
      // poin yang didapat
      "source": "upvote_received"
      // source xp (asal xp)
    }
  }
}
```

### Remove Upvote

```http
DELETE /questions/:question_id/upvote
Authorization: Bearer {token}
```

Response: 200

```json
{
  "data": {
    "question": {
      "id": 123,
      "upvote_count": 15
    }
  }
}
```

### Validate Question (Presenter Only)

```http
PATCH /questions/:question_id/validate
Authorization: Bearer {token}
```

Request Body:

```json
{
  "status": "highlighted"
  // atau "answered", default "pending"
}
```

Response: 200

```json
{
  "data": {
    "question": {
      "id": 123,
      "status": "highlighted",
      "is_validated_by_presenter": true
    },
    "xp_awarded": {
      "participant_id": 42,
      // user/participant yang mendapatkan xp
      "points": 25,
      // point yang didapat
      "new_total": 48
      // total kalkulasi baru
    }
  }
}
```

## Poll

### Create Poll

```http
POST /rooms/:room_id/polls
Authorization: Bearer {token}
```

Request Body:

```json
{
  "question": "string",
  "options": [
    "string",
    "string",
    "string",
    "string"
  ]
}
```

Response: 201

```json
{
  "data": {
    "poll": {
      "id": 5,
      "room_id": 10,
      "question": "string",
      "status": "active",
      "created_at": "2024-11-11T10:30:00Z",
      "options": [
        {
          "id": 1,
          "poll_id": 5,
          "option_text": "string",
          "vote_count": 0,
          "order": 1
        },
        {
          "id": 2,
          "poll_id": 5,
          "option_text": "string",
          "vote_count": 0,
          "order": 2
        },
        {
          "id": 3,
          "poll_id": 5,
          "option_text": "string",
          "vote_count": 0,
          "order": 3
        },
        {
          "id": 4,
          "poll_id": 5,
          "option_text": "string",
          "vote_count": 0,
          "order": 4
        }
      ]
    }
  }
}
```

### Get Active Poll

```http
GET /rooms/:room_id/polls/active
Authorization: Bearer {token}
```

Response: 200

```json
{
  "data": {
    "polls": [
      {
        "id": 5,
        "question": "string",
        "status": "active",
        "total_votes": 35,
        "created_at": "2024-11-11T10:30:00Z",
        "options": [
          {
            "id": 1,
            "option_text": "string",
            "vote_count": 15,
            "order": 1
          },
          {
            "id": 2,
            "option_text": "string",
            "vote_count": 12,
            "order": 2
          },
          {
            "id": 3,
            "option_text": "string",
            "vote_count": 5,
            "order": 3
          },
          {
            "id": 4,
            "option_text": "string",
            "vote_count": 3,
            "order": 4
          }
        ]
      }
    ]
  }
}
```

### Submit Vote

```http
POST /polls/:poll_id/vote
Authorization: Bearer {token}
```

Request Body:

```json
{
  "option_id": 2
}
```

Response: 200

```json
{
  "data": {
    "response": {
      "id": 789,
      "poll_id": 5,
      "participant_id": 42,
      "poll_option_id": 2,
      "created_at": "2024-11-11T10:35:00Z"
    },
    "updated_results": {
      "poll_id": 5,
      "total_votes": 36,
      "options": [
        {
          "id": 1,
          "vote_count": 15,
          "percentage": 41.67
        },
        {
          "id": 2,
          "vote_count": 13,
          "percentage": 36.11
        },
        {
          "id": 3,
          "vote_count": 5,
          "percentage": 13.89
        },
        {
          "id": 4,
          "vote_count": 3,
          "percentage": 8.33
        }
      ]
    },
    "xp_earned": {
      "points": 5,
      "new_total": 15
    }
  }
}
```

### Close Poll (Presenter Only)

```http
PATCH /polls/:poll_id/close
Authorization: Bearer {token}
```

Response: 200

```json
{
  "data": {
    "poll": {
      "id": 5,
      "status": "closed",
      "closed_at": "2024-11-11T10:45:00Z",
      "final_results": {
        "total_votes": 36,
        "options": [
          {
            "id": 1,
            "option_text": "string",
            "vote_count": 15,
            "percentage": 41.67
          },
          {
            "id": 2,
            "option_text": "string",
            "vote_count": 13,
            "percentage": 36.11
          }
        ]
      }
    }
  }
}
```

### Get Poll History

```http
GET /rooms/:room_id/polls
Authorization: Bearer {token}
```

Query Parameters:

- `status` (active/closed/all) - default: all
- `limit` (default: 10)

Response: 200

```json
{
  "data": {
    "polls": [
      {
        "id": 5,
        "question": "string?",
        "status": "closed",
        "total_votes": 36,
        "created_at": "2024-11-11T10:30:00Z",
        "closed_at": "2024-11-11T10:45:00Z"
      },
      {
        "id": 4,
        "question": "string?",
        "status": "closed",
        "total_votes": 42,
        "created_at": "2024-11-11T09:50:00Z",
        "closed_at": "2024-11-11T09:55:00Z"
      }
    ],
    "total": 2
  }
}
```

## Message

### Send Message

```http
POST /rooms/:room_id/messages
Authorization: Bearer {token}
```

Request Body:

```json
{
  "content": "string"
}
```

Response: 201

```json
{
  "data": {
    "message": {
      "id": 9999,
      "room_id": 10,
      "participant": {
        "id": 42,
        "display_name": "string"
      },
      "content": "string",
      "created_at": "2024-11-11T10:50:00Z"
    }
  }
}
```

### Get Messages

```http
GET /rooms/:room_id/messages
Authorization: Bearer {token}
```

Response: 200

```json
{
  "data": {
    "messages": [
      {
        "id": 9999,
        "participant": {
          "id": 42,
          "display_name": "string"
        },
        "content": "string",
        "created_at": "2024-11-11T10:50:00Z"
      },
      {
        "id": 9998,
        "participant": {
          "id": 43,
          "display_name": "string"
        },
        "content": "string",
        "created_at": "2024-11-11T10:49:30Z"
      }
    ]
  }
}
```

## XP

### Get Leaderboard

```http
GET /rooms/:room_id/leaderboard
Authorization: Bearer {token}
```

Query Parameters:

- `limit` (default: 10)

```json
{
  "data": {
    "leaderboard": [
      {
        "rank": 1,
        "participant": {
          "id": 42,
          "display_name": "string",
          "is_anonymous": false
        },
        "xp_score": 95
      },
      {
        "rank": 2,
        "participant": {
          "id": 45,
          "display_name": "string",
          "is_anonymous": true
        },
        "xp_score": 78
      },
      {
        "rank": 3,
        "participant": {
          "id": 50,
          "display_name": "string",
          "is_anonymous": false
        },
        "xp_score": 65
      }
    ],
    "my_rank": {
      "rank": 5,
      "xp_score": 48
    },
    "total_participants": 45
  }
}
```

### Get My XP Details

```http
GET /rooms/:room_id/my_xp
Authorization: Bearer {token}
```

Response: 200

```json
{
  "data": {
    "participant": {
      "id": 42,
      "display_name": "string",
      "xp_score": 95,
      "rank": 1
    },
    "xp_breakdown": {
      "poll": {
        "count": 3,
        "total_xp": 15
      },
      "question_created": {
        "count": 4,
        "total_xp": 40
      },
      "upvote_received": {
        "count": 8,
        "total_xp": 24
      },
      "presenter_validated": {
        "count": 2,
        "total_xp": 50
      }
    },
    "recent_transactions": [
      {
        "id": 789,
        "points": 25,
        "source_type": "presenter_validated",
        "source_detail": "string",
        "created_at": "2024-11-11T10:40:00Z"
      },
      {
        "id": 788,
        "points": 3,
        "source_type": "upvote_received",
        "source_detail": "string",
        "created_at": "2024-11-11T10:35:00Z"
      }
    ]
  }
}
```

### Get Xp Transactions History (Log)

```http
GET /rooms/:room_id/xp_transactions
Authorization: Bearer {token}
```

Query Parameters:

- `participant_id` (optional, for presenter to view specific user)
- `limit` (default: 20)
- `offset` (default: 0)

```json
{
  "data": {
    "transactions": [
      {
        "id": 789,
        "participant": {
          "id": 42,
          "display_name": "string"
        },
        "points": 25,
        "source_type": "presenter_validated",
        "source_id": 123,
        "created_at": "2024-11-11T10:40:00Z"
      }
    ],
    "total": 15,
    "summary": {
      "total_xp": 95,
      "total_transactions": 15
    }
  }
}
```