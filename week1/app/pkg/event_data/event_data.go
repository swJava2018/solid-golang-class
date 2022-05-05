package pipelines

import (
	"context"
	"encoding/json"
	"errors"
	"event-data-pipeline/pkg/config"
	"event-data-pipeline/pkg/consumers"
	"event-data-pipeline/pkg/logger"
	"event-data-pipeline/pkg/pipelines"
	"event-data-pipeline/pkg/processors"
	"event-data-pipeline/pkg/sys"
	"os"
	"sync"
)

type EventDataPipeline struct {
	p       *pipelines.Pipeline
	Configs []*config.PipelineCfg
	signal  *sys.Signal
}

func NewEventDataPipeline(cfg config.Config, signal *sys.Signal) (*EventDataPipeline, error) {
	ec := &EventDataPipeline{
		signal: signal,
	}

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

func (e *EventDataPipeline) RunCollector(cfg config.Config) error {
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

	err := e.runCollector(e.Configs)
	if err != nil {
		return err
	}
	return nil
}

func (e *EventDataPipeline) runCollector(collectorConfs []*config.PipelineCfg) error {

	var wg sync.WaitGroup
	// context with canel for graceful shutdown
	ctx, cancelFunc := context.WithCancel(context.Background())
	defer cancelFunc()

	// 2. for each collector...
	for _, conf := range collectorConfs {
		wg.Add(1)

		// Source
		consumer, err := consumers.CreateConsumer(conf.Consumer.Name, conf.Consumer.Config)
		if err != nil {
			logger.Errorf("%v", err)
			return err
		}

		// Muliple Processors
		_processors := make([]processors.Processor, len(conf.Processors))

		// Multiple Sink
		_storage := make([]storage.Storage, len(*p.Storage))
		for i, s := range *p.Storage {
			logger.Debugf("storage[%d]: %v", i, s.Type)
			_storage[i], err = storage.CreateStorage(s.Type, s.Config)
			if err != nil {
				logger.Errorf("%v", err)
				return err
			}
		}

		// e. create event logger collector
		// collector := .NewCollector(consumer, _processors)

		// f. run event logger collector
		go collector.Start(ctx, &wg)
	}
	e.signal.ReceiveShutDown() // blocks until shutdown signal
	logger.Infof("\n*********************************\nGraceful shutdown signal received\n*********************************")
	cancelFunc() // Signal cancellation to context.Context
	wg.Wait()
	go e.signal.SendDone(true) // send shutdown done signal to receiver
	logger.Infof("\n*********************************\nGraceful shutdown completed      \n*********************************")
	return nil
}

func (e *EventDataPipeline) WaitCollectorShutDown() {
	e.signal.ReceiveDone()
}

func (e *EventDataPipeline) ShutdownCollector(signal os.Signal) {
	e.signal.SendShutDown(signal)
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
