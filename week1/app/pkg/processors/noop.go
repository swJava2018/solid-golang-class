package processors

import (
	"context"
	"event-data-pipeline/pkg/payloads"
)

func init() {
	Register("noop", NewNoopProcessor)
}

type NoopProcessor struct {
	Validator
}

func NewNoopProcessor(config jsonObj) Processor {
	validator := NewValidator()
	return &NoopProcessor{*validator}
}

// Processor return Payload as it is
func (n *NoopProcessor) Process(ctx context.Context, p payloads.Payload) (payloads.Payload, error) {
	n.Validate(ctx, p)
	return p, nil
}
