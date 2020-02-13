package kafka

import (
	"fmt"
	"github.com/Shopify/sarama"
	"github.com/YasiruR/ktool-backend/database"
	"github.com/YasiruR/ktool-backend/log"
	"github.com/google/uuid"
	traceable_context "github.com/pickme-go/traceable-context"
	"strconv"
)

var (
	ClusterList 			[]KCluster
	SelectedClusterList 	[]KCluster
)

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

func InitAllClusters() {
	ctx := traceable_context.WithUUID(uuid.New())
	clusterList, err := database.GetAllClusters(ctx)
	if err != nil {
		log.Logger.Fatal("initializing clusters failed")
	}

	clusterLoop:
	for _, cluster := range clusterList {
		var brokerList []string
		var clustClient KCluster
		clustClient.ClusterID = cluster.ID
		clustClient.ClusterName = cluster.ClusterName

		for _, broker := range cluster.Brokers {
			addr := broker.Host + strconv.Itoa(broker.Port)
			brokerList = append(brokerList, addr)
		}

		client, err := InitClient(ctx, brokerList)
		if err != nil {
			log.Logger.Error("client could not be initialized for cluster", cluster.ClusterName, err)
			ClusterList = append(ClusterList, clustClient)
			continue
		}

		saramaBrokers := client.Brokers()
		clustClient.Brokers = saramaBrokers

		saramaConsumer, err := InitClusterConfig(ctx, cluster.ClusterName, brokerList, "")
		if err != nil {
			log.Logger.Error("cluster config could not be initialized for cluster", cluster.ClusterName, err)
			ClusterList = append(ClusterList, clustClient)
			continue
		}

		topics, err := GetTopicList(ctx, saramaConsumer)
		if err != nil {
			log.Logger.Error("topic list could not be fetched from ")
			ClusterList = append(ClusterList, clustClient)
			continue
		}

		for _, topic := range topics {
			var clusterTopic KTopic
			clusterTopic.Name = topic
			clusterTopic.Partitions, err = saramaConsumer.Partitions(topic)
			if err != nil {
				log.Logger.Error(fmt.Sprintf("partitions could not be fetched for %v topic in %v cluster", topic, cluster.ClusterName), err)
				ClusterList = append(ClusterList, clustClient)
				continue clusterLoop
			}
			clustClient.Topics = append(clustClient.Topics, clusterTopic)
		}

		clustClient.Consumer = saramaConsumer
		clustClient.Client = client
		clustClient.Available = true
		ClusterList = append(ClusterList, clustClient)
	}

	log.Logger.Trace("cluster initialization completed")
}
