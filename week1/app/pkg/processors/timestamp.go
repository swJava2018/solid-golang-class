package processors

import (
	"context"
	"event-data-pipeline/pkg/payloads"
)

type TimeStampProcessor struct {
}

func NewProcessTimeProcessor() Processor {

	return &TimeStampProcessor{}
}

func (pt *TimeStampProcessor) Process(ctx context.Context, p payloads.Payload) (payloads.Payload, error) {
	// timestamp
	return p, nil
}
