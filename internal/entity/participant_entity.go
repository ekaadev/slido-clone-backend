package entity

import "time"

type Participant struct {
	ID          uint      `gorm:"column:id;primaryKey;autoIncrement"`
	RoomID      uint      `gorm:"column:room_id;not null;index:idx_participants_room"`
	UserID      *uint     `gorm:"column:user_id;index"` // NULL for anonymous
	DisplayName string    `gorm:"column:display_name;type:varchar(100);not null"`
	XPScore     int       `gorm:"column:xp_score;type:int unsigned;default:0;not null;index:idx_participants_xp_score"`
	IsAnonymous *bool     `gorm:"column:is_anonymous;default:true;not null"`
	JoinedAt    time.Time `gorm:"column:joined_at;autoCreateTime;not null;index:idx_participants_joined_at"`

	// Relationships
	Room           Room            `gorm:"foreignKey:RoomID;references:ID;constraint:OnDelete:CASCADE"`
	User           *User           `gorm:"foreignKey:UserID;references:ID;constraint:OnDelete:SET NULL"`
	Questions      []Question      `gorm:"foreignKey:ParticipantID;references:ID;constraint:OnDelete:CASCADE"`
	Votes          []Vote          `gorm:"foreignKey:ParticipantID;references:ID;constraint:OnDelete:CASCADE"`
	PollResponses  []PollResponse  `gorm:"foreignKey:ParticipantID;references:ID;constraint:OnDelete:CASCADE"`
	Messages       []Message       `gorm:"foreignKey:ParticipantID;references:ID;constraint:OnDelete:CASCADE"`
	XPTransactions []XPTransaction `gorm:"foreignKey:ParticipantID;references:ID;constraint:OnDelete:CASCADE"`
}

func (p *Participant) TableName() string {
	return "participants"
}
