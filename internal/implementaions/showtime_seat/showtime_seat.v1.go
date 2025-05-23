package showtimeseat

import (
	"bmt_showtime_service/db/sqlc"
	"bmt_showtime_service/dto/request"
	"bmt_showtime_service/global"
	"bmt_showtime_service/internal/services"
	"context"
	"fmt"
	"net/http"
	"time"
)

type ShowtimeSeatService struct {
	SqlStore    sqlc.IStore
	RedisClient services.IRedis
}

const (
	two_days = 60 * 24 * 2
)

// GetAllShowtimeSeatsByShowtimeId implements services.IShowtimeSeat.
func (s *ShowtimeSeatService) GetAllShowtimeSeatsByShowtimeId(ctx context.Context, showtimeId int32) (interface{}, int, error) {
	showDate, err := s.SqlStore.GetShowdateByShowtimeId(ctx, showtimeId)
	if err != nil {
		return nil, http.StatusInternalServerError, fmt.Errorf("failed to get show date: %v", err)
	}

	if !showDate.Valid {
		return nil, http.StatusBadRequest, fmt.Errorf("show date is invalid")
	}

	today := time.Now().Truncate(24 * time.Hour)
	showDateTime := showDate.Time.Truncate(24 * time.Hour)
	if showDateTime.Before(today) {
		return nil, http.StatusBadRequest, fmt.Errorf("cannot get showtime seats for past date (%s)", showDateTime.Format("2006-01-02"))
	}

	var seats []sqlc.ShowtimeSeat
	var key string = fmt.Sprintf("%s%d::%s", global.SHOWTIME_SEATS, showtimeId, showDateTime.Format("2006-01-02"))

	err = s.RedisClient.Get(key, &seats)
	if err != nil {
		if err.Error() == fmt.Sprintf("key %s does not exist", key) {
			seats, err := s.SqlStore.GetAllShowtimeSeatsByShowtimeId(ctx, showtimeId)
			if err != nil {
				return nil, http.StatusInternalServerError, fmt.Errorf("failed to get all showtime seats: %w", err)
			}

			savingErr := s.RedisClient.Save(key, &seats, two_days)
			if savingErr != nil {
				return nil, http.StatusInternalServerError, fmt.Errorf("warning: failed to save to Redis: %v", savingErr)
			}

			return seats, http.StatusOK, nil
		}

		return nil, http.StatusInternalServerError, fmt.Errorf("getting value occur an error: %w", err)
	}

	return seats, http.StatusOK, nil
}

// UpdateShowtimeSeatStatus implements services.IShowtimeSeat.
func (s *ShowtimeSeatService) UpdateShowtimeSeatStatus(ctx context.Context, arg request.UpdateShowtimeSeatStatusReq) (int, error) {
	panic("unimplemented")
}

func NewShowtimeSeatService(
	sqlStore sqlc.IStore,
	redisClient services.IRedis) services.IShowtimeSeat {
	return &ShowtimeSeatService{
		SqlStore:    sqlStore,
		RedisClient: redisClient,
	}
}
