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
	payload       payloads.Payload
}

func NewKafkaConsumerClient(config jsonObj) Consumer {

	// Read config into KafkaClientConfig struct
	var kcCfg KafkaClientConfig
	_json, err := json.Marshal(config)
	if err != nil {
		logger.Panicf(err.Error())
	}
	json.Unmarshal(_json, &kcCfg)

	// create a new Consumer concrete type - KafkaConsumerClient
	client := &KafkaConsumerClient{
		// pass ConsumerOptions only
		kafkaConsumer: kafka.NewKafkaConsumer(kcCfg.Topic, kcCfg.ConsumerOptions),
	}
	err = client.InitClient()
	if err != nil {
		logger.Panicf(err.Error())
	}
	return client

}

func (kc *KafkaConsumerClient) InitClient() error {
	var err error

	err = kc.Create()
	if err != nil {
		return err
	}
	err = kc.CreateAdmin()
	if err != nil {
		return err
	}
	err = kc.GetPartitions()
	if err != nil {
		return err
	}
	return nil
}

// Consumer 인터페이스 구현
func (kc *KafkaConsumerClient) Consume(ctx context.Context, stream chan interface{}, errc chan error, shutdown chan bool) error {

	// Read
	kc.kafkaConsumer.Read(ctx, stream, errc, shutdown)

	return err
}

// kafkaConsumer 인스턴스 생성
func (kc *KafkaConsumerClient) Create() error {
	var err error

	// Create Kafka Consumer
	err = kc.kafkaConsumer.Create()
	if err != nil {
		return err
	}
	return err
}

// Consumer 인터페이스 구현
func (kc *KafkaConsumerClient) CreateAdmin() error {
	// create admin client from the consumer created
	err := kc.kafkaConsumer.CreateAdmin()
	if err != nil {
		return err
	}

	return nil
}

// Consumer 인터페이스 구현
func (kc *KafkaConsumerClient) GetPartitions() error {

	// get partitions to read
	err := kc.kafkaConsumer.GetPartitions()
	if err != nil {
		return err
	}

	return nil
}

// Source 인터페이스 구현
// 다음 Payload 가 있는지 확인
func (kc *KafkaConsumerClient) Next(context.Context) bool {
	// POll 을 통해 다음 레코드가 있으면 인스턴스에 저장.
	return true
}

// Source 인터페이스 구현
// TODO: 설정 값에 따라 사용할 페이로드 타입을 결정 일단은 UsersPayload로 진행
func (kc *KafkaConsumerClient) Payload() payloads.Payload {
	return &payloads.UsersPayload{}
}

// Source 인터페이스 구현
func (kc *KafkaConsumerClient) Error() error {
	return nil
}
