package usecase

import (
	"slido-clone-backend/internal/entity"
	"slido-clone-backend/internal/repository"

	"github.com/go-playground/validator/v10"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type XPTransactionUseCase struct {
	DB                      *gorm.DB
	Validate                *validator.Validate
	Log                     *logrus.Logger
	XPTransactionRepository *repository.XPTransactionRepository
}

func NewXPTransactionUseCase(db *gorm.DB, validate *validator.Validate, log *logrus.Logger, xpTransactionRepository *repository.XPTransactionRepository) *XPTransactionUseCase {
	return &XPTransactionUseCase{
		DB:                      db,
		Validate:                validate,
		Log:                     log,
		XPTransactionRepository: xpTransactionRepository,
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
