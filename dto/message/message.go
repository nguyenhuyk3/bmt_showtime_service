package message

type SeatIsBookedMsg struct {
	ShowtimeSeatId int32  `json:"showtime_seat_id" binding:"required"`
	Status         string `json:"status" binding:"required"`
	BookedBy       string `json:"booked_by" binding:"required"`
}
