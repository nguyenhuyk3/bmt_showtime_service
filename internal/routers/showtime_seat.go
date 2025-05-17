package routers

import (
	"bmt_showtime_service/db/sqlc"
	"bmt_showtime_service/global"
	"bmt_showtime_service/internal/controllers"
	"bmt_showtime_service/internal/implementaions/redis"
	showtimeseat "bmt_showtime_service/internal/implementaions/showtime_seat"

	"github.com/gin-gonic/gin"
)

type ShowtimeSeatRouter struct {
}

func (sr *ShowtimeSeatRouter) InitShowtimeSeatRouter(router *gin.RouterGroup) {
	sqlStore := sqlc.NewStore(global.Postgresql)
	redisClient := redis.NewRedisClient()
	showtimeSeatService := showtimeseat.NewShowtimeSeatService(sqlStore, redisClient)
	showtimeSeatController := controllers.NewShowtimeSeatController(showtimeSeatService)

	showtimeSeatPublicRouter := router.Group("/showtime_seat")
	{
		showtimePublicRouter := showtimeSeatPublicRouter.Group("/public")
		{
			showtimePublicRouter.GET("/get_all/", showtimeSeatController.GetAllShowtimeSeatsByShowtimeId)
		}
	}
}
