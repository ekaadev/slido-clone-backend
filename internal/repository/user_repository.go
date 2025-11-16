package repository

import (
	"errors"
	"slido-clone-backend/internal/entity"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type UserRepository struct {
	Repository[entity.User]
	Log *logrus.Logger
}

// NewUserRepository create new instance of UserRepository
func NewUserRepository(log *logrus.Logger) *UserRepository {
	return &UserRepository{
		Log: log,
	}
}

// FindByEmailOrUsername more efficient to find user by email or username
func (r *UserRepository) FindByEmailOrUsername(db *gorm.DB, email string, username string) (*entity.User, error) {
	var user entity.User
	err := db.Where("email = ? OR username = ?", email, username).First(&user).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}

	return &user, err
}

// FindByUsername find user by username
func (r *UserRepository) FindByUsername(db *gorm.DB, username string) (*entity.User, error) {
	var user entity.User
	err := db.Where("username = ?", username).First(&user).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}

	return &user, err
}
