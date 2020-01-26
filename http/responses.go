package http

type clusterInfo struct {
	Id 					int 		`json:"id"`
	ClusterName			string 		`json:"cluster_name"`
	KafkaVersion		float64		`json:"kafka_version"`
	ActiveControllers	int			`json:"active_controllers"`
}
