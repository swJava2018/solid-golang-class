package consumers

import "context"

func init() {

	Register("rabbitmq", NewKafkaConsumerClient)

}

type RabbitMQConsumerClient struct {
}

func NewRabbitMQConsumerClient() Consumer {

	return &RabbitMQConsumerClient{}
}

// Consume implements Consumer
func (*RabbitMQConsumerClient) Consume(ctx context.Context, stream chan interface{}, errc chan error, shutdown chan bool) error {
	// DO SOMETHING
	return nil
}
