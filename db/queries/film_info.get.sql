-- name: IsFilmIdExist :one
SELECT EXISTS (
    SELECT 1 FROM "film_infos" WHERE film_id = $1
) AS EXISTS;

-- name: GetDuration :one
SELECT duration
FROM film_infos
WHERE film_id = $1;

