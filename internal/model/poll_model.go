package model

import "time"

// ========================================
// Request Models
// ========================================

// CreatePollRequest request untuk membuat poll baru
type CreatePollRequest struct {
	RoomID      uint     `json:"-" validate:"required,min=1"`
	PresenterID uint     `json:"-" validate:"required,min=1"`
	Question    string   `json:"question" validate:"required,min=1,max=500"`
	Options     []string `json:"options" validate:"required,min=2,max=10,dive,min=1,max=255"`
}

// GetActivePollsRequest request untuk mendapatkan active polls
type GetActivePollsRequest struct {
	RoomID        uint `json:"-" validate:"required,min=1"`
	ParticipantID uint `json:"-" validate:"required,min=1"`
}

// GetPollHistoryRequest request untuk mendapatkan poll history
type GetPollHistoryRequest struct {
	RoomID uint   `json:"-" validate:"required,min=1"`
	Status string `json:"status" validate:"omitempty,oneof=active closed all"`
	Limit  int    `json:"limit" validate:"omitempty,min=1,max=100"`
}

// SubmitPollVoteRequest request untuk submit vote pada poll
type SubmitPollVoteRequest struct {
	PollID        uint `json:"-" validate:"required,min=1"`
	ParticipantID uint `json:"-" validate:"required,min=1"`
	RoomID        uint `json:"-" validate:"required,min=1"`
	OptionID      uint `json:"option_id" validate:"required,min=1"`
}

// ClosePollRequest request untuk menutup poll
type ClosePollRequest struct {
	PollID      uint `json:"-" validate:"required,min=1"`
	PresenterID uint `json:"-" validate:"required,min=1"`
}

// ========================================
// Response Models
// ========================================

// PollOptionResponse response untuk single poll option
type PollOptionResponse struct {
	ID         uint    `json:"id"`
	PollID     uint    `json:"poll_id,omitempty"`
	OptionText string  `json:"option_text"`
	VoteCount  int     `json:"vote_count"`
	Order      int     `json:"order"`
	Percentage float64 `json:"percentage,omitempty"`
}

// PollResponse response untuk single poll
type PollResponse struct {
	ID          uint                 `json:"id"`
	RoomID      uint                 `json:"room_id,omitempty"`
	Question    string               `json:"question"`
	Status      string               `json:"status"`
	TotalVotes  int                  `json:"total_votes,omitempty"`
	CreatedAt   time.Time            `json:"created_at"`
	ActivatedAt *time.Time           `json:"activated_at,omitempty"`
	ClosedAt    *time.Time           `json:"closed_at,omitempty"`
	Options     []PollOptionResponse `json:"options,omitempty"`
	HasVoted    bool                 `json:"has_voted,omitempty"`
	MyVoteID    *uint                `json:"my_vote_id,omitempty"`
}

// CreatePollResponse response setelah membuat poll
type CreatePollResponse struct {
	Poll PollResponse `json:"poll"`
}

// GetActivePollsResponse response untuk active polls
type GetActivePollsResponse struct {
	Polls []PollResponse `json:"polls"`
}

// PollHistoryResponse response untuk poll history
type PollHistoryResponse struct {
	Polls []PollResponse `json:"polls"`
	Total int64          `json:"total"`
}

// PollResponseResponse response untuk single poll response (vote)
type PollResponseResponse struct {
	ID            uint      `json:"id"`
	PollID        uint      `json:"poll_id"`
	ParticipantID uint      `json:"participant_id"`
	PollOptionID  uint      `json:"poll_option_id"`
	CreatedAt     time.Time `json:"created_at"`
}

// UpdatedPollResultsResponse response untuk updated poll results setelah vote
type UpdatedPollResultsResponse struct {
	PollID     uint                 `json:"poll_id"`
	TotalVotes int                  `json:"total_votes"`
	Options    []PollOptionResponse `json:"options"`
}

// SubmitPollVoteResponse response setelah submit vote
type SubmitPollVoteResponse struct {
	Response       PollResponseResponse       `json:"response"`
	UpdatedResults UpdatedPollResultsResponse `json:"updated_results"`
	XPEarned       *XPEarned                  `json:"xp_earned,omitempty"`
}

// FinalPollResultsResponse hasil akhir poll
type FinalPollResultsResponse struct {
	TotalVotes int                  `json:"total_votes"`
	Options    []PollOptionResponse `json:"options"`
}

// ClosePollResponse response setelah close poll
type ClosePollResponse struct {
	Poll struct {
		ID           uint                     `json:"id"`
		Status       string                   `json:"status"`
		ClosedAt     *time.Time               `json:"closed_at"`
		FinalResults FinalPollResultsResponse `json:"final_results"`
	} `json:"poll"`
}
