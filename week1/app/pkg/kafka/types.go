package kafka

type PartitionsResponse struct {
	Topic      string    `json:"topic"`
	Partitions []float64 `json:"partitions"`
}
