package request

type AddShowtimeRequest struct {
	FilmId       int32 `json:"film_id" binding:"required"`
	AuditoriumId int32 `json:"auditorium_id" binding:"required"`
	ChangedBy    string
	StartTime    string `json:"start_time" binding:"required"`
	EndTime      string `json:"end_time" binding:"required"`
}
