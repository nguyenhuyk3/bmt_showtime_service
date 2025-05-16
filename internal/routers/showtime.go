package routers

import (
	"bmt_showtime_service/internal/controllers"
	"bmt_showtime_service/internal/implementaions/redis"
	"bmt_showtime_service/internal/implementaions/showtime"
	"bmt_showtime_service/internal/middlewares"

	"github.com/gin-gonic/gin"
)

type ShowtimeRouter struct {
}

func (sr *ShowtimeRouter) InitShowtimeRouter(router *gin.RouterGroup) {
	redisClient := redis.NewRedisClient()
	showtimeService := showtime.NewShowtimeService(redisClient)
	showtimeController := controllers.NewShowtimeController(showtimeService)
	getFromHeaderMiddleware := middlewares.NewGetFromHeaderMiddleware()

	showtimePublicRouter := router.Group("/showtime")
	{
		adminShowtimePrivateRouter := showtimePublicRouter.Group("/admin")
		{
			adminShowtimePrivateRouter.POST("/add",
				getFromHeaderMiddleware.GetEmailFromHeader(),
				showtimeController.AddShowtime)
		}

		showtimePublicRouter := showtimePublicRouter.Group("/public")
		{
			showtimePublicRouter.GET("/get/:showtime_id", showtimeController.GetShowTime)
			showtimePublicRouter.GET("/get_all_showtimes", showtimeController.GetAllShowTimesByFilmIdInOneDate)
		}
	}
}
