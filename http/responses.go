package http

//-------------------cluster--------------------------//

type clusterRes struct {
	Clusters 	[]clusterInfo	`json:"clusters"`
}

type clusterInfo struct {
	Id 					int 		`json:"id"`
	ClusterName			string 		`json:"cluster_name"`
	Brokers 			[]string	`json:"brokers"`
	Topics 				[]topic 	`json:"topics"`
	Available 			bool		`json:"available"`
	Connected			bool		`json:"connected"`
}

type topic struct {
	Name 		string		`json:"name"`
	Partitions 	[]int32 	`json:"partitions"`
}

type errorMessage struct {
	Mesg 	string		`json:"mesg"`
}

//-----------------add user req-------------------------//

type userRes struct {
	Token 		string		`json:"token"`
}

////-----------------broker overview---------------------//
//
//type brokerOverviewRes struct {
//	TotalBrokers 			int			`json:"total_brokers"`
//	TotalPartitions 		int			`json:"total_partitions"`
//	TotalReplicas 			int			`json:"total_replicas"`
//	TotalProductionRate		float64		`json:"total_partition_rate"`
//	TotalConsumptionRate	float64		`json:"total_consumption_rate"`
//	ActiveControllerID 		string		`json:"active_controller_id"`
//	ZookeeperAvail			bool		`json:"zookeeper_avail"`
//	KafkaVersion 			string		`json:"kafka_version"`
//	Brokers 				[]broker	`json:"brokers"`
//}
//
//type broker struct {
//	IncomingByteRate       float64 		`json:"incoming_byte_rate"`
//	RequestRate            int64 		`json:"request_rate"`
//	RequestSize            int64		`json:"request_size"`
//	RequestLatency         int64		`json:"request_latency"`
//	OutgoingByteRate       float64		`json:"outgoing_byte_rate"`
//	ResponseRate           float64		`json:"response_rate"`
//	ResponseSize           int64		`json:"response_size"`
//	BrokerIncomingByteRate float64		`json:"broker_incoming_byte_rate"`
//	BrokerRequestRate      float64		`json:"broker_request_rate"`
//	BrokerRequestSize      int64		`json:"broker_request_size"`
//	BrokerRequestLatency   int64		`json:"broker_request_latency"`
//	BrokerOutgoingByteRate float64		`json:"broker_outgoing_byte_rate"`
//	BrokerResponseRate     float64		`json:"broker_response_rate"`
//	BrokerResponseSize     int64		`json:"broker_response_size"`
//}