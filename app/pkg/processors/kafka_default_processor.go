package processors

import (
	"context"
	"event-data-pipeline/pkg/payloads"
)

var _ Processor = new(KafkaDefaultProcessor)

type KafkaDefaultProcessor struct {
	KafkaMetaInjector
}

func init() {
	Register("kafka_default_processor", NewKafkaDefaultProcessor)
}

func NewKafkaDefaultProcessor(config jsonObj) Processor {
	return &KafkaDefaultProcessor{
		KafkaMetaInjector{},
	}
}

func (k *KafkaDefaultProcessor) Process(ctx context.Context, p payloads.Payload) (payloads.Payload, error) {

	//KafkaMetaInject struct embedding
	p, err := k.KafkaMetaInjector.Process(ctx, p)
	if err != nil {
		return nil, err
	}

	return p, nil
}
