package usecase

import (
	"context"
	"crypto/rand"
	"math/big"
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

// Get usecase untuk mendapatkan detail room berdasarkan room code
func (c *RoomUseCase) Get(ctx context.Context, request *model.GetRoomRequestByRoomCode) (*model.GetRoomDetailResponse, error) {
	// begin transaction
	tx := c.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	// validate get room request
	err := c.Validate.Struct(request)
	if err != nil {
		c.Log.Warnf("Invalid request body: %+v", err)
		return nil, fiber.ErrBadRequest
	}

	// logic to get room by room code
	// call repository to find room by room code
	existingRoom, err := c.RoomRepository.FindByRoomCode(tx, request.RoomCode)
	if err != nil {
		c.Log.Warnf("Failed to find room by room code: %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	if existingRoom == nil {
		return nil, fiber.ErrNotFound
	}

	// commit transaction
	if err = tx.Commit().Error; err != nil {
		c.Log.Warnf("Failed to commit transaction: %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	// return room detail response
	return converter.RoomToGetRoomDetailResponse(existingRoom), nil
}

// UpdateToClosed usecase untuk mengupdate room berdasarkan id
// hanya bisa diubah statusnya menjadi closed
// TODO: berpotensi untuk digunakan dalam websocket "on close room"
func (c *RoomUseCase) UpdateToClosed(ctx context.Context, request *model.UpdateToCloseRoomRequestByID) (*model.UpdateToCloseRoomResponse, error) {
	// begin transaction
	tx := c.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	// validate update room request
	err := c.Validate.Struct(request)
	if err != nil {
		c.Log.Warnf("Invalid request body: %+v", err)
		return nil, fiber.ErrBadRequest
	}

	// logic to update room by id
	room, err := c.RoomRepository.FindById(tx, request.RoomID, request.PresenterID)
	if err != nil {
		c.Log.Warnf("Failed to find room by id: %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	if room == nil {
		return nil, fiber.ErrNotFound
	}

	// update room status to closed
	room.Status = request.Status
	now := time.Now()
	room.ClosedAt = &now

	// update room in repository
	err = c.RoomRepository.Update(tx, room)
	if err != nil {
		c.Log.Warnf("Failed to update room: %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	// commit transaction
	if err = tx.Commit().Error; err != nil {
		c.Log.Warnf("Failed to commit transaction: %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	return converter.UpdateToCloseRoomToResponse(room), nil
}

// Search usecase untuk mencari room yang dimiliki oleh presenter
func (c *RoomUseCase) Search(ctx context.Context, request *model.SearchRoomsRequest) (*model.RoomListResponse, error) {
	// begin transaction
	tx := c.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	// validate search room request
	if err := c.Validate.Struct(request); err != nil {
		c.Log.Warnf("Invalid request: %+v", err)
		return nil, fiber.ErrBadRequest
	}

	// logic to search rooms by presenter id
	// call repository to search rooms
	rooms, err := c.RoomRepository.Search(tx, request.PresenterID)
	if err != nil {
		c.Log.Warnf("Failed to search rooms: %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	// commit transaction
	if err = tx.Commit().Error; err != nil {
		c.Log.Warnf("Failed to commit transaction: %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	// return room list response
	responses := make([]*model.RoomListItem, len(rooms))
	for i, room := range rooms {
		responses[i] = converter.RoomToListItemResponse(&room)
	}
	return converter.RoomsToListResponse(responses), nil
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
