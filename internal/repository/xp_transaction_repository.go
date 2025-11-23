package repository

import (
	"slido-clone-backend/internal/entity"

	"github.com/sirupsen/logrus"
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
