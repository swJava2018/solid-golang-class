package kafka_test

import (
	"context"
	"event-data-pipeline/pkg/kafka"
	"testing"
)

func TestInitKafkaConsumerClient(t *testing.T) {
	client := kafka.NewKafkaConsumerClient(nil)
	client.Read(context.TODO(), nil, nil, nil)
}
