package http

import (
	"slido-clone-backend/internal/delivery/http/middleware"
	"slido-clone-backend/internal/model"
	"slido-clone-backend/internal/usecase"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

// XPTransactionController controller untuk XP Transaction operations
type XPTransactionController struct {
	Log                  *logrus.Logger
	XPTransactionUseCase *usecase.XPTransactionUseCase
}

// NewXPTransactionController create new instance of XPTransactionController
func NewXPTransactionController(log *logrus.Logger, xpTransactionUseCase *usecase.XPTransactionUseCase) *XPTransactionController {
	return &XPTransactionController{
		Log:                  log,
		XPTransactionUseCase: xpTransactionUseCase,
	}
}

// GetTransactions handler untuk mendapatkan XP transactions history
func (c *XPTransactionController) GetTransactions(ctx *fiber.Ctx) error {
	auth := middleware.GetUser(ctx)

	// parse room_id from params
	roomIDStr := ctx.Params("room_id")
	roomIDUint64, err := strconv.ParseUint(roomIDStr, 10, 64)
	if err != nil {
		c.Log.Warnf("GetTransactions - Invalid room_id: %v", err)
		return fiber.ErrBadRequest
	}

	// create request
	request := &model.GetXPTransactionsRequest{
		RoomID:        uint(roomIDUint64),
		ParticipantID: *auth.ParticipantID,
		Limit:         ctx.QueryInt("limit", 50),
	}

	// call usecase
	response, err := c.XPTransactionUseCase.GetTransactions(ctx.UserContext(), request)
	if err != nil {
		c.Log.Warnf("GetTransactions - XPTransactionUseCase.GetTransactions error: %v", err)
		return err
	}

	return ctx.Status(fiber.StatusOK).JSON(model.WebResponse{
		Data: response,
	})
}
