package repository

import (
	"slido-clone-backend/internal/entity"

	"github.com/sirupsen/logrus"
)

type RoomRepository struct {
	Repository[entity.Room]
	Log *logrus.Logger
}

// NewRoomRepository create new instance of RoomRepository
func NewRoomRepository(log *logrus.Logger) *RoomRepository {
	return &RoomRepository{
		Log: log,
	}
}
