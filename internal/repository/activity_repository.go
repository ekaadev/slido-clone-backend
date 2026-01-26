package repository

import (
	"slido-clone-backend/internal/entity"
	"slido-clone-backend/internal/model"
	"time"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type ActivityRepository struct {
	Log *logrus.Logger
}

// NewActivityRepository create new instance of ActivityRepository
func NewActivityRepository(log *logrus.Logger) *ActivityRepository {
	return &ActivityRepository{
		Log: log,
	}
}

// GetTimelineRaw mendapatkan timeline items dengan UNION query dari messages, questions, polls
// Returns raw items (type, id, created_at) yang perlu di-enrich dengan data lengkap
func (r *ActivityRepository) GetTimelineRaw(db *gorm.DB, roomID uint, before *time.Time, after *time.Time, limit int) ([]model.TimelineRawItem, error) {
	var items []model.TimelineRawItem

	// Build UNION ALL query
	query := `
		SELECT 'message' as type, id, created_at FROM messages WHERE room_id = ? AND deleted_at IS NULL
		UNION ALL
		SELECT 'question' as type, id, created_at FROM questions WHERE room_id = ?
		UNION ALL
		SELECT 'poll' as type, id, created_at FROM polls WHERE room_id = ?
	`
	args := []interface{}{roomID, roomID, roomID}

	// Wrap dengan subquery untuk filtering dan ordering
	wrapQuery := "SELECT type, id, created_at FROM (" + query + ") AS timeline"

	// Add cursor conditions
	if before != nil {
		wrapQuery += " WHERE created_at < ?"
		args = append(args, *before)
	} else if after != nil {
		wrapQuery += " WHERE created_at > ?"
		args = append(args, *after)
	}

	// Order dan limit
	if after != nil {
		// Jika load newer, order ASC dulu lalu reverse di code
		wrapQuery += " ORDER BY created_at ASC LIMIT ?"
	} else {
		// Default atau load older: order DESC
		wrapQuery += " ORDER BY created_at DESC LIMIT ?"
	}
	args = append(args, limit+1) // +1 untuk check has_more

	err := db.Raw(wrapQuery, args...).Scan(&items).Error
	if err != nil {
		r.Log.WithField("error", err).Error("GetTimelineRaw - failed to query")
		return nil, err
	}

	return items, nil
}

// GetMessagesByIDs mendapatkan messages berdasarkan IDs
func (r *ActivityRepository) GetMessagesByIDs(db *gorm.DB, ids []uint) (map[uint]entity.Message, error) {
	var messages []entity.Message
	err := db.Preload("Participant").Where("id IN ?", ids).Find(&messages).Error
	if err != nil {
		return nil, err
	}

	result := make(map[uint]entity.Message)
	for _, m := range messages {
		result[m.ID] = m
	}
	return result, nil
}

// GetQuestionsByIDs mendapatkan questions berdasarkan IDs
func (r *ActivityRepository) GetQuestionsByIDs(db *gorm.DB, ids []uint) (map[uint]entity.Question, error) {
	var questions []entity.Question
	err := db.Preload("Participant").Where("id IN ?", ids).Find(&questions).Error
	if err != nil {
		return nil, err
	}

	result := make(map[uint]entity.Question)
	for _, q := range questions {
		result[q.ID] = q
	}
	return result, nil
}

// GetPollsByIDs mendapatkan polls berdasarkan IDs
func (r *ActivityRepository) GetPollsByIDs(db *gorm.DB, ids []uint) (map[uint]entity.Poll, error) {
	var polls []entity.Poll
	err := db.Preload("Options").Where("id IN ?", ids).Find(&polls).Error
	if err != nil {
		return nil, err
	}

	result := make(map[uint]entity.Poll)
	for _, p := range polls {
		result[p.ID] = p
	}
	return result, nil
}
