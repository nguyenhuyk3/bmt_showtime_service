package readers

import (
	"bmt_showtime_service/db/sqlc"
	"bmt_showtime_service/global"
	"bmt_showtime_service/internal/services"
	"context"
	"log"
)

type MessageBrokerReader struct {
	SqlQuery    sqlc.IStore
	RedisClient services.IRedis
	Context     context.Context
}

var topics = []string{
	global.NEW_FILM_WAS_CREATED_TOPIC,
	global.BMT_ORDER_PUBLIC_OUTBOXES,
}

func NewMessageBrokerReader(
	sqlQuery sqlc.IStore,
	redisClient services.IRedis,
) *MessageBrokerReader {
	return &MessageBrokerReader{
		SqlQuery:    sqlQuery,
		RedisClient: redisClient,
		Context:     context.Background(),
	}
}

func (m *MessageBrokerReader) InitReaders() {
	log.Printf("=============== Showtime Service is listening to messages ... ===============\n\n\n")

	for _, topic := range topics {
		go m.startReader(topic)
	}
}
