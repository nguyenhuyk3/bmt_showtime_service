package services

import (
	"bmt_showtime_service/dto/request"
	"context"
)

type IShowtime interface {
	AddShowtime(ctx context.Context, arg request.AddShowtimeRequest) (int, error)
	DeleteShowtime(ctx context.Context) (int, error)
	GetShowtime(ctx context.Context) (interface{}, int, error)
}
