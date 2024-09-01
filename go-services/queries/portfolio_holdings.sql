-- queries/portfolio_holding.sql

-- name: CreatePortfolioHolding :one
INSERT INTO portfolio_holdings (
    portfolio_id,
    asset_id,
    quantity,
    average_buy_price
) VALUES (
    $1, $2, $3, $4
)
RETURNING *;

-- name: GetPortfolioHolding :one
SELECT * FROM portfolio_holdings
WHERE id = $1 LIMIT 1;

-- name: GetPortfolioHoldingByAsset :one
SELECT * FROM portfolio_holdings
WHERE portfolio_id = $1 AND asset_id = $2 LIMIT 1;

-- name: ListPortfolioHoldings :many
SELECT ph.*, a.symbol, a.name, a.asset_type
FROM portfolio_holdings ph
JOIN assets a ON ph.asset_id = a.id
WHERE ph.portfolio_id = $1
ORDER BY a.symbol
LIMIT $2 OFFSET $3;

-- name: UpdatePortfolioHolding :one
UPDATE portfolio_holdings
SET 
    quantity = $3,
    average_buy_price = $4,
    updated_at = CURRENT_TIMESTAMP
WHERE portfolio_id = $1 AND asset_id = $2
RETURNING *;

-- name: UpsertPortfolioHolding :one
INSERT INTO portfolio_holdings (
    portfolio_id,
    asset_id,
    quantity,
    average_buy_price
) VALUES (
    $1, $2, $3, $4
)
ON CONFLICT (portfolio_id, asset_id) 
DO UPDATE SET
    quantity = portfolio_holdings.quantity + EXCLUDED.quantity,
    average_buy_price = (portfolio_holdings.average_buy_price * portfolio_holdings.quantity + EXCLUDED.average_buy_price * EXCLUDED.quantity) / (portfolio_holdings.quantity + EXCLUDED.quantity),
    updated_at = CURRENT_TIMESTAMP
RETURNING *;

-- name: DeletePortfolioHolding :exec
DELETE FROM portfolio_holdings
WHERE portfolio_id = $1 AND asset_id = $2;

-- name: GetTotalPortfolioValue :one
SELECT COALESCE(SUM(ph.quantity * a.current_price), 0) as total_value
FROM portfolio_holdings ph
JOIN assets a ON ph.asset_id = a.id
WHERE ph.portfolio_id = $1;

-- name: GetPortfolioHoldingsWithCurrentValues :many
SELECT 
    ph.*,
    a.symbol,
    a.name,
    a.asset_type,
    a.current_price,
    (ph.quantity * a.current_price) as current_value,
    ((ph.quantity * a.current_price) - (ph.quantity * ph.average_buy_price)) as unrealized_pnl
FROM portfolio_holdings ph
JOIN assets a ON ph.asset_id = a.id
WHERE ph.portfolio_id = $1
ORDER BY current_value DESC
LIMIT $2 OFFSET $3;