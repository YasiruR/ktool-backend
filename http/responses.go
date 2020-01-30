package http

type clusterInfo struct {
	Id 					int 		`json:"id"`
	ClusterName			string 		`json:"cluster_name"`
	KafkaVersion		string		`json:"kafka_version"`
}
