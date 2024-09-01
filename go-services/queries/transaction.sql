-- queries/transaction.sql

-- name: CreateTransaction :one
INSERT INTO transactions (
    id,
    debit_account_id,
    credit_account_id,
    amount,
    description
) VALUES (
    $1, $2, $3, $4, $5
)
RETURNING *;

-- name: GetTransaction :one
SELECT * FROM transactions
WHERE id = $1 LIMIT 1;

-- name: ListTransactionsByAccount :many
SELECT * FROM transactions
WHERE debit_account_id = $1 OR credit_account_id = $1
ORDER BY timestamp DESC
LIMIT $2 OFFSET $3;

-- name: ListTransactionsByTimeRange :many
SELECT * FROM transactions
WHERE timestamp BETWEEN $1 AND $2
ORDER BY timestamp DESC
LIMIT $3 OFFSET $4;