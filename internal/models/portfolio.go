package models

import (
	"gorm.io/gorm"
)

type Portfolio struct {
	gorm.Model
	UserID uint    `json:"user_id" gorm:"index"`
	Value  float64 `json:"value"`
	Assets []Asset `json:"assets" gorm:"foreignKey:PortfolioID"`
}

type Asset struct {
	gorm.Model
	PortfolioID  uint    `json:"portfolio_id" gorm:"index"`
	Symbol       string  `json:"symbol" gorm:"not null"`
	Quantity     float64 `json:"quantity"`
	CurrentPrice float64 `json:"current_price"`
}