package http

import (
	"slido-clone-backend/internal/delivery/http/middleware"
	"slido-clone-backend/internal/model"
	"slido-clone-backend/internal/usecase"

	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

type ParticipantController struct {
	Log                *logrus.Logger
	ParticipantUseCase *usecase.ParticipantUseCase
}

func NewParticipantController(log *logrus.Logger, participantUseCase *usecase.ParticipantUseCase) *ParticipantController {
	return &ParticipantController{
		Log:                log,
		ParticipantUseCase: participantUseCase,
	}
}

// Join handler untuk participant bergabung ke dalam room
func (c *ParticipantController) Join(ctx *fiber.Ctx) error {
	// get auth from locals
	auth := middleware.GetUser(ctx)

	// create model request join room
	request := &model.JoinRoomRequest{
		Username: auth.Username,
		RoomCode: ctx.Params("room_code"),
	}

	// parse body request
	if err := ctx.BodyParser(request); err != nil {
		c.Log.Warnf("request body parse failed: %v", err)
		return fiber.ErrBadRequest
	}

	// call usecase join room
	response, err := c.ParticipantUseCase.Join(ctx.UserContext(), request)
	if err != nil {
		c.Log.Warnf("request failed: %v", err)
		return err
	}

	// return response
	return ctx.Status(fiber.StatusCreated).JSON(model.WebResponse{
		Data: response,
	})
}
