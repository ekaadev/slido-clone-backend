package converter

import (
	"slido-clone-backend/internal/entity"
	"slido-clone-backend/internal/model"
)

func ParticipantToResponse(participant *entity.Participant) *model.ParticipantResponse {
	return &model.ParticipantResponse{
		ID:          participant.ID,
		RoomID:      participant.RoomID,
		DisplayName: participant.DisplayName,
		XPScore:     participant.XPScore,
		IsAnonymous: *participant.IsAnonymous,
		JoinedAt:    participant.JoinedAt,
	}
}

func ParticipantToJoinRoomResponse(participant *entity.Participant, token string) *model.JoinRoomResponse {
	return &model.JoinRoomResponse{
		Participant: *ParticipantToResponse(participant),
		Token:       token,
	}
}
