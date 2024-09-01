-- queries/portfolio.sql

-- name: CreatePortfolio :one
INSERT INTO portfolios (
    account_id,
    total_value
) VALUES (
    $1, $2
)
RETURNING *;

-- name: GetPortfolio :one
SELECT * FROM portfolios
WHERE id = $1 LIMIT 1;

-- name: GetPortfolioByAccountID :one
SELECT * FROM portfolios
WHERE account_id = $1 LIMIT 1;

-- name: UpdatePortfolioValue :one
UPDATE portfolios
SET 
    total_value = $2,
    updated_at = CURRENT_TIMESTAMP
WHERE id = $1
RETURNING *;

-- name: ListPortfolioHoldings :many
SELECT ph.*, a.symbol, a.name, a.asset_type
FROM portfolio_holdings ph
JOIN assets a ON ph.asset_id = a.id
WHERE ph.portfolio_id = $1
ORDER BY a.symbol;

-- name: UpdatePortfolioHolding :one
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
RETURNING *;