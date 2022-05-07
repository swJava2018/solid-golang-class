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

	"github.com/alexflint/go-arg"
)

func getCurDir() string {
	_, filename, _, _ := runtime.Caller(0)
	dir := path.Join(path.Dir(filename), "../../")
	err := os.Chdir(dir)
	if err != nil {
		panic(err)
	}
	return dir
}
func TestConsumerKafkaCreate(t *testing.T) {
	configPath := getCurDir() + "/test/consumers/config.json"
	os.Setenv("EDP_CONFIG", configPath)
	os.Args = nil
	arg.MustParse(&cli.Args)
	cfg := config.NewConfig()
	pipeCfgs := config.NewPipelineConfig(cfg.PipelineCfgsPath)
	for _, cfg := range pipeCfgs {
		kafkaConsumer, err := CreateConsumer(cfg.Consumer.Name, cfg.Consumer.Config)
		if err != nil {
			t.Error(err)
		}
		t.Logf("%T", kafkaConsumer)
		consumer, ok := kafkaConsumer.(*KafkaConsumerClient)
		if !ok {
			t.Error("failed to switch type to *KafkaConsumerClient")
		}
		err = consumer.kafkaConsumer.Create()
		if err != nil {
			t.Error(err)
		}

	}

}

func TestConsumerKafkaConsume(t *testing.T) {
	configPath := getCurDir() + "/test/consumers/config.json"
	os.Setenv("EDP_ENABLE_DEBUG_LOGGING", "true")
	os.Setenv("EDP_CONFIG", configPath)
	os.Args = nil
	arg.MustParse(&cli.Args)
	logger.Setup()
	cfg := config.NewConfig()
	pipeCfgs := config.NewPipelineConfig(cfg.PipelineCfgsPath)
	for _, cfg := range pipeCfgs {
		kafkaConsumer, err := CreateConsumer(cfg.Consumer.Name, cfg.Consumer.Config)
		if err != nil {
			t.Error(err)
		}
		t.Logf("%T", kafkaConsumer)
		consumer, ok := kafkaConsumer.(*KafkaConsumerClient)
		if !ok {
			t.Error("failed to switch type to *KafkaConsumerClient")
		}
		if err != nil {
			t.Error(err)
		}
		consumer.Consume(context.TODO(), nil, nil, nil)

	}

}
