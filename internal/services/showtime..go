package services

import (
	"bmt_showtime_service/dto/request"
	"context"
)

type IShowtime interface {
	AddShowtime(ctx context.Context, arg request.AddShowtimeRequest) (int, error)
	TurnOnShowtime(ctx context.Context, showtimeId int32) (int, error)
	GetShowtime(ctx context.Context, showtimeId int32) (interface{}, int, error)
	GetAllShowTimesByFilmIdInOneDate(ctx context.Context, arg request.GetAllShowTimesInOneDateRequest) (interface{}, int, error)
}
