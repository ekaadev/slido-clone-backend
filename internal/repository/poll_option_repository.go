package repository

import (
	"slido-clone-backend/internal/entity"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

// PollOptionRepository repository untuk operasi database PollOption
type PollOptionRepository struct {
	Repository[entity.PollOption]
	Log *logrus.Logger
}

// NewPollOptionRepository create new instance of PollOptionRepository
func NewPollOptionRepository(log *logrus.Logger) *PollOptionRepository {
	return &PollOptionRepository{
		Log: log,
	}
}

// CreateBatch create multiple poll options at once
func (r *PollOptionRepository) CreateBatch(db *gorm.DB, options []entity.PollOption) error {
	return db.Create(&options).Error
}

// IncrementVoteCount increment vote count untuk option
func (r *PollOptionRepository) IncrementVoteCount(db *gorm.DB, optionID uint) error {
	return db.Model(&entity.PollOption{}).Where("id = ?", optionID).
		Update("vote_count", gorm.Expr("vote_count + 1")).Error
}

// GetByPollID get all options for a poll
func (r *PollOptionRepository) GetByPollID(db *gorm.DB, pollID uint) ([]entity.PollOption, error) {
	var options []entity.PollOption
	err := db.Where("poll_id = ?", pollID).Order("poll_options.order ASC").Find(&options).Error
	return options, err
}

// GetTotalVotesByPollID get total votes untuk poll
func (r *PollOptionRepository) GetTotalVotesByPollID(db *gorm.DB, pollID uint) (int, error) {
	var total int64
	err := db.Model(&entity.PollOption{}).Where("poll_id = ?", pollID).
		Select("COALESCE(SUM(vote_count), 0)").Scan(&total).Error
	return int(total), err
}

// ValidateOptionBelongsToPoll check if option belongs to poll
func (r *PollOptionRepository) ValidateOptionBelongsToPoll(db *gorm.DB, optionID uint, pollID uint) (bool, error) {
	var count int64
	err := db.Model(&entity.PollOption{}).Where("id = ? AND poll_id = ?", optionID, pollID).Count(&count).Error
	return count > 0, err
}
