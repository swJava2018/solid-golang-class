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

	return &RabbitMQConsumerClient{}
}

// Init implements Consumer
func (rc *RabbitMQConsumerClient) Init() error {
	err := rc.Create()
	if err != nil {
		return err
	}

	return nil
}

// Consume implements Consumer
func (rc *RabbitMQConsumerClient) Consume(ctx context.Context, stream chan interface{}, errc chan error) error {
	// DO SOMETHING
	return nil
}
