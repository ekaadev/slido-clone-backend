package repository

import (
	"errors"
	"slido-clone-backend/internal/entity"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type RoomRepository struct {
	Repository[entity.Room]
	Log *logrus.Logger
}

// NewRoomRepository create new instance of RoomRepository
func NewRoomRepository(log *logrus.Logger) *RoomRepository {
	return &RoomRepository{
		Log: log,
	}
}

// FindByRoomCode find room by room code
func (r *RoomRepository) FindByRoomCode(db *gorm.DB, code string, presenterId uint) (*entity.Room, error) {
	var room entity.Room
	// with preload associations, karena butuh data yang terkait
	err := db.Preload(clause.Associations).Where("room_code = ? AND presenter_id = ?", code, presenterId).First(&room).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}

	return &room, err
}

// FindById find room by id
func (r *RoomRepository) FindById(db *gorm.DB, id uint, presenterId uint) (*entity.Room, error) {
	var room entity.Room
	err := db.Where("id = ? AND presenter_id = ?", id, presenterId).First(&room).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}

	return &room, err
}

// Search get all rooms by presenter id
func (r *RoomRepository) Search(db *gorm.DB, presenterId uint) ([]entity.Room, error) {
	var rooms []entity.Room
	err := db.Preload("Participants").Where("presenter_id = ?", presenterId).Find(&rooms).Error
	if err != nil {
		return nil, err
	}

	return rooms, nil
}
