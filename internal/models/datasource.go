package models

import (
	"gorm.io/gorm"
)

type DataSource struct {
	gorm.Model
	Name        string `json:"name" gorm:"uniqueIndex"`
	Description string `json:"description"`
	URL         string `json:"url" gorm:"not null"`
	Type        string `json:"type" gorm:"not null"` // e.g., "stock", "forex", "crypto"
}