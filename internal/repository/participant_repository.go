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

// List mencari participant berdasarkan room ID dengan pagination
func (r *ParticipantRepository) List(db *gorm.DB, roomID uint, page int, size int) ([]entity.Participant, int64, error) {
	var participants []entity.Participant
	err := db.Where("room_id = ?", roomID).Offset((page - 1) * size).Limit(size).Find(&participants).Error
	if err != nil {
		return nil, 0, err
	}

	var total int64
	err = db.Model(&entity.Participant{}).Where("room_id = ?", roomID).Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	return participants, total, nil
}
