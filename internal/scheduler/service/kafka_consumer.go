package service

import (
	"context"
	"encoding/json"
	"log"
	"fmt"
	"github.com/segmentio/kafka-go"
	"notification-system/internal/scheduler/model"
)

func ConsumeFromKafka(broker, topic string, handleFunc func(model.KafkaPayload)) {
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{broker},
		Topic:   topic,
		GroupID: "notification-consumer-group",
	})

	defer r.Close()

	for {
		m, err := r.ReadMessage(context.Background())
		if err != nil {
			log.Println("Error reading message from Kafka:", err)
			continue
		}

		var payload model.KafkaPayload
		err = json.Unmarshal(m.Value, &payload)
		if err != nil {
			log.Println("Error unmarshalling Kafka payload:", err)
			continue
		}

		handleFunc(payload)
	}
}

func HandleKafkaNotification(payload model.KafkaPayload) {
	fmt.Printf("Received Kafka Notification: %+v\n", payload)
}