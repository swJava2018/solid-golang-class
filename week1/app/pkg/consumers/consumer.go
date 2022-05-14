package consumers

import (
	"context"
	"errors"
	"event-data-pipeline/pkg/logger"
	"fmt"
	"strings"
)

type (
	jsonObj = map[string]interface{}
	jsonArr = []interface{}
)

// consumer interface
type Consumer interface {
	//read data
	Consume(ctx context.Context, stream chan interface{}, errc chan error) error
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

// CreateConsumer is a factory method that will create the named consumer
func CreateConsumer(name string, config jsonObj) (Consumer, error) {

	factory, ok := consumerFactories[name]
	if !ok {
		// Factory has not been registered.
		// Make a list of all available datastore factories for logging.
		availableConsumers := make([]string, 0)
		for k := range consumerFactories {
			availableConsumers = append(availableConsumers, k)
		}
		return nil, errors.New(fmt.Sprintf("Invalid Consumer name. Must be one of: %s", strings.Join(availableConsumers, ", ")))
	}

	// Run the factory with the configuration.
	return factory(config), nil
}
