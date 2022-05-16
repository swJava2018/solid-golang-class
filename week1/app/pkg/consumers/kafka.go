package consumers

import (
	"context"
	"encoding/json"
	"event-data-pipeline/pkg/logger"
	"event-data-pipeline/pkg/payloads"
	"event-data-pipeline/pkg/pipelines"
	"time"

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
	stream        chan interface{}
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
		kafkaConsumer: kafka.NewKafkaConsumer(kcCfg.Topic, kcCfg.ConsumerOptions),
	}
	err = client.InitClient()
	if err != nil {
		logger.Panicf(err.Error())
	}
	return client

}

// 컨슈머 생성, 파티션 정보를 가져옴.
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
func (kc *KafkaConsumerClient) Consume(ctx context.Context, stream chan interface{}, errc chan error) error {
	err := kc.kafkaConsumer.Read(ctx, stream)
	if err != nil {
		return err
	}
	return nil
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
// 특정 페이로드 타입으로 변환
func (kc *KafkaConsumerClient) Next(ctx context.Context) bool {

	//스트림으로부터 읽어오기
	for {
		select {
		// 스트림이 있을 때
		case p := <-kc.kafkaConsumer.Stream:
			data, err := json.Marshal(p)
			if err != nil {
				logger.Errorf("failed to marshall the stream data")
				return false
			}
			var kfkPayload payloads.KafkaPayload
			err = json.Unmarshal(data, &kfkPayload)
			if err != nil {
				logger.Errorf(err.Error())
			}
			kc.payload = &kfkPayload
			return true
		// Shutdown
		case <-ctx.Done():
			logger.Debugf("Context cancelled")
		// 스트림이 없을 때 Sleep 후 다시 읽기 시도.
		default:
			logger.Debugf("No Next Stream. Sleeping for 2 seconds")
			time.Sleep(2 * time.Second)
		}
	}

}

// Source 인터페이스 구현
func (kc *KafkaConsumerClient) Payload() payloads.Payload {
	return kc.payload
}

// Source 인터페이스 구현
func (kc *KafkaConsumerClient) Error() error {
	return nil
}
