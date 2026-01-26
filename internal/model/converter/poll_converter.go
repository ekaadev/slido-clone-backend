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

// PollOptionToResponseWithPercentage convert dengan percentage
func PollOptionToResponseWithPercentage(option *entity.PollOption, totalVotes int) model.PollOptionResponse {
	percentage := 0.0
	if totalVotes > 0 {
		percentage = float64(option.VoteCount) / float64(totalVotes) * 100
		// round to 2 decimal places
		percentage = float64(int(percentage*100)) / 100
	}
	return model.PollOptionResponse{
		ID:         option.ID,
		OptionText: option.OptionText,
		VoteCount:  option.VoteCount,
		Percentage: percentage,
		Order:      option.Order,
	}
}

// PollToResponse convert entity Poll to model PollResponse
func PollToResponse(poll *entity.Poll) model.PollResponse {
	options := make([]model.PollOptionResponse, len(poll.Options))
	for i, opt := range poll.Options {
		options[i] = PollOptionToResponse(&opt)
	}

	return model.PollResponse{
		ID:          poll.ID,
		RoomID:      poll.RoomID,
		Question:    poll.Question,
		Status:      poll.Status,
		CreatedAt:   poll.CreatedAt,
		ActivatedAt: poll.ActivatedAt,
		ClosedAt:    poll.ClosedAt,
		Options:     options,
	}
}

// PollToResponseWithDetails convert dengan total votes dan voted info
func PollToResponseWithDetails(poll *entity.Poll, totalVotes int, hasVoted bool, votedOption *uint) model.PollResponse {
	options := make([]model.PollOptionResponse, len(poll.Options))
	for i, opt := range poll.Options {
		options[i] = PollOptionToResponseWithPercentage(&opt, totalVotes)
	}

	return model.PollResponse{
		ID:          poll.ID,
		Question:    poll.Question,
		Status:      poll.Status,
		TotalVotes:  totalVotes,
		CreatedAt:   poll.CreatedAt,
		ActivatedAt: poll.ActivatedAt,
		ClosedAt:    poll.ClosedAt,
		Options:     options,
		HasVoted:    hasVoted,
		VotedOption: votedOption,
	}
}

// PollToCreateResponse convert untuk create response
func PollToCreateResponse(poll *entity.Poll) *model.CreatePollResponse {
	return &model.CreatePollResponse{
		Poll: PollToResponse(poll),
	}
}

// PollResponseToVoteData convert entity PollResponse ke vote response data
func PollResponseToVoteData(response *entity.PollResponse) model.PollVoteResponseData {
	return model.PollVoteResponseData{
		ID:            response.ID,
		PollID:        response.PollID,
		ParticipantID: response.ParticipantID,
		PollOptionID:  response.PollOptionID,
		CreatedAt:     response.CreatedAt,
	}
}

// PollOptionsToUpdatedResults convert options ke updated results
func PollOptionsToUpdatedResults(pollID uint, options []entity.PollOption, totalVotes int) model.UpdatedResultsResponse {
	optResponses := make([]model.PollOptionResponse, len(options))
	for i, opt := range options {
		optResponses[i] = PollOptionToResponseWithPercentage(&opt, totalVotes)
	}

	return model.UpdatedResultsResponse{
		PollID:     pollID,
		TotalVotes: totalVotes,
		Options:    optResponses,
	}
}

// PollToHistoryItem convert poll to history item
func PollToHistoryItem(poll *entity.Poll, totalVotes int) model.PollHistoryItem {
	return model.PollHistoryItem{
		ID:         poll.ID,
		Question:   poll.Question,
		Status:     poll.Status,
		TotalVotes: totalVotes,
		CreatedAt:  poll.CreatedAt,
		ClosedAt:   poll.ClosedAt,
	}
}
