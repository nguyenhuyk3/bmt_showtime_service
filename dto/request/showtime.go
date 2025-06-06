package request

type AddShowtimeReq struct {
	FilmId       int32  `json:"film_id" binding:"required"`
	AuditoriumId int32  `json:"auditorium_id" binding:"required"`
	ShowDate     string `json:"show_date" binding:"required"`
	ChangedBy    string
	// StartTime    string `json:"start_time" binding:"required"`
	// EndTime      string `json:"end_time" binding:"required"`
}

type GetAllShowtimesByFilmIdInOneDateReq struct {
	FilmId   int32  `json:"film_id" binding:"required"`
	ShowDate string `json:"show_date" binding:"required"`
}

type ReleaseShowtimeByIdReq struct {
	ShowtimeId int32 `json:"showtime_id" binding:"required"`
	ChangedBy  string
}

type GetAllShowTimesByFilmIdAndByCinemaIdAndByAuditoriumIdAndInOneDateReq struct {
	FilmId       int32  `json:"film_id" binding:"required"`
	CinemaId     int32  `json:"cinema_id" binding:"required"`
	AuditoriumId int32  `json:"auditorium_id" binding:"required"`
	ShowDate     string `json:"show_date" binding:"required"`
}
