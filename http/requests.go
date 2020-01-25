package http

type reqAddExistingCluster struct {
	ClusterName 		string 			`json:"cluster_name"`
	ClusterVersion 		float64			`json:"cluster_version"`
	ZookeeperHost 		string 			`json:"zookeeper_host"`
	ZookeeperPort 		int				`json:"zookeeper_port"`
}

type reqTestNewCluster struct {
	ZookeeperHost 		string 			`json:"zookeeper_host"`
	ZookeeperPort 		int				`json:"zookeeper_port"`
}