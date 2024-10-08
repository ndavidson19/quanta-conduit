// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.19.1
// source: executed_orders.sql

package db

import (
	"context"
)

const createExecutedOrder = `-- name: CreateExecutedOrder :one

INSERT INTO executed_orders (
    order_id,
    executed_quantity,
    executed_price
) VALUES (
    $1, $2, $3
)
RETURNING id, order_id, executed_quantity, executed_price, execution_timestamp
`

type CreateExecutedOrderParams struct {
	OrderID          int32  `json:"order_id"`
	ExecutedQuantity string `json:"executed_quantity"`
	ExecutedPrice    string `json:"executed_price"`
}

// queries/executed_order.sql
func (q *Queries) CreateExecutedOrder(ctx context.Context, arg CreateExecutedOrderParams) (ExecutedOrder, error) {
	row := q.db.QueryRowContext(ctx, createExecutedOrder, arg.OrderID, arg.ExecutedQuantity, arg.ExecutedPrice)
	var i ExecutedOrder
	err := row.Scan(
		&i.ID,
		&i.OrderID,
		&i.ExecutedQuantity,
		&i.ExecutedPrice,
		&i.ExecutionTimestamp,
	)
	return i, err
}

const deleteExecutedOrder = `-- name: DeleteExecutedOrder :exec
DELETE FROM executed_orders
WHERE id = $1
`

func (q *Queries) DeleteExecutedOrder(ctx context.Context, id int32) error {
	_, err := q.db.ExecContext(ctx, deleteExecutedOrder, id)
	return err
}

const getExecutedOrder = `-- name: GetExecutedOrder :one
SELECT id, order_id, executed_quantity, executed_price, execution_timestamp FROM executed_orders
WHERE id = $1 LIMIT 1
`

func (q *Queries) GetExecutedOrder(ctx context.Context, id int32) (ExecutedOrder, error) {
	row := q.db.QueryRowContext(ctx, getExecutedOrder, id)
	var i ExecutedOrder
	err := row.Scan(
		&i.ID,
		&i.OrderID,
		&i.ExecutedQuantity,
		&i.ExecutedPrice,
		&i.ExecutionTimestamp,
	)
	return i, err
}

const getTotalExecutedQuantity = `-- name: GetTotalExecutedQuantity :one
SELECT COALESCE(SUM(executed_quantity), 0) as total_executed_quantity
FROM executed_orders
WHERE order_id = $1
`

func (q *Queries) GetTotalExecutedQuantity(ctx context.Context, orderID int32) (interface{}, error) {
	row := q.db.QueryRowContext(ctx, getTotalExecutedQuantity, orderID)
	var total_executed_quantity interface{}
	err := row.Scan(&total_executed_quantity)
	return total_executed_quantity, err
}

const listExecutedOrdersByAccount = `-- name: ListExecutedOrdersByAccount :many
SELECT eo.id, eo.order_id, eo.executed_quantity, eo.executed_price, eo.execution_timestamp
FROM executed_orders eo
JOIN orders o ON eo.order_id = o.id
WHERE o.account_id = $1
ORDER BY eo.execution_timestamp DESC
LIMIT $2 OFFSET $3
`

type ListExecutedOrdersByAccountParams struct {
	AccountID int64 `json:"account_id"`
	Limit     int32 `json:"limit"`
	Offset    int32 `json:"offset"`
}

func (q *Queries) ListExecutedOrdersByAccount(ctx context.Context, arg ListExecutedOrdersByAccountParams) ([]ExecutedOrder, error) {
	rows, err := q.db.QueryContext(ctx, listExecutedOrdersByAccount, arg.AccountID, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []ExecutedOrder
	for rows.Next() {
		var i ExecutedOrder
		if err := rows.Scan(
			&i.ID,
			&i.OrderID,
			&i.ExecutedQuantity,
			&i.ExecutedPrice,
			&i.ExecutionTimestamp,
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

const listExecutedOrdersByOrder = `-- name: ListExecutedOrdersByOrder :many
SELECT id, order_id, executed_quantity, executed_price, execution_timestamp FROM executed_orders
WHERE order_id = $1
ORDER BY execution_timestamp
LIMIT $2 OFFSET $3
`

type ListExecutedOrdersByOrderParams struct {
	OrderID int32 `json:"order_id"`
	Limit   int32 `json:"limit"`
	Offset  int32 `json:"offset"`
}

func (q *Queries) ListExecutedOrdersByOrder(ctx context.Context, arg ListExecutedOrdersByOrderParams) ([]ExecutedOrder, error) {
	rows, err := q.db.QueryContext(ctx, listExecutedOrdersByOrder, arg.OrderID, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []ExecutedOrder
	for rows.Next() {
		var i ExecutedOrder
		if err := rows.Scan(
			&i.ID,
			&i.OrderID,
			&i.ExecutedQuantity,
			&i.ExecutedPrice,
			&i.ExecutionTimestamp,
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
