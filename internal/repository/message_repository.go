package repository

import (
	"slido-clone-backend/internal/entity"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type MessageRepository struct {
	Repository[entity.Message]
	Log *logrus.Logger
}

func NewMessageRepository(log *logrus.Logger) *MessageRepository {
	return &MessageRepository{
		Log: log,
	}
}

// CreateAndLoad membuat message baru dan memuat relasi yang diperlukan
func (r *MessageRepository) CreateAndLoad(db *gorm.DB, message *entity.Message) (*entity.Message, error) {
	err := db.Create(message).Error
	if err != nil {
		return nil, err
	}

	var createdMessage entity.Message
	err = db.Preload(clause.Associations).Where("id = ?", message.ID).Take(&createdMessage).Error
	if err != nil {
		return nil, err
	}

	return &createdMessage, nil
}

// List mencari message berdasarkan room ID dengan pagination sebelum waktu tertentu
func (r *MessageRepository) List(db *gorm.DB, roomID uint, limit int, before *int64) ([]entity.Message, error) {
	var messages []entity.Message
	query := db.Preload("Participant").Where("room_id = ?", roomID).Order("created_at DESC").Limit(limit)

	// jika ada before, ambil message sebelum waktu tersebut
	if before != nil {
		query = query.Where("id < ?", *before)
	}

	err := query.Find(&messages).Error
	if err != nil {
		return nil, err
	}

	return messages, nil
}

// CountByRoomID menghitung jumlah message dalam sebuah room
func (r *MessageRepository) CountByRoomID(db *gorm.DB, roomID uint) (int64, error) {
	var count int64
	err := db.Model(entity.Message{}).Where("room_id = ?", roomID).Count(&count).Error
	return count, err
}
