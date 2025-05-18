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
	SqlStore    sqlc.IStore
	RedisClient services.IRedis
}

const (
	time_off = 25 * time.Minute
	two_days = 60 * 24 * 2
)

// AddShowtime implements services.IShowtime.
func (s *showtimeService) AddShowtime(ctx context.Context, arg request.AddShowtimeReq) (int, error) {
	isFilmExist, err := s.SqlStore.IsFilmIdExist(ctx, arg.FilmId)
	if err != nil {
		return http.StatusInternalServerError, fmt.Errorf("failed to check film existence: %w", err)
	}
	if !isFilmExist {
		return http.StatusNotFound, fmt.Errorf("film doesn't exist")
	}

	isAuditoriumExist, err := s.SqlStore.IsAuditoriumExist(ctx, arg.AuditoriumId)
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

	if showDate.Day() == time.Now().Day() &&
		showDate.Month() == time.Now().Month() &&
		showDate.Year() == time.Now().Year() {
		return http.StatusBadRequest, fmt.Errorf("cannot add showtime for today")
	}

	var startTime time.Time

	latestShowtime, err := s.SqlStore.GetLatestShowtimeByAuditoriumId(ctx, arg.AuditoriumId)

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

	filmDuration, err := s.SqlStore.GetDuration(ctx, arg.FilmId)
	if err != nil {
		return http.StatusInternalServerError, fmt.Errorf("failed to get film duration: %w", err)
	}
	if !filmDuration.Valid {
		return http.StatusBadRequest, fmt.Errorf("film duration is invalid")
	}

	duration := time.Duration(filmDuration.Microseconds) * time.Microsecond
	rawEndTime := convertors.RoundDurationToNearestFive(duration)
	endTime := startTime.Add(rawEndTime)

	err = s.SqlStore.CreateShowTime(ctx, sqlc.CreateShowTimeParams{
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

// ReleaseShowtime implements services.IShowtime.
func (s *showtimeService) ReleaseShowtime(ctx context.Context, arg request.ReleaseShowtimeByIdReq) (int, error) {
	err := s.SqlStore.ReleaseShowtimeTran(ctx, arg)
	if err != nil {
		if errors.Is(err, global.ErrNoShowtimeExist) {
			return http.StatusNotFound, err
		}
		if errors.Is(err, global.ErrShowtimeHaveBeenReleased) {
			return http.StatusBadRequest, err
		}
		return http.StatusInternalServerError, err
	}

	showtime, err := s.SqlStore.GetShowtimeById(context.Background(), arg.ShowtimeId)
	if err != nil {
		return http.StatusInternalServerError, fmt.Errorf("failed to get showtime: %w", err)
	}

	go func() {
		showDate := showtime.ShowDate.Time.Format("2006-01-02")
		showtimes, _ := s.SqlStore.GetAllShowTimesByFilmIdInOneDate(context.Background(),
			sqlc.GetAllShowTimesByFilmIdInOneDateParams{
				FilmID: showtime.FilmID,
				ShowDate: pgtype.Date{
					Time:  showtime.ShowDate.Time,
					Valid: true,
				},
			})
		key := fmt.Sprintf("%s%d::%s", global.SHOWTIME_FILM_DATE, showtime.FilmID, showDate)

		_ = s.RedisClient.Delete(key)
		_ = s.RedisClient.Save(key, &showtimes, two_days)
	}()

	go func() {
		_ = s.RedisClient.Save(fmt.Sprintf("%s%d", global.SHOWTIME, arg.ShowtimeId), showtime, two_days)
		_ = s.RedisClient.Save(fmt.Sprintf("%s%d", global.SHOWTIME_ID, arg.ShowtimeId), arg.ShowtimeId, two_days)
	}()

	return http.StatusOK, nil
}

// GetShowtime implements services.IShowtime.
func (s *showtimeService) GetShowtime(ctx context.Context, showtimeId int32) (interface{}, int, error) {
	isShowtimeExist, err := s.SqlStore.IsShowtimeExist(ctx, showtimeId)
	if err != nil {
		return nil, http.StatusInternalServerError, fmt.Errorf("failed to check showtime existence: %w", err)
	}
	if !isShowtimeExist {
		return nil, http.StatusNotFound, fmt.Errorf("showtime with id %d does not exist", showtimeId)
	}

	var showtime sqlc.Showtimes
	var key string = fmt.Sprintf("%s%d", global.SHOWTIME, showtimeId)

	err = s.RedisClient.Get(key, &showtime)
	if err != nil {
		if err.Error() == fmt.Sprintf("key %s does not exist", key) {
			queriedShowtime, err := s.SqlStore.GetShowtimeById(ctx, showtimeId)
			if err != nil {
				if errors.Is(err, sql.ErrNoRows) {
					return nil, http.StatusNotFound, fmt.Errorf("showtime with id %d does not release", showtimeId)
				}

				return nil, http.StatusInternalServerError, fmt.Errorf("failed to get showtime: %w", err)
			}

			savingErr := s.RedisClient.Save(key, &queriedShowtime, two_days)
			if savingErr != nil {
				return nil, http.StatusInternalServerError, fmt.Errorf("warning: failed to save to Redis: %v", savingErr)
			}

			return queriedShowtime, http.StatusOK, nil
		}

		return nil, http.StatusInternalServerError, fmt.Errorf("getting value occur an error: %w", err)
	}

	return showtime, http.StatusOK, nil
}

// GetAllShowTimesByFilmIdInOneDate implements services.IShowtime.
func (s *showtimeService) GetAllShowtimesByFilmIdInOneDate(
	ctx context.Context,
	arg request.GetAllShowtimesByFilmIdInOneDateReq) (interface{}, int, error) {
	showDate, err := convertors.ConvertDateStringToTime(arg.ShowDate)
	if err != nil {
		return nil, http.StatusBadRequest, err
	}

	var key string = fmt.Sprintf("%s%d::%s", global.SHOWTIME_FILM_DATE, arg.FilmId, arg.ShowDate)
	var showtimes []sqlc.Showtimes

	err = s.RedisClient.Get(key, &showtimes)
	if err != nil {
		if err.Error() == fmt.Sprintf("key %s does not exist", key) {
			showtimes, err = s.SqlStore.GetAllShowTimesByFilmIdInOneDate(ctx,
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

			savingErr := s.RedisClient.Save(key, &showtimes, two_days)
			if savingErr != nil {
				return nil, http.StatusInternalServerError, fmt.Errorf("warning: failed to save to Redis: %v", savingErr)
			}

			return showtimes, http.StatusOK, nil
		}

		return nil, http.StatusInternalServerError, fmt.Errorf("redis error: %w", err)
	}

	return showtimes, http.StatusOK, nil
}

func NewShowtimeService(
	sqlStore sqlc.IStore,
	redisClient services.IRedis) services.IShowtime {
	return &showtimeService{
		SqlStore:    sqlStore,
		RedisClient: redisClient,
	}
}
