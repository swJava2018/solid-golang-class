package payloads

import (
	"sync"
)

var (
	_ Payload = (*pageviewPayload)(nil)

	pageviewPayloadPool = sync.Pool{
		New: func() interface{} { return new(pageviewPayload) },
	}
)

type pageviewPayload struct {
	viewtime int
	userid   string
	pageid   string
}

// Clone implements pipeline.Payload.
func (p *pageviewPayload) Clone() Payload {
	newP := pageviewPayloadPool.Get().(*pageviewPayload)
	newP.viewtime = p.viewtime
	newP.userid = p.userid
	newP.pageid = p.pageid
	return newP
}

// MarkAsProcessed implements pipeline.Payload
func (p *pageviewPayload) MarkAsProcessed() {
	p.viewtime = 0
	p.userid = ""
	p.pageid = ""
	pageviewPayloadPool.Put(p)
}
