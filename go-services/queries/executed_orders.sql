-- queries/executed_order.sql

-- name: CreateExecutedOrder :one
INSERT INTO executed_orders (
    order_id,
    executed_quantity,
    executed_price
) VALUES (
    $1, $2, $3
)
RETURNING *;

-- name: GetExecutedOrder :one
SELECT * FROM executed_orders
WHERE id = $1 LIMIT 1;

-- name: ListExecutedOrdersByOrder :many
SELECT * FROM executed_orders
WHERE order_id = $1
ORDER BY execution_timestamp
LIMIT $2 OFFSET $3;

-- name: ListExecutedOrdersByAccount :many
SELECT eo.*
FROM executed_orders eo
JOIN orders o ON eo.order_id = o.id
WHERE o.account_id = $1
ORDER BY eo.execution_timestamp DESC
LIMIT $2 OFFSET $3;

-- name: GetTotalExecutedQuantity :one
SELECT COALESCE(SUM(executed_quantity), 0) as total_executed_quantity
FROM executed_orders
WHERE order_id = $1;

-- name: DeleteExecutedOrder :exec
DELETE FROM executed_orders
WHERE id = $1;