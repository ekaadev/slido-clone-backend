package usecase

import (
	"context"
	"slido-clone-backend/internal/entity"
	"slido-clone-backend/internal/model"
	"slido-clone-backend/internal/model/converter"
	"slido-clone-backend/internal/repository"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

// XP points configuration for polls
const (
	XPPollVote = 5 // XP untuk submit vote pada poll
)

// PollUseCase usecase untuk poll operations
type PollUseCase struct {
	DB                      *gorm.DB
	Log                     *logrus.Logger
	Validator               *validator.Validate
	PollRepository          *repository.PollRepository
	RoomRepository          *repository.RoomRepository
	ParticipantRepository   *repository.ParticipantRepository
	XPTransactionRepository *repository.XPTransactionRepository
}

// NewPollUseCase create new instance of PollUseCase
func NewPollUseCase(
	db *gorm.DB,
	log *logrus.Logger,
	validate *validator.Validate,
	pollRepository *repository.PollRepository,
	roomRepository *repository.RoomRepository,
	participantRepository *repository.ParticipantRepository,
	xpTransactionRepository *repository.XPTransactionRepository,
) *PollUseCase {
	return &PollUseCase{
		DB:                      db,
		Log:                     log,
		Validator:               validate,
		PollRepository:          pollRepository,
		RoomRepository:          roomRepository,
		ParticipantRepository:   participantRepository,
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
	var room entity.Room
	if err := c.RoomRepository.FindById(tx, &room, request.RoomID); err != nil {
		if err == gorm.ErrRecordNotFound {
			c.Log.Warnf("Create - Room not found: %d", request.RoomID)
			return nil, fiber.ErrNotFound
		}
		c.Log.Errorf("Create - RoomRepository.FindById error: %v", err)
		return nil, fiber.ErrInternalServerError
	}

	// check if presenter is the owner of the room
	if room.PresenterID != request.PresenterID {
		c.Log.Warnf("Create - User %d is not the presenter of room %d", request.PresenterID, request.RoomID)
		return nil, fiber.ErrForbidden
	}

	// check room is active
	if room.Status != "active" {
		c.Log.Warnf("Create - Room %d is not active", request.RoomID)
		return nil, fiber.NewError(fiber.StatusBadRequest, "Room is not active")
	}

	// create poll entity
	poll := &entity.Poll{
		RoomID:   request.RoomID,
		Question: request.Question,
		Status:   "active", // langsung active saat dibuat
	}

	// create poll options
	options := make([]entity.PollOption, len(request.Options))
	for i, optText := range request.Options {
		options[i] = entity.PollOption{
			OptionText: optText,
			VoteCount:  0,
			Order:      i + 1,
		}
	}

	// save poll with options
	if err := c.PollRepository.CreatePollWithOptions(tx, poll, options); err != nil {
		c.Log.Errorf("Create - CreatePollWithOptions error: %v", err)
		return nil, fiber.ErrInternalServerError
	}

	if err := tx.Commit().Error; err != nil {
		c.Log.Errorf("Create - Commit error: %v", err)
		return nil, fiber.ErrInternalServerError
	}

	return converter.PollToCreateResponse(poll), nil
}

// GetActivePolls usecase untuk mendapatkan active polls di room
func (c *PollUseCase) GetActivePolls(ctx context.Context, request *model.GetActivePollsRequest) (*model.GetActivePollsResponse, error) {
	tx := c.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	// validate request
	if err := c.Validator.Struct(request); err != nil {
		c.Log.Warnf("GetActivePolls - Invalid request: %v", err)
		return nil, fiber.ErrBadRequest
	}

	// check room exists
	roomCount, err := c.RoomRepository.CountById(tx, request.RoomID)
	if err != nil {
		c.Log.Errorf("GetActivePolls - RoomRepository.CountById error: %v", err)
		return nil, fiber.ErrInternalServerError
	}
	if roomCount == 0 {
		return nil, fiber.ErrNotFound
	}

	// get active polls
	polls, err := c.PollRepository.GetActivePollsByRoomID(tx, request.RoomID)
	if err != nil {
		c.Log.Errorf("GetActivePolls - GetActivePollsByRoomID error: %v", err)
		return nil, fiber.ErrInternalServerError
	}

	// build response with participant vote info
	result := make([]model.PollResponse, len(polls))
	for i, poll := range polls {
		// calculate total votes
		totalVotes := 0
		for _, opt := range poll.Options {
			totalVotes += opt.VoteCount
		}

		// check if participant has voted
		existingResponse, err := c.PollRepository.GetPollResponseByParticipant(tx, poll.ID, request.ParticipantID)
		if err != nil {
			c.Log.Errorf("GetActivePolls - GetPollResponseByParticipant error: %v", err)
			return nil, fiber.ErrInternalServerError
		}

		hasVoted := existingResponse != nil
		var myVoteID *uint
		if hasVoted {
			myVoteID = &existingResponse.PollOptionID
		}

		pollResp := converter.PollToResponseWithOptions(&poll, totalVotes, hasVoted, myVoteID)
		result[i] = *pollResp
	}

	if err := tx.Commit().Error; err != nil {
		c.Log.Errorf("GetActivePolls - Commit error: %v", err)
		return nil, fiber.ErrInternalServerError
	}

	return &model.GetActivePollsResponse{
		Polls: result,
	}, nil
}

// GetHistory usecase untuk mendapatkan poll history di room
func (c *PollUseCase) GetHistory(ctx context.Context, request *model.GetPollHistoryRequest) (*model.PollHistoryResponse, error) {
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
	if request.Status == "" {
		request.Status = "all"
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
	polls, total, err := c.PollRepository.GetPollsByRoomID(tx, request.RoomID, request.Status, request.Limit)
	if err != nil {
		c.Log.Errorf("GetHistory - GetPollsByRoomID error: %v", err)
		return nil, fiber.ErrInternalServerError
	}

	if err := tx.Commit().Error; err != nil {
		c.Log.Errorf("GetHistory - Commit error: %v", err)
		return nil, fiber.ErrInternalServerError
	}

	return converter.PollsToHistoryResponse(polls, total), nil
}

// Vote usecase untuk submit vote pada poll
func (c *PollUseCase) Vote(ctx context.Context, request *model.SubmitPollVoteRequest) (*model.SubmitPollVoteResponse, error) {
	tx := c.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	// validate request
	if err := c.Validator.Struct(request); err != nil {
		c.Log.Warnf("Vote - Invalid request: %v", err)
		return nil, fiber.ErrBadRequest
	}

	// get poll with options
	poll, err := c.PollRepository.GetPollByIDWithOptions(tx, request.PollID)
	if err != nil {
		c.Log.Errorf("Vote - GetPollByIDWithOptions error: %v", err)
		return nil, fiber.ErrInternalServerError
	}
	if poll == nil {
		return nil, fiber.ErrNotFound
	}

	// check poll is active
	if poll.Status != "active" {
		c.Log.Warnf("Vote - Poll %d is not active", request.PollID)
		return nil, fiber.NewError(fiber.StatusBadRequest, "Poll is not active")
	}

	// check participant is in the same room as the poll
	if poll.RoomID != request.RoomID {
		c.Log.Warnf("Vote - Participant not in the same room as poll")
		return nil, fiber.ErrForbidden
	}

	// check if option belongs to the poll
	optionExists := false
	for _, opt := range poll.Options {
		if opt.ID == request.OptionID {
			optionExists = true
			break
		}
	}
	if !optionExists {
		c.Log.Warnf("Vote - Option %d does not belong to poll %d", request.OptionID, request.PollID)
		return nil, fiber.NewError(fiber.StatusBadRequest, "Invalid option")
	}

	// check if participant already voted
	existingResponse, err := c.PollRepository.GetPollResponseByParticipant(tx, request.PollID, request.ParticipantID)
	if err != nil {
		c.Log.Errorf("Vote - GetPollResponseByParticipant error: %v", err)
		return nil, fiber.ErrInternalServerError
	}
	if existingResponse != nil {
		c.Log.Warnf("Vote - Participant %d already voted on poll %d", request.ParticipantID, request.PollID)
		return nil, fiber.NewError(fiber.StatusBadRequest, "Already voted")
	}

	// create poll response
	pollResponse := &entity.PollResponse{
		PollID:        request.PollID,
		ParticipantID: request.ParticipantID,
		PollOptionID:  request.OptionID,
	}
	if err := c.PollRepository.CreatePollResponse(tx, pollResponse); err != nil {
		c.Log.Errorf("Vote - CreatePollResponse error: %v", err)
		return nil, fiber.ErrInternalServerError
	}

	// NOTE: vote_count di-increment otomatis oleh database trigger (after_poll_response_insert)
	// Jangan panggil IncrementOptionVoteCount manual karena akan double-count!

	// add XP for voting
	xpTx := &entity.XPTransaction{
		ParticipantID: request.ParticipantID,
		RoomID:        request.RoomID,
		Points:        XPPollVote,
		SourceType:    "poll",
		SourceID:      pollResponse.ID,
	}
	if err := c.XPTransactionRepository.Create(tx, xpTx); err != nil {
		c.Log.Errorf("Vote - XPTransactionRepository.Create error: %v", err)
		return nil, fiber.ErrInternalServerError
	}

	// update participant XP score
	if err := c.XPTransactionRepository.AddXP(tx, request.ParticipantID, XPPollVote); err != nil {
		c.Log.Errorf("Vote - AddXP error: %v", err)
		return nil, fiber.ErrInternalServerError
	}

	// get updated participant XP
	var participant entity.Participant
	if err := tx.Select("xp_score").Where("id = ?", request.ParticipantID).First(&participant).Error; err != nil {
		c.Log.Errorf("Vote - Get participant error: %v", err)
		return nil, fiber.ErrInternalServerError
	}

	// get updated poll with options
	updatedPoll, err := c.PollRepository.GetPollByIDWithOptions(tx, request.PollID)
	if err != nil {
		c.Log.Errorf("Vote - GetPollByIDWithOptions after vote error: %v", err)
		return nil, fiber.ErrInternalServerError
	}

	// calculate total votes
	totalVotes, err := c.PollRepository.GetTotalVotesByPollID(tx, request.PollID)
	if err != nil {
		c.Log.Errorf("Vote - GetTotalVotesByPollID error: %v", err)
		return nil, fiber.ErrInternalServerError
	}

	if err := tx.Commit().Error; err != nil {
		c.Log.Errorf("Vote - Commit error: %v", err)
		return nil, fiber.ErrInternalServerError
	}

	return converter.PollToVoteResponse(pollResponse, updatedPoll, totalVotes, XPPollVote, participant.XPScore), nil
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
	poll, err := c.PollRepository.GetPollByIDWithOptions(tx, request.PollID)
	if err != nil {
		c.Log.Errorf("Close - GetPollByIDWithOptions error: %v", err)
		return nil, fiber.ErrInternalServerError
	}
	if poll == nil {
		return nil, fiber.ErrNotFound
	}

	// get room to check presenter
	var room entity.Room
	if err := c.RoomRepository.FindById(tx, &room, poll.RoomID); err != nil {
		c.Log.Errorf("Close - RoomRepository.FindById error: %v", err)
		return nil, fiber.ErrInternalServerError
	}

	// check if user is presenter
	if room.PresenterID != request.PresenterID {
		c.Log.Warnf("Close - User %d is not the presenter of room %d", request.PresenterID, poll.RoomID)
		return nil, fiber.ErrForbidden
	}

	// check poll is active
	if poll.Status != "active" {
		c.Log.Warnf("Close - Poll %d is already closed", request.PollID)
		return nil, fiber.NewError(fiber.StatusBadRequest, "Poll is already closed")
	}

	// close poll
	if err := c.PollRepository.ClosePoll(tx, poll); err != nil {
		c.Log.Errorf("Close - ClosePoll error: %v", err)
		return nil, fiber.ErrInternalServerError
	}

	// reload poll to get updated status and closed_at
	poll, err = c.PollRepository.GetPollByIDWithOptions(tx, request.PollID)
	if err != nil {
		c.Log.Errorf("Close - GetPollByIDWithOptions after close error: %v", err)
		return nil, fiber.ErrInternalServerError
	}

	// calculate total votes
	totalVotes, err := c.PollRepository.GetTotalVotesByPollID(tx, request.PollID)
	if err != nil {
		c.Log.Errorf("Close - GetTotalVotesByPollID error: %v", err)
		return nil, fiber.ErrInternalServerError
	}

	if err := tx.Commit().Error; err != nil {
		c.Log.Errorf("Close - Commit error: %v", err)
		return nil, fiber.ErrInternalServerError
	}

	return converter.PollToCloseResponse(poll, totalVotes), nil
}
