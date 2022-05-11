package event_data

import (
	"context"
	"encoding/json"
	"errors"
	"event-data-pipeline/pkg/config"
	"event-data-pipeline/pkg/consumers"
	"event-data-pipeline/pkg/logger"
	"event-data-pipeline/pkg/pipelines"
	"event-data-pipeline/pkg/processors"
	"sync"
)

type EventDataPipeline struct {
	p        *pipelines.Pipeline
	cfgsPath string
	cfgs     []*config.PipelineCfg
}

func NewEventDataPipeline(cfg config.Config) (*EventDataPipeline, error) {
	ec := &EventDataPipeline{}

	// 설정 정보 경로 값 인스턴스에 저장.
	ec.cfgsPath = cfg.PipelineCfgsPath

	var err error
	// 제공된 경로로 부터 설정 정보를 읽어옵니다.
	ec.cfgs = config.NewPipelineConfig(cfg.PipelineCfgsPath)
	if ec.cfgs == nil {
		logger.Errorf("loaded configuration is nil")
		err = errors.New("loaded configuration is nil")
	}
	return ec, err
}

func (e *EventDataPipeline) SetCollectorRuntimeConfig(confs []*config.PipelineCfg) {
	e.cfgs = confs
}

func (e *EventDataPipeline) ValidateConfigs() error {
	// 인스턴스가 제로값인 경우 에러를 반환.
	if e == nil {
		logger.Errorf("%t is %v", e, e)
		return errors.New("EventDataPipeline instance is nil")
	}

	// 메모리에 로드된 설정 정보를 출력.
	if e.cfgs != nil {
		logger.Debugf("Loading EventDataPipeline Configs from memory: %s", ObjectToJsonString(e.cfgs))
	}

	// 메모리 상 설정 값이 비어있는 경우
	// 파일로부터 다시 읽기를 시도
	if e.cfgs == nil {
		e.cfgs = config.NewPipelineConfig(e.cfgsPath)
		logger.Infof("Loading EventDataPipeline Configs from file : %s", ObjectToJsonString(e.cfgs))
	}
	if e.cfgs == nil {
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
	for _, cfg := range e.cfgs {
		wg.Add(1)

		// 이벤트 기반 데이터를 소비하는 컨슈머 생성
		consumer, err := consumers.CreateConsumer(cfg.Consumer.Name, cfg.Consumer.Config)
		if err != nil {
			logger.Errorf("%v", err)
			return err
		}
		logger.Debugf("%v consumer created", consumer)

		// 컨슈머로 부터 데이터를 받아 처리하는 0개 이상의 프로세서 슬라이스 초기화
		stageRunners := make([]processors.Processor, len(cfg.Processors))
		for _, p := range cfg.Processors {
			processor, err := processors.CreateProcessor(p.Name, p.Config)
			if err != nil {
				return err
			}
			// 스테이지 러너에 생성된 프로세서를 등록
			stageRunners = append(stageRunners, processor)
		}

	}

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

	//Run Process
	e.p.Process(ctx, nil, nil)

	return nil
}

func (e *EventDataPipeline) GetCollectorRuntimeConfig() []*config.PipelineCfg {
	return e.cfgs
}

func ObjectToJsonString(obj interface{}) string {
	b, err := json.Marshal(obj)
	if err != nil {
		logger.Panicf("%v", err)
	}
	return string(b)
}
