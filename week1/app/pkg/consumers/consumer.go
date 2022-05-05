package consumers

import (
	"context"
	"event-data-pipeline/pkg/logger"
)

type (
	jsonObj = map[string]interface{}
	jsonArr = []interface{}
)

// consumer interface
type Consumer interface {
	//create Consumer instance
	Create() error

	//read data
	Read(ctx context.Context, stream chan interface{}, errc chan error, shutdown chan bool) error

	//delete Consumer instance
	Delete() error
}

// a factory func type to instantiate a concrete Consumer type
type ConsumerFactory func(config jsonObj) Consumer

// factories to register Consumers in init() function in each consumer
var consumerFactories = make(map[string]ConsumerFactory)

// Each consumer implementation must Register itself
func Register(name string, factory ConsumerFactory) {
	logger.Debugf("Registering consumer factory for %s", name)
	if factory == nil {
		logger.Panicf("Consumer factory %s does not exist.", name)
	}
	_, registered := consumerFactories[name]
	if registered {
		logger.Errorf("Consumer factory %s already registered. Ignoring.", name)
	}
	consumerFactories[name] = factory
}
