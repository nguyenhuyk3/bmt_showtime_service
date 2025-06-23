//go:build wireinject

package injectors

import (
	"bmt_showtime_service/internal/implementaions/nke"

	"github.com/google/wire"
)

func InitNKE() (*nke.NKE, error) {
	wire.Build(
		dbSet,
		redisSet,

		nke.NewNKE,
	)

	return &nke.NKE{}, nil
}
