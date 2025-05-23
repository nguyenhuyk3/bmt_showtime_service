package readers

import (
	"bmt_showtime_service/db/sqlc"
	"bmt_showtime_service/dto/message"
	"bmt_showtime_service/global"
	"bmt_showtime_service/utils/convertors"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/segmentio/kafka-go"
)

func (m *MessageBrokerReader) startReader(topic string) {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{
			global.Config.ServiceSetting.KafkaSetting.KafkaBroker_1,
			global.Config.ServiceSetting.KafkaSetting.KafkaBroker_2,
			global.Config.ServiceSetting.KafkaSetting.KafkaBroker_3,
		},
		GroupID:        global.SHOWTIME_SERVICE_GROUP,
		Topic:          topic,
		CommitInterval: time.Second * 5,
	})
	defer reader.Close()

	for {
		message, err := reader.ReadMessage(context.Background())
		if err != nil {
			log.Printf("error reading message: %v\n", err)
			continue
		}

		m.processMessage(topic, message.Value)
	}
}

func (m *MessageBrokerReader) processMessage(topic string, value []byte) {
	switch topic {
	case global.BMT_PRODUCT_PUBLIC_OUTBOXES:
		var productMessage message.BMTPublicOutboxesMsg
		if err := json.Unmarshal(value, &productMessage); err != nil {
			log.Printf("failed to unmarshal new film updating message: %v\n", err)
			return
		}

		switch productMessage.After.EventType {
		case global.FILM_CREATED:
			var payloadProductFilmData message.NewFilmCreationMsg
			if err := json.Unmarshal([]byte(productMessage.After.Payload), &payloadProductFilmData); err != nil {
				log.Printf("failed to parse payload (%s): %v", productMessage.After.EventType, err)
				return
			}

			m.handleFilmCreation(payloadProductFilmData)
		case global.FAB_CREATED:
			var payloadProductFABData message.NewFABCreateMsg
			if err := json.Unmarshal([]byte(productMessage.After.Payload), &payloadProductFABData); err != nil {
				log.Printf("failed to parse payload (%s): %v", productMessage.After.EventType, err)
				return
			}

			m.hanleFABCreation(payloadProductFABData)
		default:
			log.Printf("unknown event type received: %s\n", productMessage.After.EventType)
		}

	// this case will handle messages from Order service
	case global.BMT_ORDER_PUBLIC_OUTBOXES:
		var orderMessage message.BMTPublicOutboxesMsg
		if err := json.Unmarshal(value, &orderMessage); err != nil {
			log.Printf("failed to unmarshal order message: %v\n", err)
			return
		}

		switch orderMessage.After.EventType {
		// change seat status available -> reserved
		case global.ORDER_CREATED:
			var payloadOrderData message.PayloadOrderData
			if err := json.Unmarshal([]byte(orderMessage.After.Payload), &payloadOrderData); err != nil {
				log.Printf("failed to parse payload (%s): %v", orderMessage.After.EventType, err)
				return
			}

			m.handleOrderCreated(payloadOrderData)

		// change seat status reserved -> available
		case global.ORDER_FAILED:
			var payloadSubOrderData message.PayloadSubOrderData
			if err := json.Unmarshal([]byte(orderMessage.After.Payload), &payloadSubOrderData); err != nil {
				log.Printf("failed to parse payload (%s): %v", orderMessage.After.EventType, err)
				return
			}

			m.handleOrderFailed(payloadSubOrderData, global.ORDER_SUCCESS)

		// change seat status reserved -> booked
		case global.ORDER_SUCCESS:
			var payloadSubOrderData message.PayloadSubOrderData
			if err := json.Unmarshal([]byte(orderMessage.After.Payload), &payloadSubOrderData); err != nil {
				log.Printf("failed to parse payload (%s): %v", orderMessage.After.EventType, err)
				return
			}

			m.handleOrderSuccess(payloadSubOrderData, global.ORDER_SUCCESS)

		default:
			log.Printf("unknown event type received: %s\n", orderMessage.After.EventType)
		}

	default:
		log.Printf("unknown topic received: %s\n", topic)
	}
}

func (m *MessageBrokerReader) handleFilmCreation(message message.NewFilmCreationMsg) {
	duration, err := convertors.ParseDurationToPGInterval(message.Duration)
	if err != nil {
		log.Printf("an error occurre when converting to duration: %v", err)
		return
	}

	err = m.SqlQuery.CreateNewFilmId(m.Context,
		sqlc.CreateNewFilmIdParams{
			FilmID:   message.FilmId,
			Duration: duration,
		})
	if err != nil {
		log.Printf("an error occur when creating new film id (%d): %v", message.FilmId, err)
	} else {
		log.Printf("create new film id (%d) successfully", message.FilmId)
	}
}

func (m *MessageBrokerReader) hanleFABCreation(payload message.NewFABCreateMsg) {
	err := m.SqlQuery.CreateNewFABInfo(m.Context,
		sqlc.CreateNewFABInfoParams{
			FabID: payload.FABId,
			Price: payload.Price,
		})
	if err != nil {
		log.Printf("an error occur when creating new fab id (%d): %v", payload.FABId, err)
	} else {
		log.Printf("create new fab id (%d) successfully", payload.FABId)
	}
}

func (m *MessageBrokerReader) handleOrderCreated(payload message.PayloadOrderData) {
	totalPrice, err := m.SqlQuery.HandleOrderCreatedTran(m.Context, payload)
	if err != nil {
		log.Printf("%v", err)
		return
	}

	if totalPrice != 0 {
		_ = m.RedisClient.Save(fmt.Sprintf("%s%d", global.ORDER_TOTAL, payload.OrderId), gin.H{
			"total_price": totalPrice,
		}, 5)
	}
}

func (m *MessageBrokerReader) handleOrderFailed(payload message.PayloadSubOrderData, status string) {
	err := m.SqlQuery.UpdateSeatStatusTran(m.Context, payload, status)
	if err != nil {
		log.Printf("%v", err)
	} else {
		log.Printf("handle update seat status successfully (%d)- (%s)", payload.OrderId, status)
	}
}

func (m *MessageBrokerReader) handleOrderSuccess(payload message.PayloadSubOrderData, status string) {
	err := m.SqlQuery.UpdateSeatStatusTran(m.Context, payload, status)
	if err != nil {
		log.Printf("%v", err)
	} else {
		log.Printf("handle update seat status successfully (%d)- (%s)", payload.OrderId, status)
	}
}
