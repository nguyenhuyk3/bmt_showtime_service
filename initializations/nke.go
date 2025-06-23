package initializations

import (
	"bmt_showtime_service/internal/injectors"
	"context"
	"log"
)

func initNKE() {
	nke, err := injectors.InitNKE()
	if err != nil {
		log.Fatalf("failed to initialize nke: %v", err)
		return
	}

	go nke.RunSubscribingToExpiredEvents(context.Background())
}
