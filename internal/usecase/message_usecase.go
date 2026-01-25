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

type MessageUseCase struct {
	DB                    *gorm.DB
	Validate              *validator.Validate
	Log                   *logrus.Logger
	MessageRepository     *repository.MessageRepository
	RoomRepository        *repository.RoomRepository
	ParticipantRepository *repository.ParticipantRepository
	XPTransactionUseCase  *XPTransactionUseCase
}

func NewMessageUseCase(db *gorm.DB, validate *validator.Validate, log *logrus.Logger, messageRepository *repository.MessageRepository, roomRepository *repository.RoomRepository, participantRepository *repository.ParticipantRepository, xpTransactionUseCase *XPTransactionUseCase) *MessageUseCase {
	return &MessageUseCase{
		DB:                    db,
		Validate:              validate,
		Log:                   log,
		MessageRepository:     messageRepository,
		RoomRepository:        roomRepository,
		ParticipantRepository: participantRepository,
		XPTransactionUseCase:  xpTransactionUseCase,
	}
}

func (c *MessageUseCase) Send(ctx context.Context, request *model.SendMessageRequest) (*model.MessageResponse, error) {
	// begin transaction
	tx := c.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	// validate request
	if err := c.Validate.Struct(request); err != nil {
		c.Log.Errorf("Send - Validate Struct Error: %v", err)
		return nil, fiber.ErrBadRequest
	}

	// logic to create message
	// check if room exists
	roomCount, err := c.RoomRepository.CountById(tx, request.RoomID)
	if err != nil {
		c.Log.Errorf("Send - RoomRepository.CountById Error: %v", err)
		return nil, fiber.ErrInternalServerError
	}

	if roomCount == 0 {
		c.Log.Warnf("Send - Room not found with ID: %d", request.RoomID)
		return nil, fiber.ErrNotFound
	}

	// check if participant is in the room
	participantCount, err := c.ParticipantRepository.CountById(tx, request.ParticipantID)
	if err != nil {
		c.Log.Errorf("Send - ParticipantRepository.CountById Error: %v", err)
		return nil, fiber.ErrInternalServerError
	}

	if participantCount == 0 {
		c.Log.Warnf("Send - Participant with ID %d not found", request.ParticipantID)
		return nil, fiber.ErrNotFound
	}

	// create message in repository
	message := &entity.Message{
		RoomID:        request.RoomID,
		ParticipantID: request.ParticipantID,
		Content:       request.Content,
	}

	err = c.MessageRepository.Create(tx, message)
	if err != nil {
		c.Log.Errorf("Send - MessageRepository.Create failed to create message: %v", err)
		return nil, fiber.ErrInternalServerError
	}

	// load participant relation
	if err = tx.Preload("Participant").First(message, message.ID).Error; err != nil {
		c.Log.Warnf("failed to laod participant relation: %v", err)
		return nil, fiber.ErrInternalServerError
	}

	// add xp to participant for sending message
	if err = c.XPTransactionUseCase.AddXPForMessage(tx, request.RoomID, request.ParticipantID, message.ID); err != nil {
		c.Log.Warnf("failed to add xp and update participant score: %v", err)
		return nil, fiber.ErrInternalServerError
	}

	// commit transaction
	if err = tx.Commit().Error; err != nil {
		c.Log.Errorf("Send - Transaction Commit Error: %v", err)
		return nil, fiber.ErrInternalServerError
	}

	// return response
	return converter.MessageToResponse(message), nil
}

func (c *MessageUseCase) List(ctx context.Context, request *model.GetMessagesRequest) (*model.MessageListResponse, error) {
	// begin transaction
	tx := c.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	// validate request
	if err := c.Validate.Struct(request); err != nil {
		c.Log.Errorf("List - Validate Struct Error: %v", err)
		return nil, fiber.ErrBadRequest
	}

	// logic to list messages with pagination
	// check if room exists
	roomCount, err := c.RoomRepository.CountById(tx, request.RoomID)
	if err != nil {
		c.Log.Errorf("List - RoomRepository.CountById Error: %v", err)
		return nil, fiber.ErrInternalServerError
	}

	if roomCount == 0 {
		c.Log.Warnf("List - Room not found with ID: %d", request.RoomID)
		return nil, fiber.ErrNotFound
	}

	// get messages from repository
	messages, err := c.MessageRepository.List(tx, request.RoomID, request.Limit+1, request.Before)
	if err != nil {
		c.Log.Errorf("Failed to list messages: %v", err)
		return nil, fiber.ErrInternalServerError
	}

	// hasMore calculation
	hasMore := len(messages) > request.Limit
	if hasMore {
		messages = messages[:request.Limit]
	}

	// commit transaction
	if err = tx.Commit().Error; err != nil {
		c.Log.Errorf("List - Transaction Commit Error: %v", err)
		return nil, fiber.ErrInternalServerError
	}

	// return response
	return converter.MessagesToMessageListResponse(messages, hasMore), nil
}
