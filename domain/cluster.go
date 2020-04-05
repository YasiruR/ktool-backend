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
	BrokerOverview	BrokerOverview
}

type BrokerOverview struct {
	TotalBrokers 			int							`json:"total_brokers"`
	TotalPartitions 		int							`json:"total_partitions"`
	TotalTopics 			int							`json:"total_topics"`
	TotalReplicas 			int							`json:"total_replicas"`
	TotalInsyncReplicas 	int							`json:"total_insync_replicas"`
	TotalOfflineReplicas	int							`json:"total_offline_replicas"`
	TotalProductionRate		float64						`json:"total_partition_rate"`
	TotalConsumptionRate	float64						`json:"total_consumption_rate"`
	ActiveController 		string						`json:"active_controller"`
	ZookeeperAvail			bool						`json:"zookeeper_avail"`
	KafkaVersion 			string						`json:"kafka_version"`
	Brokers					map[int32]BrokerMetrics	`json:"brokers"`
}

type BrokerMetrics struct {
	IncomingByteRate       	float64 	`json:"incoming_byte_rate"`
	RequestRate            	int64 		`json:"request_rate"`
	RequestSize            	int64		`json:"request_size"`
	RequestLatency         	int64		`json:"request_latency"`
	OutgoingByteRate       	float64		`json:"outgoing_byte_rate"`
	ResponseRate           	float64		`json:"response_rate"`
	ResponseSize           	int64		`json:"response_size"`
	BrokerIncomingByteRate 	float64		`json:"broker_incoming_byte_rate"`
	BrokerRequestRate      	float64		`json:"broker_request_rate"`
	BrokerRequestSize      	int64		`json:"broker_request_size"`
	BrokerRequestLatency   	int64		`json:"broker_request_latency"`			//in ms
	BrokerOutgoingByteRate 	float64		`json:"broker_outgoing_byte_rate"`
	BrokerResponseRate     	float64		`json:"broker_response_rate"`
	BrokerResponseSize     	int64		`json:"broker_response_size"`
}