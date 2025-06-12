package routers

import (
	"bmt_showtime_service/internal/injectors"
	"log"

	"github.com/gin-gonic/gin"
)

type ShowtimeSeatRouter struct {
}

func (sr *ShowtimeSeatRouter) InitShowtimeSeatRouter(router *gin.RouterGroup) {
	showtimeSeatController, err := injectors.InitShowtimeSeatController()
	if err != nil {
		log.Fatalf("failed to initialize ShowtimeSeatController: %v", err)
		return
	}

	showtimeSeatPublicRouter := router.Group("/showtime_seat")
	{
		showtimePublicRouter := showtimeSeatPublicRouter.Group("/public")
		{
			showtimePublicRouter.GET("/get_all",
				showtimeSeatController.GetAllShowtimeSeatsByShowtimeId)
			showtimePublicRouter.GET("/get_all_from_earliest_tomorrow",
				showtimeSeatController.GetAllShowtimeSeatsFromEarliestTomorrow)
		}
	}
}
