-- queries/order.sql

-- name: CreateOrder :one
INSERT INTO orders (
    account_id,
    asset_id,
    order_type,
    side,
    quantity,
    price,
    status
) VALUES (
    $1, $2, $3, $4, $5, $6, $7
)
RETURNING *;

-- name: GetOrder :one
SELECT * FROM orders
WHERE id = $1 LIMIT 1;

-- name: ListOrdersByAccount :many
SELECT * FROM orders
WHERE account_id = $1
ORDER BY created_at DESC
LIMIT $2 OFFSET $3;

-- name: ListPendingOrders :many
SELECT * FROM orders
WHERE status = 'pending'
ORDER BY created_at
LIMIT $1 OFFSET $2;

-- name: UpdateOrderStatus :one
UPDATE orders
SET 
    status = $2,
    updated_at = CURRENT_TIMESTAMP
WHERE id = $1
RETURNING *;

-- name: CancelOrder :one
UPDATE orders
SET 
    status = 'cancelled',
    updated_at = CURRENT_TIMESTAMP
WHERE id = $1 AND status = 'pending'
RETURNING *;

-- name: DeleteOrder :exec
DELETE FROM orders
WHERE id = $1;