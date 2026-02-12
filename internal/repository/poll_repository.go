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

// GetActivePollsByRoomID retrieves all active polls for a given room ID with options
func (r *PollRepository) GetActivePollsByRoomID(db *gorm.DB, roomID uint) ([]entity.Poll, error) {
	var polls []entity.Poll
	err := db.Where("room_id = ? AND status = ?", roomID, "active").
		Preload("Options", func(db *gorm.DB) *gorm.DB {
			return db.Order("poll_options.`order` ASC")
		}).
		Order("created_at DESC").
		Find(&polls).Error
	return polls, err
}

// GetPollsByRoomID retrieves polls by room ID with optional status filter
func (r *PollRepository) GetPollsByRoomID(db *gorm.DB, roomID uint, status string, limit int) ([]entity.Poll, int64, error) {
	var polls []entity.Poll
	var total int64

	query := db.Model(&entity.Poll{}).Where("room_id = ?", roomID)

	if status != "" && status != "all" {
		query = query.Where("status = ?", status)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := query.
		Preload("Options", func(db *gorm.DB) *gorm.DB {
			return db.Order("poll_options.`order` ASC")
		}).
		Order("created_at DESC").
		Limit(limit).
		Find(&polls).Error

	return polls, total, err
}

// GetPollByIDWithOptions retrieves a poll by ID with its options
func (r *PollRepository) GetPollByIDWithOptions(db *gorm.DB, pollID uint) (*entity.Poll, error) {
	var poll entity.Poll
	err := db.Where("id = ?", pollID).
		Preload("Options", func(db *gorm.DB) *gorm.DB {
			return db.Order("poll_options.`order` ASC")
		}).
		First(&poll).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &poll, err
}

// CreatePollWithOptions creates a poll with its options in a single transaction
func (r *PollRepository) CreatePollWithOptions(db *gorm.DB, poll *entity.Poll, options []entity.PollOption) error {
	if err := db.Create(poll).Error; err != nil {
		return err
	}

	for i := range options {
		options[i].PollID = poll.ID
		options[i].Order = i + 1
	}

	if err := db.Create(&options).Error; err != nil {
		return err
	}

	poll.Options = options
	return nil
}

// GetPollOptionByID retrieves a poll option by ID
func (r *PollRepository) GetPollOptionByID(db *gorm.DB, optionID uint) (*entity.PollOption, error) {
	var option entity.PollOption
	err := db.Where("id = ?", optionID).First(&option).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &option, err
}

// GetPollResponseByParticipant retrieves a participant's response for a poll
func (r *PollRepository) GetPollResponseByParticipant(db *gorm.DB, pollID, participantID uint) (*entity.PollResponse, error) {
	var response entity.PollResponse
	err := db.Where("poll_id = ? AND participant_id = ?", pollID, participantID).First(&response).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &response, err
}

// CreatePollResponse creates a new poll response (vote)
func (r *PollRepository) CreatePollResponse(db *gorm.DB, response *entity.PollResponse) error {
	return db.Create(response).Error
}

// IncrementOptionVoteCount increments the vote count for a poll option
func (r *PollRepository) IncrementOptionVoteCount(db *gorm.DB, optionID uint) error {
	return db.Model(&entity.PollOption{}).
		Where("id = ?", optionID).
		UpdateColumn("vote_count", gorm.Expr("vote_count + ?", 1)).Error
}

// GetTotalVotesByPollID returns the total number of votes for a poll
func (r *PollRepository) GetTotalVotesByPollID(db *gorm.DB, pollID uint) (int, error) {
	var totalVotes int64
	err := db.Model(&entity.PollResponse{}).Where("poll_id = ?", pollID).Count(&totalVotes).Error
	return int(totalVotes), err
}

// ClosePoll updates poll status to closed
func (r *PollRepository) ClosePoll(db *gorm.DB, poll *entity.Poll) error {
	return db.Model(poll).Updates(map[string]interface{}{
		"status":    "closed",
		"closed_at": gorm.Expr("NOW()"),
	}).Error
}
