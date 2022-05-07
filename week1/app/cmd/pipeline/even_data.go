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
	p       *pipelines.Pipeline
	Configs []*config.PipelineCfg
}

func NewEventDataPipeline(cfg config.Config) (*EventDataPipeline, error) {
	ec := &EventDataPipeline{}

	var err error
	ec.Configs = config.NewPipelineConfig(cfg.PipelineCfgsPath)
	if ec.Configs == nil {
		logger.Errorf("loaded configuration is nil")
		err = errors.New("loaded configuration is nil")
	}
	return ec, err
}

func (e *EventDataPipeline) SetCollectorRuntimeConfig(confs []*config.PipelineCfg) {
	e.Configs = confs
}

func (e *EventDataPipeline) Run(cfg config.Config) error {
	if e == nil {
		logger.Errorf("%t is %v", e, e)
		return errors.New("Event logger Collector instance is nil")
	}
	var collectorConfs []*config.PipelineCfg

	// load from memory for runtime config modification reflection
	if e.Configs != nil {
		collectorConfs = e.Configs
		logger.Debugf("Loading CollectorConfs from memory: %s", ObjectToJsonString(collectorConfs))
	}

	// load from file if memory missing
	if collectorConfs == nil {
		collectorConfs := config.NewPipelineConfig(cfg.PipelineCfgsPath)
		logger.Infof("Loading CollectorConfs from file: %s", ObjectToJsonString(collectorConfs))
	}

	err := e.run(e.Configs)
	if err != nil {
		return err
	}
	return nil
}

func (e *EventDataPipeline) run(cfgs []*config.PipelineCfg) error {

	var wg sync.WaitGroup
	// context with canel for graceful shutdown
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
	return e.Configs
}

func ObjectToJsonString(obj interface{}) string {
	b, err := json.Marshal(obj)
	if err != nil {
		logger.Panicf("%v", err)
	}
	return string(b)
}
