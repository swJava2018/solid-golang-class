package processors

import (
	"context"
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

func (b *Noop) Process(ctx context.Context, p Payload) (Payload, error) {

	return p, nil
}
