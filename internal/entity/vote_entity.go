package entity

import "time"

type Vote struct {
	ID            uint      `gorm:"column:id;primaryKey;autoIncrement"`
	QuestionID    uint      `gorm:"column:question_id;not null;index:idx_votes_question;uniqueIndex:unique_vote_per_question"`
	ParticipantID uint      `gorm:"column:participant_id;not null;index:idx_votes_participant;uniqueIndex:unique_vote_per_question"`
	CreatedAt     time.Time `gorm:"column:created_at;autoCreateTime;not null;index:idx_votes_created_at"`

	// Relationships
	Question    Question    `gorm:"foreignKey:QuestionID;references:ID;constraint:OnDelete:CASCADE"`
	Participant Participant `gorm:"foreignKey:ParticipantID;references:ID;constraint:OnDelete:CASCADE"`
}

func (v *Vote) TableName() string {
	return "votes"
}
