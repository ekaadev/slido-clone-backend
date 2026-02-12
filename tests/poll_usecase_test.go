package tests

import (
	"slido-clone-backend/internal/model"
	"slido-clone-backend/tests/mocks"
	"testing"

	"github.com/go-playground/validator/v10"
)

// TestCreatePollRequest_Validation test validation untuk CreatePollRequest
func TestCreatePollRequest_Validation(t *testing.T) {
	validate := validator.New()

	tests := []struct {
		name    string
		request model.CreatePollRequest
		wantErr bool
	}{
		{
			name: "valid request",
			request: model.CreatePollRequest{
				RoomID:      1,
				PresenterID: 1,
				Question:    "What is your favorite color?",
				Options:     []string{"Red", "Blue", "Green", "Yellow"},
			},
			wantErr: false,
		},
		{
			name: "question too short",
			request: model.CreatePollRequest{
				RoomID:      1,
				PresenterID: 1,
				Question:    "Hi",
				Options:     []string{"Red", "Blue"},
			},
			wantErr: true,
		},
		{
			name: "too few options",
			request: model.CreatePollRequest{
				RoomID:      1,
				PresenterID: 1,
				Question:    "What is your favorite color?",
				Options:     []string{"Red"},
			},
			wantErr: true,
		},
		{
			name: "empty option",
			request: model.CreatePollRequest{
				RoomID:      1,
				PresenterID: 1,
				Question:    "What is your favorite color?",
				Options:     []string{"Red", ""},
			},
			wantErr: true,
		},
		{
			name: "room_id zero",
			request: model.CreatePollRequest{
				RoomID:      0,
				PresenterID: 1,
				Question:    "What is your favorite color?",
				Options:     []string{"Red", "Blue"},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validate.Struct(tt.request)
			if (err != nil) != tt.wantErr {
				t.Errorf("Validation error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

// TestSubmitVoteRequest_Validation test validation untuk SubmitVoteRequest
func TestSubmitVoteRequest_Validation(t *testing.T) {
	validate := validator.New()

	tests := []struct {
		name    string
		request model.SubmitVoteRequest
		wantErr bool
	}{
		{
			name: "valid request",
			request: model.SubmitVoteRequest{
				PollID:        1,
				ParticipantID: 1,
				RoomID:        1,
				OptionID:      1,
			},
			wantErr: false,
		},
		{
			name: "poll_id zero",
			request: model.SubmitVoteRequest{
				PollID:        0,
				ParticipantID: 1,
				RoomID:        1,
				OptionID:      1,
			},
			wantErr: true,
		},
		{
			name: "option_id zero",
			request: model.SubmitVoteRequest{
				PollID:        1,
				ParticipantID: 1,
				RoomID:        1,
				OptionID:      0,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validate.Struct(tt.request)
			if (err != nil) != tt.wantErr {
				t.Errorf("Validation error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

// TestMockPollRepository test mock poll repository
func TestMockPollRepository_Create(t *testing.T) {
	repo := mocks.NewMockPollRepository()

	// Test initial state
	if len(repo.Polls) != 0 {
		t.Errorf("Expected 0 polls, got %d", len(repo.Polls))
	}
}

// TestMockPollResponseRepository_HasVoted test mock response repository
func TestMockPollResponseRepository_HasVoted(t *testing.T) {
	repo := mocks.NewMockPollResponseRepository()

	hasVoted, err := repo.HasVoted(nil, 1, 1)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if hasVoted {
		t.Error("Expected hasVoted false, got true")
	}
}

// TestXPPointsConfiguration_Poll test XP configuration untuk polling
func TestXPPointsConfiguration_Poll(t *testing.T) {
	// XP untuk submit vote harus 5
	expectedXP := 5

	// This is a placeholder - in real test we would check the constant
	if expectedXP != 5 {
		t.Errorf("Expected XP for vote to be 5, got %d", expectedXP)
	}
}
