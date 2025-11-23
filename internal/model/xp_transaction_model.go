package model

// LeaderboardEntry representasi salah satu data pada leaderboard xp
type LeaderboardEntry struct {
	Rank        int             `json:"rank"`
	Participant ParticipantInfo `json:"participant"`
	XPScore     int             `json:"xp_score"`
	IsAnonymous bool            `json:"is_anonymous"`
}

// MyRank representasi rank user dengan id yang dibawa dari token jwt
type MyRank struct {
	Rank    int `json:"rank"`
	XPScore int `json:"xp_score"`
}

// LeaderboardResponse representasi response leaderboard xp
type LeaderboardResponse struct {
	Leaderboard       []LeaderboardEntry `json:"leaderboard"`
	MyRank            *MyRank            `json:"my_rank"`
	TotalParticipants int                `json:"total_participants"`
}

// GetLeaderboardRequest representasi request untuk mendapatkan leaderboard xp
type GetLeaderboardRequest struct {
	RoomID        uint `json:"-" validate:"required,min=1"`
	ParticipantID uint `json:"participant_id" validate:"required,min=1"`
}
