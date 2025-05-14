-- name: IsFilmIdExist :one
SELECT EXISTS (
    SELECT 1 FROM "film_ids" WHERE id = $1
) AS exists;
