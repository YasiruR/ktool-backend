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
	ClusterID 		int
	ClusterName 	string
	Consumer  		sarama.Consumer
	Client    		sarama.Client
	Brokers 		[]*sarama.Broker
	Topics 			[]KTopic
	Available 		bool
	BrokerOverview	ClusterOverview
}

type ClusterOverview struct {
	TotalBrokers              	int     					`json:"brokers"`
	TotalPartitions           	int     					`json:"partitions"`
	TotalTopics               	int     					`json:"topics"`
	TotalReplicas             	int     					`json:"replicas"`
	UnderReplicatedPartitions 	int     					`json:"under_replicated_partitions"`
	OfflinePartitions			int							`json:"offline_partitions"`
	OfflineReplicas				int							`json:"offline_replicas"`
	TotalOutgoingRate       	float64 					`json:"total_outgoing_rate"`		//in kb
	TotalIncomingRate      		float64 					`json:"total_incoming_rate"`		//in kb
	TotalRequestRate          	float64 					`json:"total_request_rate"`
	TotalRequestSize          	float64 					`json:"total_request_size"`
	TotalRequestLatency       	float64 					`json:"total_request_latency"`
	TotalResponseRate         	float64 					`json:"total_response_rate"`
	TotalResponseSize         	float64 					`json:"total_response_size"`
	ActiveController          	string  					`json:"active_controller"`
	ZookeeperAvail            	bool    					`json:"zookeeper_avail"`
	KafkaVersion              	string  					`json:"kafka_version"`
	Brokers						map[int32]BrokerMetrics		`json:"brokers"`
}

type BrokerMetrics struct {
	IncomingByteRate       	float64 	`json:"incoming_byte_rate"`
	RequestRate            	float64		`json:"request_rate"`
	RequestSize            	float64		`json:"request_size"`
	RequestLatency         	float64		`json:"request_latency"`		//in ms
	OutgoingByteRate       	float64		`json:"outgoing_byte_rate"`
	ResponseRate           	float64		`json:"response_rate"`
	ResponseSize           	float64		`json:"response_size"`
}