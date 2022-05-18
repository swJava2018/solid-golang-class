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
var _ Consumer = new(KafkaConsumerClient)
var _ ConsumerFactory = NewKafkaConsumerClient

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
	kafka.Consumer
	payload payloads.Payload
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
	client := &KafkaConsumerClient{}
	client.Consumer = kafka.NewKafkaConsumer(kcCfg.Topic, kcCfg.ConsumerOptions)

	return client

}

// Init implements Consumer
func (kc *KafkaConsumerClient) Init() error {
	var err error

	err = kc.CreateConsumer()
	if err != nil {
		return err
	}
	err = kc.CreateAdminConsumer()
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
	err := kc.Read(ctx, stream, errc)
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
		case p := <-kc.Stream():
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
			// logger.Debugf("No Next Stream. Sleeping for 2 seconds")
			// time.Sleep(2 * time.Second)
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
