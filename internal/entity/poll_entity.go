package entity

import "time"

type Poll struct {
	ID          uint       `gorm:"column:id;primaryKey;autoIncrement"`
	RoomID      uint       `gorm:"column:room_id;not null;index:idx_polls_room;index:idx_polls_room_status"`
	Question    string     `gorm:"column:question;type:text;not null"`
	Status      string     `gorm:"column:status;type:enum('draft','active','closed');default:'draft';not null;index:idx_polls_status;index:idx_polls_room_status"`
	CreatedAt   time.Time  `gorm:"column:created_at;autoCreateTime;not null;index:idx_polls_created_at"`
	ActivatedAt *time.Time `gorm:"column:activated_at;index:idx_polls_activated_at"`
	ClosedAt    *time.Time `gorm:"column:closed_at"`

	// Relationships
	Room          Room           `gorm:"foreignKey:RoomID;references:ID;constraint:OnDelete:CASCADE"`
	Options       []PollOption   `gorm:"foreignKey:PollID;references:ID;constraint:OnDelete:CASCADE"`
	PollResponses []PollResponse `gorm:"foreignKey:PollID;references:ID;constraint:OnDelete:CASCADE"`
}

func (p *Poll) TableName() string {
	return "polls"
}
