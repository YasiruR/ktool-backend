package domain

type Cluster struct {
	ID 					int64
	ClusterName 		string
	KafkaVersion 		string
	Zookeepers 			[]Zookeeper
	Brokers				[]Broker
	SchemaRegistry		SchemaRegistry
	ActiveControllers	int
}
