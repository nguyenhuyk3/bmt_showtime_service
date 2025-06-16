-- name: UpdateShowtimeSeatByIdAndShowtimeId :exec
UPDATE showtime_seats
SET status = $2,
    booked_by = $3,
    booked_at = CASE WHEN $2 = 'booked'::seat_statuses THEN NOW() ELSE NULL::timestamp END
WHERE seat_id = $1 AND showtime_id = $4;

-- name: UpdateShowtimeSeatByIdAndShowtimeIdSuccess :exec
UPDATE showtime_seats
SET status = $2,
    booked_at = CASE WHEN $2 = 'booked'::seat_statuses THEN NOW() ELSE NULL::timestamp END,
    updated_at = NOW()
WHERE seat_id = $1 AND showtime_id = $3;

-- name: UpdateShowtimeSeatByIdAndShowtimeIdFailed :exec
UPDATE showtime_seats
SET status = $2,
    booked_by = NULL,
    updated_at = NOW()
WHERE seat_id = $1 AND showtime_id = $3;


