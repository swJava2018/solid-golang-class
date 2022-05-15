package payloads

import (
	"sync"
)

var (
	_ Payload = (*UsersPayload)(nil)

	usersPayloadPool = sync.Pool{
		New: func() interface{} { return new(UsersPayload) },
	}
)

type UsersPayload struct {
	Userid    string `json:"userid,omitempty"`
	Offset    int    `json:"offset,omitempty"`
	Partition int    `json:"partition,omitempty"`
}

// Clone implements pipeline.Payload.
func (p *UsersPayload) Clone() Payload {
	newP := usersPayloadPool.Get().(*UsersPayload)

	return newP
}

// MarkAsProcessed implements pipeline.Payload
func (p *UsersPayload) MarkAsProcessed() {
	p.Userid = ""
	usersPayloadPool.Put(p)
}
