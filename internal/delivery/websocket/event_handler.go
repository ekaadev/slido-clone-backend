package websocket

import (
	"context"
	"encoding/json"
	"fmt"
	"slido-clone-backend/internal/model"
	"slido-clone-backend/internal/sfu"
	"slido-clone-backend/internal/usecase"
	"time"

	"github.com/pion/webrtc/v4"
)

type EventHandler struct {
	messageUseCase     *usecase.MessageUseCase
	participantUseCase *usecase.ParticipantUseCase
	questionUseCase    *usecase.QuestionUseCase
	sfuManager         *sfu.SFUManager
}

func NewEventHandler(messageUseCase *usecase.MessageUseCase, participantUseCase *usecase.ParticipantUseCase, questionUseCase *usecase.QuestionUseCase, sfuManager *sfu.SFUManager) *EventHandler {
	return &EventHandler{
		messageUseCase:     messageUseCase,
		participantUseCase: participantUseCase,
		questionUseCase:    questionUseCase,
		sfuManager:         sfuManager,
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
	// WebRTC
	case EventWebrtcOffer:
		return h.handleWebrtcOffer(client, wsMsg.Data)
	case EventWebrtcAnswer:
		return h.handleWebrtcAnswer(client, wsMsg.Data)
	case EventWebrtcCandidate:
		return h.handleWebrtcCandidate(client, wsMsg.Data)
	// Conference events
	case EventConferenceStart:
		return h.handleConferenceStart(client)
	case EventConferenceStop:
		return h.handleConferenceStop(client)
	case EventConferenceJoin:
		return h.handleConferenceJoin(client)
	case EventConferenceLeave:
		return h.handleConferenceLeave(client)
	case EventRaiseHand:
		return h.handleRaiseHand(client)
	case EventLowerHand:
		return h.handleLowerHand(client)
	case EventPromoteSpeaker:
		return h.handlePromoteSpeaker(client, wsMsg.Data)
	case EventDemoteSpeaker:
		return h.handleDemoteSpeaker(client, wsMsg.Data)
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

	// broadcast question:upvoted ke semua client di room dengan action dan participant_id
	broadcastPayload := map[string]interface{}{
		"question":       response.Question,
		"participant_id": client.participantID,
		"action":         "add",
	}
	broadcastData := WSMessage{
		Event: EventQuestionUpvoted,
		Data:  mustMarshal(broadcastPayload),
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

	// broadcast question:upvoted (dengan updated count) ke semua client di room dengan action dan participant_id
	broadcastPayload := map[string]interface{}{
		"question":       response.Question,
		"participant_id": client.participantID,
		"action":         "remove",
	}
	broadcastData := WSMessage{
		Event: EventQuestionUpvoted,
		Data:  mustMarshal(broadcastPayload),
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

func (h *EventHandler) HandleDisconnect(client *Client) {
	peerID := fmt.Sprintf("%d", client.participantID)
	// We use the sfuManager directly.
	// Make sure RemovePeer is safe to call even if peer doesn't exist.
	h.sfuManager.RemovePeer(client.roomID, peerID)
}

func (h *EventHandler) handleWebrtcOffer(client *Client, data json.RawMessage) error {
	var offerPayload struct {
		Type        string `json:"type"`
		SDP         string `json:"sdp"`
		Renegotiate bool   `json:"renegotiate"`
		Reason      string `json:"reason"`
	}
	if err := json.Unmarshal(data, &offerPayload); err != nil {
		client.hub.log.WithField("error", err).Warn("failed to parse webrtc offer")
		return err
	}

	offer := webrtc.SessionDescription{
		Type: webrtc.SDPTypeOffer,
		SDP:  offerPayload.SDP,
	}

	peerID := fmt.Sprintf("%d", client.participantID)

	signalFunc := func(payload interface{}) {
		// Prevent panic if channel is closed
		defer func() {
			if r := recover(); r != nil {
				client.hub.log.Warnf("failed to send signal: %v", r)
			}
		}()

		pl, ok := payload.(map[string]interface{})
		if !ok {
			return
		}
		typ, _ := pl["type"].(string)

		var event string
		switch typ {
		case "offer":
			event = EventWebrtcOffer
		case "answer":
			event = EventWebrtcAnswer
		case "candidate":
			event = EventWebrtcCandidate
		default:
			event = "webrtc:signal"
		}

		msg := WSMessage{
			Event: event,
			Data:  mustMarshal(payload),
		}
		client.send <- mustMarshal(msg)
	}

	room := h.sfuManager.GetRoom(client.roomID)
	existingPeer := room.GetPeer(peerID)

	// Check if this is a renegotiation (peer already exists)
	if existingPeer != nil && offerPayload.Renegotiate {
		client.hub.log.WithFields(map[string]interface{}{
			"participant_id": peerID,
			"reason":         offerPayload.Reason,
		}).Info("Handling renegotiation offer")

		// Handle renegotiation - just process the new offer without recreating peer
		return existingPeer.HandleOffer(offer)
	}

	// Create new peer if doesn't exist
	peer, err := h.sfuManager.CreatePeer(client.roomID, peerID, signalFunc)
	if err != nil {
		client.hub.log.WithField("error", err).Warn("failed to create peer")
		return err
	}

	return peer.HandleOffer(offer)
}

func (h *EventHandler) handleWebrtcAnswer(client *Client, data json.RawMessage) error {
	var answer webrtc.SessionDescription
	if err := json.Unmarshal(data, &answer); err != nil {
		client.hub.log.WithField("error", err).Warn("failed to parse webrtc answer")
		return err
	}
	peerID := fmt.Sprintf("%d", client.participantID)
	return h.sfuManager.HandleAnswer(client.roomID, peerID, answer)
}

func (h *EventHandler) handleWebrtcCandidate(client *Client, data json.RawMessage) error {
	var candidate webrtc.ICECandidateInit
	if err := json.Unmarshal(data, &candidate); err != nil {
		client.hub.log.WithField("error", err).Warn("failed to parse webrtc candidate")
		return err
	}
	peerID := fmt.Sprintf("%d", client.participantID)
	return h.sfuManager.HandleCandidate(client.roomID, peerID, candidate)
}

// mustMarshal helper untuk marshal JSON, panic jika error
func mustMarshal(v interface{}) []byte {
	data, err := json.Marshal(v)
	if err != nil {
		panic(err)
	}
	return data
}

// Conference Handlers

func (h *EventHandler) handleConferenceStart(client *Client) error {
	// Authorization: Only room owner (host) can start conference
	if !client.isRoomOwner {
		return fmt.Errorf("unauthorized: only room owner can start conference")
	}

	peerID := fmt.Sprintf("%d", client.participantID)
	room := h.sfuManager.GetRoom(client.roomID)

	if err := room.StartConference(peerID); err != nil {
		return err
	}

	// Broadcast conference started to all clients in room
	state := room.GetConferenceState()
	broadcastData := WSMessage{
		Event: EventConferenceStarted,
		Data: mustMarshal(map[string]interface{}{
			"host_id":      state.HostID,
			"is_active":    state.IsActive,
			"speakers":     state.Speakers,
			"raised_hands": state.RaisedHands,
		}),
	}
	client.hub.BroadcastToRoom(client.roomID, mustMarshal(broadcastData))
	return nil
}

func (h *EventHandler) handleConferenceStop(client *Client) error {
	// Authorization: Only room owner (host) can stop conference
	if !client.isRoomOwner {
		return fmt.Errorf("unauthorized: only room owner can stop conference")
	}

	peerID := fmt.Sprintf("%d", client.participantID)
	room := h.sfuManager.GetRoom(client.roomID)

	if err := room.StopConference(peerID); err != nil {
		return err
	}

	// Broadcast conference ended to all clients in room
	broadcastData := WSMessage{
		Event: EventConferenceEnded,
		Data:  mustMarshal(map[string]interface{}{}),
	}
	client.hub.BroadcastToRoom(client.roomID, mustMarshal(broadcastData))
	return nil
}

func (h *EventHandler) handleConferenceJoin(client *Client) error {
	peerID := fmt.Sprintf("%d", client.participantID)
	room := h.sfuManager.GetRoom(client.roomID)

	// Send current state to the joining client (including their role info)
	state := room.GetConferenceState()
	stateData := WSMessage{
		Event: EventConferenceState,
		Data: mustMarshal(map[string]interface{}{
			"host_id":       state.HostID,
			"is_active":     state.IsActive,
			"speakers":      state.Speakers,
			"raised_hands":  state.RaisedHands,
			"is_room_owner": client.isRoomOwner, // inform client their role
		}),
	}
	client.send <- mustMarshal(stateData)

	// Broadcast that someone joined (with their role info)
	broadcastData := WSMessage{
		Event: EventConferenceJoined,
		Data: mustMarshal(map[string]interface{}{
			"participant_id": peerID,
			"is_room_owner":  client.isRoomOwner,
		}),
	}
	client.hub.BroadcastToRoom(client.roomID, mustMarshal(broadcastData))
	return nil
}

func (h *EventHandler) handleConferenceLeave(client *Client) error {
	peerID := fmt.Sprintf("%d", client.participantID)

	// Broadcast that someone left
	broadcastData := WSMessage{
		Event: EventConferenceLeft,
		Data: mustMarshal(map[string]interface{}{
			"participant_id": peerID,
		}),
	}
	client.hub.BroadcastToRoom(client.roomID, mustMarshal(broadcastData))
	return nil
}

func (h *EventHandler) handleRaiseHand(client *Client) error {
	peerID := fmt.Sprintf("%d", client.participantID)
	room := h.sfuManager.GetRoom(client.roomID)

	room.RaiseHand(peerID, time.Now().Unix())

	// Broadcast hand raised
	broadcastData := WSMessage{
		Event: EventHandRaised,
		Data: mustMarshal(map[string]interface{}{
			"participant_id": peerID,
			"timestamp":      time.Now().Unix(),
		}),
	}
	client.hub.BroadcastToRoom(client.roomID, mustMarshal(broadcastData))
	return nil
}

func (h *EventHandler) handleLowerHand(client *Client) error {
	peerID := fmt.Sprintf("%d", client.participantID)
	room := h.sfuManager.GetRoom(client.roomID)

	room.LowerHand(peerID)

	// Broadcast hand lowered
	broadcastData := WSMessage{
		Event: EventHandLowered,
		Data: mustMarshal(map[string]interface{}{
			"participant_id": peerID,
		}),
	}
	client.hub.BroadcastToRoom(client.roomID, mustMarshal(broadcastData))
	return nil
}

func (h *EventHandler) handlePromoteSpeaker(client *Client, data json.RawMessage) error {
	// Authorization: Only room owner (host) can promote speakers
	if !client.isRoomOwner {
		return fmt.Errorf("unauthorized: only room owner can promote speakers")
	}

	var payload struct {
		ParticipantID string `json:"participant_id"`
	}
	if err := json.Unmarshal(data, &payload); err != nil {
		return err
	}

	hostID := fmt.Sprintf("%d", client.participantID)
	room := h.sfuManager.GetRoom(client.roomID)

	if !room.PromoteSpeaker(hostID, payload.ParticipantID) {
		return fmt.Errorf("not authorized to promote")
	}

	// Broadcast speaker promoted
	broadcastData := WSMessage{
		Event: EventSpeakerPromoted,
		Data: mustMarshal(map[string]interface{}{
			"participant_id": payload.ParticipantID,
		}),
	}
	client.hub.BroadcastToRoom(client.roomID, mustMarshal(broadcastData))
	return nil
}

func (h *EventHandler) handleDemoteSpeaker(client *Client, data json.RawMessage) error {
	// Authorization: Only room owner (host) can demote speakers
	if !client.isRoomOwner {
		return fmt.Errorf("unauthorized: only room owner can demote speakers")
	}

	var payload struct {
		ParticipantID string `json:"participant_id"`
	}
	if err := json.Unmarshal(data, &payload); err != nil {
		return err
	}

	hostID := fmt.Sprintf("%d", client.participantID)
	room := h.sfuManager.GetRoom(client.roomID)

	if !room.DemoteSpeaker(hostID, payload.ParticipantID) {
		return fmt.Errorf("not authorized to demote")
	}

	// Broadcast speaker demoted
	broadcastData := WSMessage{
		Event: EventSpeakerDemoted,
		Data: mustMarshal(map[string]interface{}{
			"participant_id": payload.ParticipantID,
		}),
	}
	client.hub.BroadcastToRoom(client.roomID, mustMarshal(broadcastData))
	return nil
}
