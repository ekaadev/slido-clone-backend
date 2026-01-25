package converter

import (
	"slido-clone-backend/internal/entity"
	"slido-clone-backend/internal/model"
)

// QuestionToResponse convert entity Question to model QuestionResponse
func QuestionToResponse(question *entity.Question) *model.QuestionResponse {
	return &model.QuestionResponse{
		ID:                     question.ID,
		RoomID:                 question.RoomID,
		ParticipantID:          question.ParticipantID,
		Content:                question.Content,
		UpvoteCount:            question.UpvoteCount,
		Status:                 question.Status,
		IsValidatedByPresenter: question.IsValidatedByPresenter,
		XPAwarded:              question.XPAwarded,
		CreatedAt:              question.CreatedAt,
	}
}

// QuestionToResponseWithParticipant convert entity Question to model QuestionResponse dengan participant info
func QuestionToResponseWithParticipant(question *entity.Question, hasVoted bool) *model.QuestionResponse {
	return &model.QuestionResponse{
		ID: question.ID,
		Participant: model.ParticipantInfo{
			ID:          question.Participant.ID,
			DisplayName: question.Participant.DisplayName,
		},
		Content:                question.Content,
		UpvoteCount:            question.UpvoteCount,
		Status:                 question.Status,
		IsValidatedByPresenter: question.IsValidatedByPresenter,
		HasVoted:               hasVoted,
		CreatedAt:              question.CreatedAt,
	}
}

// QuestionToSubmitResponse convert untuk submit response
func QuestionToSubmitResponse(question *entity.Question, xpPoints int, newTotal int) *model.SubmitQuestionResponse {
	return &model.SubmitQuestionResponse{
		Question: model.QuestionResponse{
			ID:                     question.ID,
			RoomID:                 question.RoomID,
			ParticipantID:          question.ParticipantID,
			Content:                question.Content,
			UpvoteCount:            question.UpvoteCount,
			Status:                 question.Status,
			IsValidatedByPresenter: question.IsValidatedByPresenter,
			XPAwarded:              question.XPAwarded,
			CreatedAt:              question.CreatedAt,
		},
		XPEarned: &model.XPEarned{
			Points:   xpPoints,
			NewTotal: newTotal,
		},
	}
}

// VoteToResponse convert entity Vote to model VoteResponse
func VoteToResponse(vote *entity.Vote) *model.VoteResponse {
	return &model.VoteResponse{
		ID:            vote.ID,
		QuestionID:    vote.QuestionID,
		ParticipantID: vote.ParticipantID,
		CreatedAt:     vote.CreatedAt,
	}
}

// VoteToUpvoteResponse convert untuk upvote response
func VoteToUpvoteResponse(vote *entity.Vote, upvoteCount int, recipientID uint, xpPoints int) *model.UpvoteResponse {
	return &model.UpvoteResponse{
		Vote: model.VoteResponse{
			ID:            vote.ID,
			QuestionID:    vote.QuestionID,
			ParticipantID: vote.ParticipantID,
			CreatedAt:     vote.CreatedAt,
		},
		Question: model.QuestionUpvoteInfo{
			ID:          vote.QuestionID,
			UpvoteCount: upvoteCount,
		},
		XPEarned: &model.XPEarnedForUpvote{
			RecipientParticipantID: recipientID,
			Points:                 xpPoints,
			Source:                 "upvote_received",
		},
	}
}

// QuestionToRemoveUpvoteResponse convert untuk remove upvote response
func QuestionToRemoveUpvoteResponse(questionID uint, upvoteCount int) *model.RemoveUpvoteResponse {
	return &model.RemoveUpvoteResponse{
		Question: model.QuestionUpvoteInfo{
			ID:          questionID,
			UpvoteCount: upvoteCount,
		},
	}
}

// QuestionToValidateResponse convert untuk validate response
func QuestionToValidateResponse(question *entity.Question, xpPoints int, newTotal int) *model.ValidateQuestionResponse {
	return &model.ValidateQuestionResponse{
		Question: model.QuestionValidateInfo{
			ID:                     question.ID,
			Status:                 question.Status,
			IsValidatedByPresenter: question.IsValidatedByPresenter,
		},
		XPAwarded: &model.XPAwardedInfo{
			ParticipantID: question.ParticipantID,
			Points:        xpPoints,
			NewTotal:      newTotal,
		},
	}
}
