-- name: UpdateShowtimeSeatSeatByIdAndShowtimeId :exec
UPDATE showtime_seats
SET status = $2,
    booked_by = $3,
    booked_at = CASE WHEN $2 = 'booked'::seat_statuses THEN NOW() ELSE NULL::timestamp END
WHERE seat_id = $1 AND showtime_id = $4;

