package consumers

import (
	"context"
	"encoding/json"
	"event-data-pipeline/pkg/logger"
	"fmt"

	"github.com/confluentinc/confluent-kafka-go/kafka"
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

	//configuration to create confluent-kafka-go consumer
	configMap *kafka.ConfigMap

	//confluent kafka go consumer
	kafkaConsumer *kafka.Consumer
}

func NewKafkaConsumerClient(config jsonObj) Consumer {

	// Read config into KafkaClientConfig struct
	var kcCfg KafkaClientConfig
	_json, err := json.Marshal(config)
	if err != nil {
		logger.Errorf(err.Error())
	}
	json.Unmarshal(_json, &kcCfg)

	// load Consumer Options to kafka.ConfigMap
	raw, _ := json.Marshal(kcCfg.ConsumerOptions)
	var kcm kafka.ConfigMap
	json.Unmarshal(raw, &kcm)

	client := &KafkaConsumerClient{
		configMap: &kcm,
	}
	err = client.Create()
	if err != nil {
		logger.Errorf(err.Error())
	}

	return client

}

//create KafkaClient instance
func (kc *KafkaConsumerClient) Create() error {

	if kc == nil {
		kc = &KafkaConsumerClient{}
	}
	var err error
	// create kafka consumer instance
	kc.kafkaConsumer, err = kafka.NewConsumer(kc.configMap)

	if err != nil {
		return err
	}
	logger.Debugf("Check in consumer creation: %s", kc.kafkaConsumer)
	return nil

}

func (kc *KafkaConsumerClient) Read(ctx context.Context, stream chan interface{}, errc chan error, shutdown chan bool) error {
	// create admin client
	// get partitions
	// loop through partitions
	// run go routing per partition
	/// send event data through input channel
	var err error
	for {

		ev := kc.kafkaConsumer.Poll(100)
		switch e := ev.(type) {
		case *kafka.Message:
			raw, _ := json.Marshal(e)
			logger.Debugf("Message on p[%v]: %v", e.TopicPartition.Partition, string(raw))
			return nil
		case kafka.Error:
			logger.Errorf("Error: %v: %v", e.Code(), e)
			err = e
			return err
		case kafka.PartitionEOF:
			logger.Debugf("[PartitionEOF][Consumer: %s][Topic: %v][Partition: %v][Offset: %d][Message: %v]", kc.kafkaConsumer.String(), *e.Topic, e.Partition, e.Offset, fmt.Sprintf("\"%s\"", e.Error.Error()))
		default:
			if e == nil {
				continue
			}
		}
	}

}

//delete kafkaClient instance
func (kc *KafkaConsumerClient) Delete() error {

	return nil

}
