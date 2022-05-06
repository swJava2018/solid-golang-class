package consumers

import (
	"context"
	"encoding/json"
	"event-data-pipeline/pkg/logger"

	"event-data-pipeline/pkg/kafka"
)

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
	// create admin client
	// get partitions
	// loop through partitions
	// run go routing per partition
	/// send event data through input channel
	var err error
	// ev := kc.kafkaConsumer.Poll(100)
	kc.kafkaConsumer.Read()

	// pass the event to output channel
	// which is input channel of a processor that implements StageRunner
	return err
}
