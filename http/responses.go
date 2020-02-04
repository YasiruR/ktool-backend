package http

type clusterInfo struct {
	Id 					int 		`json:"id"`
	ClusterName			string 		`json:"cluster_name"`
	KafkaVersion		string		`json:"kafka_version"`
}

type topicData struct {
	Topics 		[]string		`json:"topics"`
}