package usecase

import (
	"context"
	"slido-clone-backend/internal/entity"
	"slido-clone-backend/internal/model"
	"slido-clone-backend/internal/model/converter"
	"slido-clone-backend/internal/repository"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

// XP points configuration for polling
const (
	XPSubmitVote = 5 // XP untuk submit vote pada poll
)

// PollUseCase usecase untuk Polling operations
type PollUseCase struct {
	DB                      *gorm.DB
	Log                     *logrus.Logger
	Validator               *validator.Validate
	PollRepository          *repository.PollRepository
	PollOptionRepository    *repository.PollOptionRepository
	PollResponseRepository  *repository.PollResponseRepository
	RoomRepository          *repository.RoomRepository
	XPTransactionRepository *repository.XPTransactionRepository
}

// NewPollUseCase create new instance of PollUseCase
func NewPollUseCase(
	db *gorm.DB,
	log *logrus.Logger,
	validate *validator.Validate,
	pollRepository *repository.PollRepository,
	pollOptionRepository *repository.PollOptionRepository,
	pollResponseRepository *repository.PollResponseRepository,
	roomRepository *repository.RoomRepository,
	xpTransactionRepository *repository.XPTransactionRepository,
) *PollUseCase {
	return &PollUseCase{
		DB:                      db,
		Log:                     log,
		Validator:               validate,
		PollRepository:          pollRepository,
		PollOptionRepository:    pollOptionRepository,
		PollResponseRepository:  pollResponseRepository,
		RoomRepository:          roomRepository,
		XPTransactionRepository: xpTransactionRepository,
	}
}

// Create usecase untuk membuat poll baru (presenter only)
func (c *PollUseCase) Create(ctx context.Context, request *model.CreatePollRequest) (*model.CreatePollResponse, error) {
	tx := c.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	// validate request
	if err := c.Validator.Struct(request); err != nil {
		c.Log.Warnf("Create - Invalid request: %v", err)
		return nil, fiber.ErrBadRequest
	}

	// check room exists
	roomCount, err := c.RoomRepository.CountById(tx, request.RoomID)
	if err != nil {
		c.Log.Errorf("Create - RoomRepository.CountById error: %v", err)
		return nil, fiber.ErrInternalServerError
	}
	if roomCount == 0 {
		return nil, fiber.ErrNotFound
	}

	// create poll entity
	poll := &entity.Poll{
		RoomID:   request.RoomID,
		Question: request.Question,
		Status:   "active", // langsung active saat dibuat
	}

	// save poll
	if err = c.PollRepository.Create(tx, poll); err != nil {
		c.Log.Errorf("Create - PollRepository.Create error: %v", err)
		return nil, fiber.ErrInternalServerError
	}

	// create poll options
	options := make([]entity.PollOption, len(request.Options))
	for i, optText := range request.Options {
		options[i] = entity.PollOption{
			PollID:     poll.ID,
			OptionText: optText,
			VoteCount:  0,
			Order:      i + 1,
		}
	}

	if err = c.PollOptionRepository.CreateBatch(tx, options); err != nil {
		c.Log.Errorf("Create - PollOptionRepository.CreateBatch error: %v", err)
		return nil, fiber.ErrInternalServerError
	}

	// activate poll
	if err = c.PollRepository.Activate(tx, poll.ID); err != nil {
		c.Log.Errorf("Create - PollRepository.Activate error: %v", err)
		return nil, fiber.ErrInternalServerError
	}

	if err = tx.Commit().Error; err != nil {
		c.Log.Errorf("Create - Commit error: %v", err)
		return nil, fiber.ErrInternalServerError
	}

	// reload poll with options
	poll.Options = options

	return converter.PollToCreateResponse(poll), nil
}

// GetActive usecase untuk mendapatkan active polls di room
func (c *PollUseCase) GetActive(ctx context.Context, request *model.GetActivePollRequest) (*model.GetActivePollResponse, error) {
	tx := c.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	// validate request
	if err := c.Validator.Struct(request); err != nil {
		c.Log.Warnf("GetActive - Invalid request: %v", err)
		return nil, fiber.ErrBadRequest
	}

	// check room exists
	roomCount, err := c.RoomRepository.CountById(tx, request.RoomID)
	if err != nil {
		c.Log.Errorf("GetActive - RoomRepository.CountById error: %v", err)
		return nil, fiber.ErrInternalServerError
	}
	if roomCount == 0 {
		return nil, fiber.ErrNotFound
	}

	// get active polls
	polls, err := c.PollRepository.FindActiveByRoomID(tx, request.RoomID)
	if err != nil {
		c.Log.Errorf("GetActive - PollRepository.FindActiveByRoomID error: %v", err)
		return nil, fiber.ErrInternalServerError
	}

	// build response with voting info
	pollResponses := make([]model.PollResponse, len(polls))
	for i, poll := range polls {
		// get total votes
		totalVotes, err := c.PollOptionRepository.GetTotalVotesByPollID(tx, poll.ID)
		if err != nil {
			c.Log.Warnf("GetActive - GetTotalVotesByPollID error: %v", err)
		}

		// check if participant has voted
		hasVoted, err := c.PollResponseRepository.HasVoted(tx, poll.ID, request.ParticipantID)
		if err != nil {
			c.Log.Warnf("GetActive - HasVoted error: %v", err)
		}

		// get voted option
		var votedOption *uint
		if hasVoted {
			votedOption, _ = c.PollResponseRepository.GetVotedOptionByParticipant(tx, poll.ID, request.ParticipantID)
		}

		pollResponses[i] = converter.PollToResponseWithDetails(&poll, totalVotes, hasVoted, votedOption)
	}

	if err = tx.Commit().Error; err != nil {
		c.Log.Errorf("GetActive - Commit error: %v", err)
		return nil, fiber.ErrInternalServerError
	}

	return &model.GetActivePollResponse{
		Polls: pollResponses,
	}, nil
}

// SubmitVote usecase untuk submit vote pada poll
func (c *PollUseCase) SubmitVote(ctx context.Context, request *model.SubmitVoteRequest) (*model.SubmitVoteResponse, error) {
	tx := c.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	// validate request
	if err := c.Validator.Struct(request); err != nil {
		c.Log.Warnf("SubmitVote - Invalid request: %v", err)
		return nil, fiber.ErrBadRequest
	}

	// get poll
	poll, err := c.PollRepository.FindByIdWithOptions(tx, request.PollID)
	if err != nil {
		c.Log.Errorf("SubmitVote - PollRepository.FindByIdWithOptions error: %v", err)
		return nil, fiber.ErrInternalServerError
	}
	if poll == nil {
		return nil, fiber.ErrNotFound
	}

	// check poll is active
	if poll.Status != "active" {
		return nil, fiber.NewError(fiber.StatusBadRequest, "Poll is not active")
	}

	// check if already voted
	hasVoted, err := c.PollResponseRepository.HasVoted(tx, request.PollID, request.ParticipantID)
	if err != nil {
		c.Log.Errorf("SubmitVote - PollResponseRepository.HasVoted error: %v", err)
		return nil, fiber.ErrInternalServerError
	}
	if hasVoted {
		return nil, fiber.NewError(fiber.StatusConflict, "Already voted on this poll")
	}

	// validate option belongs to poll
	validOption, err := c.PollOptionRepository.ValidateOptionBelongsToPoll(tx, request.OptionID, request.PollID)
	if err != nil {
		c.Log.Errorf("SubmitVote - ValidateOptionBelongsToPoll error: %v", err)
		return nil, fiber.ErrInternalServerError
	}
	if !validOption {
		return nil, fiber.NewError(fiber.StatusBadRequest, "Invalid option for this poll")
	}

	// create poll response
	pollResponse := &entity.PollResponse{
		PollID:        request.PollID,
		ParticipantID: request.ParticipantID,
		PollOptionID:  request.OptionID,
	}

	if err = c.PollResponseRepository.Create(tx, pollResponse); err != nil {
		c.Log.Errorf("SubmitVote - PollResponseRepository.Create error: %v", err)
		return nil, fiber.ErrInternalServerError
	}

	// increment vote count
	if err = c.PollOptionRepository.IncrementVoteCount(tx, request.OptionID); err != nil {
		c.Log.Errorf("SubmitVote - PollOptionRepository.IncrementVoteCount error: %v", err)
		return nil, fiber.ErrInternalServerError
	}

	// add XP for voting
	xpTx := &entity.XPTransaction{
		ParticipantID: request.ParticipantID,
		RoomID:        poll.RoomID,
		Points:        XPSubmitVote,
		SourceType:    "poll_voted",
		SourceID:      pollResponse.ID,
	}
	if err = c.XPTransactionRepository.Create(tx, xpTx); err != nil {
		c.Log.Errorf("SubmitVote - XPTransactionRepository.Create error: %v", err)
		return nil, fiber.ErrInternalServerError
	}

	// update participant XP score
	if err = c.XPTransactionRepository.AddXP(tx, request.ParticipantID, XPSubmitVote); err != nil {
		c.Log.Errorf("SubmitVote - AddXP error: %v", err)
		return nil, fiber.ErrInternalServerError
	}

	// get updated options
	options, err := c.PollOptionRepository.GetByPollID(tx, request.PollID)
	if err != nil {
		c.Log.Errorf("SubmitVote - GetByPollID error: %v", err)
		return nil, fiber.ErrInternalServerError
	}

	// get total votes
	totalVotes, err := c.PollOptionRepository.GetTotalVotesByPollID(tx, request.PollID)
	if err != nil {
		c.Log.Errorf("SubmitVote - GetTotalVotesByPollID error: %v", err)
		return nil, fiber.ErrInternalServerError
	}

	// get new total XP
	var participant entity.Participant
	if err = tx.Select("xp_score").Where("id = ?", request.ParticipantID).First(&participant).Error; err != nil {
		c.Log.Errorf("SubmitVote - Get participant error: %v", err)
		return nil, fiber.ErrInternalServerError
	}

	if err = tx.Commit().Error; err != nil {
		c.Log.Errorf("SubmitVote - Commit error: %v", err)
		return nil, fiber.ErrInternalServerError
	}

	return &model.SubmitVoteResponse{
		Response:       converter.PollResponseToVoteData(pollResponse),
		UpdatedResults: converter.PollOptionsToUpdatedResults(request.PollID, options, totalVotes),
		XPEarned: &model.XPEarned{
			Points:   XPSubmitVote,
			NewTotal: participant.XPScore,
		},
	}, nil
}

// Close usecase untuk menutup poll (presenter only)
func (c *PollUseCase) Close(ctx context.Context, request *model.ClosePollRequest) (*model.ClosePollResponse, error) {
	tx := c.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	// validate request
	if err := c.Validator.Struct(request); err != nil {
		c.Log.Warnf("Close - Invalid request: %v", err)
		return nil, fiber.ErrBadRequest
	}

	// get poll with options
	poll, err := c.PollRepository.FindByIdWithOptions(tx, request.PollID)
	if err != nil {
		c.Log.Errorf("Close - PollRepository.FindByIdWithOptions error: %v", err)
		return nil, fiber.ErrInternalServerError
	}
	if poll == nil {
		return nil, fiber.ErrNotFound
	}

	// check poll is active
	if poll.Status != "active" {
		return nil, fiber.NewError(fiber.StatusBadRequest, "Poll is not active")
	}

	// close poll
	if err = c.PollRepository.Close(tx, request.PollID); err != nil {
		c.Log.Errorf("Close - PollRepository.Close error: %v", err)
		return nil, fiber.ErrInternalServerError
	}

	// get total votes
	totalVotes, err := c.PollOptionRepository.GetTotalVotesByPollID(tx, request.PollID)
	if err != nil {
		c.Log.Errorf("Close - GetTotalVotesByPollID error: %v", err)
		return nil, fiber.ErrInternalServerError
	}

	// get updated options
	options, err := c.PollOptionRepository.GetByPollID(tx, request.PollID)
	if err != nil {
		c.Log.Errorf("Close - GetByPollID error: %v", err)
		return nil, fiber.ErrInternalServerError
	}

	if err = tx.Commit().Error; err != nil {
		c.Log.Errorf("Close - Commit error: %v", err)
		return nil, fiber.ErrInternalServerError
	}

	// build response
	optResponses := make([]model.PollOptionResponse, len(options))
	for i, opt := range options {
		optResponses[i] = converter.PollOptionToResponseWithPercentage(&opt, totalVotes)
	}

	return &model.ClosePollResponse{
		Poll: struct {
			ID           uint                       `json:"id"`
			Status       string                     `json:"status"`
			ClosedAt     *time.Time                 `json:"closed_at"`
			FinalResults model.FinalResultsResponse `json:"final_results"`
		}{
			ID:       poll.ID,
			Status:   "closed",
			ClosedAt: poll.ClosedAt,
			FinalResults: model.FinalResultsResponse{
				TotalVotes: totalVotes,
				Options:    optResponses,
			},
		},
	}, nil
}

// GetHistory usecase untuk mendapatkan poll history
func (c *PollUseCase) GetHistory(ctx context.Context, request *model.GetPollHistoryRequest) (*model.GetPollHistoryResponse, error) {
	tx := c.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	// validate request
	if err := c.Validator.Struct(request); err != nil {
		c.Log.Warnf("GetHistory - Invalid request: %v", err)
		return nil, fiber.ErrBadRequest
	}

	// set defaults
	if request.Limit == 0 {
		request.Limit = 10
	}

	// check room exists
	roomCount, err := c.RoomRepository.CountById(tx, request.RoomID)
	if err != nil {
		c.Log.Errorf("GetHistory - RoomRepository.CountById error: %v", err)
		return nil, fiber.ErrInternalServerError
	}
	if roomCount == 0 {
		return nil, fiber.ErrNotFound
	}

	// get polls
	polls, err := c.PollRepository.ListByRoomID(tx, request.RoomID, request.Status, request.Limit)
	if err != nil {
		c.Log.Errorf("GetHistory - PollRepository.ListByRoomID error: %v", err)
		return nil, fiber.ErrInternalServerError
	}

	// get total count
	total, err := c.PollRepository.CountByRoomID(tx, request.RoomID, request.Status)
	if err != nil {
		c.Log.Errorf("GetHistory - PollRepository.CountByRoomID error: %v", err)
		return nil, fiber.ErrInternalServerError
	}

	// build response
	pollItems := make([]model.PollHistoryItem, len(polls))
	for i, poll := range polls {
		totalVotes, _ := c.PollOptionRepository.GetTotalVotesByPollID(tx, poll.ID)
		pollItems[i] = converter.PollToHistoryItem(&poll, totalVotes)
	}

	if err = tx.Commit().Error; err != nil {
		c.Log.Errorf("GetHistory - Commit error: %v", err)
		return nil, fiber.ErrInternalServerError
	}

	return &model.GetPollHistoryResponse{
		Polls: pollItems,
		Total: total,
	}, nil
}
