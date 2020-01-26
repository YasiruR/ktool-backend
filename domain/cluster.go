package domain

type Cluster struct {
	ID 					int
	ClusterName 		string
	KafkaVersion 		float64
	Zookeepers 			[]Zookeeper
	Brokers				[]Broker
	SchemaRegistry		SchemaRegistry
	ActiveControllers	int
}
