-- name: GetPriceOfFAB :one
SELECT price
FROM fab_infos
WHERE fab_id = $1;