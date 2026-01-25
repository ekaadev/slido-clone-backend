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

// setupParticipantUseCaseTest setup test environment for ParticipantUseCase
func setupParticipantUseCaseTest(t *testing.T) (*usecase.ParticipantUseCase, sqlmock.Sqlmock) {
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
	uc := &usecase.ParticipantUseCase{
		DB:                    gormDB,
		Log:                   log,
		Validate:              validate,
		ParticipantRepository: &repository.ParticipantRepository{Log: log},
		RoomRepository:        &repository.RoomRepository{Log: log},
		UserRepository:        &repository.UserRepository{Log: log},
		TokenUtil:             &util.TokenUtil{SecretKey: "test-secret"},
	}

	return uc, mockDB
}

// TestParticipantUseCase_Join_InvalidRequest test join room with invalid request
func TestParticipantUseCase_Join_InvalidRequest(t *testing.T) {
	uc, mockDB := setupParticipantUseCaseTest(t)

	// expect begin transaction
	mockDB.ExpectBegin()
	mockDB.ExpectRollback()

	// test with invalid room code
	request := &model.JoinRoomRequest{
		Username: "testuser",
		RoomCode: "ABC", // invalid: must be 6 characters
	}

	result, err := uc.Join(context.Background(), request)

	assert.Nil(t, result)
	assert.Error(t, err)
}

// TestParticipantUseCase_List_InvalidRequest test list participants with invalid request
func TestParticipantUseCase_List_InvalidRequest(t *testing.T) {
	uc, mockDB := setupParticipantUseCaseTest(t)

	// expect begin transaction
	mockDB.ExpectBegin()
	mockDB.ExpectRollback()

	// test with invalid page
	request := &model.ListParticipantsRequest{
		RoomID:        1,
		ParticipantID: 1,
		Page:          0, // invalid: min 1
		Size:          10,
	}

	result, total, err := uc.List(context.Background(), request)

	assert.Nil(t, result)
	assert.Equal(t, int64(0), total)
	assert.Error(t, err)
}

// TestParticipantUseCase_Leaderboard_InvalidRequest test leaderboard with invalid request
func TestParticipantUseCase_Leaderboard_InvalidRequest(t *testing.T) {
	uc, mockDB := setupParticipantUseCaseTest(t)

	// expect begin transaction
	mockDB.ExpectBegin()
	mockDB.ExpectRollback()

	// test with invalid room id
	request := &model.GetLeaderboardRequest{
		RoomID:        0, // invalid: min 1
		ParticipantID: 1,
	}

	result, err := uc.Leaderboard(context.Background(), request)

	assert.Nil(t, result)
	assert.Error(t, err)
}

// TestJoinRoomRequest_Validation test join room request validation
func TestJoinRoomRequest_Validation(t *testing.T) {
	validate := validator.New()

	tests := []struct {
		name        string
		request     model.JoinRoomRequest
		shouldError bool
	}{
		{
			name: "valid request",
			request: model.JoinRoomRequest{
				Username: "testuser",
				RoomCode: "ABC123",
			},
			shouldError: false,
		},
		{
			name: "valid request with display name",
			request: model.JoinRoomRequest{
				Username:    "testuser",
				DisplayName: "Test User",
				RoomCode:    "ABC123",
			},
			shouldError: false,
		},
		{
			name: "username too short",
			request: model.JoinRoomRequest{
				Username: "ab", // min 3
				RoomCode: "ABC123",
			},
			shouldError: true,
		},
		{
			name: "room code too short",
			request: model.JoinRoomRequest{
				Username: "testuser",
				RoomCode: "ABC", // must be 6
			},
			shouldError: true,
		},
		{
			name: "empty username",
			request: model.JoinRoomRequest{
				Username: "",
				RoomCode: "ABC123",
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

// TestListParticipantsRequest_Validation test list participants request validation
func TestListParticipantsRequest_Validation(t *testing.T) {
	validate := validator.New()

	tests := []struct {
		name        string
		request     model.ListParticipantsRequest
		shouldError bool
	}{
		{
			name: "valid request",
			request: model.ListParticipantsRequest{
				RoomID:        1,
				ParticipantID: 1,
				Page:          1,
				Size:          10,
			},
			shouldError: false,
		},
		{
			name: "page zero",
			request: model.ListParticipantsRequest{
				RoomID:        1,
				ParticipantID: 1,
				Page:          0, // min 1
				Size:          10,
			},
			shouldError: true,
		},
		{
			name: "size too large",
			request: model.ListParticipantsRequest{
				RoomID:        1,
				ParticipantID: 1,
				Page:          1,
				Size:          101, // max 100
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

// TestMockParticipantRepository_List test mock participant repository
func TestMockParticipantRepository_List(t *testing.T) {
	mockParticipantRepo := new(mocks.MockParticipantRepository)

	isAnon := false
	expectedParticipants := []entity.Participant{
		{DisplayName: "User 1", XPScore: 100, IsAnonymous: &isAnon},
		{DisplayName: "User 2", XPScore: 80, IsAnonymous: &isAnon},
	}

	mockParticipantRepo.On("List", mock.Anything, uint(1), 0, 10).Return(expectedParticipants, nil)

	participants, err := mockParticipantRepo.List(nil, 1, 0, 10)

	assert.NoError(t, err)
	assert.Len(t, participants, 2)
	assert.Equal(t, "User 1", participants[0].DisplayName)
	mockParticipantRepo.AssertExpectations(t)
}

// TestMockParticipantRepository_ListByXPScore test mock participant repository leaderboard
func TestMockParticipantRepository_ListByXPScore(t *testing.T) {
	mockParticipantRepo := new(mocks.MockParticipantRepository)

	isAnon := false
	expectedParticipants := []entity.Participant{
		{DisplayName: "User 1", XPScore: 100, IsAnonymous: &isAnon},
		{DisplayName: "User 2", XPScore: 80, IsAnonymous: &isAnon},
		{DisplayName: "User 3", XPScore: 60, IsAnonymous: &isAnon},
	}

	mockParticipantRepo.On("ListByXPScore", mock.Anything, uint(1)).Return(expectedParticipants, nil)

	participants, err := mockParticipantRepo.ListByXPScore(nil, 1)

	assert.NoError(t, err)
	assert.Len(t, participants, 3)
	// should be ordered by XP score descending
	assert.Equal(t, 100, participants[0].XPScore)
	assert.Equal(t, 80, participants[1].XPScore)
	mockParticipantRepo.AssertExpectations(t)
}

// TestMockParticipantRepository_GetRankAndScore test mock get rank and score
func TestMockParticipantRepository_GetRankAndScore(t *testing.T) {
	mockParticipantRepo := new(mocks.MockParticipantRepository)

	mockParticipantRepo.On("GetRankAndScore", mock.Anything, uint(1), uint(5)).Return(3, 75, nil)

	rank, score, err := mockParticipantRepo.GetRankAndScore(nil, 1, 5)

	assert.NoError(t, err)
	assert.Equal(t, 3, rank)
	assert.Equal(t, 75, score)
	mockParticipantRepo.AssertExpectations(t)
}
