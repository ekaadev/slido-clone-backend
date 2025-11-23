package entity

import "time"

type XPTransaction struct {
	ID            uint      `gorm:"column:id;primaryKey;autoIncrement"`
	ParticipantID uint      `gorm:"column:participant_id;not null;index:idx_xp_transactions_participant"`
	RoomID        uint      `gorm:"column:room_id;not null;index:idx_xp_transactions_room"`
	Points        int       `gorm:"column:points;not null"` // Can be positive or negative
	SourceType    string    `gorm:"column:source_type;type:enum('poll','question_created','upvote_received','presenter_validated', 'message_created');not null;index:idx_xp_transactions_source"`
	SourceID      uint      `gorm:"column:source_id;not null;index:idx_xp_transactions_source"` // Polymorphic
	CreatedAt     time.Time `gorm:"column:created_at;autoCreateTime;not null;index:idx_xp_transactions_created_at"`

	// Relationships
	Participant Participant `gorm:"foreignKey:ParticipantID;references:ID;constraint:OnDelete:CASCADE"`
	Room        Room        `gorm:"foreignKey:RoomID;references:ID;constraint:OnDelete:CASCADE"`
}

func (xp *XPTransaction) TableName() string {
	return "xp_transactions"
}
