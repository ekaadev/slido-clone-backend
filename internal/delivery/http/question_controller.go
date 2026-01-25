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

// QuestionController controller untuk Q&A operations
type QuestionController struct {
	Log             *logrus.Logger
	QuestionUseCase *usecase.QuestionUseCase
	WSHub           *websocket.Hub
}

// NewQuestionController create new instance of QuestionController
func NewQuestionController(log *logrus.Logger, questionUseCase *usecase.QuestionUseCase, wsHub *websocket.Hub) *QuestionController {
	return &QuestionController{
		Log:             log,
		QuestionUseCase: questionUseCase,
		WSHub:           wsHub,
	}
}

// Submit handler untuk submit question baru
func (c *QuestionController) Submit(ctx *fiber.Ctx) error {
	auth := middleware.GetUser(ctx)

	// parse room_id from params
	roomIDStr := ctx.Params("room_id")
	roomIDUint64, err := strconv.ParseUint(roomIDStr, 10, 64)
	if err != nil {
		c.Log.Warnf("Submit - Invalid room_id: %v", err)
		return fiber.ErrBadRequest
	}

	// create request
	request := &model.SubmitQuestionRequest{
		RoomID:        uint(roomIDUint64),
		ParticipantID: *auth.ParticipantID,
	}

	// parse body
	if err = ctx.BodyParser(request); err != nil {
		c.Log.Warnf("Submit - BodyParser error: %v", err)
		return fiber.ErrBadRequest
	}

	// call usecase
	response, err := c.QuestionUseCase.Submit(ctx.UserContext(), request)
	if err != nil {
		c.Log.Warnf("Submit - QuestionUseCase.Submit error: %v", err)
		return err
	}

	// broadcast ke websocket clients di room
	c.broadcastQuestionCreated(request.RoomID, response)

	return ctx.Status(fiber.StatusCreated).JSON(model.WebResponse{
		Data: response,
	})
}

// List handler untuk mendapatkan list questions
func (c *QuestionController) List(ctx *fiber.Ctx) error {
	auth := middleware.GetUser(ctx)

	// parse room_id from params
	roomIDStr := ctx.Params("room_id")
	roomIDUint64, err := strconv.ParseUint(roomIDStr, 10, 64)
	if err != nil {
		c.Log.Warnf("List - Invalid room_id: %v", err)
		return fiber.ErrBadRequest
	}

	// parse query params
	request := &model.GetQuestionsRequest{
		RoomID:        uint(roomIDUint64),
		ParticipantID: *auth.ParticipantID,
		Status:        ctx.Query("status"),
		SortBy:        ctx.Query("sort_by", "upvotes"),
		Limit:         ctx.QueryInt("limit", 20),
		Offset:        ctx.QueryInt("offset", 0),
	}

	// call usecase
	response, err := c.QuestionUseCase.List(ctx.UserContext(), request)
	if err != nil {
		c.Log.Warnf("List - QuestionUseCase.List error: %v", err)
		return err
	}

	return ctx.Status(fiber.StatusOK).JSON(model.WebResponse{
		Data: response,
	})
}

// Upvote handler untuk upvote question
func (c *QuestionController) Upvote(ctx *fiber.Ctx) error {
	auth := middleware.GetUser(ctx)

	// parse question_id from params
	questionIDStr := ctx.Params("question_id")
	questionIDUint64, err := strconv.ParseUint(questionIDStr, 10, 64)
	if err != nil {
		c.Log.Warnf("Upvote - Invalid question_id: %v", err)
		return fiber.ErrBadRequest
	}

	// create request
	request := &model.UpvoteRequest{
		QuestionID:    uint(questionIDUint64),
		ParticipantID: *auth.ParticipantID,
		RoomID:        *auth.RoomID,
	}

	// call usecase
	response, err := c.QuestionUseCase.Upvote(ctx.UserContext(), request)
	if err != nil {
		c.Log.Warnf("Upvote - QuestionUseCase.Upvote error: %v", err)
		return err
	}

	// broadcast ke websocket clients di room
	c.broadcastQuestionUpvoted(request.RoomID, response)

	return ctx.Status(fiber.StatusOK).JSON(model.WebResponse{
		Data: response,
	})
}

// RemoveUpvote handler untuk remove upvote
func (c *QuestionController) RemoveUpvote(ctx *fiber.Ctx) error {
	auth := middleware.GetUser(ctx)

	// parse question_id from params
	questionIDStr := ctx.Params("question_id")
	questionIDUint64, err := strconv.ParseUint(questionIDStr, 10, 64)
	if err != nil {
		c.Log.Warnf("RemoveUpvote - Invalid question_id: %v", err)
		return fiber.ErrBadRequest
	}

	// create request
	request := &model.UpvoteRequest{
		QuestionID:    uint(questionIDUint64),
		ParticipantID: *auth.ParticipantID,
		RoomID:        *auth.RoomID,
	}

	// call usecase
	response, err := c.QuestionUseCase.RemoveUpvote(ctx.UserContext(), request)
	if err != nil {
		c.Log.Warnf("RemoveUpvote - QuestionUseCase.RemoveUpvote error: %v", err)
		return err
	}

	// broadcast upvote update (menggunakan event yang sama dengan upvote)
	c.broadcastQuestionUpvoteRemoved(request.RoomID, response)

	return ctx.Status(fiber.StatusOK).JSON(model.WebResponse{
		Data: response,
	})
}

// Validate handler untuk validate question (presenter only)
func (c *QuestionController) Validate(ctx *fiber.Ctx) error {
	auth := middleware.GetUser(ctx)

	// parse question_id from params
	questionIDStr := ctx.Params("question_id")
	questionIDUint64, err := strconv.ParseUint(questionIDStr, 10, 64)
	if err != nil {
		c.Log.Warnf("Validate - Invalid question_id: %v", err)
		return fiber.ErrBadRequest
	}

	// get room_id from question (need to get from usecase)
	roomID, err := c.QuestionUseCase.QuestionRepository.GetRoomIDByQuestionID(c.QuestionUseCase.DB, uint(questionIDUint64))
	if err != nil {
		c.Log.Warnf("Validate - GetRoomIDByQuestionID error: %v", err)
		return fiber.ErrInternalServerError
	}

	// create request
	request := &model.ValidateQuestionRequest{
		QuestionID:  uint(questionIDUint64),
		PresenterID: *auth.UserID,
	}

	// parse body
	if err = ctx.BodyParser(request); err != nil {
		c.Log.Warnf("Validate - BodyParser error: %v", err)
		return fiber.ErrBadRequest
	}

	// call usecase
	response, err := c.QuestionUseCase.Validate(ctx.UserContext(), request)
	if err != nil {
		c.Log.Warnf("Validate - QuestionUseCase.Validate error: %v", err)
		return err
	}

	// broadcast ke websocket clients di room
	c.broadcastQuestionValidated(roomID, response)

	return ctx.Status(fiber.StatusOK).JSON(model.WebResponse{
		Data: response,
	})
}

// broadcastQuestionCreated broadcast question created event ke semua clients di room
func (c *QuestionController) broadcastQuestionCreated(roomID uint, response *model.SubmitQuestionResponse) {
	if c.WSHub == nil {
		return
	}
	data := websocket.WSMessage{
		Event: websocket.EventQuestionCreated,
		Data:  mustMarshalJSON(response),
	}
	c.WSHub.BroadcastToRoom(roomID, mustMarshalJSON(data))
}

// broadcastQuestionUpvoted broadcast question upvoted event ke semua clients di room
func (c *QuestionController) broadcastQuestionUpvoted(roomID uint, response *model.UpvoteResponse) {
	if c.WSHub == nil {
		return
	}
	data := websocket.WSMessage{
		Event: websocket.EventQuestionUpvoted,
		Data:  mustMarshalJSON(response),
	}
	c.WSHub.BroadcastToRoom(roomID, mustMarshalJSON(data))
}

// broadcastQuestionUpvoteRemoved broadcast upvote removed event ke semua clients di room
func (c *QuestionController) broadcastQuestionUpvoteRemoved(roomID uint, response *model.RemoveUpvoteResponse) {
	if c.WSHub == nil {
		return
	}
	// menggunakan event yang sama seperti upvoted untuk update count
	data := websocket.WSMessage{
		Event: websocket.EventQuestionUpvoted,
		Data:  mustMarshalJSON(response),
	}
	c.WSHub.BroadcastToRoom(roomID, mustMarshalJSON(data))
}

// broadcastQuestionValidated broadcast question validated event ke semua clients di room
func (c *QuestionController) broadcastQuestionValidated(roomID uint, response *model.ValidateQuestionResponse) {
	if c.WSHub == nil {
		return
	}
	data := websocket.WSMessage{
		Event: websocket.EventQuestionValidated,
		Data:  mustMarshalJSON(response),
	}
	c.WSHub.BroadcastToRoom(roomID, mustMarshalJSON(data))
}

// mustMarshalJSON helper untuk marshal JSON
func mustMarshalJSON(v interface{}) []byte {
	data, err := json.Marshal(v)
	if err != nil {
		panic(err)
	}
	return data
}
