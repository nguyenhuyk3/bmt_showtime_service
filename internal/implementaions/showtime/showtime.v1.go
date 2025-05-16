package showtime

import (
	"bmt_showtime_service/db/sqlc"
	"bmt_showtime_service/dto/request"
	"bmt_showtime_service/global"
	"bmt_showtime_service/internal/services"
	"bmt_showtime_service/utils/convertors"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
)

type showtimeService struct {
	Querier     sqlc.Querier
	RedisClient services.IRedis
}

const (
	time_off = 25 * time.Minute
)

// AddShowtime implements services.IShowtime.
func (s *showtimeService) AddShowtime(ctx context.Context, arg request.AddShowtimeRequest) (int, error) {
	isFilmExist, err := s.Querier.IsFilmIdExist(ctx, arg.FilmId)
	if err != nil {
		return http.StatusInternalServerError, fmt.Errorf("failed to check film existence: %w", err)
	}
	if !isFilmExist {
		return http.StatusNotFound, fmt.Errorf("film doesn't exist")
	}

	isAuditoriumExist, err := s.Querier.IsAuditoriumExist(ctx, arg.AuditoriumId)
	if err != nil {
		return http.StatusInternalServerError, fmt.Errorf("failed to check auditorium existence: %w", err)
	}
	if !isAuditoriumExist {
		return http.StatusNotFound, fmt.Errorf("auditorium not found for the given film")
	}

	showDate, err := convertors.ConvertDateStringToTime(arg.ShowDate)
	if err != nil {
		return http.StatusBadRequest, err
	}

	var startTime time.Time

	latestShowtime, err := s.Querier.GetLatestShowtimeByAuditoriumId(ctx, arg.AuditoriumId)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			startTime = time.Date(
				showDate.Year(), showDate.Month(), showDate.Day(),
				9, 0, 0, 0, time.Local,
			)
		} else {
			return http.StatusInternalServerError, fmt.Errorf("failed to get latest showtime: %w", err)
		}
	} else {
		if !latestShowtime.Time.IsZero() &&
			(latestShowtime.Time.Year() != showDate.Year() ||
				latestShowtime.Time.Month() != showDate.Month() ||
				latestShowtime.Time.Day() != showDate.Day()) {
			return http.StatusBadRequest, fmt.Errorf("latest showtime already crosses into the next day")
		}
		startTime = latestShowtime.Time.Add(time_off)
	}

	filmDuration, err := s.Querier.GetDuration(ctx, arg.FilmId)
	if err != nil {
		return http.StatusInternalServerError, fmt.Errorf("failed to get film duration: %w", err)
	}
	if !filmDuration.Valid {
		return http.StatusBadRequest, fmt.Errorf("film duration is invalid")
	}

	duration := time.Duration(filmDuration.Microseconds) * time.Microsecond
	rawEndTime := convertors.RoundDurationToNearestFive(duration)
	endTime := startTime.Add(rawEndTime)

	err = s.Querier.CreateShowTime(ctx, sqlc.CreateShowTimeParams{
		FilmID:       arg.FilmId,
		AuditoriumID: arg.AuditoriumId,
		ChangedBy:    arg.ChangedBy,
		ShowDate: pgtype.Date{
			Time:  showDate,
			Valid: true,
		},
		StartTime: pgtype.Timestamp{
			Time:  startTime,
			Valid: true,
		},
		EndTime: pgtype.Timestamp{
			Time:  endTime,
			Valid: true,
		},
	})
	if err != nil {
		return http.StatusInternalServerError, fmt.Errorf("failed to create showtime: %w", err)
	}

	return http.StatusOK, nil
}

// TurnOnShowtime implements services.IShowtime.
func (s *showtimeService) TurnOnShowtime(ctx context.Context, showtimeId int32) (int, error) {
	isShowtimeExist, err := s.Querier.IsShowtimeExist(ctx, showtimeId)
	if err != nil {
		return http.StatusInternalServerError, fmt.Errorf("failed to check showtime existence: %w", err)
	}
	if !isShowtimeExist {
		return http.StatusNotFound, fmt.Errorf("showtime with id %d does not exist", showtimeId)
	}

	err = s.Querier.TurnOnShowtime(ctx, showtimeId)
	if err != nil {
		return http.StatusInternalServerError, fmt.Errorf("an error occur when querying db: %w", err)
	}

	return http.StatusOK, nil
}

// GetShowtime implements services.IShowtime.
func (s *showtimeService) GetShowtime(ctx context.Context, showtimeId int32) (interface{}, int, error) {
	isShowtimeExist, err := s.Querier.IsShowtimeExist(ctx, showtimeId)
	if err != nil {
		return nil, http.StatusInternalServerError, fmt.Errorf("failed to check showtime existence: %w", err)
	}
	if !isShowtimeExist {
		return nil, http.StatusNotFound, fmt.Errorf("showtime with id %d does not exist", showtimeId)
	}

	showtime, err := s.Querier.GetShowtimeById(ctx, showtimeId)
	if err != nil {
		return nil, http.StatusInternalServerError, fmt.Errorf("failed to get showtime: %w", err)
	}

	return showtime, http.StatusOK, nil
}

// GetAllShowTimesByFilmIdInOneDate implements services.IShowtime.
func (s *showtimeService) GetAllShowTimesByFilmIdInOneDate(ctx context.Context, arg request.GetAllShowTimesInOneDateRequest) (interface{}, int, error) {
	showDate, err := convertors.ConvertDateStringToTime(arg.ShowDate)
	if err != nil {
		return nil, http.StatusBadRequest, err
	}

	showtimes, err := s.Querier.GetAllShowTimesByFilmIdInOneDate(ctx,
		sqlc.GetAllShowTimesByFilmIdInOneDateParams{
			FilmID: arg.FilmId,
			ShowDate: pgtype.Date{
				Time:  showDate,
				Valid: true,
			},
		})
	if err != nil {
		return nil, http.StatusInternalServerError, fmt.Errorf("failed to get showtimes: %w", err)
	}

	if len(showtimes) == 0 {
		return []interface{}{}, http.StatusNotFound, fmt.Errorf("no showtimes")
	}

	return showtimes, http.StatusOK, nil
}

func NewShowtimeService(redisClient services.IRedis) services.IShowtime {
	return &showtimeService{
		Querier:     sqlc.New(global.Postgresql),
		RedisClient: redisClient,
	}
}
