-- name: GetPriceOfSeatBySeatId :one
SELECT price
FROM seats
WHERE id = $1;

-- name: GetSeatById :one
SELECT * 
FROM seats
WHERE id = $1;