package routers

import (
	"bmt_showtime_service/internal/injectors"
	"log"

	"github.com/gin-gonic/gin"
)

type CinemaRouter struct{}

func (cr *CinemaRouter) InitCinemaRouter(router *gin.RouterGroup) {
	cinemaController, err := injectors.InitCinemaController()
	if err != nil {
		log.Fatalf("failed to initialize CinemaController: %v", err)
		return
	}

	cinemaRouter := router.Group("/cinema")
	{
		cinemaPublicRouter := cinemaRouter.Group("/public")
		{
			cinemaPublicRouter.GET("/get_cinemas_for_showing_film_by_film_id/:film_id",
				cinemaController.GetCinemasForShowingFilmByFilmId)
		}
	}
}
