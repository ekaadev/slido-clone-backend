package model

import (
	"time"
)

type JoinRoomRequest struct {
	Username    string `json:"username" validate:"required,min=3,max=30,alphanum"`
	DisplayName string `json:"display_name,omitempty" validate:"omitempty,min=2,max=100"`
	RoomCode    string `json:"room_code" validate:"required,len=6,alphanum"`
}

type ParticipantResponse struct {
	ID          uint      `json:"id"`
	RoomID      uint      `json:"room_id,omitempty"`
	DisplayName string    `json:"display_name"`
	XPScore     int       `json:"xp_score"`
	IsAnonymous bool      `json:"is_anonymous"`
	RoomRole    string    `json:"room_role,omitempty"` // "host" atau "audience"
	JoinedAt    time.Time `json:"joined_at,omitempty"`
}

type ListParticipantsRequest struct {
	RoomID        uint `json:"-" validate:"required,min=1"`
	ParticipantID uint `json:"-" validate:"required,min=1"`
	Page          int  `json:"page" validate:"required,min=1"`
	Size          int  `json:"size" validate:"required,min=1,max=100"`
}

type ParticipantListItem struct {
	ID          uint   `json:"id"`
	DisplayName string `json:"display_name"`
	XPScore     int    `json:"xp_score"`
	IsAnonymous bool   `json:"is_anonymous"`
}

type ParticipantListResponse struct {
	Participants []*ParticipantListItem `json:"participants"`
}

type JoinRoomResponse struct {
	Participant ParticipantResponse `json:"participant"`
	Token       string              `json:"token"`
}

type ParticipantInfo struct {
	ID          uint   `json:"id"`
	DisplayName string `json:"display_name"`
}
