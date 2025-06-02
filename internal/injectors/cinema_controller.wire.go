//go:build wireinject

package injectors

import (
	"bmt_showtime_service/internal/controllers"
	"bmt_showtime_service/internal/implementaions/cinema"
	"bmt_showtime_service/internal/injectors/provider"

	"github.com/google/wire"
)

func InitCinemaController() (*controllers.CinemaController, error) {
	wire.Build(
		provider.ProvideQueries,
		redisSet,

		cinema.NewCinemaService,
		controllers.NewCinemaController,
	)

	return &controllers.CinemaController{}, nil
}
