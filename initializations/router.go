package initializations

import (
	"bmt_showtime_service/internal/routers"

	"github.com/gin-gonic/gin"
)

func initRouter() *gin.Engine {
	r := gin.Default()

	// Routers
	showtimetRouter := routers.ShowtimeServiceRouterGroup.Showtime
	showtimeSeatRouter := routers.ShowtimeServiceRouterGroup.ShowtimeSeat
	cinemaRouter := routers.ShowtimeServiceRouterGroup.Cinema

	mainGroup := r.Group("/v1")
	{
		showtimetRouter.InitShowtimeRouter(mainGroup)
		showtimeSeatRouter.InitShowtimeSeatRouter(mainGroup)
		cinemaRouter.InitCinmaRouter(mainGroup)
	}

	return r
}
