package initializations

import (
	"bmt_showtime_service/db/sqlc"
	"bmt_showtime_service/global"
	"bmt_showtime_service/internal/implementaions/message_broker/readers"
	"bmt_showtime_service/internal/implementaions/redis"
)

func initMessageBrokerReader() {
	redisClient := redis.NewRedisClient()
	reader := readers.NewMessageBrokerReader(sqlc.New(global.Postgresql), redisClient)

	reader.InitReaders()
}
