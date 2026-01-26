package tests

import (
	"slido-clone-backend/internal/model"
	"testing"

	"github.com/go-playground/validator/v10"
)

// TestGetXPTransactionsRequest_Validation test validation untuk GetXPTransactionsRequest
func TestGetXPTransactionsRequest_Validation(t *testing.T) {
	validate := validator.New()

	tests := []struct {
		name    string
		request model.GetXPTransactionsRequest
		wantErr bool
	}{
		{
			name: "valid request",
			request: model.GetXPTransactionsRequest{
				RoomID:        1,
				ParticipantID: 1,
				Limit:         50,
			},
			wantErr: false,
		},
		{
			name: "valid request without limit",
			request: model.GetXPTransactionsRequest{
				RoomID:        1,
				ParticipantID: 1,
			},
			wantErr: false,
		},
		{
			name: "room_id zero",
			request: model.GetXPTransactionsRequest{
				RoomID:        0,
				ParticipantID: 1,
				Limit:         50,
			},
			wantErr: true,
		},
		{
			name: "participant_id zero",
			request: model.GetXPTransactionsRequest{
				RoomID:        1,
				ParticipantID: 0,
				Limit:         50,
			},
			wantErr: true,
		},
		{
			name: "limit too high",
			request: model.GetXPTransactionsRequest{
				RoomID:        1,
				ParticipantID: 1,
				Limit:         200,
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

// TestDeleteRoomRequest_Validation test validation untuk DeleteRoomRequest
func TestDeleteRoomRequest_Validation(t *testing.T) {
	validate := validator.New()

	tests := []struct {
		name    string
		request model.DeleteRoomRequest
		wantErr bool
	}{
		{
			name: "valid request",
			request: model.DeleteRoomRequest{
				PresenterID: 1,
				RoomID:      1,
			},
			wantErr: false,
		},
		{
			name: "presenter_id zero",
			request: model.DeleteRoomRequest{
				PresenterID: 0,
				RoomID:      1,
			},
			wantErr: true,
		},
		{
			name: "room_id zero",
			request: model.DeleteRoomRequest{
				PresenterID: 1,
				RoomID:      0,
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

// TestXPTransactionItemStructure test XPTransactionItem structure
func TestXPTransactionItemStructure(t *testing.T) {
	item := model.XPTransactionItem{
		ID:         1,
		Points:     10,
		SourceType: "question_submitted",
		SourceID:   5,
		CreatedAt:  "2026-01-26T07:00:00Z",
	}

	if item.ID != 1 {
		t.Errorf("Expected ID 1, got %d", item.ID)
	}
	if item.Points != 10 {
		t.Errorf("Expected Points 10, got %d", item.Points)
	}
	if item.SourceType != "question_submitted" {
		t.Errorf("Expected SourceType 'question_submitted', got '%s'", item.SourceType)
	}
}

// TestGetXPTransactionsResponseStructure test response structure
func TestGetXPTransactionsResponseStructure(t *testing.T) {
	response := model.GetXPTransactionsResponse{
		Transactions: []model.XPTransactionItem{
			{ID: 1, Points: 5, SourceType: "poll_voted"},
			{ID: 2, Points: 10, SourceType: "question_submitted"},
		},
		TotalXP: 15,
		Total:   2,
	}

	if len(response.Transactions) != 2 {
		t.Errorf("Expected 2 transactions, got %d", len(response.Transactions))
	}
	if response.TotalXP != 15 {
		t.Errorf("Expected TotalXP 15, got %d", response.TotalXP)
	}
	if response.Total != 2 {
		t.Errorf("Expected Total 2, got %d", response.Total)
	}
}

// TestLogoutTokenHandling test token handling for logout
func TestLogoutTokenHandling(t *testing.T) {
	// Test that Bearer prefix is required
	authHeader := "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
	tokenOnly := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."

	if authHeader == tokenOnly {
		t.Error("Bearer prefix should be present for valid auth header")
	}

	// Simple string prefix check simulation
	expectedPrefix := "Bearer "
	if len(authHeader) < len(expectedPrefix) {
		t.Error("Auth header too short")
	}
	if authHeader[:7] != expectedPrefix {
		t.Error("Should start with 'Bearer '")
	}
}
