package domain

import "github.com/Shopify/sarama"

type Cluster struct {
	ID                int
	ClusterName       string
	KafkaVersion      string
	Zookeepers        []Zookeeper
	Brokers           []Broker
	SchemaRegistry    SchemaRegistry
	ActiveControllers int
	ZookeeperId       int
}

type KCluster struct{
	ClusterID       int
	ClusterName     string
	Consumer        sarama.Consumer
	Client          sarama.Client
	Brokers         []*sarama.Broker
	Topics          []KTopic
	Available       bool
	ClusterOverview ClusterOverview
}

type ClusterOverview struct {
	TotalPartitions           int             		`json:"partitions"`
	TotalTopics               int             		`json:"topics"`
	TotalReplicas             int             		`json:"replicas"`
	UnderReplicatedPartitions int             		`json:"under_replicated_partitions"`
	OfflinePartitions         int             		`json:"offline_partitions"`
	OfflineReplicas           int             		`json:"offline_replicas"`
	TotalByteInRate     	  map[int64]int64 		`json:"total_byte_in_rate"`
	TotalByteOutRate	      map[int64]int64 		`json:"total_byte_out_rate"`
	ActiveController          string          		`json:"active_controller"`
	ZookeeperAvail            bool            		`json:"zookeeper_avail"`
	KafkaVersion              string          		`json:"kafka_version"`
	Brokers                   []BrokerMetrics 		`json:"brokers"`
}

type BrokerMetrics struct {
	Host 						string				`json:"host"`
	Port 						int					`json:"port"`

	NumPartitions				int					`json:"partitions"`
	NumLeaders	 				int					`json:"leaders"`
	NumActControllers			int					`json:"act_controllers"`
	OfflinePartitions			int					`json:"offline_partitions"`
	UnderReplicated				int					`json:"under_replicated"`
	MessageRate					int					`json:"message_rate"`
	IsrExpansionRate			float64				`json:"isr_expansion_rate"`		//todo could be brought into a graph
	IsrShrinkRate				float64				`json:"isr_shrink_rate"`		//todo could be brought into a graph
	TimeToSendResProduce		float64				`json:"time_to_send_res_produce"` //todo could be brought into a graph
	NetworkProcAvgIdlePercent	float64				`json:"network_proc_avg_idle_percent"`
	ReqQueueWaitingTime			float64				`json:"req_queue_waiting_time"`
	ReqWaitForFollowerProd		float64				`json:"req_wait_for_follower_prod"`
	LeaderProcessTime			float64				`json:"leader_process_time"`
	TotalReqTimeRate			float64				`json:"total_req_time_rate"`
	MaxLagBtwLeadAndRepl		float64				`json:"max_lag_btw_lead_and_repl"`
	UncleanLeadElec				float64				`json:"unclean_lead_elec"`
	FailedFetchReqRate			float64				`json:"failed_fetch_req_rate"`
	FailedProdReqRate			float64				`json:"failed_prod_req_rate"`

	ByteInRate 					map[int64]int64		`json:"byte_in_rate"`
	ByteOutRate					map[int64]int64		`json:"byte_out_rate"`
}