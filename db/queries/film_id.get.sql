-- name: IsFilmIdExist :one
SELECT EXISTS (
    SELECT 1 FROM "film_infos" WHERE id = $1
) AS exists;

-- name: GetDuration :one
SELECT duration
FROM film_infos
WHERE film_id = $1;

