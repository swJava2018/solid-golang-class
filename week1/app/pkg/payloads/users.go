package payloads

import (
	"sync"
	"time"
)

var (
	_ Payload = (*UsersPayload)(nil)

	usersPayloadPool = sync.Pool{
		New: func() interface{} { return new(UsersPayload) },
	}
)

type UsersPayload struct {
	userid       string
	regionid     string
	gender       string
	registertime time.Time
}

// Clone implements pipeline.Payload.
func (p *UsersPayload) Clone() Payload {
	newP := usersPayloadPool.Get().(*UsersPayload)

	return newP
}

// MarkAsProcessed implements pipeline.Payload
func (p *UsersPayload) MarkAsProcessed() {
	p.userid = ""
	p.gender = ""
	p.regionid = ""
	p.registertime = time.Time{}
	usersPayloadPool.Put(p)
}
