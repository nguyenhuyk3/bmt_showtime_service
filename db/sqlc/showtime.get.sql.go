// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: showtime.get.sql

package sqlc

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
)

const getAllShowTimesByFilmIdAndByCinemaIdAndInDayRange = `-- name: GetAllShowTimesByFilmIdAndByCinemaIdAndInDayRange :many
SELECT sh.id, sh.film_id, sh.auditorium_id, sh.show_date, sh.start_time, sh.end_time, sh.is_released, sh.changed_by, sh.created_at, sh.updated_at
FROM showtimes sh
JOIN auditoriums a ON sh.auditorium_id = a.id
JOIN cinemas c ON a.cinema_id = c.id
WHERE sh.is_released = TRUE
    AND sh.film_id = $1
    AND c.id = $2
    AND sh.show_date BETWEEN CURRENT_DATE AND CURRENT_DATE + INTERVAL '14 days'
ORDER BY sh.show_date, sh.start_time
`

type GetAllShowTimesByFilmIdAndByCinemaIdAndInDayRangeParams struct {
	FilmID int32 `json:"film_id"`
	ID     int32 `json:"id"`
}

func (q *Queries) GetAllShowTimesByFilmIdAndByCinemaIdAndInDayRange(ctx context.Context, arg GetAllShowTimesByFilmIdAndByCinemaIdAndInDayRangeParams) ([]Showtime, error) {
	rows, err := q.db.Query(ctx, getAllShowTimesByFilmIdAndByCinemaIdAndInDayRange, arg.FilmID, arg.ID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []Showtime{}
	for rows.Next() {
		var i Showtime
		if err := rows.Scan(
			&i.ID,
			&i.FilmID,
			&i.AuditoriumID,
			&i.ShowDate,
			&i.StartTime,
			&i.EndTime,
			&i.IsReleased,
			&i.ChangedBy,
			&i.CreatedAt,
			&i.UpdatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getAllShowTimesByFilmIdInOneDate = `-- name: GetAllShowTimesByFilmIdInOneDate :many
SELECT id, film_id, auditorium_id, show_date, start_time, end_time, is_released, changed_by, created_at, updated_at 
FROM showtimes 
WHERE film_id = $1 AND show_date = $2 AND is_released = true
`

type GetAllShowTimesByFilmIdInOneDateParams struct {
	FilmID   int32       `json:"film_id"`
	ShowDate pgtype.Date `json:"show_date"`
}

func (q *Queries) GetAllShowTimesByFilmIdInOneDate(ctx context.Context, arg GetAllShowTimesByFilmIdInOneDateParams) ([]Showtime, error) {
	rows, err := q.db.Query(ctx, getAllShowTimesByFilmIdInOneDate, arg.FilmID, arg.ShowDate)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []Showtime{}
	for rows.Next() {
		var i Showtime
		if err := rows.Scan(
			&i.ID,
			&i.FilmID,
			&i.AuditoriumID,
			&i.ShowDate,
			&i.StartTime,
			&i.EndTime,
			&i.IsReleased,
			&i.ChangedBy,
			&i.CreatedAt,
			&i.UpdatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getFilmIdsInToday = `-- name: GetFilmIdsInToday :many
SELECT DISTINCT film_id
FROM showtimes 
WHERE show_date = $1
`

func (q *Queries) GetFilmIdsInToday(ctx context.Context, showDate pgtype.Date) ([]int32, error) {
	rows, err := q.db.Query(ctx, getFilmIdsInToday, showDate)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []int32{}
	for rows.Next() {
		var film_id int32
		if err := rows.Scan(&film_id); err != nil {
			return nil, err
		}
		items = append(items, film_id)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getLatestShowtimeByAuditoriumId = `-- name: GetLatestShowtimeByAuditoriumId :one
SELECT end_time
FROM showtimes
WHERE auditorium_id = $1 AND show_date = $2
ORDER BY created_at DESC
LIMIT 1
`

type GetLatestShowtimeByAuditoriumIdParams struct {
	AuditoriumID int32       `json:"auditorium_id"`
	ShowDate     pgtype.Date `json:"show_date"`
}

func (q *Queries) GetLatestShowtimeByAuditoriumId(ctx context.Context, arg GetLatestShowtimeByAuditoriumIdParams) (pgtype.Timestamp, error) {
	row := q.db.QueryRow(ctx, getLatestShowtimeByAuditoriumId, arg.AuditoriumID, arg.ShowDate)
	var end_time pgtype.Timestamp
	err := row.Scan(&end_time)
	return end_time, err
}

const getShowdateByShowtimeId = `-- name: GetShowdateByShowtimeId :one
SELECT show_date
FROM showtimes
WHERE id = $1
`

func (q *Queries) GetShowdateByShowtimeId(ctx context.Context, id int32) (pgtype.Date, error) {
	row := q.db.QueryRow(ctx, getShowdateByShowtimeId, id)
	var show_date pgtype.Date
	err := row.Scan(&show_date)
	return show_date, err
}

const getShowtimeById = `-- name: GetShowtimeById :one
SELECT id, film_id, auditorium_id, show_date, start_time, end_time, is_released, changed_by, created_at, updated_at
FROM showtimes
WHERE id = $1 AND is_released = true
LIMIT 1
`

func (q *Queries) GetShowtimeById(ctx context.Context, id int32) (Showtime, error) {
	row := q.db.QueryRow(ctx, getShowtimeById, id)
	var i Showtime
	err := row.Scan(
		&i.ID,
		&i.FilmID,
		&i.AuditoriumID,
		&i.ShowDate,
		&i.StartTime,
		&i.EndTime,
		&i.IsReleased,
		&i.ChangedBy,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const isShowtimeExist = `-- name: IsShowtimeExist :one
SELECT EXISTS (
    SELECT 1 FROM showtimes 
    WHERE id = $1
) AS EXISTS
`

func (q *Queries) IsShowtimeExist(ctx context.Context, id int32) (bool, error) {
	row := q.db.QueryRow(ctx, isShowtimeExist, id)
	var exists bool
	err := row.Scan(&exists)
	return exists, err
}

const isShowtimeRealeased = `-- name: isShowtimeRealeased :one
SELECT is_released
FROM showtimes
WHERE id = $1
`

func (q *Queries) isShowtimeRealeased(ctx context.Context, id int32) (bool, error) {
	row := q.db.QueryRow(ctx, isShowtimeRealeased, id)
	var is_released bool
	err := row.Scan(&is_released)
	return is_released, err
}
