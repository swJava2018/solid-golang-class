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

// Hook up gocheck into the "go test" runner.
func FSTest(t *testing.T) { gc.TestingT(t) }

type FilesystemSuite struct{}

var _ = gc.Suite(&FilesystemSuite{})

func (f *FilesystemSuite) SetUpSuite(c *gc.C) {
	os.Args = nil
	os.Setenv("EDP_ENABLE_DEBUG_LOGGING", "true")
	logger.Setup()

	fmt.Println("Setting up suite: clearing fs directory...")
	err := os.RemoveAll("fs")
	c.Assert(err, gc.IsNil)
}
func (f *FilesystemSuite) TearDownSuite(c *gc.C) {
	fmt.Println("Tearind down the suit : clearing fs directory...")
	err := os.RemoveAll("fs")
	c.Assert(err, gc.IsNil)
}
func (f *FilesystemSuite) TestWrite(c *gc.C) {

	// filesystem config 오브젝트 생성
	fsCfg := make(jsonObj)
	fsCfg["path"] = "fs/"

	// filesystem storage provider 인스턴스 생성
	filesystem, err := storage_providers.CreateStorageProvider("filesystem", fsCfg)

	// 에러 체크
	c.Assert(err, gc.IsNil)

	// 페이로드 stub 생성
	payload := &esPayloadStub{"event-data-test", fmt.Sprintf("filesystem.write.test.%d", 0)}

	filesystem.Write(payload)

	dirs, err := os.ReadDir("fs/event-data-test")
	c.Assert(err, gc.IsNil)

	for _, dir := range dirs {
		c.Assert("filesystem.write.test.0", gc.Equals, dir.Name())
	}

}

func (f *FilesystemSuite) TestConcurrentWrite(c *gc.C) {

	// filesystem config 오브젝트 생성
	fsCfg := make(jsonObj)
	fsCfg["path"] = "fs/"

	// filesystem storage provider 인스턴스 생성
	filesystem, err := storage_providers.CreateStorageProvider("filesystem", fsCfg)

	// 에러 체크
	c.Assert(err, gc.IsNil)

	// 동시 트린잭션 개수
	requests := 10

	// WaitGroup 생성
	var wg sync.WaitGroup

	// 고루틴 생성
	for i := 0; i < requests; i++ {
		wg.Add(1)
		go func(idx int, wg *sync.WaitGroup) {
			// 페이로드 stub 생성
			payload := &esPayloadStub{"event-data-test-concurrent", fmt.Sprintf("filesystem.write.test.%d", idx)}
			filesystem.Write(payload)
			wg.Done()
		}(i, &wg)
	}
	wg.Wait()

	dirs, err := os.ReadDir("fs/event-data-test-concurrent")
	c.Assert(err, gc.IsNil)

	for idx, dir := range dirs {
		c.Assert(fmt.Sprintf("filesystem.write.test.%d", idx), gc.Equals, dir.Name())
	}
}

type jsonObj = map[string]interface{}

var _ payloads.Payload = new(fsPayloadStub)

type fsPayloadStub struct {
	dir      string
	filename string
}

// Clone implements payloads.Payload
func (*fsPayloadStub) Clone() payloads.Payload {
	ps := &fsPayloadStub{}
	return ps
}

// MarkAsProcessed implements payloads.Payload
func (*fsPayloadStub) MarkAsProcessed() {

}

// Out implements payloads.Payload
func (p *fsPayloadStub) Out() (string, string, []byte) {
	return p.dir, p.filename, []byte(`{}`)
}
