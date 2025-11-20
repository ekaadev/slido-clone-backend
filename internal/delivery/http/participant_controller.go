package http

import (
	"math"
	"slido-clone-backend/internal/delivery/http/middleware"
	"slido-clone-backend/internal/model"
	"slido-clone-backend/internal/usecase"
	"strconv"

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

// List handler untuk mendapatkan daftar participant dalam room
func (c *ParticipantController) List(ctx *fiber.Ctx) error {
	// get auth from locals
	auth := middleware.GetUser(ctx)

	// create a model request
	request := &model.ListParticipantsRequest{
		ParticipantID: *auth.ParticipantID,
		Page:          ctx.QueryInt("page", 1),
		Size:          ctx.QueryInt("size", 10),
	}

	roomId := ctx.Params("room_id")
	idUint64, err := strconv.ParseUint(roomId, 10, 64)
	if err != nil {
		c.Log.Warnf("Invalid room id: %s", err)
		return fiber.ErrBadRequest
	}

	request.RoomID = uint(idUint64)

	// call usecase to get list of participants
	responses, total, err := c.ParticipantUseCase.List(ctx.UserContext(), request)
	if err != nil {
		c.Log.Warnf("request failed: %v", err)
		return err
	}

	// calculate pagination
	paging := &model.PaginationResponse{
		Page:      request.Page,
		Size:      request.Size,
		TotalItem: total,
		TotalPage: int64(math.Ceil(float64(total) / float64(request.Size))),
	}

	// return response
	return ctx.Status(fiber.StatusOK).JSON(model.WebResponse{
		Data:   responses,
		Paging: paging,
	})
}
