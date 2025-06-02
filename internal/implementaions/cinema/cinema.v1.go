package cinema

import (
	"bmt_showtime_service/db/sqlc"
	"bmt_showtime_service/internal/services"
	"context"
	"fmt"
	"net/http"
)

type cinemaService struct {
	SqlStore    *sqlc.Queries
	RedisClient services.IRedis
}

// GetCinemasForShowingFilm implements services.ICinema.
func (c *cinemaService) GetCinemasForShowingFilm(ctx context.Context, filmId int32) (any, int, error) {
	cinema, err := c.SqlStore.GetCinemasForShowingFilm(ctx, filmId)
	if err != nil {
		return nil, http.StatusInternalServerError, fmt.Errorf("failed to get cinema for showing film: %w", err)
	}

	if len(cinema) == 0 {
		return nil, http.StatusNotFound, fmt.Errorf("the film was not released with id (%d)", filmId)
	}

	return cinema, http.StatusOK, nil
}

func NewCinemaService(
	sqlStore *sqlc.Queries,
	redisClient services.IRedis) services.ICinema {
	return &cinemaService{
		SqlStore:    sqlStore,
		RedisClient: redisClient,
	}
}
