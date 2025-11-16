package entity

import "time"

type Message struct {
	ID            uint      `gorm:"column:id;primaryKey;autoIncrement"`
	RoomID        uint      `gorm:"column:room_id;not null;index:idx_messages_room;index:idx_messages_room_created"`
	ParticipantID uint      `gorm:"column:participant_id;not null;index:idx_messages_participant"`
	Content       string    `gorm:"column:content;type:text;not null"`
	CreatedAt     time.Time `gorm:"column:created_at;autoCreateTime;not null;index:idx_messages_room_created"`

	// Relationships
	Room        Room        `gorm:"foreignKey:RoomID;references:ID;constraint:OnDelete:CASCADE"`
	Participant Participant `gorm:"foreignKey:ParticipantID;references:ID;constraint:OnDelete:CASCADE"`
}

func (m *Message) TableName() string {
	return "messages"
}
