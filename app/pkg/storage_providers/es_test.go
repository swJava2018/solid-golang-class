package storage_providers_test

import (
	"event-data-pipeline/pkg/logger"
	"event-data-pipeline/pkg/payloads"
	"event-data-pipeline/pkg/storage_providers"
	"fmt"
	"os"
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

var _ payloads.Payload = new(esPayloadStub)

type esPayloadStub struct {
	dir      string
	filename string
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
	return p.dir, p.filename, []byte(`{}`)
}
