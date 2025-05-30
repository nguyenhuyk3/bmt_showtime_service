package injectors

import (
	"bmt_showtime_service/db/sqlc"
	"bmt_showtime_service/internal/implementaions/redis"
	"bmt_showtime_service/internal/injectors/provider"

	"github.com/google/wire"
)

var dbSet = wire.NewSet(
	provider.ProvidePgxPool,
	provider.ProvideProductClient,
	sqlc.NewStore,
)

var redisSet = wire.NewSet(
	redis.NewRedisClient,
)

var productClientSet = wire.NewSet(
	provider.ProvideProductClient,
)
