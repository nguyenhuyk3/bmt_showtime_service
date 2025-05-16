-- name: TurnOnShowtime :exec
UPDATE showtimes
SET is_deleted = true,
    updated_at = NOW()
WHERE id = $1;