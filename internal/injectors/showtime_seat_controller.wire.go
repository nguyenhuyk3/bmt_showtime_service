//go:build wireinject

package injectors

import (
	"bmt_showtime_service/internal/controllers"
	showtimeseat "bmt_showtime_service/internal/implementaions/showtime_seat"

	"github.com/google/wire"
)

func InitShowtimeSeatController() (*controllers.ShowtimeSeatController, error) {
	wire.Build(
		dbSet,
		redisSet,

		showtimeseat.NewShowtimeSeatService,
		controllers.NewShowtimeSeatController,
	)

	return &controllers.ShowtimeSeatController{}, nil
}
