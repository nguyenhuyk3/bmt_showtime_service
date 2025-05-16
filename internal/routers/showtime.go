package routers

import (
	"bmt_showtime_service/db/sqlc"
	"bmt_showtime_service/global"
	"bmt_showtime_service/internal/controllers"
	"bmt_showtime_service/internal/implementaions/redis"
	"bmt_showtime_service/internal/implementaions/showtime"
	"bmt_showtime_service/internal/middlewares"

	"github.com/gin-gonic/gin"
)

type ShowtimeRouter struct {
}

func (sr *ShowtimeRouter) InitShowtimeRouter(router *gin.RouterGroup) {
	sqlStore := sqlc.NewStore(global.Postgresql)
	redisClient := redis.NewRedisClient()
	showtimeService := showtime.NewShowtimeService(sqlStore, redisClient)
	showtimeController := controllers.NewShowtimeController(showtimeService)
	getFromHeaderMiddleware := middlewares.NewGetFromHeaderMiddleware()

	showtimePublicRouter := router.Group("/showtime")
	{
		adminShowtimePrivateRouter := showtimePublicRouter.Group("/admin")
		{
			adminShowtimePrivateRouter.POST("/add",
				getFromHeaderMiddleware.GetEmailFromHeader(),
				showtimeController.AddShowtime)
			adminShowtimePrivateRouter.POST("/release",
				getFromHeaderMiddleware.GetEmailFromHeader(),
				showtimeController.ReleaseShowtime)
		}

		showtimePublicRouter := showtimePublicRouter.Group("/public")
		{
			showtimePublicRouter.GET("/get/:showtime_id", showtimeController.GetShowTime)
			showtimePublicRouter.GET("/get_all_showtimes", showtimeController.GetAllShowTimesByFilmIdInOneDate)
		}
	}
}
