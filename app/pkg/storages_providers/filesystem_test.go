package storages_providers_test

import (
	"event-data-pipeline/pkg/logger"
	"event-data-pipeline/pkg/payloads"
	"event-data-pipeline/pkg/storages_providers"
	"fmt"
	"os"
	"sync"
	"testing"

	gc "gopkg.in/check.v1"
)

// Hook up gocheck into the "go test" runner.
func Test(t *testing.T) { gc.TestingT(t) }

type FilesystemSuite struct{}

var _ = gc.Suite(&FilesystemSuite{})

func (f *FilesystemSuite) SetUpSuite(c *gc.C) {
	os.Args = nil
	os.Setenv("EDP_ENABLE_DEBUG_LOGGING", "true")
	logger.Setup()
}

func (f *FilesystemSuite) SetUpTest(c *gc.C) {
	fmt.Println("setup test")
}
func (f *FilesystemSuite) TearDownTest(c *gc.C) {
	fmt.Println("tear down test")
}
func (f *FilesystemSuite) TearDownSuite(c *gc.C) {
	fmt.Println("tear down suit")
}

func (f *FilesystemSuite) TestWrite(c *gc.C) {

	// filesystem config 오브젝트 생성
	fsCfg := make(jsonObj)
	fsCfg["path"] = "fs/"

	// filesystem storage provider 인스턴스 생성
	filesystem, err := storages_providers.CreateStorageProvider("filesystem", fsCfg)

	// 에러 체크
	c.Assert(err, gc.IsNil)

	// 페이로드 stub 생성
	payload := &payloadStub{0, "event-data-test"}

	filesystem.Write(payload)
}

func (f *FilesystemSuite) TestConcurrentWrite(c *gc.C) {

	// filesystem config 오브젝트 생성
	fsCfg := make(jsonObj)
	fsCfg["path"] = "fs/"

	// filesystem storage provider 인스턴스 생성
	filesystem, err := storages_providers.CreateStorageProvider("filesystem", fsCfg)

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
			payload := &payloadStub{idx, "event-data-test-concurrent"}
			filesystem.Write(payload)
			wg.Done()
		}(i, &wg)
	}
	wg.Wait()
}

type jsonObj = map[string]interface{}

var _ payloads.Payload = new(payloadStub)

type payloadStub struct {
	idx int
	dir string
}

// Clone implements payloads.Payload
func (*payloadStub) Clone() payloads.Payload {
	ps := &payloadStub{}
	return ps
}

// MarkAsProcessed implements payloads.Payload
func (*payloadStub) MarkAsProcessed() {

}

// Out implements payloads.Payload
func (p *payloadStub) Out() (string, string, []byte) {
	return p.dir, fmt.Sprintf("filesystem.storage_provider.test.%v", p.idx), []byte(`{}`)
}
