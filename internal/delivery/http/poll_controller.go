package http

import (
	"context"
	"encoding/json"
	"slido-clone-backend/internal/delivery/http/middleware"
	"slido-clone-backend/internal/delivery/websocket"
	"slido-clone-backend/internal/model"
	"slido-clone-backend/internal/usecase"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

// PollController controller untuk Poll operations
type PollController struct {
	Log                *logrus.Logger
	PollUseCase        *usecase.PollUseCase
	ParticipantUseCase *usecase.ParticipantUseCase
	WSHub              *websocket.Hub
}

// NewPollController create new instance of PollController
func NewPollController(log *logrus.Logger, pollUseCase *usecase.PollUseCase, participantUseCase *usecase.ParticipantUseCase, wsHub *websocket.Hub) *PollController {
	return &PollController{
		Log:                log,
		PollUseCase:        pollUseCase,
		ParticipantUseCase: participantUseCase,
		WSHub:              wsHub,
	}
}

// Create handler untuk membuat poll baru (presenter only)
func (c *PollController) Create(ctx *fiber.Ctx) error {
	auth := middleware.GetUser(ctx)

	// check if user is presenter
	if auth.Role != "presenter" && auth.Role != "admin" {
		c.Log.Warnf("Create - User is not presenter")
		return fiber.ErrForbidden
	}

	// parse room_id from params
	roomIDStr := ctx.Params("room_id")
	roomIDUint64, err := strconv.ParseUint(roomIDStr, 10, 64)
	if err != nil {
		c.Log.Warnf("Create - Invalid room_id: %v", err)
		return fiber.ErrBadRequest
	}

	// create request
	request := &model.CreatePollRequest{
		RoomID:      uint(roomIDUint64),
		PresenterID: *auth.UserID,
	}

	// parse body
	if err = ctx.BodyParser(request); err != nil {
		c.Log.Warnf("Create - BodyParser error: %v", err)
		return fiber.ErrBadRequest
	}

	// call usecase
	response, err := c.PollUseCase.Create(ctx.UserContext(), request)
	if err != nil {
		c.Log.Warnf("Create - PollUseCase.Create error: %v", err)
		return err
	}

	// broadcast ke websocket clients di room
	c.broadcastPollCreated(uint(roomIDUint64), response)

	return ctx.Status(fiber.StatusCreated).JSON(model.WebResponse{
		Data: response,
	})
}

// GetActive handler untuk mendapatkan active polls di room
func (c *PollController) GetActive(ctx *fiber.Ctx) error {
	auth := middleware.GetUser(ctx)

	// parse room_id from params
	roomIDStr := ctx.Params("room_id")
	roomIDUint64, err := strconv.ParseUint(roomIDStr, 10, 64)
	if err != nil {
		c.Log.Warnf("GetActive - Invalid room_id: %v", err)
		return fiber.ErrBadRequest
	}

	// create request
	request := &model.GetActivePollsRequest{
		RoomID:        uint(roomIDUint64),
		ParticipantID: *auth.ParticipantID,
	}

	// call usecase
	response, err := c.PollUseCase.GetActivePolls(ctx.UserContext(), request)
	if err != nil {
		c.Log.Warnf("GetActive - PollUseCase.GetActivePolls error: %v", err)
		return err
	}

	return ctx.Status(fiber.StatusOK).JSON(model.WebResponse{
		Data: response,
	})
}

// GetHistory handler untuk mendapatkan poll history di room
func (c *PollController) GetHistory(ctx *fiber.Ctx) error {
	// parse room_id from params
	roomIDStr := ctx.Params("room_id")
	roomIDUint64, err := strconv.ParseUint(roomIDStr, 10, 64)
	if err != nil {
		c.Log.Warnf("GetHistory - Invalid room_id: %v", err)
		return fiber.ErrBadRequest
	}

	// parse query params
	request := &model.GetPollHistoryRequest{
		RoomID: uint(roomIDUint64),
		Status: ctx.Query("status", "all"),
		Limit:  ctx.QueryInt("limit", 10),
	}

	// call usecase
	response, err := c.PollUseCase.GetHistory(ctx.UserContext(), request)
	if err != nil {
		c.Log.Warnf("GetHistory - PollUseCase.GetHistory error: %v", err)
		return err
	}

	return ctx.Status(fiber.StatusOK).JSON(model.WebResponse{
		Data: response,
	})
}

// Vote handler untuk submit vote pada poll
func (c *PollController) Vote(ctx *fiber.Ctx) error {
	auth := middleware.GetUser(ctx)

	// check if participant has joined a room
	if auth.ParticipantID == nil || auth.RoomID == nil {
		c.Log.Warnf("Vote - User has not joined a room")
		return fiber.NewError(fiber.StatusBadRequest, "You must join a room first")
	}

	// parse poll_id from params
	pollIDStr := ctx.Params("poll_id")
	pollIDUint64, err := strconv.ParseUint(pollIDStr, 10, 64)
	if err != nil {
		c.Log.Warnf("Vote - Invalid poll_id: %v", err)
		return fiber.ErrBadRequest
	}

	// create request
	request := &model.SubmitPollVoteRequest{
		PollID:        uint(pollIDUint64),
		ParticipantID: *auth.ParticipantID,
		RoomID:        *auth.RoomID,
	}

	// parse body
	if err = ctx.BodyParser(request); err != nil {
		c.Log.Warnf("Vote - BodyParser error: %v", err)
		return fiber.ErrBadRequest
	}

	// call usecase
	response, err := c.PollUseCase.Vote(ctx.UserContext(), request)
	if err != nil {
		c.Log.Warnf("Vote - PollUseCase.Vote error: %v", err)
		return err
	}

	// broadcast ke websocket clients di room
	c.broadcastPollVoted(*auth.RoomID, response)

	// broadcast leaderboard update setelah user vote poll (karena dapat XP)
	c.broadcastLeaderboardUpdate(*auth.RoomID, *auth.ParticipantID)

	return ctx.Status(fiber.StatusOK).JSON(model.WebResponse{
		Data: response,
	})
}

// Close handler untuk menutup poll (presenter only)
func (c *PollController) Close(ctx *fiber.Ctx) error {
	auth := middleware.GetUser(ctx)

	// check if user is presenter
	if auth.Role != "presenter" && auth.Role != "admin" {
		c.Log.Warnf("Close - User is not presenter")
		return fiber.ErrForbidden
	}

	// parse poll_id from params
	pollIDStr := ctx.Params("poll_id")
	pollIDUint64, err := strconv.ParseUint(pollIDStr, 10, 64)
	if err != nil {
		c.Log.Warnf("Close - Invalid poll_id: %v", err)
		return fiber.ErrBadRequest
	}

	// create request
	request := &model.ClosePollRequest{
		PollID:      uint(pollIDUint64),
		PresenterID: *auth.UserID,
	}

	// call usecase
	response, err := c.PollUseCase.Close(ctx.UserContext(), request)
	if err != nil {
		c.Log.Warnf("Close - PollUseCase.Close error: %v", err)
		return err
	}

	// broadcast ke websocket clients di room - get room_id from context
	if auth.RoomID != nil {
		c.broadcastPollClosed(*auth.RoomID, response)
	}

	return ctx.Status(fiber.StatusOK).JSON(model.WebResponse{
		Data: response,
	})
}

// ========================================
// WebSocket Broadcast Functions
// ========================================

// broadcastPollCreated broadcast event poll created ke room
func (c *PollController) broadcastPollCreated(roomID uint, response *model.CreatePollResponse) {
	if c.WSHub == nil {
		return
	}
	data := websocket.WSMessage{
		Event: websocket.EventPollCreated,
		Data:  pollMustMarshalJSON(response),
	}
	c.WSHub.BroadcastToRoom(roomID, pollMustMarshalJSON(data))
}

// broadcastPollVoted broadcast event poll voted ke room
func (c *PollController) broadcastPollVoted(roomID uint, response *model.SubmitPollVoteResponse) {
	if c.WSHub == nil {
		return
	}
	// hanya kirim updated results, tidak kirim data personal voter
	broadcastData := struct {
		UpdatedResults model.UpdatedPollResultsResponse `json:"updated_results"`
	}{
		UpdatedResults: response.UpdatedResults,
	}
	data := websocket.WSMessage{
		Event: websocket.EventPollResultsUpdate,
		Data:  pollMustMarshalJSON(broadcastData),
	}
	c.WSHub.BroadcastToRoom(roomID, pollMustMarshalJSON(data))
}

// broadcastPollClosed broadcast event poll closed ke room
func (c *PollController) broadcastPollClosed(roomID uint, response *model.ClosePollResponse) {
	if c.WSHub == nil {
		return
	}
	data := websocket.WSMessage{
		Event: websocket.EventPollClosed,
		Data:  pollMustMarshalJSON(response),
	}
	c.WSHub.BroadcastToRoom(roomID, pollMustMarshalJSON(data))
}

// broadcastLeaderboardUpdate broadcast leaderboard setelah ada perubahan XP
func (c *PollController) broadcastLeaderboardUpdate(roomID uint, participantID uint) {
	if c.WSHub == nil || c.ParticipantUseCase == nil {
		return
	}

	request := &model.GetLeaderboardRequest{
		RoomID:        roomID,
		ParticipantID: participantID,
	}

	leaderboard, err := c.ParticipantUseCase.Leaderboard(context.Background(), request)
	if err != nil {
		c.Log.Warnf("broadcastLeaderboardUpdate - error: %v", err)
		return
	}

	data := websocket.WSMessage{
		Event: websocket.EventLeaderboardUpdate,
		Data: pollMustMarshalJSON(map[string]interface{}{
			"leaderboard":        leaderboard,
			"total_participants": leaderboard.TotalParticipants,
		}),
	}
	c.WSHub.BroadcastToRoom(roomID, pollMustMarshalJSON(data))
}

// pollMustMarshalJSON helper untuk marshal JSON
func pollMustMarshalJSON(v interface{}) []byte {
	data, err := json.Marshal(v)
	if err != nil {
		panic(err)
	}
	return data
}
