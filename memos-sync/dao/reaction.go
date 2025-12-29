package dao

type Reaction struct {
	ID uint64 `gorm:"primaryKey;autoIncrement" json:"id"`

	CreatedTs    int64  `gorm:"not null;autoCreateTime:nano" json:"created_ts"`
	CreatorID    uint64 `gorm:"not null;index;uniqueIndex:uk_reaction" json:"creator_id"`
	ContentID    string `gorm:"type:text;not null;uniqueIndex:uk_reaction" json:"content_id"`
	ReactionType string `gorm:"type:text;not null;uniqueIndex:uk_reaction" json:"reaction_type"`
}

func (Reaction) TableName() string {
	return "reaction"
}
