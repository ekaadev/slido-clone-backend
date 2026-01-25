package mocks

import (
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

// MockXPTransactionUseCase mock untuk XPTransactionUseCase
type MockXPTransactionUseCase struct {
	mock.Mock
}

func (m *MockXPTransactionUseCase) AddXPForMessage(db *gorm.DB, roomID uint, participantID uint, messageID uint) error {
	args := m.Called(db, roomID, participantID, messageID)
	return args.Error(0)
}
