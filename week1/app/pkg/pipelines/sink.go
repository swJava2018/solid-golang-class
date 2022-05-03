package pipelines

import "context"

// Sink is implemented by types that can operate as the tail of a pipeline.
type Sink interface {
	// Consume processes a Payload instance that has been emitted out of
	// a Pipeline instance.
	Consume(context.Context, Payload) error
}
