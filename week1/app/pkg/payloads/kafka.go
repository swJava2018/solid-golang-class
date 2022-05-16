package payloads

import (
	"sync"
	"time"
)

var (
	// 컴파일 타임 타입 변경 체크
	_ Payload = (*KafkaPayload)(nil)

	kafkaPayloadPool = sync.Pool{
		New: func() interface{} { return new(KafkaPayload) },
	}
)

type KafkaPayload struct {
	Topic     string                 `json:"topic,omitempty"`
	Partition float64                `json:"partition,omitempty"`
	Offset    float64                `json:"offset,omitempty"`
	Key       string                 `json:"key,omitempty"`
	Value     map[string]interface{} `json:"value,omitempty"`
	Timestamp time.Time              `json:"timestamp,omitempty"`
}

// Clone implements pipeline.Payload.
func (p *KafkaPayload) Clone() Payload {
	newP := kafkaPayloadPool.Get().(*KafkaPayload)

	return newP
}

// MarkAsProcessed implements pipeline.Payload
func (p *KafkaPayload) MarkAsProcessed() {

	usersPayloadPool.Put(p)
}
