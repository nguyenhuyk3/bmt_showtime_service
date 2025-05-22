package request

type OrderSeatReq struct {
	SeatId int32 `json:"seat_id" binding:"required"`
}

type OrderFABReq struct {
	FABId    int32 `json:"fab_id" binding:"required"`
	Quantity int   `json:"quantity" binding:"required"`
}

type SubOrder struct {
	ShowtimeId int32          `json:"showtime_id" binding:"required"`
	Seats      []OrderSeatReq `json:"seats" binding:"required"`
	FABs       []OrderFABReq  `json:"fab"`
}
