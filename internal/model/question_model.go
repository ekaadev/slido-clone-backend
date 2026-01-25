package model

import "time"

// SubmitQuestionRequest request untuk submit question baru
type SubmitQuestionRequest struct {
	RoomID        uint   `json:"-" validate:"required,min=1"`
	ParticipantID uint   `json:"-" validate:"required,min=1"`
	Content       string `json:"content" validate:"required,min=1,max=1000"`
}

// GetQuestionsRequest request untuk mendapatkan list questions
type GetQuestionsRequest struct {
	RoomID        uint   `json:"-" validate:"required,min=1"`
	ParticipantID uint   `json:"-" validate:"required,min=1"`
	Status        string `json:"status" validate:"omitempty,oneof=pending answered highlighted"`
	SortBy        string `json:"sort_by" validate:"omitempty,oneof=recent upvotes validated"`
	Limit         int    `json:"limit" validate:"omitempty,min=1,max=100"`
	Offset        int    `json:"offset" validate:"omitempty,min=0"`
}

// UpvoteRequest request untuk upvote question
type UpvoteRequest struct {
	QuestionID    uint `json:"-" validate:"required,min=1"`
	ParticipantID uint `json:"-" validate:"required,min=1"`
	RoomID        uint `json:"-" validate:"required,min=1"`
}

// ValidateQuestionRequest request untuk validate question (presenter only)
type ValidateQuestionRequest struct {
	QuestionID  uint   `json:"-" validate:"required,min=1"`
	PresenterID uint   `json:"-" validate:"required,min=1"`
	Status      string `json:"status" validate:"required,oneof=answered highlighted"`
}

// QuestionResponse response untuk single question
type QuestionResponse struct {
	ID                     uint            `json:"id"`
	RoomID                 uint            `json:"room_id,omitempty"`
	ParticipantID          uint            `json:"participant_id,omitempty"`
	Participant            ParticipantInfo `json:"participant,omitempty"`
	Content                string          `json:"content"`
	UpvoteCount            int             `json:"upvote_count"`
	Status                 string          `json:"status"`
	IsValidatedByPresenter bool            `json:"is_validated_by_presenter"`
	XPAwarded              int             `json:"xp_awarded,omitempty"`
	HasVoted               bool            `json:"has_voted,omitempty"`
	CreatedAt              time.Time       `json:"created_at"`
}

// SubmitQuestionResponse response setelah submit question
type SubmitQuestionResponse struct {
	Question QuestionResponse `json:"question"`
	XPEarned *XPEarned        `json:"xp_earned,omitempty"`
}

// XPEarned representasi XP yang didapatkan
type XPEarned struct {
	Points   int `json:"points"`
	NewTotal int `json:"new_total"`
}

// QuestionListResponse response untuk list questions
type QuestionListResponse struct {
	Questions []QuestionResponse `json:"questions"`
	Paging    QuestionPaging     `json:"paging"`
}

// QuestionPaging paging info untuk questions
type QuestionPaging struct {
	Total  int64 `json:"total"`
	Limit  int   `json:"limit"`
	Offset int   `json:"offset"`
}

// VoteResponse response untuk vote
type VoteResponse struct {
	ID            uint      `json:"id"`
	QuestionID    uint      `json:"question_id"`
	ParticipantID uint      `json:"participant_id"`
	CreatedAt     time.Time `json:"created_at"`
}

// UpvoteResponse response setelah upvote
type UpvoteResponse struct {
	Vote     VoteResponse       `json:"vote"`
	Question QuestionUpvoteInfo `json:"question"`
	XPEarned *XPEarnedForUpvote `json:"xp_earned,omitempty"`
}

// QuestionUpvoteInfo info question setelah upvote
type QuestionUpvoteInfo struct {
	ID          uint `json:"id"`
	UpvoteCount int  `json:"upvote_count"`
}

// XPEarnedForUpvote XP yang didapat oleh penerima upvote
type XPEarnedForUpvote struct {
	RecipientParticipantID uint   `json:"recipient_participant_id"`
	Points                 int    `json:"points"`
	Source                 string `json:"source"`
}

// RemoveUpvoteResponse response setelah remove upvote
type RemoveUpvoteResponse struct {
	Question QuestionUpvoteInfo `json:"question"`
}

// ValidateQuestionResponse response setelah validate question
type ValidateQuestionResponse struct {
	Question  QuestionValidateInfo `json:"question"`
	XPAwarded *XPAwardedInfo       `json:"xp_awarded,omitempty"`
}

// QuestionValidateInfo info question setelah validate
type QuestionValidateInfo struct {
	ID                     uint   `json:"id"`
	Status                 string `json:"status"`
	IsValidatedByPresenter bool   `json:"is_validated_by_presenter"`
}

// XPAwardedInfo info XP yang diberikan setelah validasi
type XPAwardedInfo struct {
	ParticipantID uint `json:"participant_id"`
	Points        int  `json:"points"`
	NewTotal      int  `json:"new_total"`
}
