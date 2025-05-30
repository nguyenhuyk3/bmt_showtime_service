//go:build wireinject

package injectors

import (
	"bmt_showtime_service/internal/controllers"
	"bmt_showtime_service/internal/implementaions/showtime"

	"github.com/google/wire"
)

func InitShowtimeController() (*controllers.ShowtimeController, error) {
	wire.Build(
		dbSet,
		redisSet,
		filmClientSet,

		showtime.NewShowtimeService,
		controllers.NewShowtimeController,
	)

	return &controllers.ShowtimeController{}, nil
}
