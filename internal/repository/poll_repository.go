package repository

import (
	"errors"
	"slido-clone-backend/internal/entity"

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
