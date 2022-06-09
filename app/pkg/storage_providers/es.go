package storage_providers

import (
	"bytes"
	"context"
	"encoding/json"
	"event-data-pipeline/pkg/concur"
	"event-data-pipeline/pkg/logger"
	"event-data-pipeline/pkg/payloads"
	"event-data-pipeline/pkg/ratelimit"
	"fmt"
	"sync"
	"time"

	"golang.org/x/time/rate"

	es "github.com/elastic/go-elasticsearch/v8"
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
	DocType    string              `json:"doc_type,omitempty"`
	Refresh    bool                `json:"refresh,omitempty"`
	RateLimit  ratelimit.RateLimit `json:"rate_limit,omitempty"`
	MaxRetries int                 `json:"max_retries,omitempty"`
	Delay      int                 `json:"delay,omitempty"`
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

	workers *concur.WorkerPool
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

	es, err := es.NewClient(esConf)
	if err != nil {
		logger.Errorf("error in creating elasticsearch client: %s", err)
	}

	transport, _ := json.Marshal(es.Transport)
	logger.Debugf("Elasticsearch transport: %s", string(transport))

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
		rateLimiter: ratelimit.NewRateLimiter(escConf.RateLimit),
		maxRetries:  escConf.MaxRetries,
		delay:       escConf.Delay,
	}

	numWorkers := 1

	ec.workers = concur.NewWorkerPool("elasticsearch-workers", ec.job, numWorkers, ec.Write)
	ec.workers.Start()

	return ec
}

func (e *ElasticSearchClient) Drain(request interface{}) {
	start := time.Now()
	logger.Debugf("writing to elasticsearch record chan...")
	e.record <- request
	logger.Debugf("done writing to elasticsearch chan in %v ms...", time.Since(start).Milliseconds())

}

func (e *ElasticSearchClient) Write(payload interface{}) (int, error) {
	if payload != nil {

		index, docID, data := payload.(payloads.Payload).Out()
		e.index = index

		// documentID 메타정보 오브젝트 생성
		meta := []byte(fmt.Sprintf(`{ "index" : { "_id" : "%v" } }%s`, docID, "\n"))

		// bulk write 을 위한 개행
		data = append(data, "\n"...)

		// 메타, 데이타 오브젝트 사이즈 버퍼 할당
		e.buf.Grow(len(meta) + len(data))

		// 메타 정보 쓰기
		e.buf.Write(meta)

		// 데이터 정보 쓰기
		e.buf.Write(data)
		e.count += 1
		// if e.count >= RECORD_CNT_THRESHOLD {
		logger.Infof("elasticsearch bulk write triggered by record count")
		logger.Infof("committing batch [count: %v, queue: %v]", e.count, len(e.record))

		// pass local copy of write data to the commit
		// for async write
		buf := e.buf.Bytes()
		e.bulkWrite(index, buf)
		// }
	}
	return 0, nil
}

func (e *ElasticSearchClient) bulkWrite(index string, data []byte) (int, error) {

	retry := 0
	lck := time.Now()
	e.mu.Lock()
	logger.Infof("took %f to acquire lock", time.Since(lck).Seconds())
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
		logger.Debugf("rate limited for %f seconds", time.Since(startWait).Seconds())

		start := time.Now()
		numErrors := 0
		numIndexed := 0

		defer func() {
			dur := time.Since(start).Milliseconds()
			if numErrors > 0 {
				logger.Errorf(
					"Indexed [%v] documents with [%v] errors in %v ms (%d docs/sec)",
					numIndexed,
					numErrors,
					dur,
					int64((1000.0/float64(dur))*float64(numIndexed)),
				)
			} else if numIndexed > 0 {
				logger.Infof(
					"Successfully indexed [%v] documents in %v ms (%d docs/sec)",
					numIndexed,
					dur,
					int64((1000.0/float64(dur))*float64(numIndexed)),
				)
			} else {
				logger.Infof("numErrors: %d, numIndexed: %d", numErrors, numIndexed)
			}
		}()

		// Observe write duration in seoncds and set histogram and guage metric
		res, err := e.client.Bulk(r,
			e.client.Bulk.WithIndex(e.index),
			e.client.Bulk.WithDocumentType(e.DocumentType),
		)
		if err != nil || res == nil {
			logger.Errorf("error in bulk writing : %s", err.Error())
			retry++
			if e.maxRetries >= 0 && retry > e.maxRetries {
				err := fmt.Errorf("retry[%d] exceeded max retries[%d]", retry, e.maxRetries)
				logger.Errorf("error in bulk writing : %s", err.Error())
				return 0, err
			}
			time.Sleep(time.Duration(time.Duration(e.delay) * time.Second))
			continue
		}
		defer res.Body.Close()

		if !res.IsError() {
			var blk *bulkResponse
			if err := json.NewDecoder(res.Body).Decode(&blk); err != nil {
				logger.Errorf("Failure to to parse response body: %s", err)
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
						// ... and print the response status and error information ...
						logger.Errorf("Error: [%d]: %s: %s: %s: %s",
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
					}
				}
				return numIndexed, nil
			}
		} else {
			numErrors += e.count
			var raw map[string]interface{}
			if err := json.NewDecoder(res.Body).Decode(&raw); err != nil {
				logger.Errorf("Failure to to parse response body: %s", err)
			} else {
				logger.Printf("Error: [%d] %s: %s",
					res.StatusCode,
					raw["error"].(map[string]interface{})["type"],
					raw["error"].(map[string]interface{})["reason"],
				)
			}
			return numErrors, nil
		}
	}

}
