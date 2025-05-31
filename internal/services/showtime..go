package services

import (
	"bmt_showtime_service/dto/request"
	"context"
)

type IShowtime interface {
	AddShowtime(ctx context.Context, arg request.AddShowtimeReq) (int, error)
	ReleaseShowtime(ctx context.Context, arg request.ReleaseShowtimeByIdReq) (int, error)
	GetShowtime(ctx context.Context, showtimeId int32) (interface{}, int, error)
	GetAllShowtimesByFilmIdInOneDate(ctx context.Context, arg request.GetAllShowtimesByFilmIdInOneDateReq) (interface{}, int, error)
	GetAllFilmsCurrentlyShowing(ctx context.Context) (any, int, error)
}
