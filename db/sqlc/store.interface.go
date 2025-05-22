package sqlc

import (
	"bmt_showtime_service/dto/message"
	"bmt_showtime_service/dto/request"
	"context"
)

type IStore interface {
	Querier
	ReleaseShowtimeTran(ctx context.Context, arg request.ReleaseShowtimeByIdReq) error
	UpdateSeatStatusTran(ctx context.Context, arg message.PayloadSubOrderData, seatStatus string) error
}
