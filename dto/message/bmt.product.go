package message

type NewFilmCreationMsg struct {
	FilmId   int32  `json:"film_id" binding:"required"`
	Duration string `json:"duration" binding:"required"`
}

type NewFABCreateMsg struct {
	FABId int32 `json:"fab_id" binding:"required"`
	Price int32 `json:"price" binding:"required"`
}
