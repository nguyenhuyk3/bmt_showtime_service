-- name: CreateNewFABInfo :exec
INSERT INTO fab_infos (
    fab_id,
    price
)
VALUES (
    $1, $2
);