package converter

import (
	"slido-clone-backend/internal/entity"
	"slido-clone-backend/internal/model"
)

// PollOptionToResponse convert entity PollOption to model PollOptionResponse
func PollOptionToResponse(option *entity.PollOption) model.PollOptionResponse {
	return model.PollOptionResponse{
		ID:         option.ID,
		PollID:     option.PollID,
		OptionText: option.OptionText,
		VoteCount:  option.VoteCount,
		Order:      option.Order,
	}
}

// PollOptionToResponseWithPercentage convert entity PollOption dengan percentage
func PollOptionToResponseWithPercentage(option *entity.PollOption, totalVotes int) model.PollOptionResponse {
	var percentage float64
	if totalVotes > 0 {
		percentage = float64(option.VoteCount) / float64(totalVotes) * 100
		// Round to 2 decimal places
		percentage = float64(int(percentage*100)) / 100
	}
	return model.PollOptionResponse{
		ID:         option.ID,
		PollID:     option.PollID,
		OptionText: option.OptionText,
		VoteCount:  option.VoteCount,
		Order:      option.Order,
		Percentage: percentage,
	}
}

// PollOptionsToResponse convert slice of PollOption to slice of PollOptionResponse
func PollOptionsToResponse(options []entity.PollOption) []model.PollOptionResponse {
	result := make([]model.PollOptionResponse, len(options))
	for i, option := range options {
		result[i] = PollOptionToResponse(&option)
	}
	return result
}

// PollOptionsToResponseWithPercentage convert slice of PollOption dengan percentage
func PollOptionsToResponseWithPercentage(options []entity.PollOption, totalVotes int) []model.PollOptionResponse {
	result := make([]model.PollOptionResponse, len(options))
	for i, option := range options {
		result[i] = PollOptionToResponseWithPercentage(&option, totalVotes)
	}
	return result
}

// PollToResponse convert entity Poll to model PollResponse
func PollToResponse(poll *entity.Poll) *model.PollResponse {
	return &model.PollResponse{
		ID:          poll.ID,
		RoomID:      poll.RoomID,
		Question:    poll.Question,
		Status:      poll.Status,
		CreatedAt:   poll.CreatedAt,
		ActivatedAt: poll.ActivatedAt,
		ClosedAt:    poll.ClosedAt,
		Options:     PollOptionsToResponse(poll.Options),
	}
}

// PollToResponseWithOptions convert entity Poll dengan options yang sudah ada
func PollToResponseWithOptions(poll *entity.Poll, totalVotes int, hasVoted bool, myVoteID *uint) *model.PollResponse {
	return &model.PollResponse{
		ID:          poll.ID,
		RoomID:      poll.RoomID,
		Question:    poll.Question,
		Status:      poll.Status,
		TotalVotes:  totalVotes,
		CreatedAt:   poll.CreatedAt,
		ActivatedAt: poll.ActivatedAt,
		ClosedAt:    poll.ClosedAt,
		Options:     PollOptionsToResponseWithPercentage(poll.Options, totalVotes),
		HasVoted:    hasVoted,
		MyVoteID:    myVoteID,
	}
}

// PollToCreateResponse convert untuk create poll response
func PollToCreateResponse(poll *entity.Poll) *model.CreatePollResponse {
	return &model.CreatePollResponse{
		Poll: model.PollResponse{
			ID:        poll.ID,
			RoomID:    poll.RoomID,
			Question:  poll.Question,
			Status:    poll.Status,
			CreatedAt: poll.CreatedAt,
			Options:   PollOptionsToResponse(poll.Options),
		},
	}
}

// PollResponseToResponse convert entity PollResponse (vote) to model PollResponseResponse
func PollResponseToResponse(response *entity.PollResponse) model.PollResponseResponse {
	return model.PollResponseResponse{
		ID:            response.ID,
		PollID:        response.PollID,
		ParticipantID: response.ParticipantID,
		PollOptionID:  response.PollOptionID,
		CreatedAt:     response.CreatedAt,
	}
}

// PollToVoteResponse convert untuk submit vote response
func PollToVoteResponse(
	response *entity.PollResponse,
	poll *entity.Poll,
	totalVotes int,
	xpPoints int,
	newTotal int,
) *model.SubmitPollVoteResponse {
	return &model.SubmitPollVoteResponse{
		Response: PollResponseToResponse(response),
		UpdatedResults: model.UpdatedPollResultsResponse{
			PollID:     poll.ID,
			TotalVotes: totalVotes,
			Options:    PollOptionsToResponseWithPercentage(poll.Options, totalVotes),
		},
		XPEarned: &model.XPEarned{
			Points:   xpPoints,
			NewTotal: newTotal,
		},
	}
}

// PollToCloseResponse convert untuk close poll response
func PollToCloseResponse(poll *entity.Poll, totalVotes int) *model.ClosePollResponse {
	response := &model.ClosePollResponse{}
	response.Poll.ID = poll.ID
	response.Poll.Status = poll.Status
	response.Poll.ClosedAt = poll.ClosedAt
	response.Poll.FinalResults = model.FinalPollResultsResponse{
		TotalVotes: totalVotes,
		Options:    PollOptionsToResponseWithPercentage(poll.Options, totalVotes),
	}
	return response
}

// PollsToHistoryResponse convert untuk poll history response
func PollsToHistoryResponse(polls []entity.Poll, total int64) *model.PollHistoryResponse {
	result := make([]model.PollResponse, len(polls))
	for i, poll := range polls {
		totalVotes := 0
		for _, opt := range poll.Options {
			totalVotes += opt.VoteCount
		}
		result[i] = model.PollResponse{
			ID:         poll.ID,
			Question:   poll.Question,
			Status:     poll.Status,
			TotalVotes: totalVotes,
			CreatedAt:  poll.CreatedAt,
			ClosedAt:   poll.ClosedAt,
		}
	}
	return &model.PollHistoryResponse{
		Polls: result,
		Total: total,
	}
}
