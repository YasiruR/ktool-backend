package http

type reqAddExistingCluster struct {
	ClusterName 		string 			`json:"cluster_name"`
	KafkaVersion 		string			`json:"kafka_version"`
	ZookeeperHost 		string 			`json:"zookeeper_host"`
	ZookeeperPort 		int				`json:"zookeeper_port"`
}

type reqTestNewCluster struct {
	Host 		string 			`json:"host"`
	Port 		int				`json:"port"`
}