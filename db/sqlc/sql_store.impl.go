package sqlc

import (
	"bmt_showtime_service/dto/message"
	"bmt_showtime_service/dto/request"
	"bmt_showtime_service/global"
	"context"
	"database/sql"
	"errors"
	"fmt"

	"product"

	"github.com/jackc/pgx/v5/pgxpool"
)

type SqlStore struct {
	connPool *pgxpool.Pool
	*Queries
	ProductClient product.ProductClient
}

// ReleaseShowtimeTran implements IStore.
func (s *SqlStore) ReleaseShowtimeTran(ctx context.Context, arg request.ReleaseShowtimeByIdReq) error {
	err := s.execTran(ctx, func(q *Queries) error {
		isShowtimeRealeased, err := q.isShowtimeRealeased(ctx, arg.ShowtimeId)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return global.ErrNoShowtimeExist
			}

			return fmt.Errorf("failed to check showtime existence: %w", err)
		}
		if isShowtimeRealeased {
			return global.ErrShowtimeHaveBeenReleased
		}

		err = q.releaseShowtime(ctx, arg.ShowtimeId)
		if err != nil {
			return fmt.Errorf("an error occur when querying db: %w", err)
		}

		err = q.updateShowtime(ctx,
			updateShowtimeParams{
				ID:        arg.ShowtimeId,
				ChangedBy: arg.ChangedBy,
			})
		if err != nil {
			return fmt.Errorf("failed to update updater: %w", err)
		}

		err = q.createShowtimeSeats(ctx, arg.ShowtimeId)
		if err != nil {
			return fmt.Errorf("failed to create showtime seats: %w", err)
		}

		return nil
	})

	return err
}

// UpdateSeatStatus implements IStore.
func (s *SqlStore) UpdateSeatStatusTran(ctx context.Context, arg message.PayloadSubOrderData, seatStatus string) error {
	return s.execTran(ctx, func(q *Queries) error {
		switch seatStatus {
		case global.ORDER_FAILED:
			for _, seat := range arg.Seats {
				if err := q.UpdateShowtimeSeatByIdAndShowtimeIdFailed(ctx,
					UpdateShowtimeSeatByIdAndShowtimeIdFailedParams{
						SeatID:     seat.SeatId,
						Status:     SeatStatusesAvailable,
						ShowtimeID: arg.ShowtimeId,
					}); err != nil {
					return fmt.Errorf("failed to update seat %d for showtime %d (%s): %w", seat.SeatId, arg.ShowtimeId, seatStatus, err)
				}
			}
		case global.ORDER_SUCCESS:
			for _, seat := range arg.Seats {
				if err := q.UpdateShowtimeSeatByIdAndShowtimeIdSuccess(ctx,
					UpdateShowtimeSeatByIdAndShowtimeIdSuccessParams{
						SeatID:     seat.SeatId,
						Status:     SeatStatusesBooked,
						ShowtimeID: arg.ShowtimeId,
					}); err != nil {
					return fmt.Errorf("failed to update seat %d for showtime %d (%s): %w", seat.SeatId, arg.ShowtimeId, seatStatus, err)
				}
			}

		default:
			return fmt.Errorf("invalid seat status: %s", seatStatus)
		}

		return nil
	})
}

// HandleOrderCreatedTran implements IStore.
func (s *SqlStore) HandleOrderCreatedTran(ctx context.Context, arg message.PayloadOrderData) (int32, error) {
	var totalPrice int32 = 0

	err := s.execTran(ctx, func(q *Queries) error {
		for _, seat := range arg.Seats {
			err := q.UpdateShowtimeSeatByIdAndShowtimeId(ctx,
				UpdateShowtimeSeatByIdAndShowtimeIdParams{
					SeatID:     seat.SeatId,
					Status:     SeatStatusesReserved,
					BookedBy:   arg.OrderedBy,
					ShowtimeID: arg.ShowtimeId,
				})
			if err != nil {
				return fmt.Errorf("an error occur when updating showtime seat %d: %w", seat.SeatId, err)
			}

			seatPrice, err := q.GetPriceOfSeatBySeatId(ctx, seat.SeatId)
			if err != nil {
				return fmt.Errorf("an error occur when get price of seat by id (%d): %w", seat.SeatId, err)
			}

			// Caculate seat price
			totalPrice = totalPrice + seatPrice
		}

		if len(arg.FABs) != 0 {
			for _, fAB := range arg.FABs {
				// get price of fab from product service
				resp, err := s.ProductClient.GetPriceOfFAB(ctx,
					&product.GetPriceOfFABReq{
						FABId: fAB.FabId,
					})
				if err != nil {
					return fmt.Errorf("an error occur when get price of fab by id (%d): %w", fAB.FabId, err)
				}

				// Caculate fab price
				totalPrice = totalPrice + resp.Price*1000
			}
		}

		return nil
	})

	if err != nil {
		return -1, fmt.Errorf("failed to update seat status tran: %w", err)
	}

	return totalPrice, nil
}

func (s *SqlStore) execTran(ctx context.Context, fn func(*Queries) error) error {
	// Start transaction
	tran, err := s.connPool.Begin(ctx)
	if err != nil {
		return err
	}

	q := New(tran)
	// fn performs a series of operations down the db
	err = fn(q)
	if err != nil {
		// If an error occurs, rollback the transaction
		if rbErr := tran.Rollback(ctx); rbErr != nil {
			return fmt.Errorf("tran err: %v, rollback err: %v", err, rbErr)
		}

		return err
	}

	return tran.Commit(ctx)
}

func NewStore(
	connPool *pgxpool.Pool,
	productClient product.ProductClient) IStore {
	return &SqlStore{
		connPool:      connPool,
		Queries:       New(connPool),
		ProductClient: productClient,
	}
}
