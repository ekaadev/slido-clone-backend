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

	// WebRTC events
	EventWebrtcOffer     = "webrtc:offer"     // Client -> Server
	EventWebrtcAnswer    = "webrtc:answer"    // Server -> Client
	EventWebrtcCandidate = "webrtc:candidate" // Bidirectional

	// Conference events (Discord-style stage)
	EventConferenceStart   = "conference:start"   // Client -> Server (host only)
	EventConferenceStop    = "conference:stop"    // Client -> Server (host only)
	EventConferenceStarted = "conference:started" // Server -> Client (broadcast)
	EventConferenceEnded   = "conference:ended"   // Server -> Client (broadcast)
	EventConferenceState   = "conference:state"   // Server -> Client (current state)

	// Raise hand events
	EventRaiseHand   = "conference:raise_hand"   // Client -> Server
	EventLowerHand   = "conference:lower_hand"   // Client -> Server
	EventHandRaised  = "conference:hand_raised"  // Server -> Client (broadcast)
	EventHandLowered = "conference:hand_lowered" // Server -> Client (broadcast)

	// Promote/Demote events
	EventPromoteSpeaker  = "conference:promote"  // Client -> Server (host only)
	EventDemoteSpeaker   = "conference:demote"   // Client -> Server (host only)
	EventSpeakerPromoted = "conference:promoted" // Server -> Client (broadcast)
	EventSpeakerDemoted  = "conference:demoted"  // Server -> Client (broadcast)

	// Join conference as audience
	EventConferenceJoin   = "conference:join"   // Client -> Server
	EventConferenceLeave  = "conference:leave"  // Client -> Server
	EventConferenceJoined = "conference:joined" // Server -> Client (broadcast)
	EventConferenceLeft   = "conference:left"   // Server -> Client (broadcast)
)
