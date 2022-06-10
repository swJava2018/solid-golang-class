package storage_providers_test

import (
	"event-data-pipeline/pkg/logger"
	"event-data-pipeline/pkg/payloads"
	"event-data-pipeline/pkg/storage_providers"
	"fmt"
	"os"
	"sync"
	"testing"

	gc "gopkg.in/check.v1"
)

func Test(t *testing.T) { gc.TestingT(t) }

type ESSuite struct{}

var _ = gc.Suite(&ESSuite{})

func (e *ESSuite) SetUpSuite(c *gc.C) {
	os.Args = nil
	os.Setenv("EDP_ENABLE_DEBUG_LOGGING", "true")
	logger.Setup()

}
func (e *ESSuite) TearDownSuite(c *gc.C) {

}

func (f *ESSuite) TestWrite(c *gc.C) {

	cfgObj := make(jsonObj)
	addresses := &[]interface{}{"http://localhost:9200"}
	cfgObj["addresses"] = addresses

	// ES storage provider 인스턴스 생성
	es, err := storage_providers.CreateStorageProvider("elasticsearch", cfgObj)

	// 에러 체크
	c.Assert(err, gc.IsNil)

	// 페이로드 stub 생성
	payload := &esPayloadStub{"event-data-test", fmt.Sprintf("es.write.test.%d", 0)}

	written, err := es.Write(payload)
	c.Assert(written, gc.NotNil)
	c.Assert(err, gc.IsNil)

}

func (f *ESSuite) TestConcurrentWrite(c *gc.C) {

	cfgObj := make(jsonObj)
	addresses := &[]interface{}{"http://localhost:9200"}
	cfgObj["addresses"] = addresses

	// ES storage provider 인스턴스 생성
	es, err := storage_providers.CreateStorageProvider("elasticsearch", cfgObj)

	// 에러 체크
	c.Assert(err, gc.IsNil)

	// 동시 트린잭션 개수
	requests := 100

	// WaitGroup 생성
	var wg sync.WaitGroup

	// 고루틴 생성
	for i := 0; i < requests; i++ {
		wg.Add(1)
		go func(idx int, wg *sync.WaitGroup) {
			// 페이로드 stub 생성
			payload := &esPayloadStub{"event-data-test-concurrent", fmt.Sprintf("es.write.test.%d", idx)}
			logger.Infof("payload:%v", payload)
			written, err := es.Write(payload)
			c.Assert(err, gc.IsNil)
			if written <= 0 {
				c.Errorf("failed to write:%s", payload)
			}
			wg.Done()
		}(i, &wg)
	}
	wg.Wait()

}

var _ payloads.Payload = new(esPayloadStub)

type esPayloadStub struct {
	index string
	docID string
}

// Clone implements payloads.Payload
func (*esPayloadStub) Clone() payloads.Payload {
	ps := &esPayloadStub{}
	return ps
}

// MarkAsProcessed implements payloads.Payload
func (*esPayloadStub) MarkAsProcessed() {

}

// Out implements payloads.Payload
func (p *esPayloadStub) Out() (string, string, []byte) {
	return p.index, p.docID, []byte(`{}`)
}
