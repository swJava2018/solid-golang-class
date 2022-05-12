package consumers

import (
	"context"
	"encoding/json"
	"event-data-pipeline/pkg/logger"
	"event-data-pipeline/pkg/payloads"
	"event-data-pipeline/pkg/pipelines"

	"event-data-pipeline/pkg/kafka"
)

// compile type assertion check
var _ pipelines.Source = new(KafkaConsumerClient)

// register kafka consumer client to factory
func init() {

	Register("kafka", NewKafkaConsumerClient)

}

type KafkaClientConfig struct {
	ClientName      string  `json:"client_name,omitempty"`
	Topic           string  `json:"topic,omitempty"`
	ConsumerOptions jsonObj `json:"consumer_options,omitempty"`
}

//implements Consumer interface
type KafkaConsumerClient struct {
	kafkaConsumer *kafka.KafkaConsumer
}

func NewKafkaConsumerClient(config jsonObj) Consumer {

	// Read config into KafkaClientConfig struct
	var kcCfg KafkaClientConfig
	_json, err := json.Marshal(config)
	if err != nil {
		logger.Errorf(err.Error())
	}
	json.Unmarshal(_json, &kcCfg)

	// create a new Consumer concrete type - KafkaConsumerClient
	client := &KafkaConsumerClient{
		// pass ConsumerOptions only
		kafkaConsumer: kafka.NewKafkaConsumer(kcCfg.Topic, kcCfg.ConsumerOptions),
	}

	return client

}

func (kc *KafkaConsumerClient) Consume(ctx context.Context, stream chan interface{}, errc chan error, shutdown chan bool) error {
	var err error

	// Create Kafka Consumer
	err = kc.kafkaConsumer.Create()
	if err != nil {
		return err
	}
	// create admin client from the consumer created
	kc.kafkaConsumer.CreateAdmin()

	// get partitions to read
	err = kc.kafkaConsumer.GetPartitions()
	if err != nil {
		return err
	}
	// Read
	kc.kafkaConsumer.Read(ctx, stream, errc, shutdown)

	return err
}

func (kc *KafkaConsumerClient) Next(context.Context) bool {
	return true
}

// Payload returns the next payload to be processed.
func (kc *KafkaConsumerClient) Payload() payloads.Payload {
	return &payloads.PageviewPayload{}
}

// Error return the last error observed by the source.
func (kc *KafkaConsumerClient) Error() error {
	return nil
}
