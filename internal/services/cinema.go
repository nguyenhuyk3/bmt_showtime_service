package services

import "context"

type ICinema interface {
	GetCinemasForShowingFilm(ctx context.Context, filmId int32) (any, int, error)
}
