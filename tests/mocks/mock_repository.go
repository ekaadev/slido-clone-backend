package mocks

import (
	"slido-clone-backend/internal/entity"

	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

// MockUserRepository mock untuk UserRepository
type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) Create(db *gorm.DB, entity *entity.User) error {
	args := m.Called(db, entity)
	return args.Error(0)
}

func (m *MockUserRepository) Update(db *gorm.DB, entity *entity.User) error {
	args := m.Called(db, entity)
	return args.Error(0)
}

func (m *MockUserRepository) Delete(db *gorm.DB, entity *entity.User) error {
	args := m.Called(db, entity)
	return args.Error(0)
}

func (m *MockUserRepository) CountById(db *gorm.DB, id any) (int64, error) {
	args := m.Called(db, id)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockUserRepository) FindById(db *gorm.DB, entity *entity.User, id any) error {
	args := m.Called(db, entity, id)
	return args.Error(0)
}

func (m *MockUserRepository) FindByEmailOrUsername(db *gorm.DB, email string, username string) (*entity.User, error) {
	args := m.Called(db, email, username)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.User), args.Error(1)
}

func (m *MockUserRepository) FindByUsername(db *gorm.DB, username string) (*entity.User, error) {
	args := m.Called(db, username)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.User), args.Error(1)
}

// MockRoomRepository mock untuk RoomRepository
type MockRoomRepository struct {
	mock.Mock
}

func (m *MockRoomRepository) Create(db *gorm.DB, entity *entity.Room) error {
	args := m.Called(db, entity)
	return args.Error(0)
}

func (m *MockRoomRepository) Update(db *gorm.DB, entity *entity.Room) error {
	args := m.Called(db, entity)
	return args.Error(0)
}

func (m *MockRoomRepository) Delete(db *gorm.DB, entity *entity.Room) error {
	args := m.Called(db, entity)
	return args.Error(0)
}

func (m *MockRoomRepository) CountById(db *gorm.DB, id any) (int64, error) {
	args := m.Called(db, id)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockRoomRepository) FindById(db *gorm.DB, entity *entity.Room, id any) error {
	args := m.Called(db, entity, id)
	return args.Error(0)
}

func (m *MockRoomRepository) FindByRoomCode(db *gorm.DB, code string) (*entity.Room, error) {
	args := m.Called(db, code)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.Room), args.Error(1)
}

func (m *MockRoomRepository) FindByRoomCodeAndPresenterId(db *gorm.DB, code string, presenterId uint) (*entity.Room, error) {
	args := m.Called(db, code, presenterId)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.Room), args.Error(1)
}

func (m *MockRoomRepository) FindByIdAndPresenterId(db *gorm.DB, id uint, presenterId uint) (*entity.Room, error) {
	args := m.Called(db, id, presenterId)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.Room), args.Error(1)
}

func (m *MockRoomRepository) Search(db *gorm.DB, presenterId uint) ([]entity.Room, error) {
	args := m.Called(db, presenterId)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]entity.Room), args.Error(1)
}

// MockParticipantRepository mock untuk ParticipantRepository
type MockParticipantRepository struct {
	mock.Mock
}

func (m *MockParticipantRepository) Create(db *gorm.DB, entity *entity.Participant) error {
	args := m.Called(db, entity)
	return args.Error(0)
}

func (m *MockParticipantRepository) Update(db *gorm.DB, entity *entity.Participant) error {
	args := m.Called(db, entity)
	return args.Error(0)
}

func (m *MockParticipantRepository) Delete(db *gorm.DB, entity *entity.Participant) error {
	args := m.Called(db, entity)
	return args.Error(0)
}

func (m *MockParticipantRepository) CountById(db *gorm.DB, id any) (int64, error) {
	args := m.Called(db, id)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockParticipantRepository) FindById(db *gorm.DB, entity *entity.Participant, id any) error {
	args := m.Called(db, entity, id)
	return args.Error(0)
}

func (m *MockParticipantRepository) FindByRoomIdAndUserId(db *gorm.DB, roomId uint, userId uint) (*entity.Participant, error) {
	args := m.Called(db, roomId, userId)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.Participant), args.Error(1)
}

func (m *MockParticipantRepository) FindByRoomIdAndUserIdWithRoom(db *gorm.DB, roomId uint, userId uint) (*entity.Participant, error) {
	args := m.Called(db, roomId, userId)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.Participant), args.Error(1)
}

func (m *MockParticipantRepository) FindByRoomCodeAndUsername(db *gorm.DB, roomCode string, username string) (*entity.Participant, error) {
	args := m.Called(db, roomCode, username)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.Participant), args.Error(1)
}

func (m *MockParticipantRepository) List(db *gorm.DB, roomId uint, offset int, limit int) ([]entity.Participant, error) {
	args := m.Called(db, roomId, offset, limit)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]entity.Participant), args.Error(1)
}

func (m *MockParticipantRepository) Count(db *gorm.DB, roomId uint) (int64, error) {
	args := m.Called(db, roomId)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockParticipantRepository) ListByXPScore(db *gorm.DB, roomId uint) ([]entity.Participant, error) {
	args := m.Called(db, roomId)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]entity.Participant), args.Error(1)
}

func (m *MockParticipantRepository) GetRankAndScore(db *gorm.DB, roomId uint, participantId uint) (int, int, error) {
	args := m.Called(db, roomId, participantId)
	return args.Int(0), args.Int(1), args.Error(2)
}

// MockMessageRepository mock untuk MessageRepository
type MockMessageRepository struct {
	mock.Mock
}

func (m *MockMessageRepository) Create(db *gorm.DB, entity *entity.Message) error {
	args := m.Called(db, entity)
	return args.Error(0)
}

func (m *MockMessageRepository) Update(db *gorm.DB, entity *entity.Message) error {
	args := m.Called(db, entity)
	return args.Error(0)
}

func (m *MockMessageRepository) Delete(db *gorm.DB, entity *entity.Message) error {
	args := m.Called(db, entity)
	return args.Error(0)
}

func (m *MockMessageRepository) CountById(db *gorm.DB, id any) (int64, error) {
	args := m.Called(db, id)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockMessageRepository) FindById(db *gorm.DB, entity *entity.Message, id any) error {
	args := m.Called(db, entity, id)
	return args.Error(0)
}

func (m *MockMessageRepository) List(db *gorm.DB, roomID uint, limit int, before *int64) ([]entity.Message, error) {
	args := m.Called(db, roomID, limit, before)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]entity.Message), args.Error(1)
}
