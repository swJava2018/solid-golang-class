package payloads

import (
	"sync"
)

var (
	_ Payload = (*PageviewPayload)(nil)

	pageviewPayloadPool = sync.Pool{
		New: func() interface{} { return new(PageviewPayload) },
	}
)

type PageviewPayload struct {
	viewtime int
	userid   string
	pageid   string
}

// Clone implements pipeline.Payload.
func (p *PageviewPayload) Clone() Payload {
	newP := pageviewPayloadPool.Get().(*PageviewPayload)
	newP.viewtime = p.viewtime
	newP.userid = p.userid
	newP.pageid = p.pageid
	return newP
}

// MarkAsProcessed implements pipeline.Payload
func (p *PageviewPayload) MarkAsProcessed() {
	p.viewtime = 0
	p.userid = ""
	p.pageid = ""
	pageviewPayloadPool.Put(p)
}
