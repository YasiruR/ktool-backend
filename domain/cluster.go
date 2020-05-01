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
	TotalBrokers              	int     					`json:"brokers"`
	TotalPartitions           	int     					`json:"partitions"`
	TotalTopics               	int     					`json:"topics"`
	TotalReplicas             	int     					`json:"replicas"`
	UnderReplicatedPartitions 	int     					`json:"under_replicated_partitions"`
	OfflinePartitions			int							`json:"offline_partitions"`
	OfflineReplicas				int							`json:"offline_replicas"`
	TotalOutgoingRate       	float64 					`json:"total_outgoing_rate"`		//in kb
	TotalIncomingRate      		float64 					`json:"total_incoming_rate"`		//in kb
	TotalMesgByteInRate			map[int64]int64				`json:"total_mesg_byte_in_rate"`
	TotalMesgByteOutRate		map[int64]int64				`json:"total_mesg_byte_out_rate"`
	ActiveController          	string  					`json:"active_controller"`
	ZookeeperAvail            	bool    					`json:"zookeeper_avail"`
	KafkaVersion              	string  					`json:"kafka_version"`
	Brokers						[]BrokerMetrics				`json:"brokers"`
}

type BrokerMetrics struct {
	Host 						string						`json:"host"`
	Port 						int							`json:"port"`
	MesgInByteRate 				map[int64]int64				`json:"mesg_in_byte_rate"`
	MesgOutByteRate				map[int64]int64				`json:"mesg_out_byte_rate"`
}