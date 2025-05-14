-- name: CreateNewFilmId :exec
INSERT INTO film_ids (
    film_id
)
VALUES (
    $1
);