package dao

type MemoRelationDTO struct {
	Type           string `json:"type"`
	UID            string `json:"uid"`
	RelatedMemoUID string `json:"related_memo_uid"`
}

type MemoRelation struct {
	MemoID        uint64 `gorm:"not null;uniqueIndex:uk_memo_relation" json:"memo_id"`
	RelatedMemoID uint64 `gorm:"not null;uniqueIndex:uk_memo_relation" json:"related_memo_id"`
	Type          string `gorm:"type:text;not null;uniqueIndex:uk_memo_relation" json:"type"`
}

func (MemoRelation) TableName() string {
	return "memo_relation"
}
