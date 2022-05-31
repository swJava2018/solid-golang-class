package casters

import (
	"encoding/json"
	"event-data-pipeline/pkg/logger"

	"github.com/streadway/amqp"
)

func init() {
	Register("purchases_caster", NewPurchasesCaster)
}

func NewPurchasesCaster() Caster {
	return CasterFunc(CastPurchases)
}

func CastPurchases(meta jsonObj, message amqp.Delivery) (jsonObj, error) {
	var record = make(jsonObj)
	for k, v := range meta {
		record[k] = v
	}
	var valObj jsonObj
	err := json.Unmarshal(message.Body, &valObj)
	if err != nil {
		logger.Errorf("error in casting value object: %v", err)
		return nil, err
	}
	record["value"] = valObj
	record["timestamp"] = message.Timestamp

	return record, nil
}
