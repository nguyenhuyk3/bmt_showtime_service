package showtime

import (
	"bmt_showtime_service/db/sqlc"
	"bmt_showtime_service/dto/request"
	"bmt_showtime_service/global"
	"bmt_showtime_service/internal/services"
	"bmt_showtime_service/utils/convertors"
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/jackc/pgx/v5/pgtype"
)

type showtimeService struct {
	Querier     sqlc.Querier
	RedisClient services.IRedis
}

// AddShowtime implements services.IShowtime.
func (s *showtimeService) AddShowtime(ctx context.Context, arg request.AddShowtimeRequest) (int, error) {
	isFilmExists, err := s.Querier.IsFilmIdExist(ctx, arg.FilmId)
	if err != nil {
		return http.StatusInternalServerError, fmt.Errorf("failed to check film existence: %w", err)
	}

	if !isFilmExists {
		return http.StatusNotFound, fmt.Errorf("film doesn't exist")
	}

	isAuditoriumExist, err := s.Querier.IsAuditoriumExist(ctx, arg.FilmId)
	if err != nil {
		return http.StatusInternalServerError, fmt.Errorf("failed to check auditorium existence: %w", err)
	}

	if !isAuditoriumExist {
		return http.StatusNotFound, errors.New("auditorium not found for the given film")
	}

	startTime, err := convertors.ParseAndValidateTime(arg.StartTime)
	if err != nil {
		return http.StatusBadRequest, fmt.Errorf("invalid start time: %w", err)
	}

	endTime, err := convertors.ParseAndValidateTime(arg.EndTime)
	if err != nil {
		return http.StatusBadRequest, fmt.Errorf("invalid end time: %w", err)
	}
	if startTime.After(endTime) {
		return http.StatusBadRequest, errors.New("start time must be before end time")
	}

	err = s.Querier.CreateShowTime(ctx, sqlc.CreateShowTimeParams{
		FilmID:       arg.FilmId,
		AuditoriumID: arg.AuditoriumId,
		ChangedBy:    arg.ChangedBy,
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

// DeleteShowtime implements services.IShowtime.
func (s *showtimeService) DeleteShowtime(ctx context.Context) (int, error) {
	panic("unimplemented")
}

// GetShowtime implements services.IShowtime.
func (s *showtimeService) GetShowtime(ctx context.Context) (interface{}, int, error) {
	panic("unimplemented")
}

func NewShowtimeService(redisClient services.IRedis) services.IShowtime {
	return &showtimeService{
		Querier:     sqlc.New(global.Postgresql),
		RedisClient: redisClient,
	}
}
