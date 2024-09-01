// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.19.1

package db

import (
	"database/sql"
)

type Account struct {
	ID              int32          `json:"id"`
	UserID          sql.NullInt32  `json:"user_id"`
	AccountID       int64          `json:"account_id"`
	Balance         string         `json:"balance"`
	PortfolioValue  string         `json:"portfolio_value"`
	AccountTier     sql.NullString `json:"account_tier"`
	AccountType     sql.NullString `json:"account_type"`
	IsSystemAccount sql.NullBool   `json:"is_system_account"`
	CreatedAt       sql.NullTime   `json:"created_at"`
	UpdatedAt       sql.NullTime   `json:"updated_at"`
}

type Asset struct {
	ID        int32        `json:"id"`
	Symbol    string       `json:"symbol"`
	Name      string       `json:"name"`
	AssetType string       `json:"asset_type"`
	CreatedAt sql.NullTime `json:"created_at"`
	UpdatedAt sql.NullTime `json:"updated_at"`
}

type ExecutedOrder struct {
	ID                 int32        `json:"id"`
	OrderID            int32        `json:"order_id"`
	ExecutedQuantity   string       `json:"executed_quantity"`
	ExecutedPrice      string       `json:"executed_price"`
	ExecutionTimestamp sql.NullTime `json:"execution_timestamp"`
}

type Order struct {
	ID        int32          `json:"id"`
	AccountID int64          `json:"account_id"`
	AssetID   int32          `json:"asset_id"`
	OrderType string         `json:"order_type"`
	Side      string         `json:"side"`
	Quantity  string         `json:"quantity"`
	Price     sql.NullString `json:"price"`
	Status    string         `json:"status"`
	CreatedAt sql.NullTime   `json:"created_at"`
	UpdatedAt sql.NullTime   `json:"updated_at"`
}

type Portfolio struct {
	ID         int32        `json:"id"`
	AccountID  int64        `json:"account_id"`
	TotalValue string       `json:"total_value"`
	CreatedAt  sql.NullTime `json:"created_at"`
	UpdatedAt  sql.NullTime `json:"updated_at"`
}

type PortfolioHolding struct {
	ID              int32        `json:"id"`
	PortfolioID     int32        `json:"portfolio_id"`
	AssetID         int32        `json:"asset_id"`
	Quantity        string       `json:"quantity"`
	AverageBuyPrice string       `json:"average_buy_price"`
	CreatedAt       sql.NullTime `json:"created_at"`
	UpdatedAt       sql.NullTime `json:"updated_at"`
}

type Transaction struct {
	ID              []byte       `json:"id"`
	DebitAccountID  int64        `json:"debit_account_id"`
	CreditAccountID int64        `json:"credit_account_id"`
	Amount          int64        `json:"amount"`
	Description     string       `json:"description"`
	Timestamp       sql.NullTime `json:"timestamp"`
}

type User struct {
	ID              int32          `json:"id"`
	Email           string         `json:"email"`
	HashedPassword  sql.NullString `json:"hashed_password"`
	Provider        sql.NullString `json:"provider"`
	EmailVerified   sql.NullBool   `json:"email_verified"`
	TwoFactorSecret sql.NullString `json:"two_factor_secret"`
	CreatedAt       sql.NullTime   `json:"created_at"`
	UpdatedAt       sql.NullTime   `json:"updated_at"`
}