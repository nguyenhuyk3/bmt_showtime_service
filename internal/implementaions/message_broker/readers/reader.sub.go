package readers

import (
	"bmt_showtime_service/db/sqlc"
	"bmt_showtime_service/dto/message"
	"bmt_showtime_service/global"
	"bmt_showtime_service/utils/convertors"
	"context"
	"encoding/json"
	"log"
	"time"

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
	case global.NEW_FILM_WAS_CREATED_TOPIC:
		var message message.NewFilmCreationMsg
		if err := json.Unmarshal(value, &message); err != nil {
			log.Printf("failed to unmarshal new film creation message: %v\n", err)
			return
		}

		m.handleNewFilmCreationTopic(message)

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
		}

	// case global.BMT_PAYMENT_PUBLIC_OUTBOXES:
	// 	var paymentMessage message.BMTPublicOutboxesMsg
	// 	if err := json.Unmarshal(value, &paymentMessage); err != nil {
	// 		log.Printf("failed to unmarshal payment message: %v\n", err)
	// 		return
	// 	}

	// 	m.handlePaymentStatusEvent(paymentMessage)

	default:
		log.Printf("unknown topic received: %s\n", topic)
	}
}

func (m *MessageBrokerReader) handleNewFilmCreationTopic(message message.NewFilmCreationMsg) {
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
		log.Printf("an error occur when creating new film id: %v", err)
	} else {
		log.Println("create new film id successfully")
	}
}

func (m *MessageBrokerReader) handleOrderCreated(payload message.PayloadOrderData) {
	for _, seat := range payload.Seats {
		err := m.SqlQuery.UpdateShowtimeSeatSeatByIdAndShowtimeId(m.Context,
			sqlc.UpdateShowtimeSeatSeatByIdAndShowtimeIdParams{
				SeatID:     seat.SeatId,
				Status:     sqlc.SeatStatusesReserved,
				BookedBy:   payload.OrderedBy,
				ShowtimeID: payload.ShowtimeId,
			})

		if err != nil {
			log.Printf("an error occur when updating showtime seat %d: %v", seat.SeatId, err)
			return
		} else {
			log.Printf("update showtime seat %d with showtime id %d successfully (reserved)", seat.SeatId, payload.ShowtimeId)
		}
	}

}

func (m *MessageBrokerReader) handleOrderFailed(payload message.PayloadSubOrderData, status string) {
	err := m.SqlQuery.UpdateSeatStatusTran(m.Context, payload, status)
	if err != nil {
		log.Printf("%v", err)
	}
}

func (m *MessageBrokerReader) handleOrderSuccess(payload message.PayloadSubOrderData, status string) {
	err := m.SqlQuery.UpdateSeatStatusTran(m.Context, payload, status)
	if err != nil {
		log.Printf("%v", err)
	}
}

// func (m *MessageBrokerReader) handlePaymentStatusEvent(messageData message.BMTPublicOutboxesMsg) {
// 	var payload message.PayloadPaymentData
// 	if err := json.Unmarshal([]byte(messageData.After.Payload), &payload); err != nil {
// 		log.Printf("failed to parse payload (%s): %v", messageData.After.EventType, err)
// 		return
// 	}

// 	orderRedisKey := fmt.Sprintf("%s%d", global.ORDER, payload.OrderId)

// 	var subOrder request.SubOrder
// 	if err := m.RedisClient.Get(orderRedisKey, &subOrder); err != nil {
// 		log.Printf("failed to get data with key %s", orderRedisKey)
// 		return
// 	}

// 	var status sqlc.SeatStatuses
// 	var bookedBy *string

// 	switch messageData.After.EventType {
// 	// If payment is successful, change seat from reserved -> booked
// 	case global.PAYMENT_SUCCESS:
// 		status = sqlc.SeatStatusesBooked
// 	// If payment is failed, change seat from reserved -> available, booked_by -> empty string
// 	case global.PAYMENT_FAILED:
// 		status = sqlc.SeatStatusesAvailable
// 		empty := ""
// 		bookedBy = &empty
// 	default:
// 		log.Printf("unsupported event type: %s", messageData.After.EventType)
// 		return
// 	}

// 	for _, seat := range subOrder.Seats {
// 		param := sqlc.UpdateShowtimeSeatSeatByIdAndShowtimeIdParams{
// 			SeatID:     seat.SeatId,
// 			Status:     status,
// 			ShowtimeID: subOrder.ShowtimeId,
// 		}
// 		if bookedBy != nil {
// 			param.BookedBy = *bookedBy
// 		}

// 		if err := m.SqlQuery.UpdateShowtimeSeatSeatByIdAndShowtimeId(m.Context, param); err != nil {
// 			log.Printf("error updating seat %d: %v", seat.SeatId, err)
// 			return
// 		}

// 		log.Printf("seat %d updated for showtime %d to status %s", seat.SeatId, subOrder.ShowtimeId, status)
// 	}
// }
