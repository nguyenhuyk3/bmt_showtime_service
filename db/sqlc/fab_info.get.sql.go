// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: fab_info.get.sql

package sqlc

import (
	"context"
)

const getPriceOfFAB = `-- name: GetPriceOfFAB :one
SELECT price
FROM fab_infos
WHERE fab_id = $1
`

func (q *Queries) GetPriceOfFAB(ctx context.Context, fabID int32) (int32, error) {
	row := q.db.QueryRow(ctx, getPriceOfFAB, fabID)
	var price int32
	err := row.Scan(&price)
	return price, err
}
