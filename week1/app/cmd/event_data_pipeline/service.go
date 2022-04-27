package event_data_pipeline

import (
	"event-data-pipeline/pkg/config"
	"log"
	"net/http"
	_ "net/http/pprof"
)

type (
	jsonObj = map[string]interface{}
	jsonArr = []interface{}
)

// Run is the entrypoint for running the http server & event log collector as a service
func Run(cfg config.Config) {
	// Force garbage collection
	// go utils.GarbageCollector()

	go func() {
		log.Println(http.ListenAndServe("localhost:6060", nil))
	}()

	// create event log collector instance
	// eventCollector, err := pipeline.NewEventCollector(cfg, collectorSignal)

	// if err != nil {
	// 	log.Panicf(err.Error())
	// }

	// // Run elc second
	// err = eventCollector.RunCollector(cfg)
	// if err != nil {
	// 	log.Panicf(err.Error())
	// }

}
