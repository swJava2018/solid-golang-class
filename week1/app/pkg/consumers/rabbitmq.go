package consumers

import "context"

func init() {

	Register("rabbitmq", NewKafkaConsumerClient)

}

type RabbitMQConsumerClient struct {
}

// Init implements Consumer
func (*RabbitMQConsumerClient) Init() error {
	panic("unimplemented")
}

func NewRabbitMQConsumerClient() Consumer {

	return &RabbitMQConsumerClient{}
}

// Consume implements Consumer
func (*RabbitMQConsumerClient) Consume(ctx context.Context, stream chan interface{}, errc chan error) error {
	// DO SOMETHING
	return nil
}
