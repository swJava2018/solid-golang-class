package processors

import (
	"context"
	"event-data-pipeline/pkg/payloads"
)

type PageViewProcessor struct {
	TimeStampProcessor
}

func NewPageViewProcessor() Processor {
	tsp := NewTimeStampProcessor().(*TimeStampProcessor)
	return &PageViewProcessor{
		*tsp,
	}
}

func (pv *PageViewProcessor) Process(ctx context.Context, p payloads.Payload) (payloads.Payload, error) {
	//
	pv.TimeStampProcessor.Process(ctx, p)
	return p, nil
}
