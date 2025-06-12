-- name: GetShowtimeById :one
SELECT *
FROM showtimes
WHERE id = $1 AND is_released = true
LIMIT 1;

-- name: GetLatestShowtimeByAuditoriumId :one
SELECT end_time
FROM showtimes
WHERE auditorium_id = $1 AND show_date = $2
ORDER BY created_at DESC
LIMIT 1;

-- name: IsShowtimeExist :one
SELECT EXISTS (
    SELECT 1 FROM showtimes 
    WHERE id = $1
) AS EXISTS;

-- name: GetAllShowTimesByFilmIdInOneDate :many
SELECT * 
FROM showtimes 
WHERE film_id = $1 AND show_date = $2 AND is_released = true;

-- name: isShowtimeRealeased :one
SELECT is_released
FROM showtimes
WHERE id = $1;

-- name: GetShowdateByShowtimeId :one
SELECT show_date
FROM showtimes
WHERE id = $1;

-- name: GetFilmIdsInToday :many
SELECT DISTINCT film_id
FROM showtimes 
WHERE show_date = $1;

-- name: GetAllShowTimesByFilmIdAndByCinemaIdAndInDayRange :many
SELECT sh.*
FROM showtimes sh
JOIN auditoriums a ON sh.auditorium_id = a.id
JOIN cinemas c ON a.cinema_id = c.id
WHERE sh.is_released = TRUE
    AND sh.film_id = $1
    AND c.id = $2
    AND sh.show_date BETWEEN CURRENT_DATE AND CURRENT_DATE + INTERVAL '14 days'
ORDER BY sh.show_date, sh.start_time;



