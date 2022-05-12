package payloads

import (
	"sync"
	"time"
)

var (
	_ Payload = (*usersPayload)(nil)

	usersPayloadPool = sync.Pool{
		New: func() interface{} { return new(usersPayload) },
	}
)

type usersPayload struct {
	userid       string
	regionid     string
	gender       string
	registertime time.Time
}

// Clone implements pipeline.Payload.
func (p *usersPayload) Clone() Payload {
	newP := usersPayloadPool.Get().(*usersPayload)

	return newP
}

// MarkAsProcessed implements pipeline.Payload
func (p *usersPayload) MarkAsProcessed() {
	p.userid = ""
	p.gender = ""
	p.regionid = ""
	p.registertime = time.Time{}
	usersPayloadPool.Put(p)
}
