package usecase

import (
	"context"
	"crypto/rand"
	"math/big"
	"slido-clone-backend/internal/entity"
	"slido-clone-backend/internal/model"
	"slido-clone-backend/internal/model/converter"
	"slido-clone-backend/internal/repository"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

const charset = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

type RoomUseCase struct {
	DB             *gorm.DB
	Log            *logrus.Logger
	Validate       *validator.Validate
	RoomRepository *repository.RoomRepository
}

// NewRoomUseCase create new instance of RoomUseCase
func NewRoomUseCase(db *gorm.DB, log *logrus.Logger, validate *validator.Validate, roomRepository *repository.RoomRepository) *RoomUseCase {
	return &RoomUseCase{
		DB:             db,
		Log:            log,
		Validate:       validate,
		RoomRepository: roomRepository,
	}
}

// Create usecase untuk membuat room baru
func (c *RoomUseCase) Create(ctx context.Context, request *model.CreateRoomRequest) (*model.CreateRoomResponse, error) {
	// begin transaction
	tx := c.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	// validate request
	err := c.Validate.Struct(request)
	if err != nil {
		c.Log.Warnf("Invalid request body: %+v", err)
		return nil, fiber.ErrBadRequest
	}

	// logic to create room
	// call function generate room code
	roomCode, err := GenerateRoomCode(6)
	if err != nil {
		c.Log.Warnf("Failed to generate room code: %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	room := &entity.Room{
		RoomCode:    roomCode,
		Title:       request.Title,
		PresenterID: request.PresenterID,
	}

	// create room and call method in repository
	err = c.RoomRepository.Create(tx, room)
	if err != nil {
		c.Log.Warnf("Failed to create room: %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	// commit transaction
	if err = tx.Commit().Error; err != nil {
		c.Log.Warnf("Failed to commit transaction: %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	// return room response
	return converter.RoomToCreateRoomResponse(room), nil
}

// GenerateRoomCode generate with crypto/rand
func GenerateRoomCode(n int) (string, error) {
	result := make([]byte, n)

	// generate random
	for i := range result {
		num, err := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		if err != nil {
			return "", err
		}
		result[i] = charset[num.Int64()]
	}

	return string(result), nil
}
