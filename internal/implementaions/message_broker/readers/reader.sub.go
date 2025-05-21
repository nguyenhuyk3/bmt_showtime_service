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
	case global.BMT_ORDER_PUBLIC_OUTBOXES:
		var orderMessage message.BMTOrderPublicOutboxesMsg
		if err := json.Unmarshal(value, &orderMessage); err != nil {
			log.Printf("failed to unmarshal seat is booked message: %v\n", err)
			return
		}

		var payloadData message.PayloadData
		if err := json.Unmarshal([]byte(orderMessage.After.Payload), &payloadData); err != nil {
			log.Printf("failed to parse payload (%s): %v", global.ORDER_CREATED, err)
			return
		}

		switch orderMessage.After.EventType {
		case global.ORDER_CREATED:
			m.handleOrderCreated(payloadData)
		}

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

// func (m *MessageBrokerReader) handleSeatIsBookedTopic(message message.SeatIsBookedMsg) {
// 	var status sqlc.NullSeatStatuses

// 	err := status.Scan(message.Status)
// 	if err != nil {
// 		log.Printf("invalid status (%s): %v", message.Status, err)
// 		return
// 	}

// 	err = m.SqlQuery.UpdateShowtimeSeatById(m.Context,
// 		sqlc.UpdateShowtimeSeatByIdParams{
// 			ID:     message.ShowtimeSeatId,
// 			Status: status.SeatStatuses,
// 		})
// 	if err != nil {
// 		log.Printf("an error occur when updating with seat id (%d): %v", message.ShowtimeSeatId, err)
// 	} else {
// 		log.Printf("update seat with id (%d) successfully", message.ShowtimeSeatId)

// 		go func() {
// 			showDate, _ := m.SqlQuery.GetShowdateByShowtimeId(context.Background(), message.ShowtimeSeatId)
// 			showDateTime := showDate.Time.Truncate(24 * time.Hour)
// 			key := fmt.Sprintf("%s%d::%s", global.SHOWTIME_SEATS, message.ShowtimeSeatId, showDateTime.Format("2006-01-02"))
// 			seats, _ := m.SqlQuery.GetAllShowtimeSeatsByShowtimeId(context.Background(), message.ShowtimeSeatId)

// 			_ = m.RedisClient.Delete(key)
// 			_ = m.RedisClient.Save(key, &seats, 60*24*2)
// 		}()
// 	}
// }

func (m *MessageBrokerReader) handleOrderCreated(payload message.PayloadData) {
	for _, seat := range payload.Seats {
		err := m.SqlQuery.UpdateShowtimeSeatSeatByIdAndShowtimeId(m.Context,
			sqlc.UpdateShowtimeSeatSeatByIdAndShowtimeIdParams{
				SeatID:     seat.SeatID,
				Status:     sqlc.SeatStatusesReserved,
				BookedBy:   payload.OrderedBy,
				ShowtimeID: payload.ShowtimeId,
			})

		if err != nil {
			log.Printf("an error occur when updating showtime seat %d: %v", seat.SeatID, err)
			return
		} else {
			log.Printf("update showtime seat %d with showtime id %d successfully", seat.SeatID, payload.ShowtimeId)
		}
	}

}
