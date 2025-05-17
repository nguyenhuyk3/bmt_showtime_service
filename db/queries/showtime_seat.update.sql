-- name: UpdateShowtimeSeatById :exec
UPDATE showtime_seats
SET status = $2,
    booked_by = $3,
    booked_at = CASE WHEN $2 = 'booked' THEN NOW() ELSE NULL END
WHERE id = $1;