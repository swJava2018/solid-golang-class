package processors

import (
	"context"
	"event-data-pipeline/pkg/payloads"
)

func init() {
	Register("noop", NewNoopProcessor)
}

func NewNoopProcessor(config jsonObj) Processor {
	return &Noop{}
}

// Processor return Payload as it is
type Noop struct {
}

func (b *Noop) Process(ctx context.Context, p payloads.Payload) (payloads.Payload, error) {

	return p, nil
}
