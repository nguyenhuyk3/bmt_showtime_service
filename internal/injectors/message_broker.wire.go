//go:build wireinject

package injectors

import (
	"bmt_showtime_service/internal/implementaions/message_broker/readers"
	"bmt_showtime_service/internal/injectors/provider"

	"github.com/google/wire"
)

func InitMessageBroker() (*readers.MessageBrokerReader, error) {
	wire.Build(
		provider.ProvideQueries,
		redisSet,

		readers.NewMessageBrokerReader,
	)

	return nil, nil
}
