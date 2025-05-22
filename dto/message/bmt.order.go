package message

type PayloadOrderData struct {
	Seats      []SeatItem `json:"seats"`
	OrderedBy  string     `json:"ordered_by"`
	ShowtimeId int32      `json:"showtime_id"`
}

type PayloadSubOrderData struct {
	OrderId    int32      `json:"order_id" binding:"required"`
	ShowtimeId int32      `json:"showtime_id" binding:"required"`
	Seats      []SeatItem `json:"seats" binding:"required"`
	FABs       []FabItem  `json:"fab"`
}

type FabItem struct {
	FabId    int32 `json:"fab_id"`
	Quantity int32 `json:"quantity"`
}

type SeatItem struct {
	SeatId int32 `json:"seat_id"`
}
