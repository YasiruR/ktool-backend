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

//-----------------topic metrics-------------------------//

type topicMetricsRes struct {
	TotalMessages 		int				`json:"total_messages"`
	BytesInRate			int				`json:"bytes_in_rate"`
	BytesOutRate		int				`json:"bytes_out_rate"`
	MessageRate			int				`json:"message_rate"`
	Topics 				[]metricsTopic	`json:"topics"`
}

type metricsTopic struct {
	Name 				string				`json:"name"`
	Brokers 			[]string			`json:"brokers"`
	Partitions 			[]topicPartition	`json:"partitions"`
	WritablePartitions	int					`json:"writable_part"`
	UnderReplicatedPart	int					`json:"under_repl_part"`
	Replicas 			int					`json:"replicas"`
	InSyncReplicas		int					`json:"in_sync_repl"`
	OfflineReplicas		int					`json:"offline_repl"`
	Messages 			int					`json:"messages"`
	BytesIn				int					`json:"bytes_in"`
	BytesOut			int					`json:"bytes_out"`
	BytesRejected		int					`json:"bytes_rej"`
	ReplBytesIn			int					`json:"repl_bytes_in"`
	ReplBytesOut 		int					`json:"repl_bytes_out"`
}

type topicPartition struct {
	ID 				int		`json:"id"`
	FirstOffset 	int		`json:"first_offset"`
	LastOffset		int		`json:"last_offset"`
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