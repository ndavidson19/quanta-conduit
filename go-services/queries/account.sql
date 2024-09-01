-- queries/account.sql

-- name: CreateAccount :one
INSERT INTO accounts (
    user_id,
    account_id,
    balance,
    portfolio_value,
    account_tier,
    account_type,
    is_system_account,
    created_at,
    updated_at
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9
)
RETURNING *;

-- name: GetAccount :one
SELECT * FROM accounts
WHERE id = $1 LIMIT 1;

-- name: GetAccountByAccountID :one
SELECT * FROM accounts
WHERE account_id = $1 LIMIT 1;

-- name: ListAccountsByUserID :many
SELECT * FROM accounts
WHERE user_id = $1
ORDER BY id;

-- name: UpdateAccount :one
UPDATE accounts
SET 
    balance = COALESCE($2, balance),
    portfolio_value = COALESCE($3, portfolio_value),
    account_tier = COALESCE($4, account_tier),
    account_type = COALESCE($5, account_type),
    updated_at = CURRENT_TIMESTAMP
WHERE id = $1
RETURNING *;

-- name: DeleteAccount :exec
DELETE FROM accounts
WHERE id = $1;