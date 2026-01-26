package repository

import (
	"errors"
	"slido-clone-backend/internal/entity"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

// PollResponseRepository repository untuk operasi database PollResponse
type PollResponseRepository struct {
	Repository[entity.PollResponse]
	Log *logrus.Logger
}

// NewPollResponseRepository create new instance of PollResponseRepository
func NewPollResponseRepository(log *logrus.Logger) *PollResponseRepository {
	return &PollResponseRepository{
		Log: log,
	}
}

// HasVoted check apakah participant sudah vote pada poll
func (r *PollResponseRepository) HasVoted(db *gorm.DB, pollID uint, participantID uint) (bool, error) {
	var count int64
	err := db.Model(&entity.PollResponse{}).
		Where("poll_id = ? AND participant_id = ?", pollID, participantID).
		Count(&count).Error
	return count > 0, err
}

// FindByPollAndParticipant find response by poll and participant
func (r *PollResponseRepository) FindByPollAndParticipant(db *gorm.DB, pollID uint, participantID uint) (*entity.PollResponse, error) {
	var response entity.PollResponse
	err := db.Where("poll_id = ? AND participant_id = ?", pollID, participantID).First(&response).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &response, err
}

// CountByPollID count total responses untuk poll
func (r *PollResponseRepository) CountByPollID(db *gorm.DB, pollID uint) (int64, error) {
	var count int64
	err := db.Model(&entity.PollResponse{}).Where("poll_id = ?", pollID).Count(&count).Error
	return count, err
}

// GetVotedOptionByParticipant get option yang sudah di-vote oleh participant
func (r *PollResponseRepository) GetVotedOptionByParticipant(db *gorm.DB, pollID uint, participantID uint) (*uint, error) {
	var response entity.PollResponse
	err := db.Select("poll_option_id").
		Where("poll_id = ? AND participant_id = ?", pollID, participantID).
		First(&response).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &response.PollOptionID, nil
}
