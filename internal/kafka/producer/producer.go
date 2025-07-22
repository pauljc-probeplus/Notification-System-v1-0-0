package producer

import (
	"context"
	"log"
	"time"

	"github.com/segmentio/kafka-go"
)

type KafkaProducer struct {
	writer *kafka.Writer
}

func NewKafkaProducer(brokers []string, topic string) *KafkaProducer {
	return &KafkaProducer{
		writer: &kafka.Writer{
			Addr:     kafka.TCP(brokers...),
			Topic:    topic,
			Balancer: &kafka.LeastBytes{},
		},
	}
}

func (kp *KafkaProducer) Produce(ctx context.Context, key string, value []byte) error {
	msg := kafka.Message{
		Key:   []byte(key),
		Value: value,
		Time:  time.Now(),
	}
	err := kp.writer.WriteMessages(ctx, msg)
	if err != nil {
		log.Printf("failed to write message to kafka: %v", err)
		return err
	}
	log.Printf("message sent to kafka: key=%s", key)
	return nil
}

func (kp *KafkaProducer) Close() error {
	return kp.writer.Close()
}
