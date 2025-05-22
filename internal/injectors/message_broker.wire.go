//go:build wireinject

package injectors

import (
	"bmt_showtime_service/internal/implementaions/message_broker/readers"

	"github.com/google/wire"
)

func InitMessageBroker() (*readers.MessageBrokerReader, error) {
	wire.Build(
		dbSet,
		redisSet,

		readers.NewMessageBrokerReader,
	)

	return nil, nil
}
