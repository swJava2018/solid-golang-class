package concur

type Job struct {
	index string
	docID string
	body  []byte
}

func NewJob(key string, path string, data []byte) *Job {
	if key == "" || path == "" {
		return nil
	}
	if data == nil {
		return nil
	}
	return &Job{
		index: key,
		docID: path,
		body:  data,
	}
}
