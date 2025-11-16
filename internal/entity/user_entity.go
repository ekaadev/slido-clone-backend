package entity

import "time"

type User struct {
	ID           uint      `gorm:"column:id;primaryKey;autoIncrement"`
	Username     string    `gorm:"column:username;type:varchar(100);uniqueIndex;not null"`
	Email        string    `gorm:"column:email;type:varchar(255);uniqueIndex;not null"`
	PasswordHash string    `gorm:"column:password_hash;type:varchar(255);not null"`
	Role         string    `gorm:"column:role;type:enum('presenter','admin');default:'presenter';not null"`
	CreatedAt    time.Time `gorm:"column:created_at;autoCreateTime;not null"`
	UpdatedAt    time.Time `gorm:"column:updated_at;autoCreateTime;autoUpdateTime;not null"`

	// Relationships
	Rooms []Room `gorm:"foreignKey:PresenterID;references:ID"`
}

func (u *User) TableName() string {
	return "users"
}
