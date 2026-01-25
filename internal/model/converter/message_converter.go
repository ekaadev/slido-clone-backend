package converter

import (
	"slido-clone-backend/internal/entity"
	"slido-clone-backend/internal/model"
)

// MessageToResponse convert entity Message to model MessageResponse
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

// MessagesToMessageListResponse convert list of entity Message to model MessageListResponse
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
