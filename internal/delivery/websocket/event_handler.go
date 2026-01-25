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
	questionUseCase    *usecase.QuestionUseCase
}

func NewEventHandler(messageUseCase *usecase.MessageUseCase, participantUseCase *usecase.ParticipantUseCase, questionUseCase *usecase.QuestionUseCase) *EventHandler {
	return &EventHandler{
		messageUseCase:     messageUseCase,
		participantUseCase: participantUseCase,
		questionUseCase:    questionUseCase,
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
	// Q&A events
	case EventQuestionSubmit:
		return h.handleQuestionSubmit(client, wsMsg.Data)
	case EventQuestionUpvote:
		return h.handleQuestionUpvote(client, wsMsg.Data)
	case EventQuestionRemoveUpvote:
		return h.handleQuestionRemoveUpvote(client, wsMsg.Data)
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
		Data:  mustMarshal(response),
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

// handleQuestionSubmit handle submit question via websocket
func (h *EventHandler) handleQuestionSubmit(client *Client, data json.RawMessage) error {
	// parse payload
	var payload struct {
		Content string `json:"content"`
	}
	if err := json.Unmarshal(data, &payload); err != nil {
		client.hub.log.WithField("error", err).Warn("failed to parse question payload")
		return err
	}

	// create request untuk usecase
	request := &model.SubmitQuestionRequest{
		RoomID:        client.roomID,
		ParticipantID: client.participantID,
		Content:       payload.Content,
	}

	// call usecase
	response, err := h.questionUseCase.Submit(context.Background(), request)
	if err != nil {
		client.hub.log.WithField("error", err).Warn("failed to submit question")
		return err
	}

	// broadcast question:created ke semua client di room
	broadcastData := WSMessage{
		Event: EventQuestionCreated,
		Data:  mustMarshal(response),
	}
	client.hub.BroadcastToRoom(client.roomID, mustMarshal(broadcastData))
	h.broadcastLeaderboardUpdate(client)
	return nil
}

// handleQuestionUpvote handle upvote question via websocket
func (h *EventHandler) handleQuestionUpvote(client *Client, data json.RawMessage) error {
	// parse payload
	var payload struct {
		QuestionID uint `json:"question_id"`
	}
	if err := json.Unmarshal(data, &payload); err != nil {
		client.hub.log.WithField("error", err).Warn("failed to parse upvote payload")
		return err
	}

	// create request untuk usecase
	request := &model.UpvoteRequest{
		QuestionID:    payload.QuestionID,
		ParticipantID: client.participantID,
		RoomID:        client.roomID,
	}

	// call usecase
	response, err := h.questionUseCase.Upvote(context.Background(), request)
	if err != nil {
		client.hub.log.WithField("error", err).Warn("failed to upvote question")
		return err
	}

	// broadcast question:upvoted ke semua client di room
	broadcastData := WSMessage{
		Event: EventQuestionUpvoted,
		Data:  mustMarshal(response),
	}
	client.hub.BroadcastToRoom(client.roomID, mustMarshal(broadcastData))
	h.broadcastLeaderboardUpdate(client)
	return nil
}

// handleQuestionRemoveUpvote handle remove upvote via websocket
func (h *EventHandler) handleQuestionRemoveUpvote(client *Client, data json.RawMessage) error {
	// parse payload
	var payload struct {
		QuestionID uint `json:"question_id"`
	}
	if err := json.Unmarshal(data, &payload); err != nil {
		client.hub.log.WithField("error", err).Warn("failed to parse remove upvote payload")
		return err
	}

	// create request untuk usecase
	request := &model.UpvoteRequest{
		QuestionID:    payload.QuestionID,
		ParticipantID: client.participantID,
		RoomID:        client.roomID,
	}

	// call usecase
	response, err := h.questionUseCase.RemoveUpvote(context.Background(), request)
	if err != nil {
		client.hub.log.WithField("error", err).Warn("failed to remove upvote")
		return err
	}

	// broadcast question:upvoted (dengan updated count) ke semua client di room
	broadcastData := WSMessage{
		Event: EventQuestionUpvoted,
		Data:  mustMarshal(response),
	}
	client.hub.BroadcastToRoom(client.roomID, mustMarshal(broadcastData))
	h.broadcastLeaderboardUpdate(client)
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
