package event_data

import (
	"event-data-pipeline/pkg/consumers"
	"testing"
)

func TestCreateConsumer(t *testing.T) {
	consumer, err := consumers.CreateConsumer("kafka", nil)
	if err != nil || consumer == nil {
		t.Error(err)
	}

}
