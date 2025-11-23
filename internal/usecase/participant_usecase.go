package usecase

import (
	"context"
	"slido-clone-backend/internal/entity"
	"slido-clone-backend/internal/model"
	"slido-clone-backend/internal/model/converter"
	"slido-clone-backend/internal/repository"
	"slido-clone-backend/internal/util"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type ParticipantUseCase struct {
	DB                    *gorm.DB
	Log                   *logrus.Logger
	Validate              *validator.Validate
	ParticipantRepository *repository.ParticipantRepository
	RoomRepository        *repository.RoomRepository
	UserRepository        *repository.UserRepository
	TokenUtil             *util.TokenUtil
}

func NewParticipantUseCase(db *gorm.DB, log *logrus.Logger, validate *validator.Validate, participantRepository *repository.ParticipantRepository, roomRepository *repository.RoomRepository, userRepository *repository.UserRepository, tokenUtil *util.TokenUtil) *ParticipantUseCase {
	return &ParticipantUseCase{
		DB:                    db,
		Log:                   log,
		Validate:              validate,
		ParticipantRepository: participantRepository,
		RoomRepository:        roomRepository,
		UserRepository:        userRepository,
		TokenUtil:             tokenUtil,
	}
}

// Join usecase digunakan untuk participant bergabung ke dalam room
func (c *ParticipantUseCase) Join(ctx context.Context, request *model.JoinRoomRequest) (*model.JoinRoomResponse, error) {
	// transaction begin
	tx := c.DB.WithContext(ctx).Begin()
	defer tx.Commit()

	// validate request
	if err := c.Validate.Struct(request); err != nil {
		c.Log.Warnf("Invalid request body: %+v", err)
		return nil, fiber.ErrBadRequest
	}

	// logic to join room
	// check room exist
	roomExisting, err := c.RoomRepository.FindByRoomCode(tx, request.RoomCode)
	if err != nil {
		c.Log.Warnf("Failed to find room by room code: %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	if roomExisting == nil {
		c.Log.Warnf("Room not found with room code: %s", request.RoomCode)
		return nil, fiber.ErrNotFound
	}

	if roomExisting.Status == "closed" {
		c.Log.Warnf("Room is closed with room code: %s", request.RoomCode)
		return nil, fiber.ErrBadRequest
	}

	// check user registered is exist
	userExisting, err := c.UserRepository.FindByUsername(tx, request.Username)
	if err != nil {
		c.Log.Warnf("Failed to find user by username: %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	if userExisting == nil {
		c.Log.Warnf("User not found with username: %s", request.Username)
		return nil, fiber.ErrUnauthorized
	}

	// check participant already join room
	participantExisting, err := c.ParticipantRepository.FindByRoomIDAndUserID(tx, roomExisting.ID, userExisting.ID)
	if err != nil {
		c.Log.Warnf("Failed to find participant by room ID and user ID: %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	// if participant already joined, return existing participant
	if participantExisting != nil {
		c.Log.Warnf("Participant already joined room with room code: %s", request.RoomCode)

		if err = tx.Commit().Error; err != nil {
			c.Log.Warnf("Failed to commit transaction: %+v", err)
			return nil, fiber.ErrInternalServerError
		}

		// Determine role based on presenter status
		role := "participant"
		if roomExisting.PresenterID == userExisting.ID {
			role = "presenter"
		}

		// Generate new token
		token, err := c.TokenUtil.CreateToken(ctx, &model.Auth{
			UserID:        &userExisting.ID,
			ParticipantID: &participantExisting.ID,
			RoomID:        &roomExisting.ID,
			Username:      userExisting.Username,
			DisplayName:   participantExisting.DisplayName,
			Email:         userExisting.Email,
			Role:          role,
			IsAnonymous:   *participantExisting.IsAnonymous,
		})
		if err != nil {
			c.Log.Warnf("Failed to create token: %+v", err)
			return nil, fiber.ErrInternalServerError
		}

		return converter.ParticipantToJoinRoomResponse(participantExisting, token), nil
	}

	anon := false
	// create participant entity
	participant := &entity.Participant{
		RoomID:      roomExisting.ID,
		UserID:      &userExisting.ID,
		DisplayName: userExisting.Username,
		IsAnonymous: &anon,
	}

	// check is display name provided
	if request.DisplayName != "" {
		participant.DisplayName = request.DisplayName
	}

	// create participant in repository
	if err = c.ParticipantRepository.Create(tx, participant); err != nil {
		c.Log.Warnf("Failed to create participant: %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	// transaction commit
	if err = tx.Commit().Error; err != nil {
		c.Log.Warnf("Failed to commit transaction: %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	// generate new token (update jwt token)
	token, err := c.TokenUtil.CreateToken(ctx, &model.Auth{
		UserID:        &userExisting.ID,
		ParticipantID: &participant.ID,
		RoomID:        &roomExisting.ID,
		Username:      userExisting.Username,
		DisplayName:   participant.DisplayName,
		Email:         userExisting.Email,
		Role:          "participant",
		IsAnonymous:   *participant.IsAnonymous,
	})
	if err != nil {
		c.Log.Warnf("Failed to create token: %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	// return response
	return converter.ParticipantToJoinRoomResponse(participant, token), nil
}

// List usecase digunakan untuk mencari participant dalam room
func (c *ParticipantUseCase) List(ctx context.Context, request *model.ListParticipantsRequest) (*model.ParticipantListResponse, int64, error) {
	// begin transaction
	tx := c.DB.WithContext(ctx).Begin()
	defer tx.Commit()

	// validate request
	if err := c.Validate.Struct(request); err != nil {
		c.Log.Warnf("Invalid request body: %+v", err)
		return nil, 0, fiber.ErrBadRequest
	}

	// logic to list participants
	// check room exist
	total, err := c.RoomRepository.CountById(tx, request.RoomID)
	if err != nil {
		c.Log.Warnf("Failed to count room by ID: %+v", err)
		return nil, 0, fiber.ErrInternalServerError
	}

	if total == 0 {
		c.Log.Warnf("Room not found with ID: %d", request.RoomID)
		return nil, 0, fiber.ErrNotFound
	}

	// list participants from repository
	participants, total, err := c.ParticipantRepository.List(tx, request.RoomID, request.Page, request.Size)
	if err != nil {
		c.Log.Warnf("Failed to list participants: %+v", err)
		return nil, 0, fiber.ErrInternalServerError
	}

	// commit transaction
	if err = tx.Commit().Error; err != nil {
		c.Log.Warnf("Failed to commit transaction: %+v", err)
		return nil, 0, fiber.ErrInternalServerError
	}

	// return response
	responses := make([]*model.ParticipantListItem, len(participants))
	for i, participant := range participants {
		responses[i] = converter.ParticipantToListItem(&participant)
	}

	return converter.ParticipantsToListResponse(responses), total, nil
}

func (c *ParticipantUseCase) Leaderboard(ctx context.Context, request *model.GetLeaderboardRequest) (*model.LeaderboardResponse, error) {
	// begin transaction
	tx := c.DB.WithContext(ctx).Begin()
	defer tx.Commit()

	// validate request
	if err := c.Validate.Struct(request); err != nil {
		c.Log.Warnf("invalid request body: %+v", err)
		return nil, fiber.ErrBadRequest
	}

	// logic to get leaderboard
	// check room exist
	roomCount, err := c.RoomRepository.CountById(tx, request.RoomID)
	if err != nil {
		c.Log.Warnf("failed to count room by id: %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	if roomCount == 0 {
		c.Log.Warnf("room not found with id: %d", request.RoomID)
		return nil, fiber.ErrNotFound
	}

	// get total participants in the room
	participantCountByRoomID, err := c.ParticipantRepository.CountByRoomID(tx, request.RoomID)
	if err != nil {
		c.Log.Warnf("failed to count participants by room id: %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	// get leaderboard from repository
	leaderboard, err := c.ParticipantRepository.ListLeaderboard(tx, request.RoomID)
	if err != nil {
		c.Log.Warnf("failed to list leaderboard: %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	// get my rank (from participant ID)
	rank, xpScore, err := c.ParticipantRepository.GetRankAndScoreByParticipantID(tx, request.RoomID, request.ParticipantID)

	myRank := &model.MyRank{
		Rank:    int(rank),
		XPScore: int(xpScore),
	}

	// commit transaction
	if err = tx.Commit().Error; err != nil {
		c.Log.Warnf("failed to commit transaction: %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	// return response
	return converter.ParticipantsToLeaderboardResponse(leaderboard, myRank, int(participantCountByRoomID)), nil
}
