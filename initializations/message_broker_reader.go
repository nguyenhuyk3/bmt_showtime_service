package initializations

import (
	"bmt_showtime_service/internal/injectors"
	"log"
)

func initMessageBrokerReader() {
	reader, err := injectors.InitMessageBroker()
	if err != nil {
		log.Fatalf("an error occur when initiallizating SHOWTIME READERS: %v", err)
	}

	reader.InitReaders()
}
