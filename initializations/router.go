package initializations

import (
	"bmt_showtime_service/internal/routers"

	"github.com/gin-gonic/gin"
)

func initRouter() *gin.Engine {
	r := gin.Default()
	// Routers
	showtimetRouter := routers.ShowtimeServiceRouterGroup.Showtime

	mainGroup := r.Group("/v1")
	{
		showtimetRouter.InitShowtimeRouter(mainGroup)
	}

	return r
}
