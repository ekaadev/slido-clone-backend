package repository

import (
	"errors"
	"slido-clone-backend/internal/entity"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type ParticipantRepository struct {
	Repository[entity.Participant]
	Log *logrus.Logger
}

func NewParticipantRepository(log *logrus.Logger) *ParticipantRepository {
	return &ParticipantRepository{
		Log: log,
	}
}

func (r *ParticipantRepository) FindByRoomIDAndUserID(db *gorm.DB, roomID uint, userID uint) (*entity.Participant, error) {
	var participant entity.Participant
	err := db.Where("room_id = ? AND user_id = ?", roomID, userID).First(&participant).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}

	return &participant, err
}
