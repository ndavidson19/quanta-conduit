-- queries/asset.sql

-- name: CreateAsset :one
INSERT INTO assets (
    symbol,
    name,
    asset_type
) VALUES (
    $1, $2, $3
)
RETURNING *;

-- name: GetAsset :one
SELECT * FROM assets
WHERE id = $1 LIMIT 1;

-- name: GetAssetBySymbol :one
SELECT * FROM assets
WHERE symbol = $1 LIMIT 1;

-- name: ListAssets :many
SELECT * FROM assets
ORDER BY symbol
LIMIT $1 OFFSET $2;

-- name: UpdateAsset :one
UPDATE assets
SET 
    name = COALESCE($2, name),
    asset_type = COALESCE($3, asset_type),
    updated_at = CURRENT_TIMESTAMP
WHERE id = $1
RETURNING *;

-- name: DeleteAsset :exec
DELETE FROM assets
WHERE id = $1;