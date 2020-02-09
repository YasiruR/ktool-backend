package http

type addExistingCluster struct {
	ClusterName 		string 			`json:"cluster_name"`
	KafkaVersion 		string			`json:"kafka_version"`
	Brokers 			[]server		`json:"brokers"`
}

type server struct {
	Host 		string 			`json:"host"`
	Port 		int				`json:"port"`
}

type connectToCluster struct {
	ClusterID 	int				`json:"cluster_id"`
	Brokers 	[]string		`json:"brokers"`
}