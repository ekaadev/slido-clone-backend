package entity

import "time"

type Room struct {
	ID          uint       `gorm:"column:id;primaryKey;autoIncrement"`
	RoomCode    string     `gorm:"column:room_code;type:varchar(20);uniqueIndex;not null"`
	Title       string     `gorm:"column:title;type:varchar(255);not null"`
	PresenterID uint       `gorm:"column:presenter_id;not null;index"`
	Status      string     `gorm:"column:status;type:enum('active','closed');default:'active';not null;index"`
	CreatedAt   time.Time  `gorm:"column:created_at;autoCreateTime;not null;index:idx_rooms_created_at"`
	ClosedAt    *time.Time `gorm:"column:closed_at"`

	// Relationships
	Presenter      User            `gorm:"foreignKey:PresenterID;references:ID;constraint:OnDelete:CASCADE"`
	Participants   []Participant   `gorm:"foreignKey:RoomID;references:ID;constraint:OnDelete:CASCADE"`
	Questions      []Question      `gorm:"foreignKey:RoomID;references:ID;constraint:OnDelete:CASCADE"`
	Polls          []Poll          `gorm:"foreignKey:RoomID;references:ID;constraint:OnDelete:CASCADE"`
	Messages       []Message       `gorm:"foreignKey:RoomID;references:ID;constraint:OnDelete:CASCADE"`
	XPTransactions []XPTransaction `gorm:"foreignKey:RoomID;references:ID;constraint:OnDelete:CASCADE"`
}

func (r *Room) TableName() string {
	return "rooms"
}
