-- name: GetAllShowtimeSeatsByShowtimeId :many
SELECT *
FROM showtime_seats
WHERE showtime_id = $1;