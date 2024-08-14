package models

import (
	"time"
	"gorm.io/gorm"
)

type Trade struct {
	gorm.Model
	UserID    uint      `json:"user_id" gorm:"index"`
	StrategyID uint     `json:"strategy_id" gorm:"index"`
	Symbol    string    `json:"symbol" gorm:"index"`
	Type      string    `json:"type" gorm:"index"` // "buy" or "sell"
	Quantity  float64   `json:"quantity"`
	Price     float64   `json:"price"`
	Status    string    `json:"status" gorm:"index"` // "pending", "executed", "cancelled"
	ExecutedAt *time.Time `json:"executed_at"`
}