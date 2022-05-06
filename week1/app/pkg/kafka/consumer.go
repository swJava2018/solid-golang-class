package kafka

import (
	"context"
	"encoding/json"
	"event-data-pipeline/pkg/logger"
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

func (kc *KafkaConsumer) Read(ctx context.Context, stream chan interface{}, errc chan error, shutdown chan bool) ([]byte, error) {

	// Create an consumer instance initially
	kc.Create()

	// create admin client from the consumer created
	ac := NewAdminClient(kc.topic, kc.kafkaConsumer)

	// get partitions from admin client
	prttnsResps, err := ac.GetPartitions()

	if err != nil {
		return nil, err
	}
	// loop through partitions
	for _, p := range prttnsResps.Partitions {

		// Copy outer KafkaConsumer
		c := kc.Copy()
		// Instantiate inner kafka.Consumer
		c.Create()
		// Assign partition to the consumer
		c.AssignPartition(int(p))

		// spin up go routine per partition

	}
	// get partitions
	// loop through partitions
	// run go routing per partition
	/// send event data through input channel
	var message []byte
	ev := kc.kafkaConsumer.Poll(100)
	switch e := ev.(type) {
	case *kafka.Message:
		message, err = json.Marshal(e)
		logger.Debugf("Message on p[%v]: %v", e.TopicPartition.Partition, string(message))
	case kafka.Error:
		err = e
		logger.Errorf("Error: %v: %v", e.Code(), e)
	case kafka.PartitionEOF:
		logger.Debugf("[PartitionEOF][Consumer: %s][Topic: %v][Partition: %v][Offset: %d][Message: %v]", kc.kafkaConsumer.String(), *e.Topic, e.Partition, e.Offset, fmt.Sprintf("\"%s\"", e.Error.Error()))
	}
	return message, err
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
func (kc *KafkaConsumer) Delete() error {

	return nil

}
