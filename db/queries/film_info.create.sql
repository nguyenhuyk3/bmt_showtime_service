-- name: CreateNewFilmId :exec
INSERT INTO film_infos (
    film_id,
    duration
)
VALUES (
    $1, $2
);