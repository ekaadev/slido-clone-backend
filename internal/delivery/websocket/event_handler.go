package websocket

import (
	"context"
	"encoding/json"
	"slido-clone-backend/internal/model"
	"slido-clone-backend/internal/usecase"
)

type EventHandler struct {
	messageUseCase     *usecase.MessageUseCase
	participantUseCase *usecase.ParticipantUseCase
}

func NewEventHandler(messageUseCase *usecase.MessageUseCase, participantUseCase *usecase.ParticipantUseCase) *EventHandler {
	return &EventHandler{
		messageUseCase:     messageUseCase,
		participantUseCase: participantUseCase,
	}
}

// HandleMessage process incoming websocket messages
func (h *EventHandler) HandleMessage(client *Client, data []byte) error {
	var wsMsg WSMessage
	if err := json.Unmarshal(data, &wsMsg); err != nil {
		client.hub.log.WithField("error", err).Warn("failed to parse ws message")
		return err
	}

	// route ke handler berdasarkan event type
	switch wsMsg.Event {
	case EventMessageSend:
		return h.handleMessageSend(client, wsMsg.Data)
	case EventChatTyping:
		return h.handleChatTyping(client, wsMsg.Data)
	case EventLeaderboardRequest:
		return h.handleLeaderboardRequest(client, wsMsg.Data)
	default:
		client.hub.log.WithField("event", wsMsg.Event).Warn("unknown event")
		return nil
	}
}

func (h *EventHandler) handleMessageSend(client *Client, data json.RawMessage) error {
	// parse payload
	var payload struct {
		Content string `json:"content"`
	}
	if err := json.Unmarshal(data, &payload); err != nil {
		client.hub.log.WithField("error", err).Warn("failed to parse message payload")
		return err
	}

	// create request untuk usecase
	request := &model.SendMessageRequest{
		RoomID:        client.roomID,
		ParticipantID: client.participantID,
		Content:       payload.Content,
	}

	// panggil usecase
	response, err := h.messageUseCase.Send(context.Background(), request)
	if err != nil {
		client.hub.log.WithField("error", err).Warn("failed to send message")
		return err
	}

	// broadcast message ke semua client di room
	broadcastData := WSMessage{
		Event: EventMessageSend,
		Data:  mustMarshal(response.Message),
	}
	client.hub.BroadcastToRoom(client.roomID, mustMarshal(broadcastData))
	h.broadcastLeaderboardUpdate(client)
	return nil
}

// handleChatTyping handle typing indicatior
func (h *EventHandler) handleChatTyping(client *Client, data json.RawMessage) error {
	// parse payload
	var payload struct {
		IsTyping bool `json:"is_typing"`
	}
	if err := json.Unmarshal(data, &payload); err != nil {
		client.hub.log.WithField("error", err).Warn("failed to parse message payload")
		return err
	}

	// broadcast typing status ke semua client di room (kecuali sender)
	typingData := WSMessage{
		Event: EventChatTyping,
		Data: mustMarshal(map[string]interface{}{
			"participant_id": client.participantID,
			"is_typing":      payload.IsTyping,
		}),
	}

	// broadcast ke room (implmenetasi nanti jika perlu exclude sender)
	client.hub.BroadcastToRoom(client.roomID, mustMarshal(typingData))
	return nil
}

// handleLeaderboardRequest handle request leaderboard
func (h *EventHandler) handleLeaderboardRequest(client *Client, data json.RawMessage) error {
	request := &model.GetLeaderboardRequest{
		RoomID:        client.roomID,
		ParticipantID: client.participantID,
	}

	leaderboard, err := h.participantUseCase.Leaderboard(context.Background(), request)
	if err != nil {
		client.hub.log.WithField("error", err).Warn("failed to get leaderboard")
		return err
	}

	responseData := WSMessage{
		Event: EventLeaderboardUpdate,
		Data:  mustMarshal(leaderboard),
	}

	client.send <- mustMarshal(responseData)
	return nil
}

func (h *EventHandler) broadcastLeaderboardUpdate(client *Client) {
	request := &model.GetLeaderboardRequest{
		RoomID:        client.roomID,
		ParticipantID: client.participantID,
	}

	leaderboard, err := h.participantUseCase.Leaderboard(context.Background(), request)
	if err != nil {
		client.hub.log.WithField("error", err).Warn("failed to get leaderboard")
		return
	}

	leaderboardData := WSMessage{
		Event: EventLeaderboardUpdate,
		Data: mustMarshal(map[string]interface{}{
			"leaderboard":        leaderboard,
			"total_participants": leaderboard.TotalParticipants,
		}),
	}

	client.hub.BroadcastToRoom(client.roomID, mustMarshal(leaderboardData))
}

// mustMarshal helper untuk marshal JSON, panic jika error
func mustMarshal(v interface{}) []byte {
	data, err := json.Marshal(v)
	if err != nil {
		panic(err)
	}
	return data
}
