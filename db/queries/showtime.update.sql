-- name: ReleaseShowtime :exec
UPDATE showtimes
SET is_deleted = !is_deleted,
    updated_at = NOW()
WHERE id = $1;

-- name: UpdateShowtime :exec
UPDATE showtimes
SET changed_by = $1,
    updated_at = NOW()
WHERE id = $1;