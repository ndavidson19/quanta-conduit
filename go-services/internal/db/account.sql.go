// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.19.1
// source: account.sql

package db

import (
	"context"
	"database/sql"
)

const createAccount = `-- name: CreateAccount :one

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
RETURNING id, user_id, account_id, balance, portfolio_value, account_tier, account_type, is_system_account, created_at, updated_at
`

type CreateAccountParams struct {
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

// queries/account.sql
func (q *Queries) CreateAccount(ctx context.Context, arg CreateAccountParams) (Account, error) {
	row := q.db.QueryRowContext(ctx, createAccount,
		arg.UserID,
		arg.AccountID,
		arg.Balance,
		arg.PortfolioValue,
		arg.AccountTier,
		arg.AccountType,
		arg.IsSystemAccount,
		arg.CreatedAt,
		arg.UpdatedAt,
	)
	var i Account
	err := row.Scan(
		&i.ID,
		&i.UserID,
		&i.AccountID,
		&i.Balance,
		&i.PortfolioValue,
		&i.AccountTier,
		&i.AccountType,
		&i.IsSystemAccount,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const deleteAccount = `-- name: DeleteAccount :exec
DELETE FROM accounts
WHERE id = $1
`

func (q *Queries) DeleteAccount(ctx context.Context, id int32) error {
	_, err := q.db.ExecContext(ctx, deleteAccount, id)
	return err
}

const getAccount = `-- name: GetAccount :one
SELECT id, user_id, account_id, balance, portfolio_value, account_tier, account_type, is_system_account, created_at, updated_at FROM accounts
WHERE id = $1 LIMIT 1
`

func (q *Queries) GetAccount(ctx context.Context, id int32) (Account, error) {
	row := q.db.QueryRowContext(ctx, getAccount, id)
	var i Account
	err := row.Scan(
		&i.ID,
		&i.UserID,
		&i.AccountID,
		&i.Balance,
		&i.PortfolioValue,
		&i.AccountTier,
		&i.AccountType,
		&i.IsSystemAccount,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const getAccountByAccountID = `-- name: GetAccountByAccountID :one
SELECT id, user_id, account_id, balance, portfolio_value, account_tier, account_type, is_system_account, created_at, updated_at FROM accounts
WHERE account_id = $1 LIMIT 1
`

func (q *Queries) GetAccountByAccountID(ctx context.Context, accountID int64) (Account, error) {
	row := q.db.QueryRowContext(ctx, getAccountByAccountID, accountID)
	var i Account
	err := row.Scan(
		&i.ID,
		&i.UserID,
		&i.AccountID,
		&i.Balance,
		&i.PortfolioValue,
		&i.AccountTier,
		&i.AccountType,
		&i.IsSystemAccount,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const listAccountsByUserID = `-- name: ListAccountsByUserID :many
SELECT id, user_id, account_id, balance, portfolio_value, account_tier, account_type, is_system_account, created_at, updated_at FROM accounts
WHERE user_id = $1
ORDER BY id
`

func (q *Queries) ListAccountsByUserID(ctx context.Context, userID sql.NullInt32) ([]Account, error) {
	rows, err := q.db.QueryContext(ctx, listAccountsByUserID, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Account
	for rows.Next() {
		var i Account
		if err := rows.Scan(
			&i.ID,
			&i.UserID,
			&i.AccountID,
			&i.Balance,
			&i.PortfolioValue,
			&i.AccountTier,
			&i.AccountType,
			&i.IsSystemAccount,
			&i.CreatedAt,
			&i.UpdatedAt,
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

const updateAccount = `-- name: UpdateAccount :one
UPDATE accounts
SET 
    balance = COALESCE($2, balance),
    portfolio_value = COALESCE($3, portfolio_value),
    account_tier = COALESCE($4, account_tier),
    account_type = COALESCE($5, account_type),
    updated_at = CURRENT_TIMESTAMP
WHERE id = $1
RETURNING id, user_id, account_id, balance, portfolio_value, account_tier, account_type, is_system_account, created_at, updated_at
`

type UpdateAccountParams struct {
	ID             int32          `json:"id"`
	Balance        string         `json:"balance"`
	PortfolioValue string         `json:"portfolio_value"`
	AccountTier    sql.NullString `json:"account_tier"`
	AccountType    sql.NullString `json:"account_type"`
}

func (q *Queries) UpdateAccount(ctx context.Context, arg UpdateAccountParams) (Account, error) {
	row := q.db.QueryRowContext(ctx, updateAccount,
		arg.ID,
		arg.Balance,
		arg.PortfolioValue,
		arg.AccountTier,
		arg.AccountType,
	)
	var i Account
	err := row.Scan(
		&i.ID,
		&i.UserID,
		&i.AccountID,
		&i.Balance,
		&i.PortfolioValue,
		&i.AccountTier,
		&i.AccountType,
		&i.IsSystemAccount,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}