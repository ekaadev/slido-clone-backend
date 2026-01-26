package websocket

import "encoding/json"

type WSMessage struct {
	Event string          `json:"event"`
	Data  json.RawMessage `json:"data"`
}

// Event types constants
const (
	// Room events
	EventRoomJoin     = "room:join"
	EventRoomUserJoin = "room:user_joined"
	EventRoomUserLeft = "room:user_left"
	EventRoomClosed   = "room:closed"
	EventRoomAnnounce = "room:announce" // Server -> Client (broadcast from presenter)

	// Message events
	EventMessageSend = "message:send" // Client -> Server
	EventMessageNew  = "message:new"  // Server -> Client (broadcast)
	EventChatTyping  = "chat:typing"  // Bidirectional

	// Question events
	EventQuestionSubmit       = "question:submit"        // Client -> Server
	EventQuestionUpvote       = "question:upvote"        // Client -> Server
	EventQuestionRemoveUpvote = "question:remove_upvote" // Client -> Server
	EventQuestionCreated      = "question:created"       // Server -> Client (broadcast)
	EventQuestionUpvoted      = "question:upvoted"       // Server -> Client (broadcast)
	EventQuestionValidated    = "question:validated"     // Server -> Client (broadcast)

	// Poll events
	EventPollVote          = "poll:vote"            // Client -> Server
	EventPollCreated       = "poll:created"         // Server -> Client (broadcast)
	EventPollResultsUpdate = "poll:results_updated" // Server -> Client (broadcast)
	EventPollClosed        = "poll:closed"          // Server -> Client (broadcast)

	// Leaderboard events
	EventLeaderboardUpdate  = "leaderboard:updated" // Server -> Client
	EventXPAwarded          = "xp:awarded"
	EventLeaderboardRequest = "leaderboard:request" // Client -> Server
)
