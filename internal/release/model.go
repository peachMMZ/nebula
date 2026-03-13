package release

import "nebula/types"

type Release struct {
	ID        uint   `gorm:"primaryKey"`
	AppID     string `gorm:"not null"`
	Tag       string `gorm:"not null;index:idx_app_tag,unique"` // Git 标签,如 v1.0.0
	Version   string `gorm:"not null"`                          // 版本号,如 1.0.0
	Notes     string
	Channel   string
	PubDate   types.JSONTime
	CreatedAt types.JSONTime
	UpdatedAt types.JSONTime
}
