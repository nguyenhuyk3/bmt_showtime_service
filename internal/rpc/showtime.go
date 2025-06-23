package rpc

import (
	"bmt_showtime_service/db/sqlc"
	"context"
	"fmt"
	rpc_showtime "showtime"
	"strings"
)

type ShowtimeRPCServer struct {
	SqlStore sqlc.Queries
	rpc_showtime.UnimplementedShowtimeServer
}

// GetSomeInformationForTicket implements showtime.ShowtimeServer.
func (s *ShowtimeRPCServer) GetSomeInformationForTicket(ctx context.Context,
	arg *rpc_showtime.GetSomeInformationForTicketReq) (*rpc_showtime.GetSomeInformationForTicketRes, error) {
	cinema, err := s.SqlStore.GetCinemaByShowtimeId(ctx, arg.ShowtimeId)
	if err != nil {
		return nil, fmt.Errorf("failed to get cinema with showtime id (%d): %w", arg.ShowtimeId, err)
	}

	seatNumbers := []string{}

	for _, seatId := range arg.SeatIds {
		seat, err := s.SqlStore.GetSeatById(ctx, seatId)
		if err != nil {
			return nil, fmt.Errorf("failed to get seat number with seat id (%d): %w", seatId, err)
		}

		seatNumbers = append(seatNumbers, seat.SeatNumber)
	}

	showtime, err := s.SqlStore.GetShowtimeById(ctx, arg.ShowtimeId)
	if err != nil {
		return nil, fmt.Errorf("failed to get show time with showtime id (%d): %w", arg.ShowtimeId, err)
	}

	return &rpc_showtime.GetSomeInformationForTicketRes{
		CinemaName: cinema.Name,
		City:       string(cinema.City),
		Location:   cinema.Location,
		RoomName:   cinema.Roomname,
		ShowDate:   showtime.ShowDate.Time.Format("2006-01-02"),
		StartTime:  showtime.StartTime.Time.Format("15:04"),
		Seats:      strings.Join(seatNumbers, ", "),
		FilmId:     showtime.FilmID,
	}, nil
}

func NewShowtimeRPCServer(
	sqlStore sqlc.Queries,
) rpc_showtime.ShowtimeServer {
	return &ShowtimeRPCServer{
		SqlStore: sqlStore,
	}
}
