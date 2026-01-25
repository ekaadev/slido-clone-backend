package repository

import (
	"errors"
	"slido-clone-backend/internal/entity"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

// VoteRepository repository untuk operasi database Vote
type VoteRepository struct {
	Repository[entity.Vote]
	Log *logrus.Logger
}

// NewVoteRepository create new instance of VoteRepository
func NewVoteRepository(log *logrus.Logger) *VoteRepository {
	return &VoteRepository{
		Log: log,
	}
}

// FindByQuestionAndParticipant find vote berdasarkan question dan participant
func (r *VoteRepository) FindByQuestionAndParticipant(db *gorm.DB, questionID uint, participantID uint) (*entity.Vote, error) {
	var vote entity.Vote
	err := db.Where("question_id = ? AND participant_id = ?", questionID, participantID).First(&vote).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &vote, err
}

// HasVoted check apakah participant sudah vote question
func (r *VoteRepository) HasVoted(db *gorm.DB, questionID uint, participantID uint) (bool, error) {
	var count int64
	err := db.Model(&entity.Vote{}).Where("question_id = ? AND participant_id = ?", questionID, participantID).Count(&count).Error
	return count > 0, err
}

// GetVotedQuestionIDs mendapatkan list question IDs yang sudah di-vote oleh participant
func (r *VoteRepository) GetVotedQuestionIDs(db *gorm.DB, participantID uint, questionIDs []uint) (map[uint]bool, error) {
	var votes []entity.Vote
	err := db.Select("question_id").Where("participant_id = ? AND question_id IN ?", participantID, questionIDs).Find(&votes).Error
	if err != nil {
		return nil, err
	}

	result := make(map[uint]bool)
	for _, v := range votes {
		result[v.QuestionID] = true
	}
	return result, nil
}

// DeleteByQuestionAndParticipant delete vote berdasarkan question dan participant
func (r *VoteRepository) DeleteByQuestionAndParticipant(db *gorm.DB, questionID uint, participantID uint) error {
	return db.Where("question_id = ? AND participant_id = ?", questionID, participantID).Delete(&entity.Vote{}).Error
}
