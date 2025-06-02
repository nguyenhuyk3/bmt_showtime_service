-- name: GetCinemasForShowingFilm :many
SELECT DISTINCT c.*
FROM showtimes sh
JOIN auditoriums a ON sh.auditorium_id = a.id
JOIN cinemas c ON a.cinema_id = c.id
WHERE sh.film_id = $1
    AND sh.show_date >= CURRENT_DATE
    AND sh.is_released = true
    AND c.is_released = true
    AND a.is_released = true;
