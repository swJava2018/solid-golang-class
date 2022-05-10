package kafka

import (
	"context"
	"encoding/json"
	"event-data-pipeline/pkg/logger"
	"event-data-pipeline/pkg/pipelines"
	"fmt"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

type Consumer interface {
	Create() error
	AssignPartition(partition int) error
	Read(ctx context.Context, stream chan interface{}, errc chan error, shutdown chan bool) ([]byte, error)
	CommitOffset(partition int, offset int) error
	GetOffsetRange(partition int) (int, int, error)
	GetCommittedOffset(partition int) (int, error)
}

type KafkaConsumer struct {
	// topic to consume from
	topic string

	//configuration to create confluent-kafka-go consumer
	configMap *kafka.ConfigMap

	//confluent kafka go consumer
	kafkaConsumer *kafka.Consumer

	//adminClient to get partitions
	adminClient *AdminClient

	//partitions response
	partitions *PartitionsResponse
}

func NewKafkaConsumer(topic string, config jsonObj) *KafkaConsumer {

	// load Consumer Options to kafka.ConfigMap
	raw, _ := json.Marshal(config)
	var kcm kafka.ConfigMap
	json.Unmarshal(raw, &kcm)

	// create a new KafkaConsumer with configMap fed in
	kafkaConsumer := &KafkaConsumer{
		topic:     topic,
		configMap: &kcm,
	}
	return kafkaConsumer

}

//create KafkaClient instance
func (kc *KafkaConsumer) Create() error {

	if kc == nil {
		kc = &KafkaConsumer{configMap: kc.configMap}
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

func (kc *KafkaConsumer) CreateAdmin() {
	kc.adminClient = NewAdminClient(kc.topic, kc.kafkaConsumer)
}

func (kc *KafkaConsumer) GetPartitions() error {
	var err error
	// get partitions from admin client
	kc.partitions, err = kc.adminClient.GetPartitions()
	if err != nil {
		return err
	}
	return nil
}

func (kc *KafkaConsumer) Read(ctx context.Context, stream chan interface{}, errc chan error, shutdown chan bool) error {

	// loop through partitions
	for _, p := range kc.partitions.Partitions {

		// Copy outer KafkaConsumer
		ckc := kc.Copy()
		// Instantiate inner kafka.Consumer
		ckc.Create()
		// Assign partition to the consumer
		ckc.AssignPartition(int(p))

		// spin up go routine per partition
		ckc.Poll(stream, errc)
	}
	return nil
}

//Copy KafkaConsumer instance
func (kc *KafkaConsumer) Copy() *KafkaConsumer {
	return &KafkaConsumer{
		topic:     kc.topic,
		configMap: kc.configMap,
	}
}

func (kc *KafkaConsumer) AssignPartition(partition int) error {

	var partitions []kafka.TopicPartition

	tp := NewTopicPartition(kc.topic, partition)
	partitions = append(partitions, *tp)

	err := kc.kafkaConsumer.Assign(partitions)
	if err != nil {
		return err
	}
	return err
}

//delete kafkaClient instance
func (kc *KafkaConsumer) Poll(stream chan<- interface{}, errc chan<- error) {

	//Next
	//Payload
	//Error
	for {
		var message []byte
		var err error
		ev := kc.kafkaConsumer.Poll(100)
		switch e := ev.(type) {
		case *kafka.Message:
			message, err = json.Marshal(e)
			if err != nil {
				logger.Errorf(e.String())
				errc <- err
				continue
			}
			// write message to the stream
			stream <- message
			logger.Debugf("Message on p[%v]: %v", e.TopicPartition.Partition, string(message))
		case kafka.Error:
			//TODO: check switch type
			errc <- ev.(error)
			logger.Errorf("Error: %v: %v", e.Code(), e)
		case kafka.PartitionEOF:
			logger.Debugf("[PartitionEOF][Consumer: %s][Topic: %v][Partition: %v][Offset: %d][Message: %v]", kc.kafkaConsumer.String(), *e.Topic, e.Partition, e.Offset, fmt.Sprintf("\"%s\"", e.Error.Error()))
		}
	}
}

// Check if Next Payload exists
func (kc *KafkaConsumer) Next(context.Context) bool {

	return false

}

// Get Payload
func (kc *KafkaConsumer) Payload() pipelines.Payload {

	return false

}

//delete kafkaClient instance
func (kc *KafkaConsumer) Delete() error {

	return nil

}
