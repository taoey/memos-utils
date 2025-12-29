package dao

type Resource struct {
	ID        uint64 `gorm:"primaryKey;autoIncrement" json:"id"`
	UID       string `gorm:"type:text;not null;uniqueIndex" json:"uid"`
	CreatorID uint64 `gorm:"not null;index" json:"creator_id"`

	CreatedTs int64 `gorm:"not null;autoCreateTime:nano" json:"created_ts"`
	UpdatedTs int64 `gorm:"not null;autoUpdateTime:nano" json:"updated_ts"`

	Filename string `gorm:"type:text;not null;default:''" json:"filename"`
	Blob     []byte `gorm:"type:blob" json:"blob,omitempty"`

	Type string `gorm:"type:text;not null;default:''" json:"type"`
	Size int64  `gorm:"not null;default:0" json:"size"`

	MemoID uint64 `gorm:"index" json:"memo_id,omitempty"`

	StorageType string `gorm:"type:text;not null;default:''" json:"storage_type"`
	Reference   string `gorm:"type:text;not null;default:''" json:"reference"`

	Payload string `gorm:"type:text;not null;default:'{}'" json:"payload"`
}

func (Resource) TableName() string {
	return "resource"
}

type SlaveMemoResource struct {
	Resource
	MemosUid string `json:"memos_uid"`
}
