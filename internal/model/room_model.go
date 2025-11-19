package model

import "time"

type CreateRoomRequest struct {
	Title       string `json:"title" validate:"required,min=3,max=255"`
	PresenterID uint   `json:"presenter_id" validate:"required,min=1"`
}

type CloseRoomRequest struct {
	Status string `json:"status" validate:"required,oneof=closed"`
}

type RoomResponse struct {
	ID          uint       `json:"id"`
	RoomCode    string     `json:"room_code"`
	Title       string     `json:"title"`
	PresenterID uint       `json:"presenter_id"`
	Status      string     `json:"status"`
	CreatedAt   time.Time  `json:"created_at"`
	ClosedAt    *time.Time `json:"closed_at,omitempty"`
}

type CreateRoomResponse struct {
	Room RoomResponse `json:"room"`
}

type RoomDetailResponse struct {
	ID        uint          `json:"id"`
	RoomCode  string        `json:"room_code"`
	Title     string        `json:"title"`
	Status    string        `json:"status"`
	Presenter PresenterInfo `json:"presenter"`
	Stats     RoomStats     `json:"stats"`
	CreatedAt time.Time     `json:"created_at"`
}

type GetRoomDetailResponse struct {
	Room RoomDetailResponse `json:"room"`
}

type PresenterInfo struct {
	ID       uint   `json:"id"`
	Username string `json:"username"`
}

type RoomStats struct {
	TotalParticipants int   `json:"total_participants"`
	TotalQuestions    int   `json:"total_questions"`
	TotalPolls        int   `json:"total_polls"`
	ActivePollID      *uint `json:"active_poll_id,omitempty"`
}

type RoomListItem struct {
	ID                uint      `json:"id"`
	RoomCode          string    `json:"room_code"`
	Title             string    `json:"title"`
	Status            string    `json:"status"`
	ParticipantsCount int       `json:"participants_count"`
	CreatedAt         time.Time `json:"created_at"`
}

type RoomListResponse struct {
	RoomListItem []*RoomListItem `json:"rooms"`
}

type GetRoomRequestByRoomCode struct {
	RoomCode string `json:"-" validate:"required,len=6,alphanum"`
}

type UpdateToCloseRoomRequestByID struct {
	PresenterID uint   `json:"-" validate:"required,min=1"`
	RoomID      uint   `json:"-" validate:"required,min=1"`
	Status      string `json:"status" validate:"required,oneof=closed"`
}

type UpdateToCloseRoom struct {
	ID       uint       `json:"id"`
	Status   string     `json:"status"`
	ClosedAt *time.Time `json:"closed_at,omitempty"`
}

type UpdateToCloseRoomResponse struct {
	Room UpdateToCloseRoom `json:"room"`
}

type SearchRoomsRequest struct {
	PresenterID uint `json:"-" validate:"required,min=1"`
}
