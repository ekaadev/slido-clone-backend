package model

import "time"

// ActivityType enum untuk tipe aktivitas dalam timeline
const (
	ActivityTypeMessage      = "message"
	ActivityTypeQuestion     = "question"
	ActivityTypePoll         = "poll"
	ActivityTypeAnnouncement = "announcement"
)

// TimelineItem represents single item in timeline
type TimelineItem struct {
	Type      string      `json:"type"`       // message, question, poll, announcement
	ID        uint        `json:"id"`         // ID dari item asli
	CreatedAt time.Time   `json:"created_at"` // waktu pembuatan
	Data      interface{} `json:"data"`       // data lengkap sesuai tipe
}

// GetTimelineRequest request dengan cursor-based pagination
type GetTimelineRequest struct {
	RoomID        uint   `json:"-" validate:"required,min=1"`
	ParticipantID uint   `json:"-"`      // untuk context user (optional)
	Before        string `json:"before"` // cursor: RFC3339 timestamp untuk load older items
	After         string `json:"after"`  // cursor: RFC3339 timestamp untuk load newer items
	Limit         int    `json:"limit" validate:"omitempty,min=1,max=100"`
}

// GetTimelineResponse response dengan pagination info
type GetTimelineResponse struct {
	Items    []TimelineItem `json:"items"`
	HasMore  bool           `json:"has_more"`            // apakah ada items lagi yang lebih lama
	OldestAt *string        `json:"oldest_at,omitempty"` // cursor untuk load older
	NewestAt *string        `json:"newest_at,omitempty"` // cursor untuk load newer
}

// TimelineRawItem untuk hasil query UNION (sebelum di-enrich dengan data lengkap)
type TimelineRawItem struct {
	Type      string    `json:"type"`
	ID        uint      `json:"id"`
	CreatedAt time.Time `json:"created_at"`
}

// MessageTimelineData data message untuk timeline
type MessageTimelineData struct {
	Content     string          `json:"content"`
	Participant ParticipantInfo `json:"participant"`
}

// QuestionTimelineData data question untuk timeline
type QuestionTimelineData struct {
	Content     string          `json:"content"`
	Participant ParticipantInfo `json:"participant"`
	UpvoteCount int             `json:"upvote_count"`
	IsValidated bool            `json:"is_validated"`
	Status      string          `json:"status"`
}

// PollTimelineData data poll untuk timeline
type PollTimelineData struct {
	Question   string               `json:"question"`
	Status     string               `json:"status"`
	Options    []PollOptionResponse `json:"options"`
	TotalVotes int                  `json:"total_votes"`
}

// AnnouncementTimelineData data announcement untuk timeline
type AnnouncementTimelineData struct {
	Message string `json:"message"`
}
