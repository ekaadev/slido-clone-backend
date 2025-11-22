package http

import (
	"slido-clone-backend/internal/delivery/http/middleware"
	"slido-clone-backend/internal/model"
	"slido-clone-backend/internal/usecase"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

type MessageController struct {
	Log            *logrus.Logger
	MessageUseCase *usecase.MessageUseCase
}

func NewMessageController(log *logrus.Logger, messageUseCase *usecase.MessageUseCase) *MessageController {
	return &MessageController{
		Log:            log,
		MessageUseCase: messageUseCase,
	}
}

func (c *MessageController) Send(ctx *fiber.Ctx) error {
	// get auth
	auth := middleware.GetUser(ctx)

	// parse from query param
	roomIDStr := ctx.Params("room_id")
	roomIDUint64, err := strconv.ParseUint(roomIDStr, 10, 64)
	if err != nil {
		c.Log.Warnf("Failed to parse to uint: %v", err)
		return fiber.ErrBadRequest
	}

	// create model request
	request := &model.SendMessageRequest{
		RoomID:        uint(roomIDUint64),
		ParticipantID: *auth.ParticipantID,
	}

	// parse body
	if err = ctx.BodyParser(request); err != nil {
		c.Log.Warnf("Failed to parse : %v", err)
		return fiber.ErrBadRequest
	}

	// call usecase
	response, err := c.MessageUseCase.Send(ctx.UserContext(), request)
	if err != nil {
		c.Log.Warnf("Failed to send message: %v", err)
		return err
	}

	// return response
	return ctx.Status(fiber.StatusCreated).JSON(model.WebResponse{
		Data: response,
	})
}

// List handler untuk mendapatkan list messages
func (c *MessageController) List(ctx *fiber.Ctx) error {
	// get auth from context
	auth := middleware.GetUser(ctx)

	// parse room_id from params
	roomIDStr := ctx.Params("room_id")
	roomIDUint64, err := strconv.ParseUint(roomIDStr, 10, 64)
	if err != nil {
		c.Log.Warnf("Invalid room_id: %s", err)
		return fiber.ErrBadRequest
	}

	// parse query params
	limit := ctx.QueryInt("limit", 50)
	var before *int64
	if beforeStr := ctx.Query("before"); beforeStr != "" {
		beforeInt64, err := strconv.ParseInt(beforeStr, 10, 64)
		if err == nil {
			before = &beforeInt64
		}
	}

	// create request
	request := &model.GetMessagesRequest{
		RoomID:        uint(roomIDUint64),
		ParticipantID: *auth.ParticipantID,
		Limit:         limit,
		Before:        before,
	}

	// call usecase
	response, err := c.MessageUseCase.List(ctx.UserContext(), request)
	if err != nil {
		c.Log.Warnf("Failed to list messages: %s", err)
		return err
	}

	// return response
	return ctx.Status(fiber.StatusOK).JSON(model.WebResponse{
		Data: response,
	})
}
