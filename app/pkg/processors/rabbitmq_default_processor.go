package processors

import (
	"context"
	"event-data-pipeline/pkg/payloads"
)

// 컴파일 타임 인터페이스 타입 체크
var _ Processor = new(RabbitMQDefaultProcessor)

func init() {
	// TODO: 2주차 과제입니다.
}

type RabbitMQDefaultProcessor struct {
	// TODO: 2주차 과제입니다.
}

func init() {
	// TODO: 2주차 과제입니다.
}

func NewRabbitMQDefaultProcessor(config jsonObj) Processor {
	// TODO: 2주차 과제입니다.
	return nil
}

func (k *RabbitMQDefaultProcessor) Process(ctx context.Context, p payloads.Payload) (payloads.Payload, error) {
	// TODO: 2주차 과제입니다.

	return p, nil
}
