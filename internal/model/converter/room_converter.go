package converter

import (
	"slido-clone-backend/internal/entity"
	"slido-clone-backend/internal/model"
)

// RoomToResponse convert entity Room to model RoomResponse
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

// RoomToCreateRoomResponse convert entity Room and participantID to model CreateRoomResponse
func RoomToCreateRoomResponse(room *entity.Room, participantID uint) *model.CreateRoomResponse {
	return &model.CreateRoomResponse{
		Room:          *RoomToResponse(room),
		ParticipantID: participantID,
	}
}

// RoomToDetailResponse convert entity Room with relations to model RoomDetailResponse
func RoomToDetailResponse(room *entity.Room) *model.RoomDetailResponse {
	return &model.RoomDetailResponse{
		ID:       room.ID,
		RoomCode: room.RoomCode,
		Title:    room.Title,
		Status:   room.Status,
		Presenter: model.PresenterInfo{
			ID:       room.Presenter.ID,
			Username: room.Presenter.Username,
		},
		Stats: model.RoomStats{
			TotalParticipants: len(room.Participants),
			TotalQuestions:    len(room.Questions),
			TotalPolls:        len(room.Polls),
			// TODO: implement logic to get active poll ID
		},
		CreatedAt: room.CreatedAt,
	}
}

// RoomToUpdateToCloseResponse convert entity Room to model UpdateToCloseRoom for close room response
func RoomToUpdateToCloseResponse(room *entity.Room) *model.UpdateToCloseRoom {
	return &model.UpdateToCloseRoom{
		ID:       room.ID,
		Status:   room.Status,
		ClosedAt: room.ClosedAt,
	}
}

// RoomToListItemResponse convert entity Room to model RoomListItem
func RoomToListItemResponse(room *entity.Room) *model.RoomListItem {
	return &model.RoomListItem{
		ID:                room.ID,
		RoomCode:          room.RoomCode,
		Title:             room.Title,
		Status:            room.Status,
		ParticipantsCount: len(room.Participants),
		CreatedAt:         room.CreatedAt,
	}
}

// RoomsToListResponse wrap list of RoomListItem to RoomListResponse
func RoomsToListResponse(roomsList []*model.RoomListItem) *model.RoomListResponse {
	return &model.RoomListResponse{
		RoomListItem: roomsList,
	}
}
