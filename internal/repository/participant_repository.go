package repository

import (
	"slido-clone-backend/internal/entity"

	"github.com/sirupsen/logrus"
)

type ParticipantRepository struct {
	Repository[entity.Participant]
	Log *logrus.Logger
}

func NewParticipantRepository(log *logrus.Logger) *ParticipantRepository {
	return &ParticipantRepository{
		Log: log,
	}
}
