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

// setupMessageUseCaseTest setup test environment for MessageUseCase
func setupMessageUseCaseTest(t *testing.T) (*usecase.MessageUseCase, sqlmock.Sqlmock) {
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

	// create xp transaction usecase (nil for simple tests)
	xpUC := &usecase.XPTransactionUseCase{
		Log: log,
	}

	// create usecase
	uc := &usecase.MessageUseCase{
		DB:                    gormDB,
		Validate:              validate,
		Log:                   log,
		MessageRepository:     &repository.MessageRepository{Log: log},
		RoomRepository:        &repository.RoomRepository{Log: log},
		ParticipantRepository: &repository.ParticipantRepository{Log: log},
		XPTransactionUseCase:  xpUC,
	}

	return uc, mockDB
}

// TestMessageUseCase_Send_InvalidRequest test send message with invalid request
func TestMessageUseCase_Send_InvalidRequest(t *testing.T) {
	uc, mockDB := setupMessageUseCaseTest(t)

	// expect begin transaction
	mockDB.ExpectBegin()
	mockDB.ExpectRollback()

	// test with empty content
	request := &model.SendMessageRequest{
		RoomID:        1,
		ParticipantID: 1,
		Content:       "", // invalid: required
	}

	result, err := uc.Send(context.Background(), request)

	assert.Nil(t, result)
	assert.Error(t, err)
}

// TestMessageUseCase_List_InvalidRequest test list messages with invalid request
func TestMessageUseCase_List_InvalidRequest(t *testing.T) {
	uc, mockDB := setupMessageUseCaseTest(t)

	// expect begin transaction
	mockDB.ExpectBegin()
	mockDB.ExpectRollback()

	// test with room id 0
	request := &model.GetMessagesRequest{
		RoomID:        0, // invalid: min 1
		ParticipantID: 1,
		Limit:         20,
	}

	result, err := uc.List(context.Background(), request)

	assert.Nil(t, result)
	assert.Error(t, err)
}

// TestSendMessageRequest_Validation test send message request validation
func TestSendMessageRequest_Validation(t *testing.T) {
	validate := validator.New()

	tests := []struct {
		name        string
		request     model.SendMessageRequest
		shouldError bool
	}{
		{
			name: "valid request",
			request: model.SendMessageRequest{
				RoomID:        1,
				ParticipantID: 1,
				Content:       "Hello, world!",
			},
			shouldError: false,
		},
		{
			name: "empty content",
			request: model.SendMessageRequest{
				RoomID:        1,
				ParticipantID: 1,
				Content:       "",
			},
			shouldError: true,
		},
		{
			name: "content too long",
			request: model.SendMessageRequest{
				RoomID:        1,
				ParticipantID: 1,
				Content:       string(make([]byte, 1001)), // max 1000
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

// TestGetMessagesRequest_Validation test get messages request validation
func TestGetMessagesRequest_Validation(t *testing.T) {
	validate := validator.New()

	tests := []struct {
		name        string
		request     model.GetMessagesRequest
		shouldError bool
	}{
		{
			name: "valid request",
			request: model.GetMessagesRequest{
				RoomID:        1,
				ParticipantID: 1,
				Limit:         20,
			},
			shouldError: false,
		},
		{
			name: "limit too high",
			request: model.GetMessagesRequest{
				RoomID:        1,
				ParticipantID: 1,
				Limit:         101, // max 100
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

// TestMockMessageRepository_List test mock message repository
func TestMockMessageRepository_List(t *testing.T) {
	mockMessageRepo := new(mocks.MockMessageRepository)

	participant := entity.Participant{
		DisplayName: "Test User",
	}
	participant.ID = 1

	expectedMessages := []entity.Message{
		{
			RoomID:        1,
			ParticipantID: 1,
			Content:       "Hello",
			Participant:   participant,
		},
		{
			RoomID:        1,
			ParticipantID: 2,
			Content:       "World",
			Participant:   participant,
		},
	}

	mockMessageRepo.On("List", mock.Anything, uint(1), 20, (*int64)(nil)).Return(expectedMessages, nil)

	messages, err := mockMessageRepo.List(nil, 1, 20, nil)

	assert.NoError(t, err)
	assert.Len(t, messages, 2)
	assert.Equal(t, "Hello", messages[0].Content)
	mockMessageRepo.AssertExpectations(t)
}

// TestMockMessageRepository_Create test mock message repository create
func TestMockMessageRepository_Create(t *testing.T) {
	mockMessageRepo := new(mocks.MockMessageRepository)

	message := &entity.Message{
		RoomID:        1,
		ParticipantID: 1,
		Content:       "Test message",
	}

	mockMessageRepo.On("Create", mock.Anything, message).Return(nil)

	err := mockMessageRepo.Create(nil, message)

	assert.NoError(t, err)
	mockMessageRepo.AssertExpectations(t)
}
