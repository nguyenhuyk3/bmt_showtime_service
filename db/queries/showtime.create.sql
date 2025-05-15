-- name: CreateShowTime :exec
INSERT INTO showtimes (
    film_id,
    auditorium_id,
    changed_by,
    show_date,
    start_time,
    end_time
)
VALUES (
    $1, $2, $3, $4, $5, $6
);