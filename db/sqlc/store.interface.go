package sqlc

import (
	"bmt_showtime_service/dto/request"
	"context"
)

type IStore interface {
	Querier
	ReleaseShowtimeTran(ctx context.Context, arg request.ReleaseShowtimeByIdReq) error
}
