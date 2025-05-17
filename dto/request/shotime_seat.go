package request

type UpdateShowtimeSeatStatusReq struct {
	BookedBy string `json:"booked_by" binding:"required"`
}
