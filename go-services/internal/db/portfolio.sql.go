// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.19.1
// source: portfolio.sql

package db

import (
	"context"
	"database/sql"
)

const createPortfolio = `-- name: CreatePortfolio :one

INSERT INTO portfolios (
    account_id,
    total_value
) VALUES (
    $1, $2
)
RETURNING id, account_id, total_value, created_at, updated_at
`

type CreatePortfolioParams struct {
	AccountID  int64  `json:"account_id"`
	TotalValue string `json:"total_value"`
}

// queries/portfolio.sql
func (q *Queries) CreatePortfolio(ctx context.Context, arg CreatePortfolioParams) (Portfolio, error) {
	row := q.db.QueryRowContext(ctx, createPortfolio, arg.AccountID, arg.TotalValue)
	var i Portfolio
	err := row.Scan(
		&i.ID,
		&i.AccountID,
		&i.TotalValue,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const getPortfolio = `-- name: GetPortfolio :one
SELECT id, account_id, total_value, created_at, updated_at FROM portfolios
WHERE id = $1 LIMIT 1
`

func (q *Queries) GetPortfolio(ctx context.Context, id int32) (Portfolio, error) {
	row := q.db.QueryRowContext(ctx, getPortfolio, id)
	var i Portfolio
	err := row.Scan(
		&i.ID,
		&i.AccountID,
		&i.TotalValue,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const getPortfolioByAccountID = `-- name: GetPortfolioByAccountID :one
SELECT id, account_id, total_value, created_at, updated_at FROM portfolios
WHERE account_id = $1 LIMIT 1
`

func (q *Queries) GetPortfolioByAccountID(ctx context.Context, accountID int64) (Portfolio, error) {
	row := q.db.QueryRowContext(ctx, getPortfolioByAccountID, accountID)
	var i Portfolio
	err := row.Scan(
		&i.ID,
		&i.AccountID,
		&i.TotalValue,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const listPortfolioHoldings = `-- name: ListPortfolioHoldings :many
SELECT ph.id, ph.portfolio_id, ph.asset_id, ph.quantity, ph.average_buy_price, ph.created_at, ph.updated_at, a.symbol, a.name, a.asset_type
FROM portfolio_holdings ph
JOIN assets a ON ph.asset_id = a.id
WHERE ph.portfolio_id = $1
ORDER BY a.symbol
`

type ListPortfolioHoldingsRow struct {
	ID              int32        `json:"id"`
	PortfolioID     int32        `json:"portfolio_id"`
	AssetID         int32        `json:"asset_id"`
	Quantity        string       `json:"quantity"`
	AverageBuyPrice string       `json:"average_buy_price"`
	CreatedAt       sql.NullTime `json:"created_at"`
	UpdatedAt       sql.NullTime `json:"updated_at"`
	Symbol          string       `json:"symbol"`
	Name            string       `json:"name"`
	AssetType       string       `json:"asset_type"`
}

func (q *Queries) ListPortfolioHoldings(ctx context.Context, portfolioID int32) ([]ListPortfolioHoldingsRow, error) {
	rows, err := q.db.QueryContext(ctx, listPortfolioHoldings, portfolioID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []ListPortfolioHoldingsRow
	for rows.Next() {
		var i ListPortfolioHoldingsRow
		if err := rows.Scan(
			&i.ID,
			&i.PortfolioID,
			&i.AssetID,
			&i.Quantity,
			&i.AverageBuyPrice,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.Symbol,
			&i.Name,
			&i.AssetType,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const updatePortfolioHolding = `-- name: UpdatePortfolioHolding :one
INSERT INTO portfolio_holdings (
    portfolio_id, asset_id, quantity, average_buy_price
) VALUES (
    $1, $2, $3, $4
)
ON CONFLICT (portfolio_id, asset_id) 
DO UPDATE SET
    quantity = portfolio_holdings.quantity + EXCLUDED.quantity,
    average_buy_price = (portfolio_holdings.average_buy_price * portfolio_holdings.quantity + EXCLUDED.average_buy_price * EXCLUDED.quantity) / (portfolio_holdings.quantity + EXCLUDED.quantity),
    updated_at = CURRENT_TIMESTAMP
RETURNING id, portfolio_id, asset_id, quantity, average_buy_price, created_at, updated_at
`

type UpdatePortfolioHoldingParams struct {
	PortfolioID     int32  `json:"portfolio_id"`
	AssetID         int32  `json:"asset_id"`
	Quantity        string `json:"quantity"`
	AverageBuyPrice string `json:"average_buy_price"`
}

func (q *Queries) UpdatePortfolioHolding(ctx context.Context, arg UpdatePortfolioHoldingParams) (PortfolioHolding, error) {
	row := q.db.QueryRowContext(ctx, updatePortfolioHolding,
		arg.PortfolioID,
		arg.AssetID,
		arg.Quantity,
		arg.AverageBuyPrice,
	)
	var i PortfolioHolding
	err := row.Scan(
		&i.ID,
		&i.PortfolioID,
		&i.AssetID,
		&i.Quantity,
		&i.AverageBuyPrice,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const updatePortfolioValue = `-- name: UpdatePortfolioValue :one
UPDATE portfolios
SET 
    total_value = $2,
    updated_at = CURRENT_TIMESTAMP
WHERE id = $1
RETURNING id, account_id, total_value, created_at, updated_at
`

type UpdatePortfolioValueParams struct {
	ID         int32  `json:"id"`
	TotalValue string `json:"total_value"`
}

func (q *Queries) UpdatePortfolioValue(ctx context.Context, arg UpdatePortfolioValueParams) (Portfolio, error) {
	row := q.db.QueryRowContext(ctx, updatePortfolioValue, arg.ID, arg.TotalValue)
	var i Portfolio
	err := row.Scan(
		&i.ID,
		&i.AccountID,
		&i.TotalValue,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}
