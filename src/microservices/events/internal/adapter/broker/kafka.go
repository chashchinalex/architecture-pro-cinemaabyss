package broker

import (
	"context"
	"fmt"

	"events/core/ports"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

type KafkaEventProducer struct {
	producer *kafka.Producer
}

var _ ports.EventProducer = &KafkaEventProducer{}

func NewKafkaEventProducer(brokers string) (*KafkaEventProducer, error) {
	p, err := kafka.NewProducer(&kafka.ConfigMap{
		"bootstrap.servers": brokers,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create Kafka producer: %v", err)
	}
	return &KafkaEventProducer{producer: p}, nil
}

func (k *KafkaEventProducer) SendEvent(ctx context.Context, topic string, _ []byte, value []byte) error {
	return k.producer.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
		Value:          value,
	}, nil)
}

func (k *KafkaEventProducer) Close() {
	k.producer.Close()
}
