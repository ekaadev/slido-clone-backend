package converter

import (
	"slido-clone-backend/internal/entity"
	"slido-clone-backend/internal/model"
)

// ParticipantToResponse convert entity Participant to model ParticipantResponse
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

// ParticipantToResponseWithRole convert entity Participant to model ParticipantResponse with room role
func ParticipantToResponseWithRole(participant *entity.Participant, isRoomOwner bool) *model.ParticipantResponse {
	roomRole := "audience"
	if isRoomOwner {
		roomRole = "host"
	}
	return &model.ParticipantResponse{
		ID:          participant.ID,
		RoomID:      participant.RoomID,
		DisplayName: participant.DisplayName,
		XPScore:     participant.XPScore,
		IsAnonymous: *participant.IsAnonymous,
		RoomRole:    roomRole,
		JoinedAt:    participant.JoinedAt,
	}
}

// ParticipantToJoinRoomResponse convert entity Participant and token to model JoinRoomResponse
func ParticipantToJoinRoomResponse(participant *entity.Participant, token string) *model.JoinRoomResponse {
	return &model.JoinRoomResponse{
		Participant: *ParticipantToResponse(participant),
		Token:       token,
	}
}

// ParticipantToJoinRoomResponseWithRole convert entity Participant and token to model JoinRoomResponse with room role
func ParticipantToJoinRoomResponseWithRole(participant *entity.Participant, token string, isRoomOwner bool) *model.JoinRoomResponse {
	return &model.JoinRoomResponse{
		Participant: *ParticipantToResponseWithRole(participant, isRoomOwner),
		Token:       token,
	}
}

// ParticipantToListItem convert entity Participant to model ParticipantListItem
func ParticipantToListItem(participant *entity.Participant) *model.ParticipantListItem {
	return &model.ParticipantListItem{
		ID:          participant.ID,
		DisplayName: participant.DisplayName,
		XPScore:     participant.XPScore,
		IsAnonymous: *participant.IsAnonymous,
	}
}

// ParticipantsToListResponse wrap list of ParticipantListItem to ParticipantListResponse
func ParticipantsToListResponse(participants []*model.ParticipantListItem) *model.ParticipantListResponse {
	return &model.ParticipantListResponse{
		Participants: participants,
	}
}

// ParticipantToInfo convert entity Participant to model ParticipantInfo
func ParticipantToInfo(participant *entity.Participant) *model.ParticipantInfo {
	return &model.ParticipantInfo{
		ID:          participant.ID,
		DisplayName: participant.DisplayName,
	}
}

// ParticipantToLeaderboardEntry convert entity Participant and rank to model LeaderboardEntry
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

// ParticipantsToLeaderboardResponse convert list of Participant to LeaderboardResponse with myRank info
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
