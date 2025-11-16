package entity

type PollOption struct {
	ID         uint   `gorm:"column:id;primaryKey;autoIncrement"`
	PollID     uint   `gorm:"column:poll_id;not null;index:idx_poll_options_poll;index:idx_poll_options_order"`
	OptionText string `gorm:"column:option_text;type:varchar(255);not null"`
	VoteCount  int    `gorm:"column:vote_count;type:int unsigned;default:0;not null"` // Denormalized
	Order      int    `gorm:"column:order;type:tinyint unsigned;not null;index:idx_poll_options_order"`

	// Relationships
	Poll          Poll           `gorm:"foreignKey:PollID;references:ID;constraint:OnDelete:CASCADE"`
	PollResponses []PollResponse `gorm:"foreignKey:PollOptionID;references:ID;constraint:OnDelete:CASCADE"`
}

func (po *PollOption) TableName() string {
	return "poll_options"
}
