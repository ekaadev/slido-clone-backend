package model

import "time"

type SendMessageRequest struct {
	RoomID        uint   `json:"-" validate:"required,min=1"`
	ParticipantID uint   `json:"-" validate:"required,min=1"`
	Content       string `json:"content" validate:"required,min=1,max=1000"`
}

type MessageResponse struct {
	ID          uint            `json:"id"`
	RoomID      uint            `json:"room_id,omitempty"`
	Participant ParticipantInfo `json:"participant"`
	Content     string          `json:"content"`
	CreatedAt   time.Time       `json:"created_at"`
}

type SendMessageResponse struct {
	Message MessageResponse `json:"message"`
}

type GetMessagesRequest struct {
	RoomID        uint   `json:"-" validate:"required,min=1"`
	ParticipantID uint   `json:"-" validate:"required,min=1"`
	Limit         int    `json:"limit" validate:"omitempty,min=1,max=100"`
	Before        *int64 `json:"before" validate:"omitempty,required"`
}

type MessageListResponse struct {
	Messages []MessageResponse `json:"messages"`
	HasMore  bool              `json:"has_more"`
}

type MessageParticipantInfo struct {
	ID          uint   `json:"id"`
	DisplayName string `json:"display_name"`
}
