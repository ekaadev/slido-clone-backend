package tests

import (
	"context"
	"slido-clone-backend/internal/entity"
	"slido-clone-backend/internal/model"
	"slido-clone-backend/internal/repository"
	"slido-clone-backend/internal/usecase"
	"slido-clone-backend/internal/util"
	"slido-clone-backend/tests/mocks"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/go-playground/validator/v10"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// setupUserUseCaseTest setup test environment for UserUseCase
func setupUserUseCaseTest(t *testing.T) (*usecase.UserUseCase, sqlmock.Sqlmock, *mocks.MockUserRepository, *mocks.MockParticipantRepository, *mocks.MockRoomRepository, *mocks.MockTokenUtil) {
	// create sqlmock
	db, mockDB, err := sqlmock.New()
	assert.NoError(t, err)

	// create gorm db
	dialector := mysql.New(mysql.Config{
		Conn:                      db,
		SkipInitializeWithVersion: true,
	})
	gormDB, err := gorm.Open(dialector, &gorm.Config{})
	assert.NoError(t, err)

	// create mocks
	mockUserRepo := new(mocks.MockUserRepository)
	mockParticipantRepo := new(mocks.MockParticipantRepository)
	mockRoomRepo := new(mocks.MockRoomRepository)
	mockTokenUtil := new(mocks.MockTokenUtil)

	// create real validator
	validate := validator.New()

	// create logger
	log := logrus.New()
	log.SetLevel(logrus.PanicLevel) // suppress logs during testing

	// create usecase with mocks - note: we need to adapt this
	uc := &usecase.UserUseCase{
		DB:                    gormDB,
		Log:                   log,
		Validate:              validate,
		UserRepository:        &repository.UserRepository{Log: log},
		ParticipantRepository: &repository.ParticipantRepository{Log: log},
		RoomRepository:        &repository.RoomRepository{Log: log},
		TokenUtil:             &util.TokenUtil{SecretKey: "test-secret"},
	}

	return uc, mockDB, mockUserRepo, mockParticipantRepo, mockRoomRepo, mockTokenUtil
}

// TestUserUseCase_Create_InvalidRequest test create user with invalid request
func TestUserUseCase_Create_InvalidRequest(t *testing.T) {
	uc, mockDB, _, _, _, _ := setupUserUseCaseTest(t)

	// expect begin transaction
	mockDB.ExpectBegin()
	mockDB.ExpectRollback()

	// test with empty request
	request := &model.RegisterUserRequest{
		Username: "", // invalid: required
		Email:    "test@example.com",
		Password: "password123",
		Role:     "presenter",
	}

	result, err := uc.Create(context.Background(), request)

	assert.Nil(t, result)
	assert.Error(t, err)
}

// TestUserUseCase_Login_InvalidRequest test login with invalid request
func TestUserUseCase_Login_InvalidRequest(t *testing.T) {
	uc, mockDB, _, _, _, _ := setupUserUseCaseTest(t)

	// expect begin transaction
	mockDB.ExpectBegin()
	mockDB.ExpectRollback()

	// test with empty request
	request := &model.LoginUserRequest{
		Username: "", // invalid: required
		Password: "password123",
	}

	result, err := uc.Login(context.Background(), request)

	assert.Nil(t, result)
	assert.Error(t, err)
}

// TestUserUseCase_Anon_InvalidRequest test anonymous join with invalid request
func TestUserUseCase_Anon_InvalidRequest(t *testing.T) {
	uc, mockDB, _, _, _, _ := setupUserUseCaseTest(t)

	// expect begin transaction
	mockDB.ExpectBegin()
	mockDB.ExpectRollback()

	// test with invalid room code (not 6 characters)
	request := &model.AnonymousUserRequest{
		RoomCode:    "ABC", // invalid: must be 6 characters
		DisplayName: "Test User",
	}

	result, err := uc.Anon(context.Background(), request)

	assert.Nil(t, result)
	assert.Error(t, err)
}

// TestUserUseCase_Anon_DisplayNameRequired test anonymous join without display name
func TestUserUseCase_Anon_DisplayNameRequired(t *testing.T) {
	uc, mockDB, _, _, _, _ := setupUserUseCaseTest(t)

	// expect begin transaction
	mockDB.ExpectBegin()
	mockDB.ExpectRollback()

	// test without display name
	request := &model.AnonymousUserRequest{
		RoomCode:    "ABC123",
		DisplayName: "", // invalid: required
	}

	result, err := uc.Anon(context.Background(), request)

	assert.Nil(t, result)
	assert.Error(t, err)
}

// setupUserUseCaseMockTest setup test with mock repositories
func setupUserUseCaseMockTest(t *testing.T) (*gorm.DB, sqlmock.Sqlmock, *mocks.MockUserRepository, *mocks.MockParticipantRepository, *mocks.MockRoomRepository, *mocks.MockTokenUtil, *validator.Validate, *logrus.Logger) {
	// create sqlmock
	db, mockDB, err := sqlmock.New()
	assert.NoError(t, err)

	// create gorm db
	dialector := mysql.New(mysql.Config{
		Conn:                      db,
		SkipInitializeWithVersion: true,
	})
	gormDB, err := gorm.Open(dialector, &gorm.Config{})
	assert.NoError(t, err)

	// create mocks
	mockUserRepo := new(mocks.MockUserRepository)
	mockParticipantRepo := new(mocks.MockParticipantRepository)
	mockRoomRepo := new(mocks.MockRoomRepository)
	mockTokenUtil := new(mocks.MockTokenUtil)

	// create real validator
	validate := validator.New()

	// create logger
	log := logrus.New()
	log.SetLevel(logrus.PanicLevel)

	return gormDB, mockDB, mockUserRepo, mockParticipantRepo, mockRoomRepo, mockTokenUtil, validate, log
}

// TestRoomCode_Validation test room code validation
func TestRoomCode_Validation(t *testing.T) {
	validate := validator.New()

	tests := []struct {
		name        string
		request     model.AnonymousUserRequest
		shouldError bool
	}{
		{
			name: "valid room code",
			request: model.AnonymousUserRequest{
				RoomCode:    "ABC123",
				DisplayName: "Test User",
			},
			shouldError: false,
		},
		{
			name: "room code too short",
			request: model.AnonymousUserRequest{
				RoomCode:    "ABC",
				DisplayName: "Test User",
			},
			shouldError: true,
		},
		{
			name: "room code too long",
			request: model.AnonymousUserRequest{
				RoomCode:    "ABC12345",
				DisplayName: "Test User",
			},
			shouldError: true,
		},
		{
			name: "display name too short",
			request: model.AnonymousUserRequest{
				RoomCode:    "ABC123",
				DisplayName: "AB", // min 3
			},
			shouldError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validate.Struct(tt.request)
			if tt.shouldError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

// TestRegisterRequest_Validation test register request validation
func TestRegisterRequest_Validation(t *testing.T) {
	validate := validator.New()

	tests := []struct {
		name        string
		request     model.RegisterUserRequest
		shouldError bool
	}{
		{
			name: "valid request",
			request: model.RegisterUserRequest{
				Username: "testuser",
				Email:    "test@example.com",
				Password: "password123",
				Role:     "presenter",
			},
			shouldError: false,
		},
		{
			name: "username too short",
			request: model.RegisterUserRequest{
				Username: "ab", // min 3
				Email:    "test@example.com",
				Password: "password123",
				Role:     "presenter",
			},
			shouldError: true,
		},
		{
			name: "invalid email",
			request: model.RegisterUserRequest{
				Username: "testuser",
				Email:    "invalid-email",
				Password: "password123",
				Role:     "presenter",
			},
			shouldError: true,
		},
		{
			name: "password too short",
			request: model.RegisterUserRequest{
				Username: "testuser",
				Email:    "test@example.com",
				Password: "pass", // min 8
				Role:     "presenter",
			},
			shouldError: true,
		},
		{
			name: "invalid role",
			request: model.RegisterUserRequest{
				Username: "testuser",
				Email:    "test@example.com",
				Password: "password123",
				Role:     "invalid", // must be presenter or admin
			},
			shouldError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validate.Struct(tt.request)
			if tt.shouldError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

// TestMockTokenUtil_CreateToken test mock token util
func TestMockTokenUtil_CreateToken(t *testing.T) {
	mockTokenUtil := new(mocks.MockTokenUtil)

	auth := &model.Auth{
		Username: "testuser",
		Role:     "presenter",
	}

	expectedToken := "mock-jwt-token"
	mockTokenUtil.On("CreateToken", mock.Anything, auth).Return(expectedToken, nil)

	token, err := mockTokenUtil.CreateToken(context.Background(), auth)

	assert.NoError(t, err)
	assert.Equal(t, expectedToken, token)
	mockTokenUtil.AssertExpectations(t)
}

// TestMockRoomRepository_FindByRoomCode test mock room repository
func TestMockRoomRepository_FindByRoomCode(t *testing.T) {
	mockRoomRepo := new(mocks.MockRoomRepository)

	expectedRoom := &entity.Room{
		RoomCode: "ABC123",
		Title:    "Test Room",
		Status:   "active",
	}
	expectedRoom.ID = 1

	mockRoomRepo.On("FindByRoomCode", mock.Anything, "ABC123").Return(expectedRoom, nil)

	room, err := mockRoomRepo.FindByRoomCode(nil, "ABC123")

	assert.NoError(t, err)
	assert.NotNil(t, room)
	assert.Equal(t, "ABC123", room.RoomCode)
	mockRoomRepo.AssertExpectations(t)
}

// TestMockRoomRepository_FindByRoomCode_NotFound test room not found
func TestMockRoomRepository_FindByRoomCode_NotFound(t *testing.T) {
	mockRoomRepo := new(mocks.MockRoomRepository)

	mockRoomRepo.On("FindByRoomCode", mock.Anything, "NOTFND").Return(nil, nil)

	room, err := mockRoomRepo.FindByRoomCode(nil, "NOTFND")

	assert.NoError(t, err)
	assert.Nil(t, room)
	mockRoomRepo.AssertExpectations(t)
}
