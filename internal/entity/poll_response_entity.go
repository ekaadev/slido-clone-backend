package entity

import "time"

type PollResponse struct {
	ID            uint      `gorm:"column:id;primaryKey;autoIncrement"`
	PollID        uint      `gorm:"column:poll_id;not null;index:idx_poll_responses_poll;uniqueIndex:unique_poll_response"`
	ParticipantID uint      `gorm:"column:participant_id;not null;index:idx_poll_responses_participant;uniqueIndex:unique_poll_response"`
	PollOptionID  uint      `gorm:"column:poll_option_id;not null;index:idx_poll_responses_option"`
	CreatedAt     time.Time `gorm:"column:created_at;autoCreateTime;not null;index:idx_poll_responses_created_at"`

	// Relationships
	Poll        Poll        `gorm:"foreignKey:PollID;references:ID;constraint:OnDelete:CASCADE"`
	Participant Participant `gorm:"foreignKey:ParticipantID;references:ID;constraint:OnDelete:CASCADE"`
	PollOption  PollOption  `gorm:"foreignKey:PollOptionID;references:ID;constraint:OnDelete:CASCADE"`
}

func (pr *PollResponse) TableName() string {
	return "poll_responses"
}
