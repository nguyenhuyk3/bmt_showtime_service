package services

import "context"

type ICinema interface {
	GetCinemasForShowingFilmByFilmId(ctx context.Context, filmId int32) (any, int, error)
}
