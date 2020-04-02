package domain

import (
	"github.com/Shopify/sarama"
)

var LoggedInUsers []User

type User struct {
	Id 					int
	Username			string
	Token				string
	AccessLevel 		int
	ConnectedClusters	[]KCluster
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

type KTopic struct {
	Name 		string
	Partitions 	[]int32
}
