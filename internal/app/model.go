package app

import "nebula/types"

type App struct {
	ID          string         `gorm:"primaryKey" json:"id"`
	Name        string         `gorm:"not null" json:"name"`
	Description string         `json:"description"`
	CreatedAt   types.JSONTime `json:"createdAt"`
	UpdatedAt   types.JSONTime `json:"updatedAt"`
}
