package converter

import (
	"slido-clone-backend/internal/entity"
	"slido-clone-backend/internal/model"
)

func MessageToResponse(message *entity.Message) *model.MessageResponse {
	return &model.MessageResponse{
		ID:     message.ID,
		RoomID: message.RoomID,
		Participant: model.ParticipantInfo{
			ID:          message.Participant.ID,
			DisplayName: message.Participant.DisplayName,
		},
		Content:   message.Content,
		CreatedAt: message.CreatedAt,
	}
}

func MessageToSendResponse(message *entity.Message) *model.SendMessageResponse {
	return &model.SendMessageResponse{
		Message: *MessageToResponse(message),
	}
}

func MessagesToMessageListResponse(messages []entity.Message, hasMore bool) *model.MessageListResponse {
	responses := make([]model.MessageResponse, len(messages))
	for i, message := range messages {
		responses[i] = *MessageToResponse(&message)
	}

	return &model.MessageListResponse{
		Messages: responses,
		HasMore:  hasMore,
	}
}
