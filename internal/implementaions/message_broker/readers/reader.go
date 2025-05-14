package readers

import (
	"bmt_showtime_service/db/sqlc"
	"bmt_showtime_service/global"
	"context"
	"log"
)

type MessageBrokerReader struct {
	SqlQuery sqlc.Querier
	Context  context.Context
}

var topics = []string{
	global.NEW_FILM_WAS_CREATED_TOPIC,
}

func NewMessageBrokerReader(
	sqlQuery *sqlc.Queries,
) *MessageBrokerReader {
	return &MessageBrokerReader{
		SqlQuery: sqlQuery,
	}
}

func (m *MessageBrokerReader) InitReaders() {
	log.Printf("=============== Showtime Service is listening to messages about new film creation ... ===============\n\n\n")

	for _, topic := range topics {
		go m.startReader(topic)
	}
}
