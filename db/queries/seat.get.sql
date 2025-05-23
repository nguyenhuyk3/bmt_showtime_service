-- name: GetPriceOfSeatBySeatIdAndShowtimeId :one
SELECT price
FROM seats
WHERE id = $1;