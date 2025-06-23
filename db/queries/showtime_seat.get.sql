-- name: GetAllShowtimeSeatsByShowtimeId :many
SELECT  
    ss.id,
    ss.showtime_id,
    ss.seat_id,
    s.seat_type,
    s.seat_number,
    s.price,
    ss.status,
    ss.booked_by,
    ss.created_at,
    ss.booked_at
FROM showtime_seats ss
JOIN seats s ON s.id = ss.seat_id
WHERE ss.showtime_id = $1;

-- name: GetAllShowtimeSeatsFromEarliestTomorrow :many
WITH next_showtime AS (
    SELECT st.id AS showtime_id
    FROM showtimes st
    WHERE st.film_id = $1
        AND st.show_date >= CURRENT_DATE + INTERVAL '1 day'
    ORDER BY st.show_date, st.start_time
    LIMIT 1
)
SELECT 
    ss.id,
    ss.showtime_id,
    ss.seat_id,
    s.seat_type,
    s.seat_number,
    s.price,
    ss.status,
    ss.booked_by,
    ss.created_at,
    ss.booked_at
FROM showtime_seats ss
JOIN next_showtime ns ON ss.showtime_id = ns.showtime_id
JOIN seats s ON ss.seat_id = s.id
ORDER BY ss.seat_id;



