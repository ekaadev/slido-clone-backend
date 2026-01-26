package model

import "time"

// CreatePollRequest request untuk membuat poll baru
type CreatePollRequest struct {
	RoomID      uint     `json:"-" validate:"required,min=1"`
	PresenterID uint     `json:"-" validate:"required,min=1"`
	Question    string   `json:"question" validate:"required,min=5,max=500"`
	Options     []string `json:"options" validate:"required,min=2,max=10,dive,required,min=1,max=255"`
}

// SubmitVoteRequest request untuk submit vote pada poll
type SubmitVoteRequest struct {
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

// GetPollHistoryRequest request untuk mendapatkan history poll
type GetPollHistoryRequest struct {
	RoomID        uint   `json:"-" validate:"required,min=1"`
	ParticipantID uint   `json:"-" validate:"required,min=1"`
	Status        string `json:"status" validate:"omitempty,oneof=active closed all"`
	Limit         int    `json:"limit" validate:"omitempty,min=1,max=50"`
}

// GetActivePollRequest request untuk mendapatkan active poll
type GetActivePollRequest struct {
	RoomID        uint `json:"-" validate:"required,min=1"`
	ParticipantID uint `json:"-" validate:"required,min=1"`
}

// PollOptionResponse response untuk single poll option
type PollOptionResponse struct {
	ID         uint    `json:"id"`
	PollID     uint    `json:"poll_id,omitempty"`
	OptionText string  `json:"option_text"`
	VoteCount  int     `json:"vote_count"`
	Percentage float64 `json:"percentage,omitempty"`
	Order      int     `json:"order"`
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
	VotedOption *uint                `json:"voted_option,omitempty"`
}

// CreatePollResponse response setelah create poll
type CreatePollResponse struct {
	Poll PollResponse `json:"poll"`
}

// GetActivePollResponse response untuk active polls
type GetActivePollResponse struct {
	Polls []PollResponse `json:"polls"`
}

// PollVoteResponseData response untuk vote
type PollVoteResponseData struct {
	ID            uint      `json:"id"`
	PollID        uint      `json:"poll_id"`
	ParticipantID uint      `json:"participant_id"`
	PollOptionID  uint      `json:"poll_option_id"`
	CreatedAt     time.Time `json:"created_at"`
}

// UpdatedResultsResponse updated results setelah vote
type UpdatedResultsResponse struct {
	PollID     uint                 `json:"poll_id"`
	TotalVotes int                  `json:"total_votes"`
	Options    []PollOptionResponse `json:"options"`
}

// SubmitVoteResponse response setelah submit vote
type SubmitVoteResponse struct {
	Response       PollVoteResponseData   `json:"response"`
	UpdatedResults UpdatedResultsResponse `json:"updated_results"`
	XPEarned       *XPEarned              `json:"xp_earned,omitempty"`
}

// FinalResultsResponse final results saat poll di-close
type FinalResultsResponse struct {
	TotalVotes int                  `json:"total_votes"`
	Options    []PollOptionResponse `json:"options"`
}

// ClosePollResponse response setelah close poll
type ClosePollResponse struct {
	Poll struct {
		ID           uint                 `json:"id"`
		Status       string               `json:"status"`
		ClosedAt     *time.Time           `json:"closed_at"`
		FinalResults FinalResultsResponse `json:"final_results"`
	} `json:"poll"`
}

// PollHistoryItem item untuk poll history
type PollHistoryItem struct {
	ID         uint       `json:"id"`
	Question   string     `json:"question"`
	Status     string     `json:"status"`
	TotalVotes int        `json:"total_votes"`
	CreatedAt  time.Time  `json:"created_at"`
	ClosedAt   *time.Time `json:"closed_at,omitempty"`
}

// GetPollHistoryResponse response untuk poll history
type GetPollHistoryResponse struct {
	Polls []PollHistoryItem `json:"polls"`
	Total int64             `json:"total"`
}
