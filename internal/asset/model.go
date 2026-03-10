package asset

import "time"

type Asset struct {
	ID          uint   `gorm:"primaryKey" json:"id"`
	ReleaseID   uint   `gorm:"not null" json:"releaseId"`
	Platform    string `json:"platform"`
	Arch        string `json:"arch"`
	StoragePath string `json:"-"` // 存储路径，不在 API 中返回

	URL       string `gorm:"not null" json:"url"`
	Signature string `json:"signature"`
	Checksum  string `json:"checksum"`

	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}
