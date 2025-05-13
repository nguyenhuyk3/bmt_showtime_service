-- name: IsAuditoriumExist :one
SELECT EXISTS (
    SELECT 1 FROM "auditoriums" WHERE id = $1
) AS exists;
