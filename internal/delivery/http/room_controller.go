package http

import (
	"slido-clone-backend/internal/delivery/http/middleware"
	"slido-clone-backend/internal/model"
	"slido-clone-backend/internal/usecase"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

type RoomController struct {
	Log         *logrus.Logger
	RoomUseCase *usecase.RoomUseCase
}

// NewRoomController create new instance of RoomController
func NewRoomController(log *logrus.Logger, roomUseCase *usecase.RoomUseCase) *RoomController {
	return &RoomController{
		Log:         log,
		RoomUseCase: roomUseCase,
	}
}

// Create handler yang digunakan untuk create room baru (call usecase create room)
func (c *RoomController) Create(ctx *fiber.Ctx) error {
	// get user from locals
	auth := middleware.GetUser(ctx)

	// create model create room request
	request := new(model.CreateRoomRequest)

	// parsing
	if err := ctx.BodyParser(request); err != nil {
		c.Log.Warnf("Failed to parse body: %s", err)
		return fiber.ErrBadRequest
	}
	request.PresenterID = *auth.UserID

	// call usecase to create room
	response, err := c.RoomUseCase.Create(ctx.UserContext(), request)
	if err != nil {
		c.Log.Warnf("Failed to create room: %s", err)
		return err
	}

	// return response
	return ctx.Status(fiber.StatusCreated).JSON(model.WebResponse{
		Data: response,
	})
}

// Get handler yang digunakan untuk mencari room berdasarkan room code
func (c *RoomController) Get(ctx *fiber.Ctx) error {
	// create model get room request
	// get room code from params
	// assign value to model request
	request := &model.GetRoomRequestByRoomCode{
		RoomCode: ctx.Params("room_code"),
	}

	// call usecase to get room
	response, err := c.RoomUseCase.Get(ctx.UserContext(), request)
	if err != nil {
		c.Log.Warnf("Failed to get room: %s", err)
		return err
	}

	// return response
	return ctx.Status(fiber.StatusOK).JSON(model.WebResponse{
		Data: response,
	})
}

// UpdateToClosed handler yang digunakan untuk mengupdate status room menjadi closed
func (c *RoomController) UpdateToClosed(ctx *fiber.Ctx) error {
	// get user from locals
	auth := middleware.GetUser(ctx)

	request := new(model.UpdateToCloseRoomRequestByID)

	// parsing body payload
	if err := ctx.BodyParser(request); err != nil {
		c.Log.Warnf("Failed to parse body: %s", err)
		return fiber.ErrBadRequest
	}

	roomId := ctx.Params("room_id")
	idUint64, err := strconv.ParseUint(roomId, 10, 64)
	if err != nil {
		c.Log.Warnf("Invalid room id: %s", err)
		return fiber.ErrBadRequest
	}

	request.PresenterID = *auth.UserID
	request.RoomID = uint(idUint64)

	// call usecase to update room status
	response, err := c.RoomUseCase.UpdateToClosed(ctx.UserContext(), request)
	if err != nil {
		c.Log.Warnf("Failed to update room to closed: %s", err)
		return err
	}

	// return response
	return ctx.Status(fiber.StatusOK).JSON(model.WebResponse{
		Data: response,
	})
}

// Search handler yang digunakan untuk mencari semua room berdasarkan presenter id
func (c *RoomController) Search(ctx *fiber.Ctx) error {
	// get user from locals
	auth := middleware.GetUser(ctx)

	// create model search room request
	request := &model.SearchRoomsRequest{
		PresenterID: *auth.UserID,
	}

	// call usecase to search room
	response, err := c.RoomUseCase.Search(ctx.UserContext(), request)
	if err != nil {
		c.Log.Warnf("Failed to search rooms: %s", err)
		return err
	}

	// return response
	return ctx.Status(fiber.StatusOK).JSON(model.WebResponse{
		Data: response,
	})
}
