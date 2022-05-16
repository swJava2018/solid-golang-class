package processors

import (
	"context"
	"event-data-pipeline/pkg/payloads"
)

var _ Processor = new(ProcessorFunc)
var _ ProcessorFunc = NormalizeKafkaPayload

func init() {
	Register("normalize_kafka_payload", NewNormalizeKafkaPayloadProcessor)
}

func NewNormalizeKafkaPayloadProcessor(config jsonObj) Processor {

	return ProcessorFunc(NormalizeKafkaPayload)
}

func NormalizeKafkaPayload(context.Context, payloads.Payload) (payloads.Payload, error) {
	return nil, nil
}
