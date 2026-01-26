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

// GetXPTransactionsRequest request untuk get XP transactions
type GetXPTransactionsRequest struct {
	RoomID        uint `json:"-" validate:"required,min=1"`
	ParticipantID uint `json:"-" validate:"required,min=1"`
	Limit         int  `json:"limit" validate:"omitempty,min=1,max=100"`
}

// XPTransactionItem single XP transaction item
type XPTransactionItem struct {
	ID         uint   `json:"id"`
	Points     int    `json:"points"`
	SourceType string `json:"source_type"`
	SourceID   uint   `json:"source_id,omitempty"`
	CreatedAt  string `json:"created_at"`
}

// GetXPTransactionsResponse response untuk get XP transactions
type GetXPTransactionsResponse struct {
	Transactions []XPTransactionItem `json:"transactions"`
	TotalXP      int                 `json:"total_xp"`
	Total        int64               `json:"total"`
}
