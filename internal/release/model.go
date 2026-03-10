package release

import "time"

type Release struct {
	ID        uint   `gorm:"primaryKey"`
	AppID     string `gorm:"not null"`
	Version   string `gorm:"not null"`
	Notes     string
	Channel   string
	PubDate   time.Time
	CreatedAt time.Time
	UpdatedAt time.Time
}
