package processors

import (
	"context"
	"event-data-pipeline/pkg/payloads"
)

func init() {
	// TODO: 2주차 과제입니다.
}

type RabbitMQMetaInjector struct {
	// TODO: 2주차 과제입니다.
}

func NewRabbitMQMetaInjector(config jsonObj) Processor {
	// TODO: 2주차 과제입니다.
	return nil
}

// Process implements Processor
func (*RabbitMQMetaInjector) Process(ctx context.Context, p payloads.Payload) (payloads.Payload, error) {
	// TODO: 2주차 과제입니다.

	return p, nil
}
