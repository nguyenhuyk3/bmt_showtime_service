package showtime

import (
	"bmt_showtime_service/db/sqlc"
	"bmt_showtime_service/dto/request"
	"bmt_showtime_service/global"
	"bmt_showtime_service/internal/services"
	"bmt_showtime_service/utils/convertors"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
)

type showtimeService struct {
	Querier     sqlc.Querier
	RedisClient services.IRedis
}

const (
	time_off = 25 * time.Minute
)

/*
	Các bước thực hiện hàm AddShowtime
	Bước 1: Kiểm tra xem phim và rạp chiếu có tồn tại hay không
	Bước 2: Kiểm tra thời gian bắt đầu chiếu có hợp lệ hay không
		Nếu arg.StartTime == ""
			Bước 2.1: Kiểm tra thời gian chiếu gần nhất của phòng arg.AuditoriumId
			Bước 2.2: StartTime = Thời gian chiếu gần nhất + 20 phút
		Nếu arg.StartTime != ""
			Bước 2.1: Thì StartTime hợp lệ
	Bước 3: Kiểm tra thời gian kết thúc có hợp lệ hay không
		Nếu arg.EndTime == ""
			Bước 3.1: EndTime bằng StartTime + Duration
		Nếu arg.EndTIme != ""
			Bước 3.1: EndTime hợp lệ
	Bước 4: Thêm mới showtime
*/
// AddShowtime implements services.IShowtime.
func (s *showtimeService) AddShowtime(ctx context.Context, arg request.AddShowtimeRequest) (int, error) {
	isFilmExist, err := s.Querier.IsFilmIdExist(ctx, arg.FilmId)
	if err != nil {
		return http.StatusInternalServerError, fmt.Errorf("failed to check film existence: %w", err)
	}
	if !isFilmExist {
		return http.StatusNotFound, fmt.Errorf("film doesn't exist")
	}

	isAuditoriumExist, err := s.Querier.IsAuditoriumExist(ctx, arg.AuditoriumId)
	if err != nil {
		return http.StatusInternalServerError, fmt.Errorf("failed to check auditorium existence: %w", err)
	}
	if !isAuditoriumExist {
		return http.StatusNotFound, fmt.Errorf("auditorium not found for the given film")
	}

	showDate, err := convertors.ConvertDateStringToTime(arg.ShowDate)
	if err != nil {
		return http.StatusBadRequest, err
	}

	var startTime time.Time

	latestShowTime, err := s.Querier.GetLatestShowtimeByAuditoriumId(ctx, arg.AuditoriumId)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			startTime = time.Date(
				showDate.Year(), showDate.Month(), showDate.Day(),
				9, 0, 0, 0, time.Local,
			)
		} else {
			return http.StatusInternalServerError, fmt.Errorf("failed to get latest showtime: %w", err)
		}
	} else {
		startTime = latestShowTime.Time.Add(time_off)
	}

	filmDuration, err := s.Querier.GetDuration(ctx, arg.FilmId)
	if err != nil {
		return http.StatusInternalServerError, fmt.Errorf("failed to get film duration: %w", err)
	}
	if !filmDuration.Valid {
		return http.StatusBadRequest, fmt.Errorf("film duration is invalid")
	}

	duration := time.Duration(filmDuration.Microseconds) * time.Microsecond
	rawEndTime := convertors.RoundDurationToNearestFive(duration)
	endTime := startTime.Add(rawEndTime)

	err = s.Querier.CreateShowTime(ctx, sqlc.CreateShowTimeParams{
		FilmID:       arg.FilmId,
		AuditoriumID: arg.AuditoriumId,
		ChangedBy:    arg.ChangedBy,
		ShowDate: pgtype.Date{
			Time:  showDate,
			Valid: true,
		},
		StartTime: pgtype.Timestamp{
			Time:  startTime,
			Valid: true,
		},
		EndTime: pgtype.Timestamp{
			Time:  endTime,
			Valid: true,
		},
	})
	if err != nil {
		return http.StatusInternalServerError, fmt.Errorf("failed to create showtime: %w", err)
	}

	return http.StatusOK, nil
}

// DeleteShowtime implements services.IShowtime.
func (s *showtimeService) DeleteShowtime(ctx context.Context) (int, error) {
	panic("unimplemented")
}

// GetShowtime implements services.IShowtime.
func (s *showtimeService) GetShowtime(ctx context.Context) (interface{}, int, error) {
	panic("unimplemented")
}

func NewShowtimeService(redisClient services.IRedis) services.IShowtime {
	return &showtimeService{
		Querier:     sqlc.New(global.Postgresql),
		RedisClient: redisClient,
	}
}
