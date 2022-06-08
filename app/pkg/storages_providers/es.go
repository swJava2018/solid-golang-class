package storages_providers

import (
	"bytes"
	"context"
	"encoding/json"
	"event-data-pipeline/pkg/concur"
	"fmt"
	"strconv"
	"sync"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"golang.org/x/time/rate"
	"nuance.xaas-logging.event-log-collector/pkg/pipeline/ratelimiter"

	es "github.com/elastic/go-elasticsearch/v8"
	log "nuance.xaas-logging.event-log-collector/pkg/pipeline/logging"
)

const (
	REMOTE_SERVICE_ES    = "ElasticSearch"
	TICKER_TIMEOUT_MS    = 30000
	RECORD_CNT_THRESHOLD = 1000
)

type bulkResponse struct {
	Errors bool `json:"errors"`
	Items  []struct {
		Index struct {
			ID     string `json:"_id"`
			Result string `json:"result"`
			Status int    `json:"status"`
			Error  struct {
				Type   string `json:"type"`
				Reason string `json:"reason"`
				Cause  struct {
					Type   string `json:"type"`
					Reason string `json:"reason"`
				} `json:"caused_by"`
			} `json:"error"`
		} `json:"index"`
	} `json:"items"`
}

type ElasticSearchClientConfig struct {
	DocType    string                `json:"doc_type,omitempty"`
	Refresh    bool                  `json:"refresh,omitempty"`
	RateLimit  ratelimiter.RateLimit `json:"rate_limit,omitempty"`
	MaxRetries int                   `json:"max_retries,omitempty"`
	Delay      int                   `json:"delay,omitempty"`
}
type ElasticSearchClient struct {
	client       *es.Client
	DocumentType string
	index        string
	Refresh      bool

	buf    bytes.Buffer
	count  int
	record chan interface{}
	mu     sync.Mutex
	ticker *time.Ticker
	done   chan bool

	workers *workerpool.WorkerPool
	job     chan interface{}

	rateLimiter *rate.Limiter
	maxRetries  int
	delay       int

	recordCount int
}

func init() {
	Register("elasticsearch", NewElasticSearchClient)
}

func NewElasticSearchClient(config jsonObj) StorageProvider {
	// 바이트 변환
	cfgData, _ := json.Marshal(config)
	var escConf ElasticSearchClientConfig
	var esConf es.Config
	json.Unmarshal(cfgData, &escConf)
	json.Unmarshal(cfgData, &esConf)

	if escConf.DocType == "" {
		escConf.DocType = "mix3-record"
	}

	es, err := es.NewClient(esConf)
	if err != nil {
		log.Errorf("error in creating elasticsearch client: %s", err)
	}

	transport, _ := json.Marshal(es.Transport)
	log.Debugf("Elasticsearch transport: %s", string(transport))

	ec := &ElasticSearchClient{
		client:       es,
		DocumentType: escConf.DocType,
		Refresh:      escConf.Refresh,

		job:         make(chan interface{}),
		record:      make(chan interface{}),
		done:        make(chan bool),
		ticker:      time.NewTicker(TICKER_TIMEOUT_MS * time.Millisecond),
		count:       0,
		recordCount: 0,
		rateLimiter: ratelimiter.NewRateLimiter(escConf.RateLimit),
		maxRetries:  escConf.MaxRetries,
		delay:       escConf.Delay,
	}

	numWorkers := 1

	ec.workers = concur.NewWorkerPool("elasticsearch-workers", ec.job, numWorkers, ec.storageHandler)
	ec.workers.Start()
	go ec.commitHandler()

	return ec
}

// implement storage Write interface method which is called from processor
func (e *ElasticSearchClient) Write(key string, path string, data []byte) (int, error) {

	return len(data), nil
}

func (e *ElasticSearchClient) Drain(request interface{}) {
	start := time.Now()
	log.Debugf("writing to elasticsearch record chan...")
	e.record <- request
	log.Debugf("done writing to elasticsearch chan in %v ms...", time.Since(start).Milliseconds())

}

func (e *ElasticSearchClient) Write(payload interface{}) {
	lastCommit := time.Now()
	for {
		select {
		case record := <-e.record:
			if record != nil {
				job := record.(*Job)
				meta := []byte(fmt.Sprintf(`{ "index" : { "_id" : "%v" } }%s`, job.docID, "\n"))
				data := append(job.body, "\n"...)
				e.mu.Lock()
				e.buf.Grow(len(meta) + len(data))
				e.buf.Write(meta)
				e.buf.Write(data)
				e.count += 1
				if e.count >= RECORD_CNT_THRESHOLD {
					log.Infof("elasticsearch bulk write triggered by record count")
					log.Infof("committing batch [count: %v, queue: %v]", e.count, len(e.record))
					go e.commit()
					lastCommit = time.Now()
				}
				e.mu.Unlock()
			}
		case <-e.ticker.C:
			log.Infof("elasticsearch bulk write triggered by elapsed time")
			if e.count == 0 {
				log.Infof("no records to commit")
			} else if time.Since(lastCommit).Milliseconds() >= TICKER_TIMEOUT_MS {
				go e.commit()
				lastCommit = time.Now()
			}
		case <-e.done:
			return
		}
	}

}

func (e *ElasticSearchClient) Shutdown() {
	e.workers.Stop()
	e.ticker.Stop()
	e.done <- true
	log.Infof("closed job queue")
}

func (e *ElasticSearchClient) Read(key string, path string) ([]byte, error) {
	return nil, fmt.Errorf("not Implemented")
}

func (e *ElasticSearchClient) commit() (int, error) {

	retry := 0
	lck := time.Now()
	e.mu.Lock()
	log.Infof("took %f to acquire lock", time.Since(lck).Seconds())
	b := make([]byte, e.buf.Len())
	e.buf.Read(b)
	r := bytes.NewReader(b)
	e.buf.Reset()
	e.recordCount += e.count
	e.count = 0
	e.mu.Unlock()

	for {
		ctx := context.Background()
		startWait := time.Now()
		// rate limiting ...
		e.rateLimiter.Wait(ctx)
		log.Debugf("rate limited for %f seconds", time.Since(startWait).Seconds())

		start := time.Now()
		numErrors := 0
		numIndexed := 0

		defer func() {
			dur := time.Since(start).Milliseconds()
			if numErrors > 0 {
				log.Errorf(
					"Indexed [%v] documents with [%v] errors in %v ms (%d docs/sec)",
					numIndexed,
					numErrors,
					dur,
					int64((1000.0/float64(dur))*float64(numIndexed)),
				)
			} else if numIndexed > 0 {
				log.Infof(
					"Successfully indexed [%v] documents in %v ms (%d docs/sec)",
					numIndexed,
					dur,
					int64((1000.0/float64(dur))*float64(numIndexed)),
				)
			} else {
				log.Infof("numErrors: %d, numIndexed: %d", numErrors, numIndexed)
			}
		}()

		// Observe write duration in seoncds and set histogram and guage metric
		timer := prometheus.NewTimer(*e.WriteDurationHistogram)
		res, err := e.client.Bulk(r,
			e.client.Bulk.WithIndex(e.index),
			e.client.Bulk.WithDocumentType(e.DocumentType),
		)
		if err != nil || res == nil {
			log.Errorf("error in bulk writing : %s", err.Error())
			e.RemoteConnectionGauge.WithLabelValues(REMOTE_SERVICE_ES).Set(0)
			retry++
			if e.maxRetries >= 0 && retry > e.maxRetries {
				err := fmt.Errorf("retry[%d] exceeded max retries[%d]", retry, e.maxRetries)
				e.WriteTotalCounter.WithLabelValues("0", err.Error()).Add(float64(e.count))
				log.Errorf("error in bulk writing : %s", err.Error())
				return 0, err
			}
			time.Sleep(time.Duration(time.Duration(e.delay) * time.Second))
			continue
		}
		e.RemoteConnectionGauge.WithLabelValues(REMOTE_SERVICE_ES).Set(1)
		defer res.Body.Close()
		e.WriteDurationGauge.WithLabelValues(REMOTE_SERVICE_ES).Set(timer.ObserveDuration().Seconds())

		if !res.IsError() {
			var blk *bulkResponse
			if err := json.NewDecoder(res.Body).Decode(&blk); err != nil {
				log.Errorf("Failure to to parse response body: %s", err)
				return 0, err
			} else {
				// If the whole request failed, print error and mark all documents as failed
				for _, d := range blk.Items {
					// ... so for any HTTP status above 201 ...
					//
					if d.Index.Status > 201 {
						// ... increment the error counter ...
						//
						numErrors++
						e.WriteTotalCounter.WithLabelValues(strconv.Itoa(d.Index.Status), d.Index.Error.Reason).Inc()
						// ... and print the response status and error information ...
						log.Errorf("Error: [%d]: %s: %s: %s: %s",
							d.Index.Status,
							d.Index.Error.Type,
							d.Index.Error.Reason,
							d.Index.Error.Cause.Type,
							d.Index.Error.Cause.Reason,
						)
					} else {
						// ... otherwise increase the success counter.
						//
						numIndexed++
						e.WriteTotalCounter.WithLabelValues(strconv.Itoa(d.Index.Status), res.Status()).Inc()
					}
				}
				return numIndexed, nil
			}
		} else {
			numErrors += e.count
			e.WriteTotalCounter.WithLabelValues(strconv.Itoa(res.StatusCode), res.Status()).Add(float64(numErrors))
			var raw map[string]interface{}
			if err := json.NewDecoder(res.Body).Decode(&raw); err != nil {
				log.Errorf("Failure to to parse response body: %s", err)
			} else {
				log.Printf("Error: [%d] %s: %s",
					res.StatusCode,
					raw["error"].(map[string]interface{})["type"],
					raw["error"].(map[string]interface{})["reason"],
				)
			}
			return numErrors, nil
		}
	}

}
