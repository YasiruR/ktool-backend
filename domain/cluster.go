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
	ClusterID 	int
	ClusterName string
	Consumer  	sarama.Consumer
	Client    	sarama.Client
	Brokers 	[]*sarama.Broker
	Topics 		[]KTopic
	Available 	bool
}

//type BrokerOverview struct {
//	TotalBrokers 			int
//	TotalProductionRate		float64
//	TotalConsumptionRate	float64
//	ActiveController		string
//	ZookeeperAvail			bool
//	TotalPartitions 		int
//	TotalReplicas 			int
//	Brokers 				[]Broker
//}