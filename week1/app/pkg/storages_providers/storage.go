package storages_providers

import (
	"event-data-pipeline/pkg/logger"
	"fmt"
	"strings"
)

type (
	jsonObj = map[string]interface{}
	jsonArr = []interface{}
)
type Storage interface {
	Write(key string, path string, data []byte) (int, error)
}

type StorageFactory func(config jsonObj) Storage

var storageFactories = make(map[string]StorageFactory)

// Each storage implementation must Register itself
func Register(name string, factory StorageFactory) {
	logger.Debugf("Registering storage factory for %s", name)
	if factory == nil {
		logger.Panicf("Storage factory %s does not exist.", name)
	}
	_, registered := storageFactories[name]
	if registered {
		logger.Errorf("Storage factory %s already registered. Ignoring.", name)
	}
	storageFactories[name] = factory
}

// CreateStorage is a factory method that will create the named storage
func CreateStorage(name string, config jsonObj) (Storage, error) {

	factory, ok := storageFactories[name]
	if !ok {
		// Factory has not been registered.
		// Make a list of all available datastore factories for logging.
		availableStorages := make([]string, 0)
		for k := range storageFactories {
			availableStorages = append(availableStorages, k)
		}
		return nil, fmt.Errorf("invalid Storage name. Must be one of: %s", strings.Join(availableStorages, ", "))
	}

	// Run the factory with the configuration.
	return factory(config), nil
}
