-- name: IsAuditoriumExist :one
SELECT EXISTS (
    SELECT 1 FROM "auditoriums" WHERE id = $1 AND is_released = true
) AS exists;

-- name: GetAuditoriumByShowtimeId :one
SELECT a.*
FROM showtimes sh
JOIN auditoriums a ON sh.auditorium_id = a.id
WHERE sh.id = $1;
