package processors

import (
	"context"
	"event-data-pipeline/pkg/payloads"
)

// 컴파일 타임 인터페이스 타입 체크
var _ Processor = new(RabbitMQDefaultProcessor)

func init() {
	Register("rabbitmq_default", NewRabbitMQDefaultProcessor)

}

type RabbitMQDefaultProcessor struct {
	Validator
	RabbitMQMetaInjector
}

func NewRabbitMQDefaultProcessor(config jsonObj) Processor {
	p := &RabbitMQDefaultProcessor{
		Validator{},
		RabbitMQMetaInjector{},
	}
	return p
}

func (k *RabbitMQDefaultProcessor) Process(ctx context.Context, p payloads.Payload) (payloads.Payload, error) {
	err := k.Validate(ctx, p)
	if err != nil {
		return nil, err
	}
	p, err = k.RabbitMQMetaInjector.Process(ctx, p)
	if err != nil {
		return nil, err
	}
	return p, nil
}
