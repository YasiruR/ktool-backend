package http

type clusterRes struct {
	Clusters 	[]clusterInfo	`json:"clusters"`
}

type clusterInfo struct {
	Id 					int 		`json:"id"`
	ClusterName			string 		`json:"cluster_name"`
	KafkaVersion		string		`json:"kafka_version"`
	Brokers 			[]broker	`json:"brokers"`
}

type topicData struct {
	Topics 		[]string		`json:"topics"`
}

type broker struct {
	Host 		string			`json:"host"`
	Port 		int				`json:"port"`
}

type errorMessage struct {
	Mesg 	string		`json:"mesg"`
}