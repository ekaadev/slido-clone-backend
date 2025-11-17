package converter

import (
	"slido-clone-backend/internal/entity"
	"slido-clone-backend/internal/model"
)

func RoomToResponse(room *entity.Room) *model.RoomResponse {
	return &model.RoomResponse{
		ID:          room.ID,
		RoomCode:    room.RoomCode,
		Title:       room.Title,
		PresenterID: room.PresenterID,
		Status:      room.Status,
		CreatedAt:   room.CreatedAt,
		ClosedAt:    room.ClosedAt,
	}
}

func RoomToCreateRoomResponse(room *entity.Room) *model.CreateRoomResponse {
	return &model.CreateRoomResponse{
		Room: *RoomToResponse(room),
	}
}
