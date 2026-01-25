package mocks

import (
	"slido-clone-backend/internal/entity"

	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

// MockQuestionRepository mock untuk QuestionRepository
type MockQuestionRepository struct {
	mock.Mock
}

func (m *MockQuestionRepository) Create(db *gorm.DB, entity *entity.Question) error {
	args := m.Called(db, entity)
	return args.Error(0)
}

func (m *MockQuestionRepository) Update(db *gorm.DB, entity *entity.Question) error {
	args := m.Called(db, entity)
	return args.Error(0)
}

func (m *MockQuestionRepository) Delete(db *gorm.DB, entity *entity.Question) error {
	args := m.Called(db, entity)
	return args.Error(0)
}

func (m *MockQuestionRepository) CountById(db *gorm.DB, id any) (int64, error) {
	args := m.Called(db, id)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockQuestionRepository) FindById(db *gorm.DB, entity *entity.Question, id any) error {
	args := m.Called(db, entity, id)
	return args.Error(0)
}

func (m *MockQuestionRepository) FindByIdWithParticipant(db *gorm.DB, id uint) (*entity.Question, error) {
	args := m.Called(db, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.Question), args.Error(1)
}

func (m *MockQuestionRepository) List(db *gorm.DB, roomID uint, status string, sortBy string, limit int, offset int) ([]entity.Question, error) {
	args := m.Called(db, roomID, status, sortBy, limit, offset)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]entity.Question), args.Error(1)
}

func (m *MockQuestionRepository) Count(db *gorm.DB, roomID uint, status string) (int64, error) {
	args := m.Called(db, roomID, status)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockQuestionRepository) IncrementUpvoteCount(db *gorm.DB, questionID uint) error {
	args := m.Called(db, questionID)
	return args.Error(0)
}

func (m *MockQuestionRepository) DecrementUpvoteCount(db *gorm.DB, questionID uint) error {
	args := m.Called(db, questionID)
	return args.Error(0)
}

func (m *MockQuestionRepository) UpdateValidation(db *gorm.DB, questionID uint, status string, xpAwarded int) error {
	args := m.Called(db, questionID, status, xpAwarded)
	return args.Error(0)
}

func (m *MockQuestionRepository) GetRoomIDByQuestionID(db *gorm.DB, questionID uint) (uint, error) {
	args := m.Called(db, questionID)
	return args.Get(0).(uint), args.Error(1)
}

// MockVoteRepository mock untuk VoteRepository
type MockVoteRepository struct {
	mock.Mock
}

func (m *MockVoteRepository) Create(db *gorm.DB, entity *entity.Vote) error {
	args := m.Called(db, entity)
	return args.Error(0)
}

func (m *MockVoteRepository) Update(db *gorm.DB, entity *entity.Vote) error {
	args := m.Called(db, entity)
	return args.Error(0)
}

func (m *MockVoteRepository) Delete(db *gorm.DB, entity *entity.Vote) error {
	args := m.Called(db, entity)
	return args.Error(0)
}

func (m *MockVoteRepository) CountById(db *gorm.DB, id any) (int64, error) {
	args := m.Called(db, id)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockVoteRepository) FindById(db *gorm.DB, entity *entity.Vote, id any) error {
	args := m.Called(db, entity, id)
	return args.Error(0)
}

func (m *MockVoteRepository) FindByQuestionAndParticipant(db *gorm.DB, questionID uint, participantID uint) (*entity.Vote, error) {
	args := m.Called(db, questionID, participantID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.Vote), args.Error(1)
}

func (m *MockVoteRepository) HasVoted(db *gorm.DB, questionID uint, participantID uint) (bool, error) {
	args := m.Called(db, questionID, participantID)
	return args.Bool(0), args.Error(1)
}

func (m *MockVoteRepository) GetVotedQuestionIDs(db *gorm.DB, participantID uint, questionIDs []uint) (map[uint]bool, error) {
	args := m.Called(db, participantID, questionIDs)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(map[uint]bool), args.Error(1)
}

func (m *MockVoteRepository) DeleteByQuestionAndParticipant(db *gorm.DB, questionID uint, participantID uint) error {
	args := m.Called(db, questionID, participantID)
	return args.Error(0)
}
