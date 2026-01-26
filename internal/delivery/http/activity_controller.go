package http

import (
	"slido-clone-backend/internal/delivery/http/middleware"
	"slido-clone-backend/internal/model"
	"slido-clone-backend/internal/usecase"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

// ActivityController controller untuk timeline/activity operations
type ActivityController struct {
	Log             *logrus.Logger
	ActivityUseCase *usecase.ActivityUseCase
}

// NewActivityController create new instance of ActivityController
func NewActivityController(log *logrus.Logger, activityUseCase *usecase.ActivityUseCase) *ActivityController {
	return &ActivityController{
		Log:             log,
		ActivityUseCase: activityUseCase,
	}
}

// GetTimeline handler untuk mendapatkan unified timeline
// Query params: before, after (RFC3339 timestamp), limit
func (c *ActivityController) GetTimeline(ctx *fiber.Ctx) error {
	auth := middleware.GetUser(ctx)

	// parse room_id from params
	roomIDStr := ctx.Params("room_id")
	roomIDUint64, err := strconv.ParseUint(roomIDStr, 10, 64)
	if err != nil {
		c.Log.Warnf("GetTimeline - Invalid room_id: %v", err)
		return fiber.ErrBadRequest
	}

	// create request dengan query params
	request := &model.GetTimelineRequest{
		RoomID:        uint(roomIDUint64),
		ParticipantID: 0,
		Before:        ctx.Query("before"),
		After:         ctx.Query("after"),
		Limit:         ctx.QueryInt("limit", 50),
	}

	// set participant_id jika ada
	if auth.ParticipantID != nil {
		request.ParticipantID = *auth.ParticipantID
	}

	// call usecase
	response, err := c.ActivityUseCase.GetTimeline(ctx.UserContext(), request)
	if err != nil {
		c.Log.Warnf("GetTimeline - ActivityUseCase.GetTimeline error: %v", err)
		return err
	}

	return ctx.Status(fiber.StatusOK).JSON(model.WebResponse{
		Data: response,
	})
}
