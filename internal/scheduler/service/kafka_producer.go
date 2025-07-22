package service

import (
	"context"

	"github.com/segmentio/kafka-go"
)


func publishToKafka(topic, payload string) error {
	writer := &kafka.Writer{
		Addr:     kafka.TCP("localhost:9092"),
		Topic:    topic,
		Balancer: &kafka.LeastBytes{},
	}
	defer writer.Close()

	msg := kafka.Message{
		Value: []byte(payload),
	}
	return writer.WriteMessages(context.Background(), msg)
}
