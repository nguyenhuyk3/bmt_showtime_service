package routers

import (
	"bmt_showtime_service/internal/injectors"
	"bmt_showtime_service/internal/middlewares"
	"log"

	"github.com/gin-gonic/gin"
)

type ShowtimeRouter struct {
}

func (sr *ShowtimeRouter) InitShowtimeRouter(router *gin.RouterGroup) {
	showtimeController, err := injectors.InitShowtimeController()
	if err != nil {
		log.Fatalf("failed to initialize ShowtimeController: %v", err)
		return
	}

	getFromHeaderMiddleware := middlewares.NewGetFromHeaderMiddleware()

	showtimeRouter := router.Group("/showtime")
	{
		adminShowtimePrivateRouter := showtimeRouter.Group("/admin")
		{
			adminShowtimePrivateRouter.POST("/add",
				getFromHeaderMiddleware.GetEmailFromHeader(),
				showtimeController.AddShowtime)
			adminShowtimePrivateRouter.POST("/release",
				getFromHeaderMiddleware.GetEmailFromHeader(),
				showtimeController.ReleaseShowtime)
		}

		showtimePublicRouter := showtimeRouter.Group("/public")
		{
			showtimePublicRouter.GET("/get/:showtime_id", showtimeController.GetShowTime)
			showtimePublicRouter.GET("/get_all_showtimes", showtimeController.GetAllShowTimesByFilmIdInOneDate)
			showtimePublicRouter.GET("/get_all_film_currently_showing", showtimeController.GetAllFilmsCurrentlyShowing)
		}
	}
}
