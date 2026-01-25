package tests

import (
	"context"
	"slido-clone-backend/internal/entity"
	"slido-clone-backend/internal/model"
	"slido-clone-backend/internal/repository"
	"slido-clone-backend/internal/usecase"
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

// setupRoomUseCaseTest setup test environment for RoomUseCase
func setupRoomUseCaseTest(t *testing.T) (*usecase.RoomUseCase, sqlmock.Sqlmock) {
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

	// create real validator
	validate := validator.New()

	// create logger
	log := logrus.New()
	log.SetLevel(logrus.PanicLevel)

	// create usecase
	uc := &usecase.RoomUseCase{
		DB:                    gormDB,
		Log:                   log,
		Validate:              validate,
		RoomRepository:        &repository.RoomRepository{Log: log},
		ParticipantRepository: &repository.ParticipantRepository{Log: log},
	}

	return uc, mockDB
}

// TestRoomUseCase_Create_InvalidRequest test create room with invalid request
func TestRoomUseCase_Create_InvalidRequest(t *testing.T) {
	uc, mockDB := setupRoomUseCaseTest(t)

	// expect begin transaction
	mockDB.ExpectBegin()
	mockDB.ExpectRollback()

	// test with empty title
	request := &model.CreateRoomRequest{
		Title:       "", // invalid: required
		PresenterID: 1,
	}

	result, err := uc.Create(context.Background(), request)

	assert.Nil(t, result)
	assert.Error(t, err)
}

// TestRoomUseCase_Create_TitleTooShort test create room with title too short
func TestRoomUseCase_Create_TitleTooShort(t *testing.T) {
	uc, mockDB := setupRoomUseCaseTest(t)

	// expect begin transaction
	mockDB.ExpectBegin()
	mockDB.ExpectRollback()

	// test with short title
	request := &model.CreateRoomRequest{
		Title:       "AB", // invalid: min 3
		PresenterID: 1,
	}

	result, err := uc.Create(context.Background(), request)

	assert.Nil(t, result)
	assert.Error(t, err)
}

// TestRoomUseCase_Get_InvalidRoomCode test get room with invalid room code
func TestRoomUseCase_Get_InvalidRoomCode(t *testing.T) {
	uc, mockDB := setupRoomUseCaseTest(t)

	// expect begin transaction
	mockDB.ExpectBegin()
	mockDB.ExpectRollback()

	// test with invalid room code
	request := &model.GetRoomRequestByRoomCode{
		RoomCode: "ABC", // invalid: must be 6 characters
	}

	result, err := uc.Get(context.Background(), request)

	assert.Nil(t, result)
	assert.Error(t, err)
}

// TestRoomUseCase_UpdateToClosed_InvalidRequest test update to closed with invalid request
func TestRoomUseCase_UpdateToClosed_InvalidRequest(t *testing.T) {
	uc, mockDB := setupRoomUseCaseTest(t)

	// expect begin transaction
	mockDB.ExpectBegin()
	mockDB.ExpectRollback()

	// test with invalid status
	request := &model.UpdateToCloseRoomRequestByID{
		PresenterID: 1,
		RoomID:      1,
		Status:      "active", // invalid: must be "closed"
	}

	result, err := uc.UpdateToClosed(context.Background(), request)

	assert.Nil(t, result)
	assert.Error(t, err)
}

// TestCreateRoomRequest_Validation test create room request validation
func TestCreateRoomRequest_Validation(t *testing.T) {
	validate := validator.New()

	tests := []struct {
		name        string
		request     model.CreateRoomRequest
		shouldError bool
	}{
		{
			name: "valid request",
			request: model.CreateRoomRequest{
				Title:       "Test Room",
				PresenterID: 1,
			},
			shouldError: false,
		},
		{
			name: "title too short",
			request: model.CreateRoomRequest{
				Title:       "AB", // min 3
				PresenterID: 1,
			},
			shouldError: true,
		},
		{
			name: "empty title",
			request: model.CreateRoomRequest{
				Title:       "",
				PresenterID: 1,
			},
			shouldError: true,
		},
		{
			name: "presenter id zero",
			request: model.CreateRoomRequest{
				Title:       "Test Room",
				PresenterID: 0, // min 1
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

// TestCloseRoomRequest_Validation test close room request validation
func TestCloseRoomRequest_Validation(t *testing.T) {
	validate := validator.New()

	tests := []struct {
		name        string
		request     model.CloseRoomRequest
		shouldError bool
	}{
		{
			name: "valid closed status",
			request: model.CloseRoomRequest{
				Status: "closed",
			},
			shouldError: false,
		},
		{
			name: "invalid status",
			request: model.CloseRoomRequest{
				Status: "active", // must be "closed"
			},
			shouldError: true,
		},
		{
			name: "empty status",
			request: model.CloseRoomRequest{
				Status: "",
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

// TestMockRoomRepository_Search test mock room repository search
func TestMockRoomRepository_Search(t *testing.T) {
	mockRoomRepo := new(mocks.MockRoomRepository)

	expectedRooms := []entity.Room{
		{RoomCode: "ABC123", Title: "Room 1", Status: "active"},
		{RoomCode: "DEF456", Title: "Room 2", Status: "closed"},
	}

	mockRoomRepo.On("Search", mock.Anything, uint(1)).Return(expectedRooms, nil)

	rooms, err := mockRoomRepo.Search(nil, 1)

	assert.NoError(t, err)
	assert.Len(t, rooms, 2)
	assert.Equal(t, "ABC123", rooms[0].RoomCode)
	mockRoomRepo.AssertExpectations(t)
}

// TestMockRoomRepository_FindByIdAndPresenterId test mock room repository
func TestMockRoomRepository_FindByIdAndPresenterId(t *testing.T) {
	mockRoomRepo := new(mocks.MockRoomRepository)

	expectedRoom := &entity.Room{
		RoomCode: "ABC123",
		Title:    "Test Room",
		Status:   "active",
	}
	expectedRoom.ID = 1

	mockRoomRepo.On("FindByIdAndPresenterId", mock.Anything, uint(1), uint(1)).Return(expectedRoom, nil)

	room, err := mockRoomRepo.FindByIdAndPresenterId(nil, 1, 1)

	assert.NoError(t, err)
	assert.NotNil(t, room)
	assert.Equal(t, uint(1), room.ID)
	mockRoomRepo.AssertExpectations(t)
}

// TestGenerateRoomCode test room code generation
func TestGenerateRoomCode(t *testing.T) {
	code, err := usecase.GenerateRoomCode(6)

	assert.NoError(t, err)
	assert.Len(t, code, 6)

	// verify code contains only valid characters
	validChars := "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	for _, c := range code {
		assert.Contains(t, validChars, string(c))
	}
}

// TestGenerateRoomCode_UniqueEachTime test that generated codes are unique
func TestGenerateRoomCode_UniqueEachTime(t *testing.T) {
	codes := make(map[string]bool)

	// generate 100 codes and check for uniqueness
	for i := 0; i < 100; i++ {
		code, err := usecase.GenerateRoomCode(6)
		assert.NoError(t, err)
		codes[code] = true
	}

	// should have high uniqueness (allow some collisions but expect most to be unique)
	assert.Greater(t, len(codes), 90)
}
