package messages

type NewFilmCreationTopic struct {
	FilmId   int32  `json:"film_id" binding:"required"`
	Duration string `json:"duraton" binding:"required"`
}
