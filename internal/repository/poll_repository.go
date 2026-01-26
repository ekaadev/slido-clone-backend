package repository

import (
	"errors"
	"slido-clone-backend/internal/entity"
	"time"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type PollRepository struct {
	Repository[entity.Poll]
	Log *logrus.Logger
}

func NewPollRepository(log *logrus.Logger) *PollRepository {
	return &PollRepository{
		Log: log,
	}
}

// GetActivePollByRoomID retrieves the active poll for a given room ID.
// 1 room hanya bisa punya 1 active poll pada satu waktu
func (r *PollRepository) GetActivePollByRoomID(db *gorm.DB, poll entity.Poll) (*entity.Poll, error) {
	var activePoll entity.Poll
	err := db.Where("room_id = ? AND status = ?", poll.RoomID, "active").First(&activePoll).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}

	return &activePoll, err
}

// FindByIdWithOptions find poll by id dengan preload options
func (r *PollRepository) FindByIdWithOptions(db *gorm.DB, id uint) (*entity.Poll, error) {
	var poll entity.Poll
	err := db.Preload("Options", func(db *gorm.DB) *gorm.DB {
		return db.Order("poll_options.order ASC")
	}).Where("id = ?", id).First(&poll).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &poll, err
}

// FindActiveByRoomID find active polls di room dengan options
func (r *PollRepository) FindActiveByRoomID(db *gorm.DB, roomID uint) ([]entity.Poll, error) {
	var polls []entity.Poll
	err := db.Preload("Options", func(db *gorm.DB) *gorm.DB {
		return db.Order("poll_options.order ASC")
	}).Where("room_id = ? AND status = ?", roomID, "active").
		Order("created_at DESC").Find(&polls).Error
	return polls, err
}

// ListByRoomID get poll history berdasarkan room
func (r *PollRepository) ListByRoomID(db *gorm.DB, roomID uint, status string, limit int) ([]entity.Poll, error) {
	var polls []entity.Poll

	query := db.Where("room_id = ?", roomID)

	// filter by status
	if status != "" && status != "all" {
		query = query.Where("status = ?", status)
	}

	// limit
	if limit > 0 {
		query = query.Limit(limit)
	}

	err := query.Order("created_at DESC").Find(&polls).Error
	return polls, err
}

// CountByRoomID count polls berdasarkan room dan status
func (r *PollRepository) CountByRoomID(db *gorm.DB, roomID uint, status string) (int64, error) {
	var count int64
	query := db.Model(&entity.Poll{}).Where("room_id = ?", roomID)

	if status != "" && status != "all" {
		query = query.Where("status = ?", status)
	}

	err := query.Count(&count).Error
	return count, err
}

// Close close poll
func (r *PollRepository) Close(db *gorm.DB, pollID uint) error {
	now := time.Now()
	return db.Model(&entity.Poll{}).Where("id = ?", pollID).
		Updates(map[string]interface{}{
			"status":    "closed",
			"closed_at": now,
		}).Error
}

// Activate activate poll (set status to active)
func (r *PollRepository) Activate(db *gorm.DB, pollID uint) error {
	now := time.Now()
	return db.Model(&entity.Poll{}).Where("id = ?", pollID).
		Updates(map[string]interface{}{
			"status":       "active",
			"activated_at": now,
		}).Error
}

// GetRoomIDByPollID get room id dari poll
func (r *PollRepository) GetRoomIDByPollID(db *gorm.DB, pollID uint) (uint, error) {
	var poll entity.Poll
	err := db.Select("room_id").Where("id = ?", pollID).First(&poll).Error
	if err != nil {
		return 0, err
	}
	return poll.RoomID, nil
}
