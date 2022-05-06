package kafka

type PartitionsResponse struct {
	Topic      string    `json:"topic"`
	Partitions []float64 `json:"partitions"`
}

type Consumer interface {
	CommitOffset(partition int, offset int) error
	GetOffsetRange(partition int) (int, int, error)
	GetCommittedOffset(partition int) (int, error)
}
