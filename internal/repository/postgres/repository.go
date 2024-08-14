package postgres

import (
	"context"

	"conduit/internal/models"

	"gorm.io/gorm"
)

type Repository struct {
	db *gorm.DB
}

func (r *Repository) GetStrategy(ctx context.Context, id uint) (*models.Strategy, error) {
	var strategy models.Strategy
	err := r.db.WithContext(ctx).First(&strategy, id).Error
	return &strategy, err
}

// Implement other methods with context...
func (r *Repository) CreateStrategy(ctx context.Context, strategy *models.Strategy) error {
	return r.db.WithContext(ctx).Create(strategy).Error
}

func (r *Repository) UpdateStrategy(ctx context.Context, strategy *models.Strategy) error {
	return r.db.WithContext(ctx).Save(strategy).Error
}

func (r *Repository) DeleteStrategy(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&models.Strategy{}, id).Error
}

func (r *Repository) ListStrategies(ctx context.Context, userID uint) ([]models.Strategy, error) {
	var strategies []models.Strategy
	err := r.db.WithContext(ctx).Where("user_id = ?", userID).Find(&strategies).Error
	return strategies, err
}

func (r *Repository) CreateTrade(ctx context.Context, trade *models.Trade) error {
	return r.db.WithContext(ctx).Create(trade).Error
}

func (r *Repository) GetTrade(ctx context.Context, id uint) (*models.Trade, error) {
	var trade models.Trade
	err := r.db.WithContext(ctx).First(&trade, id).Error
	return &trade, err
}

func (r *Repository) ListTradesByUser(ctx context.Context, userID uint) ([]models.Trade, error) {
	var trades []models.Trade
	err := r.db.WithContext(ctx).Where("user_id = ?", userID).Find(&trades).Error
	return trades, err
}

func (r *Repository) GetPortfolio(ctx context.Context, userID uint) (*models.Portfolio, error) {
	var portfolio models.Portfolio
	err := r.db.WithContext(ctx).Where("user_id = ?", userID).First(&portfolio).Error
	return &portfolio, err
}
