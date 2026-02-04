package model

import "github.com/golang-jwt/jwt/v5"

type Auth struct {
	// Core identification
	UserID        *uint `json:"user_id,omitempty"`        // NULL for anonymous
	ParticipantID *uint `json:"participant_id,omitempty"` // Set when joined room
	RoomID        *uint `json:"room_id,omitempty"`        // Current room context

	// User info
	Username    string `json:"username,omitempty"`     // Empty for anonymous
	DisplayName string `json:"display_name,omitempty"` // For anonymous users
	Email       string `json:"email,omitempty"`        // Empty for anonymous
	Role        string `json:"role"`                   // "presenter" | "admin" | "anonymous"

	// Flags
	IsAnonymous bool `json:"is_anonymous"`
	IsRoomOwner bool `json:"is_room_owner"` // true jika user adalah pembuat room (host)

	// Standard claims
	jwt.RegisteredClaims
}

// Update struct for updating Auth information
