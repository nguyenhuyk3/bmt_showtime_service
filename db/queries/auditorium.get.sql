-- name: IsAuditoriumExist :one
SELECT EXISTS (
    SELECT 1 FROM "auditoriums" WHERE id = $1 AND is_released = true
) AS exists;
