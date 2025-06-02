package response

type FilmCurrentlyShowing struct {
	FilmId    int32  `json:"film_id"`
	Title     string `json:"title"`
	PosterUrl string `json:"poster_url"`
	Genres    string `json:"genres"`
	Duration  string `json:"duration"`
}
