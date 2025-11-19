package model

import "time"

type JoinRoomRequest struct {
	DisplayName string `json:"display_name,omitempty" validate:"omitempty,min=2,max=100"`
}

type ParticipantResponse struct {
	ID          uint      `json:"id"`
	RoomID      uint      `json:"room_id,omitempty"`
	DisplayName string    `json:"display_name"`
	XPScore     int       `json:"xp_score"`
	IsAnonymous bool      `json:"is_anonymous"`
	JoinedAt    time.Time `json:"joined_at,omitempty"`
}

type ParticipantListItem struct {
	ID          uint   `json:"id"`
	DisplayName string `json:"display_name"`
	XPScore     int    `json:"xp_score"`
	IsAnonymous bool   `json:"is_anonymous"`
}
