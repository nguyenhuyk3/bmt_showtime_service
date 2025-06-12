package services

import (
	"bmt_showtime_service/dto/request"
	"context"
)

type IShowtimeSeat interface {
	GetAllShowtimeSeatsByShowtimeId(ctx context.Context,
		showtimeId int32) (interface{}, int, error)
	UpdateShowtimeSeatStatus(ctx context.Context,
		arg request.UpdateShowtimeSeatStatusReq) (int, error)
	GetShowtimeSeatsFromEarliestTomorrow(ctx context.Context,
		arg request.GetShowtimeSeatsFromEarliestTomorrowReq) (any, int, error)
}
