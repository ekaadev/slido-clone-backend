package usecase

import (
	"context"
	"slido-clone-backend/internal/model"
	"slido-clone-backend/internal/repository"
	"sort"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type ActivityUseCase struct {
	DB                 *gorm.DB
	Log                *logrus.Logger
	Validate           *validator.Validate
	ActivityRepository *repository.ActivityRepository
	RoomRepository     *repository.RoomRepository
}

// NewActivityUseCase create new instance of ActivityUseCase
func NewActivityUseCase(db *gorm.DB, log *logrus.Logger, validate *validator.Validate, activityRepository *repository.ActivityRepository, roomRepository *repository.RoomRepository) *ActivityUseCase {
	return &ActivityUseCase{
		DB:                 db,
		Log:                log,
		Validate:           validate,
		ActivityRepository: activityRepository,
		RoomRepository:     roomRepository,
	}
}

// GetTimeline usecase untuk mendapatkan unified timeline
func (c *ActivityUseCase) GetTimeline(ctx context.Context, request *model.GetTimelineRequest) (*model.GetTimelineResponse, error) {
	tx := c.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	// validate request
	if err := c.Validate.Struct(request); err != nil {
		c.Log.Warnf("GetTimeline - Invalid request: %v", err)
		return nil, fiber.ErrBadRequest
	}

	// set default limit
	if request.Limit == 0 {
		request.Limit = 50
	}

	// check room exists
	roomCount, err := c.RoomRepository.CountById(tx, request.RoomID)
	if err != nil {
		c.Log.Errorf("GetTimeline - RoomRepository.CountById error: %v", err)
		return nil, fiber.ErrInternalServerError
	}
	if roomCount == 0 {
		return nil, fiber.ErrNotFound
	}

	// parse cursors
	var beforeTime, afterTime *time.Time
	if request.Before != "" {
		t, err := time.Parse(time.RFC3339, request.Before)
		if err != nil {
			c.Log.Warnf("GetTimeline - Invalid before cursor: %v", err)
			return nil, fiber.ErrBadRequest
		}
		beforeTime = &t
	}
	if request.After != "" {
		t, err := time.Parse(time.RFC3339, request.After)
		if err != nil {
			c.Log.Warnf("GetTimeline - Invalid after cursor: %v", err)
			return nil, fiber.ErrBadRequest
		}
		afterTime = &t
	}

	// get raw timeline items
	rawItems, err := c.ActivityRepository.GetTimelineRaw(tx, request.RoomID, beforeTime, afterTime, request.Limit)
	if err != nil {
		c.Log.Errorf("GetTimeline - GetTimelineRaw error: %v", err)
		return nil, fiber.ErrInternalServerError
	}

	// check has_more
	hasMore := len(rawItems) > request.Limit
	if hasMore {
		rawItems = rawItems[:request.Limit]
	}

	// jika load newer (after), reverse order agar newest first
	if afterTime != nil {
		for i, j := 0, len(rawItems)-1; i < j; i, j = i+1, j-1 {
			rawItems[i], rawItems[j] = rawItems[j], rawItems[i]
		}
	}

	// collect IDs per type untuk batch fetch
	messageIDs := make([]uint, 0)
	questionIDs := make([]uint, 0)
	pollIDs := make([]uint, 0)

	for _, item := range rawItems {
		switch item.Type {
		case model.ActivityTypeMessage:
			messageIDs = append(messageIDs, item.ID)
		case model.ActivityTypeQuestion:
			questionIDs = append(questionIDs, item.ID)
		case model.ActivityTypePoll:
			pollIDs = append(pollIDs, item.ID)
		}
	}

	// batch fetch data
	messagesMap, _ := c.ActivityRepository.GetMessagesByIDs(tx, messageIDs)
	questionsMap, _ := c.ActivityRepository.GetQuestionsByIDs(tx, questionIDs)
	pollsMap, _ := c.ActivityRepository.GetPollsByIDs(tx, pollIDs)

	// build timeline items with full data
	items := make([]model.TimelineItem, 0, len(rawItems))
	for _, raw := range rawItems {
		item := model.TimelineItem{
			Type:      raw.Type,
			ID:        raw.ID,
			CreatedAt: raw.CreatedAt,
		}

		switch raw.Type {
		case model.ActivityTypeMessage:
			if msg, ok := messagesMap[raw.ID]; ok {
				item.Data = model.MessageTimelineData{
					Content: msg.Content,
					Participant: model.ParticipantInfo{
						ID:          msg.Participant.ID,
						DisplayName: msg.Participant.DisplayName,
					},
				}
			}
		case model.ActivityTypeQuestion:
			if q, ok := questionsMap[raw.ID]; ok {
				item.Data = model.QuestionTimelineData{
					Content: q.Content,
					Participant: model.ParticipantInfo{
						ID:          q.Participant.ID,
						DisplayName: q.Participant.DisplayName,
					},
					UpvoteCount: q.UpvoteCount,
					IsValidated: q.IsValidatedByPresenter,
					Status:      q.Status,
				}
			}
		case model.ActivityTypePoll:
			if p, ok := pollsMap[raw.ID]; ok {
				options := make([]model.PollOptionResponse, len(p.Options))
				totalVotes := 0
				for i, opt := range p.Options {
					options[i] = model.PollOptionResponse{
						ID:         opt.ID,
						OptionText: opt.OptionText,
						VoteCount:  opt.VoteCount,
					}
					totalVotes += opt.VoteCount
				}
				item.Data = model.PollTimelineData{
					Question:   p.Question,
					Status:     p.Status,
					Options:    options,
					TotalVotes: totalVotes,
				}
			}
		}

		items = append(items, item)
	}

	// sort by created_at DESC (newest first)
	sort.Slice(items, func(i, j int) bool {
		return items[i].CreatedAt.After(items[j].CreatedAt)
	})

	if err = tx.Commit().Error; err != nil {
		c.Log.Errorf("GetTimeline - Commit error: %v", err)
		return nil, fiber.ErrInternalServerError
	}

	// build response
	response := &model.GetTimelineResponse{
		Items:   items,
		HasMore: hasMore,
	}

	// set cursors
	if len(items) > 0 {
		oldest := items[len(items)-1].CreatedAt.Format(time.RFC3339)
		newest := items[0].CreatedAt.Format(time.RFC3339)
		response.OldestAt = &oldest
		response.NewestAt = &newest
	}

	return response, nil
}
