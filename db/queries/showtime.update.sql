-- name: releaseShowtime :exec
UPDATE showtimes
SET is_released = NOT is_released,
    updated_at = NOW()
WHERE id = $1;


-- name: updateShowtime :exec
UPDATE showtimes
SET changed_by = $2,
    updated_at = NOW()
WHERE id = $1;