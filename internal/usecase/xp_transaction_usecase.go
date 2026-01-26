package usecase

import (
	"context"
	"slido-clone-backend/internal/entity"
	"slido-clone-backend/internal/model"
	"slido-clone-backend/internal/repository"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type XPTransactionUseCase struct {
	DB                      *gorm.DB
	Validate                *validator.Validate
	Log                     *logrus.Logger
	XPTransactionRepository *repository.XPTransactionRepository
	RoomRepository          *repository.RoomRepository
}

func NewXPTransactionUseCase(db *gorm.DB, validate *validator.Validate, log *logrus.Logger, xpTransactionRepository *repository.XPTransactionRepository, roomRepository *repository.RoomRepository) *XPTransactionUseCase {
	return &XPTransactionUseCase{
		DB:                      db,
		Validate:                validate,
		Log:                     log,
		XPTransactionRepository: xpTransactionRepository,
		RoomRepository:          roomRepository,
	}
}

func (c *XPTransactionUseCase) AddXPForMessage(tx *gorm.DB, roomID, participantID, messageID uint) error {
	xpPoint := 1 // poin XP untuk setiap pesan yang dikirim

	xp := &entity.XPTransaction{
		ParticipantID: participantID,
		RoomID:        roomID,
		Points:        xpPoint,
		SourceType:    "message_created",
		SourceID:      messageID,
	}

	if err := c.XPTransactionRepository.Create(tx, xp); err != nil {
		return err
	}

	// update participant score
	if err := c.XPTransactionRepository.AddXP(tx, participantID, xpPoint); err != nil {
		return err
	}

	return nil
}

// GetTransactions usecase untuk mendapatkan XP transactions history
func (c *XPTransactionUseCase) GetTransactions(ctx context.Context, request *model.GetXPTransactionsRequest) (*model.GetXPTransactionsResponse, error) {
	tx := c.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	// validate request
	if err := c.Validate.Struct(request); err != nil {
		c.Log.Warnf("GetTransactions - Invalid request: %v", err)
		return nil, fiber.ErrBadRequest
	}

	// set default limit
	if request.Limit == 0 {
		request.Limit = 50
	}

	// check room exists
	roomCount, err := c.RoomRepository.CountById(tx, request.RoomID)
	if err != nil {
		c.Log.Errorf("GetTransactions - RoomRepository.CountById error: %v", err)
		return nil, fiber.ErrInternalServerError
	}
	if roomCount == 0 {
		return nil, fiber.ErrNotFound
	}

	// get transactions
	transactions, err := c.XPTransactionRepository.FindByRoomAndParticipant(tx, request.RoomID, request.ParticipantID, request.Limit)
	if err != nil {
		c.Log.Errorf("GetTransactions - FindByRoomAndParticipant error: %v", err)
		return nil, fiber.ErrInternalServerError
	}

	// get total count
	total, err := c.XPTransactionRepository.CountByRoomAndParticipant(tx, request.RoomID, request.ParticipantID)
	if err != nil {
		c.Log.Errorf("GetTransactions - CountByRoomAndParticipant error: %v", err)
		return nil, fiber.ErrInternalServerError
	}

	// get total XP
	totalXP, err := c.XPTransactionRepository.GetTotalXPByParticipant(tx, request.ParticipantID)
	if err != nil {
		c.Log.Warnf("GetTransactions - GetTotalXPByParticipant error: %v", err)
	}

	if err = tx.Commit().Error; err != nil {
		c.Log.Errorf("GetTransactions - Commit error: %v", err)
		return nil, fiber.ErrInternalServerError
	}

	// build response
	items := make([]model.XPTransactionItem, len(transactions))
	for i, t := range transactions {
		items[i] = model.XPTransactionItem{
			ID:         t.ID,
			Points:     t.Points,
			SourceType: t.SourceType,
			SourceID:   t.SourceID,
			CreatedAt:  t.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		}
	}

	return &model.GetXPTransactionsResponse{
		Transactions: items,
		TotalXP:      totalXP,
		Total:        total,
	}, nil
}
