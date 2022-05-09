package pipeline

import (
	"context"
	"encoding/json"
	"errors"
	"event-data-pipeline/pkg/config"
	"event-data-pipeline/pkg/consumers"
	"event-data-pipeline/pkg/logger"
	"event-data-pipeline/pkg/pipelines"
	"sync"
)

type EventDataPipeline struct {
	p        *pipelines.Pipeline
	CfgsPath string
	Cfgs     []*config.PipelineCfg
}

func NewEventDataPipeline(cfg config.Config) (*EventDataPipeline, error) {
	ec := &EventDataPipeline{}

	// 설정 정보 경로 값 인스턴스에 저장.
	ec.CfgsPath = cfg.PipelineCfgsPath

	var err error
	// 제공된 경로로 부터 설정 정보를 읽어옵니다.
	ec.Cfgs = config.NewPipelineConfig(cfg.PipelineCfgsPath)
	if ec.Cfgs == nil {
		logger.Errorf("loaded configuration is nil")
		err = errors.New("loaded configuration is nil")
	}
	return ec, err
}

func (e *EventDataPipeline) SetCollectorRuntimeConfig(confs []*config.PipelineCfg) {
	e.Cfgs = confs
}

func (e *EventDataPipeline) ValidateConfigs() error {
	// 인스턴스가 제로값인 경우 에러를 반환.
	if e == nil {
		logger.Errorf("%t is %v", e, e)
		return errors.New("EventDataPipeline instance is nil")
	}

	// 메모리에 로드된 설정 정보를 출력.
	if e.Cfgs != nil {
		logger.Debugf("Loading EventDataPipeline Configs from memory: %s", ObjectToJsonString(e.Cfgs))
	}

	// 메모리 상 설정 값이 비어있는 경우
	// 파일로부터 다시 읽기를 시도
	if e.Cfgs == nil {
		e.Cfgs = config.NewPipelineConfig(e.CfgsPath)
		logger.Infof("Loading EventDataPipeline Configs from file : %s", ObjectToJsonString(e.Cfgs))
	}
	if e.Cfgs == nil {
		return errors.New("did not pass configs validation.")
	}
	return nil
}

// 파이프라인을 구동하는 메소드
func (e *EventDataPipeline) Run() error {

	// Goroutine 실행 후 대기를 위한 WaiterGroup
	var wg sync.WaitGroup

	// Graceful Shutdown 을 위한 Context, CancelFunction
	ctx, cancelFunc := context.WithCancel(context.Background())
	defer cancelFunc()

	// loop through PipelineConfigs
	for _, cfg := range cfgs {
		wg.Add(1)
		consumer, err := consumers.CreateConsumer(cfg.Consumer.Name, cfg.Consumer.Config)
		if err != nil {
			logger.Errorf("%v", err)
			return err
		}

	}

	// for _, conf := range collectorConfs {
	// 	wg.Add(1)

	// 	// Source

	// Muliple Processors
	// _processors := make([]processors.Processor, len(conf.Processors))

	// Multiple Sink
	// _storage := make([]storage.Storage, len(*p.Storage))
	// for i, s := range *p.Storage {
	// 	logger.Debugf("storage[%d]: %v", i, s.Type)
	// 	_storage[i], err = storage.CreateStorage(s.Type, s.Config)
	// 	if err != nil {
	// 		logger.Errorf("%v", err)
	// 		return err
	// 	}
	// }

	// e. create event logger collector
	// collector := .NewCollector(consumer, _processors)

	// f. run event logger collector
	// go collector.Start(ctx, &wg)
	// }
	// e.signal.ReceiveShutDown() // blocks until shutdown signal
	// logger.Infof("\n*********************************\nGraceful shutdown signal received\n*********************************")
	// cancelFunc() // Signal cancellation to context.Context
	// wg.Wait()
	// go e.signal.SendDone(true) // send shutdown done signal to receiver
	// logger.Infof("\n*********************************\nGraceful shutdown completed      \n*********************************")
	return nil
}

func (e *EventDataPipeline) GetCollectorRuntimeConfig() []*config.PipelineCfg {
	return e.Cfgs
}

func ObjectToJsonString(obj interface{}) string {
	b, err := json.Marshal(obj)
	if err != nil {
		logger.Panicf("%v", err)
	}
	return string(b)
}
