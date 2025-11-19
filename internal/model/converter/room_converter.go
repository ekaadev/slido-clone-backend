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

func RoomToGetRoomDetailResponse(room *entity.Room) *model.GetRoomDetailResponse {
	return &model.GetRoomDetailResponse{
		Room: *RoomToDetailResponse(room),
	}
}

func RoomToUpdateToCloseResponse(room *entity.Room) *model.UpdateToCloseRoom {
	return &model.UpdateToCloseRoom{
		ID:       room.ID,
		Status:   room.Status,
		ClosedAt: room.ClosedAt,
	}
}

func UpdateToCloseRoomToResponse(update *entity.Room) *model.UpdateToCloseRoomResponse {
	return &model.UpdateToCloseRoomResponse{
		Room: *RoomToUpdateToCloseResponse(update),
	}
}

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

func RoomsToListResponse(roomsList []*model.RoomListItem) *model.RoomListResponse {
	return &model.RoomListResponse{
		RoomListItem: roomsList,
	}
}
