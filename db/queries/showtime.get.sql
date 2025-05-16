-- name: GetShowtimeById :one
SELECT *
FROM showtimes
WHERE id = $1
LIMIT 1;

-- name: GetLatestShowtimeByAuditoriumId :one
SELECT end_time
FROM showtimes
WHERE auditorium_id = $1
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
WHERE film_id = $1 AND show_date = $2;
