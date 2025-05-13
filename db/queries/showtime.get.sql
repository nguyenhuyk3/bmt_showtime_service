-- name: GetShowtimeById :one
SELECT *
FROM showtimes
WHERE id = $1
LIMIT 1;




