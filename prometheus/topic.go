package prometheus

var PromTopicMap map[int]map[string]topicMetrics

type topicMetrics struct {
	Brokers 		[]string
	Messages 		int
	BytesIn 		int
	BytesOut 		int
	BytesRejected	int
	ReplBytesIn		int
	ReplBytesOut	int
}

