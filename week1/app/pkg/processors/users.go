package processors

import (
	"context"
	"event-data-pipeline/pkg/payloads"
)

type UsersProcessor struct {
	TimeStampProcessor
}

func NewUsersProcessor() Processor {
	tsp := NewTimeStampProcessor().(*TimeStampProcessor)
	return &UsersProcessor{
		*tsp,
	}
}

func (b *UsersProcessor) Process(ctx context.Context, p payloads.Payload) (payloads.Payload, error) {
	b.TimeStampProcessor.Process(ctx, p)
	return p, nil
}
