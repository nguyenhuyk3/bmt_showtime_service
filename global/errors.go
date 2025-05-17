package global

import "errors"

var (
	ErrNoShowtimeExist          = errors.New("showtime id does not exist")
	ErrShowtimeHaveBeenReleased = errors.New("showtime have been released")
)
