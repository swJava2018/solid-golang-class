package processors

import (
	"context"
	"event-data-pipeline/pkg/payloads"
)

var _ Processor = new(KafkaDefaultProcessor)

type KafkaDefaultProcessor struct {
	Validator
	KafkaMetaInjector
}

func init() {
	Register("kafka_default", NewKafkaDefaultProcessor)
}

func NewKafkaDefaultProcessor(config jsonObj) Processor {
	return &KafkaDefaultProcessor{
		Validator{},
		KafkaMetaInjector{},
	}
}

func (k *KafkaDefaultProcessor) Process(ctx context.Context, p payloads.Payload) (payloads.Payload, error) {
	//Validator struct embedding
	err := k.Validate(ctx, p)
	if err != nil {
		return nil, err
	}
	//KafkaMetaInject struct embedding
	p, err = k.KafkaMetaInjector.Process(ctx, p)
	if err != nil {
		return nil, err
	}

	return p, nil
}
