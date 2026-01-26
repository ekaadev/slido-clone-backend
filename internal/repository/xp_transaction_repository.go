package repository

import (
	"slido-clone-backend/internal/entity"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type XPTransactionRepository struct {
	Repository[entity.XPTransaction]
	Log *logrus.Logger
}

func NewXPTransactionRepository(log *logrus.Logger) *XPTransactionRepository {
	return &XPTransactionRepository{
		Log: log,
	}
}

func (r *XPTransactionRepository) AddXP(tx *gorm.DB, participantID uint, points int) error {
	return tx.Model(&entity.Participant{}).Where("id = ?", participantID).Update("xp_score", gorm.Expr("xp_score + ?", points)).Error
}

// FindByRoomAndParticipant get XP transactions untuk participant di room
func (r *XPTransactionRepository) FindByRoomAndParticipant(db *gorm.DB, roomID uint, participantID uint, limit int) ([]entity.XPTransaction, error) {
	var transactions []entity.XPTransaction
	query := db.Where("room_id = ? AND participant_id = ?", roomID, participantID).
		Order("created_at DESC")
	if limit > 0 {
		query = query.Limit(limit)
	}
	err := query.Find(&transactions).Error
	return transactions, err
}

// CountByRoomAndParticipant count total XP transactions
func (r *XPTransactionRepository) CountByRoomAndParticipant(db *gorm.DB, roomID uint, participantID uint) (int64, error) {
	var count int64
	err := db.Model(&entity.XPTransaction{}).
		Where("room_id = ? AND participant_id = ?", roomID, participantID).
		Count(&count).Error
	return count, err
}

// GetTotalXPByParticipant get total XP score dari participant
func (r *XPTransactionRepository) GetTotalXPByParticipant(db *gorm.DB, participantID uint) (int, error) {
	var participant entity.Participant
	err := db.Select("xp_score").Where("id = ?", participantID).First(&participant).Error
	if err != nil {
		return 0, err
	}
	return participant.XPScore, nil
}
