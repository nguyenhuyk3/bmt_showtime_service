package readers

import (
	"bmt_showtime_service/dto/messages"
	"bmt_showtime_service/global"
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
		var message messages.NewFilmCreationTopic
		if err := json.Unmarshal(value, &message); err != nil {
			log.Printf("failed to unmarshal image message: %v\n", err)
			return
		}

		m.handleNewFilmCreationTopic(message)
	default:
		log.Printf("unknown topic received: %s\n", topic)
	}
}

func (m *MessageBrokerReader) handleNewFilmCreationTopic(message messages.NewFilmCreationTopic) {
	err := m.SqlQuery.CreateNewFilmId(m.Context, message.FilmId)
	if err != nil {
		log.Printf("an error occur when creating new film id: %v", err)
	} else {
		log.Println("create new film id successfully")
	}
}
