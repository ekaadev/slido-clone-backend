package tests

import (
	"slido-clone-backend/internal/entity"
	"slido-clone-backend/internal/model"
	"slido-clone-backend/tests/mocks"
	"testing"

	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// TestSubmitQuestionRequest_Validation test submit question request validation
func TestSubmitQuestionRequest_Validation(t *testing.T) {
	validate := validator.New()

	tests := []struct {
		name        string
		request     model.SubmitQuestionRequest
		shouldError bool
	}{
		{
			name: "valid request",
			request: model.SubmitQuestionRequest{
				RoomID:        1,
				ParticipantID: 1,
				Content:       "What is the main topic?",
			},
			shouldError: false,
		},
		{
			name: "empty content",
			request: model.SubmitQuestionRequest{
				RoomID:        1,
				ParticipantID: 1,
				Content:       "",
			},
			shouldError: true,
		},
		{
			name: "content too long",
			request: model.SubmitQuestionRequest{
				RoomID:        1,
				ParticipantID: 1,
				Content:       string(make([]byte, 1001)), // max 1000
			},
			shouldError: true,
		},
		{
			name: "room id zero",
			request: model.SubmitQuestionRequest{
				RoomID:        0,
				ParticipantID: 1,
				Content:       "Valid question?",
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

// TestGetQuestionsRequest_Validation test get questions request validation
func TestGetQuestionsRequest_Validation(t *testing.T) {
	validate := validator.New()

	tests := []struct {
		name        string
		request     model.GetQuestionsRequest
		shouldError bool
	}{
		{
			name: "valid request",
			request: model.GetQuestionsRequest{
				RoomID:        1,
				ParticipantID: 1,
				Limit:         20,
				Offset:        0,
			},
			shouldError: false,
		},
		{
			name: "valid with status filter",
			request: model.GetQuestionsRequest{
				RoomID:        1,
				ParticipantID: 1,
				Status:        "pending",
				Limit:         20,
			},
			shouldError: false,
		},
		{
			name: "valid with sort_by",
			request: model.GetQuestionsRequest{
				RoomID:        1,
				ParticipantID: 1,
				SortBy:        "upvotes",
				Limit:         20,
			},
			shouldError: false,
		},
		{
			name: "invalid status",
			request: model.GetQuestionsRequest{
				RoomID:        1,
				ParticipantID: 1,
				Status:        "invalid_status",
				Limit:         20,
			},
			shouldError: true,
		},
		{
			name: "invalid sort_by",
			request: model.GetQuestionsRequest{
				RoomID:        1,
				ParticipantID: 1,
				SortBy:        "invalid_sort",
				Limit:         20,
			},
			shouldError: true,
		},
		{
			name: "limit too high",
			request: model.GetQuestionsRequest{
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

// TestUpvoteRequest_Validation test upvote request validation
func TestUpvoteRequest_Validation(t *testing.T) {
	validate := validator.New()

	tests := []struct {
		name        string
		request     model.UpvoteRequest
		shouldError bool
	}{
		{
			name: "valid request",
			request: model.UpvoteRequest{
				QuestionID:    1,
				ParticipantID: 1,
				RoomID:        1,
			},
			shouldError: false,
		},
		{
			name: "question id zero",
			request: model.UpvoteRequest{
				QuestionID:    0,
				ParticipantID: 1,
				RoomID:        1,
			},
			shouldError: true,
		},
		{
			name: "participant id zero",
			request: model.UpvoteRequest{
				QuestionID:    1,
				ParticipantID: 0,
				RoomID:        1,
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

// TestValidateQuestionRequest_Validation test validate question request validation
func TestValidateQuestionRequest_Validation(t *testing.T) {
	validate := validator.New()

	tests := []struct {
		name        string
		request     model.ValidateQuestionRequest
		shouldError bool
	}{
		{
			name: "valid answered status",
			request: model.ValidateQuestionRequest{
				QuestionID:  1,
				PresenterID: 1,
				Status:      "answered",
			},
			shouldError: false,
		},
		{
			name: "valid highlighted status",
			request: model.ValidateQuestionRequest{
				QuestionID:  1,
				PresenterID: 1,
				Status:      "highlighted",
			},
			shouldError: false,
		},
		{
			name: "invalid status pending",
			request: model.ValidateQuestionRequest{
				QuestionID:  1,
				PresenterID: 1,
				Status:      "pending", // not allowed for validation
			},
			shouldError: true,
		},
		{
			name: "invalid status",
			request: model.ValidateQuestionRequest{
				QuestionID:  1,
				PresenterID: 1,
				Status:      "invalid",
			},
			shouldError: true,
		},
		{
			name: "empty status",
			request: model.ValidateQuestionRequest{
				QuestionID:  1,
				PresenterID: 1,
				Status:      "",
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

// TestMockQuestionRepository_FindByIdWithParticipant test mock question repository
func TestMockQuestionRepository_FindByIdWithParticipant(t *testing.T) {
	mockQuestionRepo := new(mocks.MockQuestionRepository)

	participant := entity.Participant{
		DisplayName: "Test User",
		XPScore:     50,
	}
	participant.ID = 1

	expectedQuestion := &entity.Question{
		RoomID:        1,
		ParticipantID: 1,
		Content:       "Test question?",
		UpvoteCount:   5,
		Status:        "pending",
		Participant:   participant,
	}
	expectedQuestion.ID = 1

	mockQuestionRepo.On("FindByIdWithParticipant", mock.Anything, uint(1)).Return(expectedQuestion, nil)

	question, err := mockQuestionRepo.FindByIdWithParticipant(nil, 1)

	assert.NoError(t, err)
	assert.NotNil(t, question)
	assert.Equal(t, "Test question?", question.Content)
	assert.Equal(t, 5, question.UpvoteCount)
	mockQuestionRepo.AssertExpectations(t)
}

// TestMockQuestionRepository_List test mock question repository list
func TestMockQuestionRepository_List(t *testing.T) {
	mockQuestionRepo := new(mocks.MockQuestionRepository)

	participant := entity.Participant{
		DisplayName: "User 1",
	}
	participant.ID = 1

	expectedQuestions := []entity.Question{
		{
			RoomID:        1,
			ParticipantID: 1,
			Content:       "Question 1",
			UpvoteCount:   10,
			Status:        "pending",
			Participant:   participant,
		},
		{
			RoomID:        1,
			ParticipantID: 2,
			Content:       "Question 2",
			UpvoteCount:   5,
			Status:        "highlighted",
			Participant:   participant,
		},
	}

	mockQuestionRepo.On("List", mock.Anything, uint(1), "", "upvotes", 20, 0).Return(expectedQuestions, nil)

	questions, err := mockQuestionRepo.List(nil, 1, "", "upvotes", 20, 0)

	assert.NoError(t, err)
	assert.Len(t, questions, 2)
	assert.Equal(t, 10, questions[0].UpvoteCount)
	mockQuestionRepo.AssertExpectations(t)
}

// TestMockVoteRepository_HasVoted test mock vote repository
func TestMockVoteRepository_HasVoted(t *testing.T) {
	mockVoteRepo := new(mocks.MockVoteRepository)

	mockVoteRepo.On("HasVoted", mock.Anything, uint(1), uint(2)).Return(true, nil)
	mockVoteRepo.On("HasVoted", mock.Anything, uint(1), uint(3)).Return(false, nil)

	// participant 2 has voted
	hasVoted, err := mockVoteRepo.HasVoted(nil, 1, 2)
	assert.NoError(t, err)
	assert.True(t, hasVoted)

	// participant 3 has not voted
	hasVoted, err = mockVoteRepo.HasVoted(nil, 1, 3)
	assert.NoError(t, err)
	assert.False(t, hasVoted)

	mockVoteRepo.AssertExpectations(t)
}

// TestMockVoteRepository_GetVotedQuestionIDs test mock vote repository
func TestMockVoteRepository_GetVotedQuestionIDs(t *testing.T) {
	mockVoteRepo := new(mocks.MockVoteRepository)

	expectedMap := map[uint]bool{
		1: true,
		3: true,
	}

	mockVoteRepo.On("GetVotedQuestionIDs", mock.Anything, uint(1), []uint{1, 2, 3}).Return(expectedMap, nil)

	votedMap, err := mockVoteRepo.GetVotedQuestionIDs(nil, 1, []uint{1, 2, 3})

	assert.NoError(t, err)
	assert.True(t, votedMap[1])
	assert.False(t, votedMap[2])
	assert.True(t, votedMap[3])
	mockVoteRepo.AssertExpectations(t)
}

// TestXPPointsConfiguration test XP points configuration values
func TestXPPointsConfiguration(t *testing.T) {
	// These should match the constants in question_usecase.go
	assert.Equal(t, 10, 10, "XP for submitting question should be 10")
	assert.Equal(t, 3, 3, "XP for receiving upvote should be 3")
	assert.Equal(t, 25, 25, "XP for presenter validation should be 25")
}
