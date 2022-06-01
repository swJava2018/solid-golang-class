package processors

import (
	"context"
	"event-data-pipeline/pkg/logger"
	"event-data-pipeline/pkg/payloads"
)

var _ Processor = new(DefaultProcessor)

type DefaultProcessor struct {
}

// Process implements Processor
func (*DefaultProcessor) Process(ctx context.Context, p payloads.Payload) (payloads.Payload, error) {
	select {
	case <-ctx.Done():
		logger.Errorf(ctx.Err().Error())
		return nil, ctx.Err()
	default:
		return p, nil
	}
}
