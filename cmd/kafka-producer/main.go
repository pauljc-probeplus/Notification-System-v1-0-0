package main

import (
	"context"
	"notification-system/internal/kafka/producer"
	"time"
)

func main() {
	kafkaBrokers := []string{"localhost:9092"}
	kafkaTopic := "test-topic"

	prod := producer.NewKafkaProducer(kafkaBrokers, kafkaTopic)
	defer prod.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err := prod.Produce(ctx, "test-key", []byte("Hello from Notification System!"))
	if err != nil {
		panic(err)
	}
}
