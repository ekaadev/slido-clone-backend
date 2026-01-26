package http

import (
	"encoding/json"
	"slido-clone-backend/internal/delivery/http/middleware"
	"slido-clone-backend/internal/delivery/websocket"
	"slido-clone-backend/internal/model"
	"slido-clone-backend/internal/usecase"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

// PollController controller untuk Polling operations
type PollController struct {
	Log         *logrus.Logger
	PollUseCase *usecase.PollUseCase
	WSHub       *websocket.Hub
}

// NewPollController create new instance of PollController
func NewPollController(log *logrus.Logger, pollUseCase *usecase.PollUseCase, wsHub *websocket.Hub) *PollController {
	return &PollController{
		Log:         log,
		PollUseCase: pollUseCase,
		WSHub:       wsHub,
	}
}

// Create handler untuk membuat poll baru (presenter only)
func (c *PollController) Create(ctx *fiber.Ctx) error {
	auth := middleware.GetUser(ctx)

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

// GetActive handler untuk mendapatkan active polls
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
	request := &model.GetActivePollRequest{
		RoomID:        uint(roomIDUint64),
		ParticipantID: *auth.ParticipantID,
	}

	// call usecase
	response, err := c.PollUseCase.GetActive(ctx.UserContext(), request)
	if err != nil {
		c.Log.Warnf("GetActive - PollUseCase.GetActive error: %v", err)
		return err
	}

	return ctx.Status(fiber.StatusOK).JSON(model.WebResponse{
		Data: response,
	})
}

// SubmitVote handler untuk submit vote
func (c *PollController) SubmitVote(ctx *fiber.Ctx) error {
	auth := middleware.GetUser(ctx)

	// parse poll_id from params
	pollIDStr := ctx.Params("poll_id")
	pollIDUint64, err := strconv.ParseUint(pollIDStr, 10, 64)
	if err != nil {
		c.Log.Warnf("SubmitVote - Invalid poll_id: %v", err)
		return fiber.ErrBadRequest
	}

	// get room_id from poll
	roomID, err := c.PollUseCase.PollRepository.GetRoomIDByPollID(c.PollUseCase.DB, uint(pollIDUint64))
	if err != nil {
		c.Log.Warnf("SubmitVote - GetRoomIDByPollID error: %v", err)
		return fiber.ErrInternalServerError
	}

	// create request
	request := &model.SubmitVoteRequest{
		PollID:        uint(pollIDUint64),
		ParticipantID: *auth.ParticipantID,
		RoomID:        roomID,
	}

	// parse body
	if err = ctx.BodyParser(request); err != nil {
		c.Log.Warnf("SubmitVote - BodyParser error: %v", err)
		return fiber.ErrBadRequest
	}

	// call usecase
	response, err := c.PollUseCase.SubmitVote(ctx.UserContext(), request)
	if err != nil {
		c.Log.Warnf("SubmitVote - PollUseCase.SubmitVote error: %v", err)
		return err
	}

	// broadcast ke websocket clients di room
	c.broadcastPollResultsUpdated(roomID, response)

	return ctx.Status(fiber.StatusOK).JSON(model.WebResponse{
		Data: response,
	})
}

// Close handler untuk menutup poll (presenter only)
func (c *PollController) Close(ctx *fiber.Ctx) error {
	auth := middleware.GetUser(ctx)

	// parse poll_id from params
	pollIDStr := ctx.Params("poll_id")
	pollIDUint64, err := strconv.ParseUint(pollIDStr, 10, 64)
	if err != nil {
		c.Log.Warnf("Close - Invalid poll_id: %v", err)
		return fiber.ErrBadRequest
	}

	// get room_id from poll
	roomID, err := c.PollUseCase.PollRepository.GetRoomIDByPollID(c.PollUseCase.DB, uint(pollIDUint64))
	if err != nil {
		c.Log.Warnf("Close - GetRoomIDByPollID error: %v", err)
		return fiber.ErrInternalServerError
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

	// broadcast ke websocket clients di room
	c.broadcastPollClosed(roomID, response)

	return ctx.Status(fiber.StatusOK).JSON(model.WebResponse{
		Data: response,
	})
}

// GetHistory handler untuk mendapatkan poll history
func (c *PollController) GetHistory(ctx *fiber.Ctx) error {
	auth := middleware.GetUser(ctx)

	// parse room_id from params
	roomIDStr := ctx.Params("room_id")
	roomIDUint64, err := strconv.ParseUint(roomIDStr, 10, 64)
	if err != nil {
		c.Log.Warnf("GetHistory - Invalid room_id: %v", err)
		return fiber.ErrBadRequest
	}

	// create request
	request := &model.GetPollHistoryRequest{
		RoomID:        uint(roomIDUint64),
		ParticipantID: *auth.ParticipantID,
		Status:        ctx.Query("status", "all"),
		Limit:         ctx.QueryInt("limit", 10),
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

// broadcastPollCreated broadcast poll created event
func (c *PollController) broadcastPollCreated(roomID uint, response *model.CreatePollResponse) {
	if c.WSHub == nil {
		return
	}
	data := websocket.WSMessage{
		Event: websocket.EventPollCreated,
		Data:  mustMarshalPollJSON(response),
	}
	c.WSHub.BroadcastToRoom(roomID, mustMarshalPollJSON(data))
}

// broadcastPollResultsUpdated broadcast poll results updated event
func (c *PollController) broadcastPollResultsUpdated(roomID uint, response *model.SubmitVoteResponse) {
	if c.WSHub == nil {
		return
	}
	data := websocket.WSMessage{
		Event: websocket.EventPollResultsUpdate,
		Data:  mustMarshalPollJSON(response.UpdatedResults),
	}
	c.WSHub.BroadcastToRoom(roomID, mustMarshalPollJSON(data))
}

// broadcastPollClosed broadcast poll closed event
func (c *PollController) broadcastPollClosed(roomID uint, response *model.ClosePollResponse) {
	if c.WSHub == nil {
		return
	}
	data := websocket.WSMessage{
		Event: websocket.EventPollClosed,
		Data:  mustMarshalPollJSON(response),
	}
	c.WSHub.BroadcastToRoom(roomID, mustMarshalPollJSON(data))
}

// mustMarshalPollJSON helper untuk marshal JSON
func mustMarshalPollJSON(v interface{}) []byte {
	data, err := json.Marshal(v)
	if err != nil {
		panic(err)
	}
	return data
}
