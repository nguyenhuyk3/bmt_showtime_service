package sqlc

import (
	"bmt_showtime_service/dto/request"
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

type SqlStore struct {
	connPool *pgxpool.Pool
	*Queries
}

// ReleaseShowtimeTran implements IStore.
func (s *SqlStore) ReleaseShowtimeTran(ctx context.Context, arg request.ReleaseShowtimeByIdReq) error {
	err := s.execTran(ctx, func(q *Queries) error {
		isShowtimeDeleted, err := q.isShowtimeRealeased(ctx, arg.ShowtimeId)
		if err != nil {
			return fmt.Errorf("failed to check showtime existence: %w", err)
		}
		if isShowtimeDeleted {
			return fmt.Errorf("showtime have been released")
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

func NewStore(connPool *pgxpool.Pool) IStore {
	return &SqlStore{
		connPool: connPool,
		Queries:  New(connPool),
	}
}
