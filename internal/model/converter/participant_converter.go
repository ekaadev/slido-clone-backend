package converter

import (
	"slido-clone-backend/internal/entity"
	"slido-clone-backend/internal/model"
)

func ParticipantToResponse(participant *entity.Participant) *model.ParticipantResponse {
	return &model.ParticipantResponse{
		ID:          participant.ID,
		RoomID:      participant.RoomID,
		DisplayName: participant.DisplayName,
		XPScore:     participant.XPScore,
		IsAnonymous: *participant.IsAnonymous,
		JoinedAt:    participant.JoinedAt,
	}
}

func ParticipantToJoinRoomResponse(participant *entity.Participant, token string) *model.JoinRoomResponse {
	return &model.JoinRoomResponse{
		Participant: *ParticipantToResponse(participant),
		Token:       token,
	}
}

func ParticipantToListItem(participant *entity.Participant) *model.ParticipantListItem {
	return &model.ParticipantListItem{
		ID:          participant.ID,
		DisplayName: participant.DisplayName,
		XPScore:     participant.XPScore,
		IsAnonymous: *participant.IsAnonymous,
	}
}

func ParticipantsToListResponse(participants []*model.ParticipantListItem) *model.ParticipantListResponse {
	return &model.ParticipantListResponse{
		Participants: participants,
	}
}

func ParticipantToInfo(participant *entity.Participant) *model.ParticipantInfo {
	return &model.ParticipantInfo{
		ID:          participant.ID,
		DisplayName: participant.DisplayName,
	}
}

func ParticipantToLeaderboardEntry(participant *entity.Participant, rank int) *model.LeaderboardEntry {
	return &model.LeaderboardEntry{
		Rank: rank,
		Participant: model.ParticipantInfo{
			ID:          participant.ID,
			DisplayName: participant.DisplayName,
		},
		XPScore:     participant.XPScore,
		IsAnonymous: *participant.IsAnonymous,
	}
}

func ParticipantsToLeaderboardResponse(participants []entity.Participant, myRank *model.MyRank, totalParticipants int) *model.LeaderboardResponse {
	leaderboard := make([]model.LeaderboardEntry, len(participants))
	for i, participant := range participants {
		leaderboard[i] = *ParticipantToLeaderboardEntry(&participant, i+1)
	}

	return &model.LeaderboardResponse{
		Leaderboard:       leaderboard,
		MyRank:            myRank,
		TotalParticipants: totalParticipants,
	}
}
