-- name: createShowtimeSeats :exec
INSERT INTO showtime_seats (showtime_id, seat_id, status, created_at)
SELECT
    $1 AS showtime_id,
    s.id AS seat_id,
    'available'::seat_statuses AS status,
    now()
FROM
    showtimes sh
JOIN
    seats s ON s.auditorium_id = sh.auditorium_id
WHERE
    sh.id = $1;
