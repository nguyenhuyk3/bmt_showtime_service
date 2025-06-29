package initializations

import (
	"bmt_showtime_service/global"
	"fmt"
)

func Run() {
	loadConfigs()
	initPostgreSql()
	initRedis()
	initMessageBrokerReader()

	go initRPC()

	r := initRouter()

	r.Run(fmt.Sprintf("0.0.0.0:%s", global.Config.Server.ServerPort))
}
