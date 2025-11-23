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
