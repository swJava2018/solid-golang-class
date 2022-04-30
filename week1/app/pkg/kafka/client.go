package kafka

import (
	"context"
	"encoding/json"
	"event-data-pipeline/pkg/logger"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

type (
	jsonObj = map[string]interface{}
)

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

func NewKafkaConsumerClient(config jsonObj) *KafkaConsumerClient {

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
		logger.Fatalf(err.Error())
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
	return nil

}

func (kc *KafkaConsumerClient) Read(ctx context.Context, stream chan interface{}, errc chan error, shutdown chan bool) error {

	return nil
}

//delete kafkaClient instance
func (kc *KafkaConsumerClient) Delete() error {

	return nil

}
