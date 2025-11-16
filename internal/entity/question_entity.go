package entity

import "time"

type Question struct {
	ID                     uint      `gorm:"column:id;primaryKey;autoIncrement"`
	RoomID                 uint      `gorm:"column:room_id;not null;index:idx_questions_room"`
	ParticipantID          uint      `gorm:"column:participant_id;not null;index:idx_questions_participant"`
	Content                string    `gorm:"column:content;type:text;not null"`
	UpvoteCount            int       `gorm:"column:upvote_count;type:int unsigned;default:0;not null;index:idx_questions_upvote_count"`
	Status                 string    `gorm:"column:status;type:enum('pending','answered','highlighted');default:'pending';not null;index:idx_questions_status"`
	IsValidatedByPresenter bool      `gorm:"column:is_validated_by_presenter;default:false;not null;index:idx_questions_validated"`
	XPAwarded              int       `gorm:"column:xp_awarded;type:int unsigned;default:0;not null"`
	CreatedAt              time.Time `gorm:"column:created_at;autoCreateTime;not null;index:idx_questions_created_at"`

	// Relationships
	Room        Room        `gorm:"foreignKey:RoomID;references:ID;constraint:OnDelete:CASCADE"`
	Participant Participant `gorm:"foreignKey:ParticipantID;references:ID;constraint:OnDelete:CASCADE"`
	Votes       []Vote      `gorm:"foreignKey:QuestionID;references:ID;constraint:OnDelete:CASCADE"`
}

func (q *Question) TableName() string {
	return "questions"
}
