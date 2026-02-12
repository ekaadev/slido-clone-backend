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

// XP points configuration
const (
	XPSubmitQuestion    = 10 // XP untuk submit question
	XPReceiveUpvote     = 3  // XP untuk menerima upvote
	XPPresenterValidate = 25 // XP untuk validasi presenter
)

// QuestionUseCase usecase untuk Q&A operations
type QuestionUseCase struct {
	DB                      *gorm.DB
	Log                     *logrus.Logger
	Validator               *validator.Validate
	QuestionRepository      *repository.QuestionRepository
	VoteRepository          *repository.VoteRepository
	RoomRepository          *repository.RoomRepository
	ParticipantRepository   *repository.ParticipantRepository
	XPTransactionRepository *repository.XPTransactionRepository
}

// NewQuestionUseCase create new instance of QuestionUseCase
func NewQuestionUseCase(
	db *gorm.DB,
	log *logrus.Logger,
	validate *validator.Validate,
	questionRepository *repository.QuestionRepository,
	voteRepository *repository.VoteRepository,
	roomRepository *repository.RoomRepository,
	participantRepository *repository.ParticipantRepository,
	xpTransactionRepository *repository.XPTransactionRepository,
) *QuestionUseCase {
	return &QuestionUseCase{
		DB:                      db,
		Log:                     log,
		Validator:               validate,
		QuestionRepository:      questionRepository,
		VoteRepository:          voteRepository,
		RoomRepository:          roomRepository,
		ParticipantRepository:   participantRepository,
		XPTransactionRepository: xpTransactionRepository,
	}
}

// Submit usecase untuk submit question baru
func (c *QuestionUseCase) Submit(ctx context.Context, request *model.SubmitQuestionRequest) (*model.SubmitQuestionResponse, error) {
	tx := c.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	// validate request
	if err := c.Validator.Struct(request); err != nil {
		c.Log.Warnf("Submit - Invalid request: %v", err)
		return nil, fiber.ErrBadRequest
	}

	// check room exists and active
	roomCount, err := c.RoomRepository.CountById(tx, request.RoomID)
	if err != nil {
		c.Log.Errorf("Submit - RoomRepository.CountById error: %v", err)
		return nil, fiber.ErrInternalServerError
	}
	if roomCount == 0 {
		return nil, fiber.ErrNotFound
	}

	// create question entity
	question := &entity.Question{
		RoomID:        request.RoomID,
		ParticipantID: request.ParticipantID,
		Content:       request.Content,
		XPAwarded:     XPSubmitQuestion,
	}

	// save question
	if err = c.QuestionRepository.Create(tx, question); err != nil {
		c.Log.Errorf("Submit - QuestionRepository.Create error: %v", err)
		return nil, fiber.ErrInternalServerError
	}

	// add XP for submitting question
	xpTx := &entity.XPTransaction{
		ParticipantID: request.ParticipantID,
		RoomID:        request.RoomID,
		Points:        XPSubmitQuestion,
		SourceType:    "question_created",
		SourceID:      question.ID,
	}
	if err = c.XPTransactionRepository.Create(tx, xpTx); err != nil {
		c.Log.Errorf("Submit - XPTransactionRepository.Create error: %v", err)
		return nil, fiber.ErrInternalServerError
	}

	// update participant XP score
	if err = c.XPTransactionRepository.AddXP(tx, request.ParticipantID, XPSubmitQuestion); err != nil {
		c.Log.Errorf("Submit - AddXP error: %v", err)
		return nil, fiber.ErrInternalServerError
	}

	// get new total XP
	var participant entity.Participant
	if err = tx.Select("xp_score").Where("id = ?", request.ParticipantID).First(&participant).Error; err != nil {
		c.Log.Errorf("Submit - Get participant error: %v", err)
		return nil, fiber.ErrInternalServerError
	}

	if err = tx.Commit().Error; err != nil {
		c.Log.Errorf("Submit - Commit error: %v", err)
		return nil, fiber.ErrInternalServerError
	}

	return converter.QuestionToSubmitResponse(question, XPSubmitQuestion, participant.XPScore), nil
}

// List usecase untuk mendapatkan list questions
func (c *QuestionUseCase) List(ctx context.Context, request *model.GetQuestionsRequest) (*model.QuestionListResponse, error) {
	tx := c.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	// validate request
	if err := c.Validator.Struct(request); err != nil {
		c.Log.Warnf("List - Invalid request: %v", err)
		return nil, fiber.ErrBadRequest
	}

	// set defaults
	if request.Limit == 0 {
		request.Limit = 20
	}
	if request.SortBy == "" {
		request.SortBy = "upvotes"
	}

	// check room exists
	roomCount, err := c.RoomRepository.CountById(tx, request.RoomID)
	if err != nil {
		c.Log.Errorf("List - RoomRepository.CountById error: %v", err)
		return nil, fiber.ErrInternalServerError
	}
	if roomCount == 0 {
		return nil, fiber.ErrNotFound
	}

	// get questions
	questions, err := c.QuestionRepository.List(tx, request.RoomID, request.Status, request.SortBy, request.Limit, request.Offset)
	if err != nil {
		c.Log.Errorf("List - QuestionRepository.List error: %v", err)
		return nil, fiber.ErrInternalServerError
	}

	// get total count
	total, err := c.QuestionRepository.Count(tx, request.RoomID, request.Status)
	if err != nil {
		c.Log.Errorf("List - QuestionRepository.Count error: %v", err)
		return nil, fiber.ErrInternalServerError
	}

	// get voted question IDs for current participant
	questionIDs := make([]uint, len(questions))
	for i, q := range questions {
		questionIDs[i] = q.ID
	}

	votedMap := make(map[uint]bool)
	if len(questionIDs) > 0 {
		votedMap, err = c.VoteRepository.GetVotedQuestionIDs(tx, request.ParticipantID, questionIDs)
		if err != nil {
			c.Log.Warnf("List - VoteRepository.GetVotedQuestionIDs error: %v", err)
		}
	}

	if err = tx.Commit().Error; err != nil {
		c.Log.Errorf("List - Commit error: %v", err)
		return nil, fiber.ErrInternalServerError
	}

	// convert to response
	responses := make([]model.QuestionResponse, len(questions))
	for i, q := range questions {
		responses[i] = *converter.QuestionToResponseWithParticipant(&q, votedMap[q.ID])
	}

	return &model.QuestionListResponse{
		Questions: responses,
		Paging: model.QuestionPaging{
			Total:  total,
			Limit:  request.Limit,
			Offset: request.Offset,
		},
	}, nil
}

// Upvote usecase untuk upvote question
func (c *QuestionUseCase) Upvote(ctx context.Context, request *model.UpvoteRequest) (*model.UpvoteResponse, error) {
	tx := c.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	// validate request
	if err := c.Validator.Struct(request); err != nil {
		c.Log.Warnf("Upvote - Invalid request: %v", err)
		return nil, fiber.ErrBadRequest
	}

	// get question
	question, err := c.QuestionRepository.FindByIdWithParticipant(tx, request.QuestionID)
	if err != nil {
		c.Log.Errorf("Upvote - QuestionRepository.FindByIdWithParticipant error: %v", err)
		return nil, fiber.ErrInternalServerError
	}
	if question == nil {
		return nil, fiber.ErrNotFound
	}

	// check if already voted
	hasVoted, err := c.VoteRepository.HasVoted(tx, request.QuestionID, request.ParticipantID)
	if err != nil {
		c.Log.Errorf("Upvote - VoteRepository.HasVoted error: %v", err)
		return nil, fiber.ErrInternalServerError
	}
	if hasVoted {
		return nil, fiber.NewError(fiber.StatusConflict, "Already voted")
	}

	// cannot vote own question
	if question.ParticipantID == request.ParticipantID {
		return nil, fiber.NewError(fiber.StatusBadRequest, "Cannot upvote own question")
	}

	// create vote
	vote := &entity.Vote{
		QuestionID:    request.QuestionID,
		ParticipantID: request.ParticipantID,
	}
	if err = c.VoteRepository.Create(tx, vote); err != nil {
		c.Log.Errorf("Upvote - VoteRepository.Create error: %v", err)
		return nil, fiber.ErrInternalServerError
	}

	// NOTE: upvote_count increment sudah di-handle oleh database trigger (after_vote_insert)
	// Jangan manual increment disini untuk menghindari double increment

	// add XP to question owner (recipient)
	xpTx := &entity.XPTransaction{
		ParticipantID: question.ParticipantID,
		RoomID:        question.RoomID,
		Points:        XPReceiveUpvote,
		SourceType:    "upvote_received",
		SourceID:      vote.ID,
	}
	if err = c.XPTransactionRepository.Create(tx, xpTx); err != nil {
		c.Log.Errorf("Upvote - XPTransactionRepository.Create error: %v", err)
		return nil, fiber.ErrInternalServerError
	}

	// update recipient XP score
	if err = c.XPTransactionRepository.AddXP(tx, question.ParticipantID, XPReceiveUpvote); err != nil {
		c.Log.Errorf("Upvote - AddXP error: %v", err)
		return nil, fiber.ErrInternalServerError
	}

	if err = tx.Commit().Error; err != nil {
		c.Log.Errorf("Upvote - Commit error: %v", err)
		return nil, fiber.ErrInternalServerError
	}

	return converter.VoteToUpvoteResponse(vote, question.UpvoteCount+1, question.ParticipantID, XPReceiveUpvote), nil
}

// RemoveUpvote usecase untuk remove upvote
func (c *QuestionUseCase) RemoveUpvote(ctx context.Context, request *model.UpvoteRequest) (*model.RemoveUpvoteResponse, error) {
	tx := c.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	// validate request
	if err := c.Validator.Struct(request); err != nil {
		c.Log.Warnf("RemoveUpvote - Invalid request: %v", err)
		return nil, fiber.ErrBadRequest
	}

	// get question
	question, err := c.QuestionRepository.FindByIdWithParticipant(tx, request.QuestionID)
	if err != nil {
		c.Log.Errorf("RemoveUpvote - QuestionRepository.FindByIdWithParticipant error: %v", err)
		return nil, fiber.ErrInternalServerError
	}
	if question == nil {
		return nil, fiber.ErrNotFound
	}

	// check if voted
	vote, err := c.VoteRepository.FindByQuestionAndParticipant(tx, request.QuestionID, request.ParticipantID)
	if err != nil {
		c.Log.Errorf("RemoveUpvote - VoteRepository.FindByQuestionAndParticipant error: %v", err)
		return nil, fiber.ErrInternalServerError
	}
	if vote == nil {
		return nil, fiber.NewError(fiber.StatusNotFound, "Vote not found")
	}

	// delete vote
	if err = c.VoteRepository.DeleteByQuestionAndParticipant(tx, request.QuestionID, request.ParticipantID); err != nil {
		c.Log.Errorf("RemoveUpvote - DeleteByQuestionAndParticipant error: %v", err)
		return nil, fiber.ErrInternalServerError
	}

	// NOTE: upvote_count decrement sudah di-handle oleh database trigger (after_vote_delete)
	// Jangan manual decrement disini untuk menghindari double decrement

	// remove XP from question owner (optional: bisa diabaikan untuk simplifikasi)
	if err = c.XPTransactionRepository.AddXP(tx, question.ParticipantID, -XPReceiveUpvote); err != nil {
		c.Log.Warnf("RemoveUpvote - RemoveXP error (ignored): %v", err)
	}

	if err = tx.Commit().Error; err != nil {
		c.Log.Errorf("RemoveUpvote - Commit error: %v", err)
		return nil, fiber.ErrInternalServerError
	}

	newCount := question.UpvoteCount - 1
	if newCount < 0 {
		newCount = 0
	}

	return converter.QuestionToRemoveUpvoteResponse(request.QuestionID, newCount), nil
}

// Validate usecase untuk validate question (presenter only)
func (c *QuestionUseCase) Validate(ctx context.Context, request *model.ValidateQuestionRequest) (*model.ValidateQuestionResponse, error) {
	tx := c.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	// validate request
	if err := c.Validator.Struct(request); err != nil {
		c.Log.Warnf("Validate - Invalid request: %v", err)
		return nil, fiber.ErrBadRequest
	}

	// get question
	question, err := c.QuestionRepository.FindByIdWithParticipant(tx, request.QuestionID)
	if err != nil {
		c.Log.Errorf("Validate - QuestionRepository.FindByIdWithParticipant error: %v", err)
		return nil, fiber.ErrInternalServerError
	}
	if question == nil {
		return nil, fiber.ErrNotFound
	}

	// check if already validated
	if question.IsValidatedByPresenter {
		return nil, fiber.NewError(fiber.StatusConflict, "Question already validated")
	}

	// update validation status
	if err = c.QuestionRepository.UpdateValidation(tx, request.QuestionID, request.Status, XPPresenterValidate); err != nil {
		c.Log.Errorf("Validate - UpdateValidation error: %v", err)
		return nil, fiber.ErrInternalServerError
	}

	// add XP to question owner
	xpTx := &entity.XPTransaction{
		ParticipantID: question.ParticipantID,
		RoomID:        question.RoomID,
		Points:        XPPresenterValidate,
		SourceType:    "presenter_validated",
		SourceID:      question.ID,
	}
	if err = c.XPTransactionRepository.Create(tx, xpTx); err != nil {
		c.Log.Errorf("Validate - XPTransactionRepository.Create error: %v", err)
		return nil, fiber.ErrInternalServerError
	}

	// update participant XP score
	if err = c.XPTransactionRepository.AddXP(tx, question.ParticipantID, XPPresenterValidate); err != nil {
		c.Log.Errorf("Validate - AddXP error: %v", err)
		return nil, fiber.ErrInternalServerError
	}

	// get new total XP
	var participant entity.Participant
	if err = tx.Select("xp_score").Where("id = ?", question.ParticipantID).First(&participant).Error; err != nil {
		c.Log.Errorf("Validate - Get participant error: %v", err)
		return nil, fiber.ErrInternalServerError
	}

	if err = tx.Commit().Error; err != nil {
		c.Log.Errorf("Validate - Commit error: %v", err)
		return nil, fiber.ErrInternalServerError
	}

	// update question for response
	question.Status = request.Status
	question.IsValidatedByPresenter = true

	return converter.QuestionToValidateResponse(question, XPPresenterValidate, participant.XPScore), nil
}
