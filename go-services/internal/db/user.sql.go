// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.19.1
// source: user.sql

package db

import (
	"context"
	"database/sql"
)

const createUser = `-- name: CreateUser :one

INSERT INTO users (
    email,
    hashed_password,
    provider,
    email_verified,
    two_factor_secret
) VALUES (
    $1, $2, $3, $4, $5
)
RETURNING id, email, hashed_password, provider, email_verified, two_factor_secret, created_at, updated_at
`

type CreateUserParams struct {
	Email           string         `json:"email"`
	HashedPassword  sql.NullString `json:"hashed_password"`
	Provider        sql.NullString `json:"provider"`
	EmailVerified   sql.NullBool   `json:"email_verified"`
	TwoFactorSecret sql.NullString `json:"two_factor_secret"`
}

// queries/user.sql
func (q *Queries) CreateUser(ctx context.Context, arg CreateUserParams) (User, error) {
	row := q.db.QueryRowContext(ctx, createUser,
		arg.Email,
		arg.HashedPassword,
		arg.Provider,
		arg.EmailVerified,
		arg.TwoFactorSecret,
	)
	var i User
	err := row.Scan(
		&i.ID,
		&i.Email,
		&i.HashedPassword,
		&i.Provider,
		&i.EmailVerified,
		&i.TwoFactorSecret,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const deleteUser = `-- name: DeleteUser :exec
DELETE FROM users
WHERE id = $1
`

func (q *Queries) DeleteUser(ctx context.Context, id int32) error {
	_, err := q.db.ExecContext(ctx, deleteUser, id)
	return err
}

const getUser = `-- name: GetUser :one
SELECT id, email, hashed_password, provider, email_verified, two_factor_secret, created_at, updated_at FROM users
WHERE id = $1 LIMIT 1
`

func (q *Queries) GetUser(ctx context.Context, id int32) (User, error) {
	row := q.db.QueryRowContext(ctx, getUser, id)
	var i User
	err := row.Scan(
		&i.ID,
		&i.Email,
		&i.HashedPassword,
		&i.Provider,
		&i.EmailVerified,
		&i.TwoFactorSecret,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const getUserByEmail = `-- name: GetUserByEmail :one
SELECT id, email, hashed_password, provider, email_verified, two_factor_secret, created_at, updated_at FROM users
WHERE email = $1 LIMIT 1
`

func (q *Queries) GetUserByEmail(ctx context.Context, email string) (User, error) {
	row := q.db.QueryRowContext(ctx, getUserByEmail, email)
	var i User
	err := row.Scan(
		&i.ID,
		&i.Email,
		&i.HashedPassword,
		&i.Provider,
		&i.EmailVerified,
		&i.TwoFactorSecret,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const listUsers = `-- name: ListUsers :many
SELECT id, email, hashed_password, provider, email_verified, two_factor_secret, created_at, updated_at FROM users
ORDER BY id
`

func (q *Queries) ListUsers(ctx context.Context) ([]User, error) {
	rows, err := q.db.QueryContext(ctx, listUsers)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []User
	for rows.Next() {
		var i User
		if err := rows.Scan(
			&i.ID,
			&i.Email,
			&i.HashedPassword,
			&i.Provider,
			&i.EmailVerified,
			&i.TwoFactorSecret,
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

const updateUser = `-- name: UpdateUser :one
UPDATE users
SET 
    email = COALESCE($2, email),
    hashed_password = COALESCE($3, hashed_password),
    provider = COALESCE($4, provider),
    email_verified = COALESCE($5, email_verified),
    two_factor_secret = COALESCE($6, two_factor_secret),
    updated_at = CURRENT_TIMESTAMP
WHERE id = $1
RETURNING id, email, hashed_password, provider, email_verified, two_factor_secret, created_at, updated_at
`

type UpdateUserParams struct {
	ID              int32          `json:"id"`
	Email           string         `json:"email"`
	HashedPassword  sql.NullString `json:"hashed_password"`
	Provider        sql.NullString `json:"provider"`
	EmailVerified   sql.NullBool   `json:"email_verified"`
	TwoFactorSecret sql.NullString `json:"two_factor_secret"`
}

func (q *Queries) UpdateUser(ctx context.Context, arg UpdateUserParams) (User, error) {
	row := q.db.QueryRowContext(ctx, updateUser,
		arg.ID,
		arg.Email,
		arg.HashedPassword,
		arg.Provider,
		arg.EmailVerified,
		arg.TwoFactorSecret,
	)
	var i User
	err := row.Scan(
		&i.ID,
		&i.Email,
		&i.HashedPassword,
		&i.Provider,
		&i.EmailVerified,
		&i.TwoFactorSecret,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}
