package initializations

import (
	"bmt_showtime_service/db/sqlc"
	"bmt_showtime_service/global"
	"bmt_showtime_service/internal/implementaions/message_broker/readers"
)

func initMessageBrokerReader() {
	reader := readers.NewMessageBrokerReader(sqlc.New(global.Postgresql))

	reader.InitReaders()
}
