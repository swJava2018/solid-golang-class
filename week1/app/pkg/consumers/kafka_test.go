package consumers

import (
	"context"
	"event-data-pipeline/pkg/cli"
	"event-data-pipeline/pkg/config"
	"event-data-pipeline/pkg/logger"
	"os"
	"path"
	"runtime"
	"testing"
	"time"

	"github.com/alexflint/go-arg"
)

func TestKafkaConsumerClient_Consume(t *testing.T) {
	configPath := getCurDir() + "/test/consumers/config.json"
	os.Setenv("EDP_ENABLE_DEBUG_LOGGING", "true")
	os.Setenv("EDP_CONFIG", configPath)
	os.Args = nil
	arg.MustParse(&cli.Args)
	logger.Setup()
	cfg := config.NewConfig()
	pipeCfgs := config.NewPipelineConfig(cfg.PipelineCfgsPath)
	stream := make(chan interface{})
	errCh := make(chan error)

	for _, cfg := range pipeCfgs {
		cfg.Consumer.Config["stream"] = stream
		cfg.Consumer.Config["errch"] = errCh
		kafkaConsumer, err := CreateConsumer(cfg.Consumer.Name, cfg.Consumer.Config)
		if err != nil {
			t.Error(err)
		}
		err = kafkaConsumer.Init()
		if err != nil {
			t.Error(err)
		}
		kafkaConsumer.Consume(context.TODO())
	}

	for {
		select {
		case data := <-stream:
			t.Logf("data: %v", data)
			return
		case err := <-errCh:
			t.Logf("err: %v", err)
		default:
			time.Sleep(1 * time.Second)
		}
	}

}

func getCurDir() string {
	_, filename, _, _ := runtime.Caller(0)
	dir := path.Join(path.Dir(filename), "../../")
	err := os.Chdir(dir)
	if err != nil {
		panic(err)
	}
	return dir
}
