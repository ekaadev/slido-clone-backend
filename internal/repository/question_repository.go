package repository

import (
	"errors"
	"slido-clone-backend/internal/entity"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

// QuestionRepository repository untuk operasi database Question
type QuestionRepository struct {
	Repository[entity.Question]
	Log *logrus.Logger
}

// NewQuestionRepository create new instance of QuestionRepository
func NewQuestionRepository(log *logrus.Logger) *QuestionRepository {
	return &QuestionRepository{
		Log: log,
	}
}

// FindByIdWithParticipant find question by id dengan preload participant
func (r *QuestionRepository) FindByIdWithParticipant(db *gorm.DB, id uint) (*entity.Question, error) {
	var question entity.Question
	err := db.Preload("Participant").Where("id = ?", id).First(&question).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &question, err
}

// List mendapatkan list questions dengan filter dan sorting
func (r *QuestionRepository) List(db *gorm.DB, roomID uint, status string, sortBy string, limit int, offset int) ([]entity.Question, error) {
	var questions []entity.Question

	query := db.Preload("Participant").Where("room_id = ?", roomID)

	// filter by status
	if status != "" {
		query = query.Where("status = ?", status)
	}

	// sorting
	switch sortBy {
	case "recent":
		query = query.Order("created_at DESC")
	case "validated":
		query = query.Order("is_validated_by_presenter DESC, upvote_count DESC, created_at DESC")
	default: // upvotes (default)
		query = query.Order("upvote_count DESC, created_at DESC")
	}

	// pagination
	if limit > 0 {
		query = query.Limit(limit)
	}
	if offset > 0 {
		query = query.Offset(offset)
	}

	err := query.Find(&questions).Error
	return questions, err
}

// Count menghitung total questions berdasarkan room dan status
func (r *QuestionRepository) Count(db *gorm.DB, roomID uint, status string) (int64, error) {
	var count int64
	query := db.Model(&entity.Question{}).Where("room_id = ?", roomID)

	if status != "" {
		query = query.Where("status = ?", status)
	}

	err := query.Count(&count).Error
	return count, err
}

// IncrementUpvoteCount increment upvote count
func (r *QuestionRepository) IncrementUpvoteCount(db *gorm.DB, questionID uint) error {
	return db.Model(&entity.Question{}).Where("id = ?", questionID).
		UpdateColumn("upvote_count", gorm.Expr("upvote_count + ?", 1)).Error
}

// DecrementUpvoteCount decrement upvote count
func (r *QuestionRepository) DecrementUpvoteCount(db *gorm.DB, questionID uint) error {
	return db.Model(&entity.Question{}).Where("id = ?", questionID).
		UpdateColumn("upvote_count", gorm.Expr("upvote_count - ?", 1)).Error
}

// UpdateValidation update validation status
func (r *QuestionRepository) UpdateValidation(db *gorm.DB, questionID uint, status string, xpAwarded int) error {
	return db.Model(&entity.Question{}).Where("id = ?", questionID).
		Updates(map[string]interface{}{
			"status":                    status,
			"is_validated_by_presenter": true,
			"xp_awarded":                gorm.Expr("xp_awarded + ?", xpAwarded),
		}).Error
}

// GetRoomIDByQuestionID get room id dari question
func (r *QuestionRepository) GetRoomIDByQuestionID(db *gorm.DB, questionID uint) (uint, error) {
	var question entity.Question
	err := db.Select("room_id").Where("id = ?", questionID).First(&question).Error
	if err != nil {
		return 0, err
	}
	return question.RoomID, nil
}
