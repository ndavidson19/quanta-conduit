-- queries/user.sql

-- name: CreateUser :one
INSERT INTO users (
    email,
    hashed_password,
    provider,
    email_verified,
    two_factor_secret
) VALUES (
    $1, $2, $3, $4, $5
)
RETURNING *;

-- name: GetUser :one
SELECT * FROM users
WHERE id = $1 LIMIT 1;

-- name: GetUserByEmail :one
SELECT * FROM users
WHERE email = $1 LIMIT 1;

-- name: ListUsers :many
SELECT * FROM users
ORDER BY id;

-- name: UpdateUser :one
UPDATE users
SET 
    email = COALESCE($2, email),
    hashed_password = COALESCE($3, hashed_password),
    provider = COALESCE($4, provider),
    email_verified = COALESCE($5, email_verified),
    two_factor_secret = COALESCE($6, two_factor_secret),
    updated_at = CURRENT_TIMESTAMP
WHERE id = $1
RETURNING *;

-- name: DeleteUser :exec
DELETE FROM users
WHERE id = $1;