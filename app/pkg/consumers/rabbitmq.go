package consumers

import (
	"context"
	"event-data-pipeline/pkg/rabbitmq"
	"event-data-pipeline/pkg/sources"
)

// compile type assertion check
var _ Consumer = new(RabbitMQConsumerClient)
var _ ConsumerFactory = NewRabbitMQConsumerClient

func init() {

	Register("rabbitmq", NewKafkaConsumerClient)

}

type RabbitMQConsumerClient struct {
	rabbitmq.Consumer
	sources.Source
}

func NewRabbitMQConsumerClient(config jsonObj) Consumer {
	//TODO: 1주차 과제입니다.
	return client
}

// Init implements Consumer
func (rc *RabbitMQConsumerClient) Init() error {
	//TODO: 1주차 과제입니다.

	return nil
}

// Consume implements Consumer
func (rc *RabbitMQConsumerClient) Consume(ctx context.Context) error {
	//TODO: 1주차 과제입니다.
	return nil
}
