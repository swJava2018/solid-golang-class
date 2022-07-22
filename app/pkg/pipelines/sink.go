package pipelines

import (
	"context"
	"event-data-pipeline/pkg/payloads"
)

type Sink interface {
	Drain(context.Context, payloads.Payload) error
}
