package http

type reqAddExistingCluster struct {
	ClusterName 		string 			`json:"cluster_name"`
	KafkaVersion 		float64			`json:"kafka_version"`
	ZookeeperHost 		string 			`json:"zookeeper_host"`
	ZookeeperPort 		int				`json:"zookeeper_port"`
}

type reqTestNewCluster struct {
	Host 		string 			`json:"zookeeper_host"`
	Port 		int				`json:"zookeeper_port"`
}