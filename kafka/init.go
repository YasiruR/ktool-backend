package kafka

import (
	"fmt"
	"github.com/YasiruR/ktool-backend/database"
	"github.com/YasiruR/ktool-backend/domain"
	"github.com/YasiruR/ktool-backend/log"
	"github.com/google/uuid"
	traceable_context "github.com/pickme-go/traceable-context"
	"strconv"
)

var (
	ClusterList 			[]domain.KCluster
)

func InitAllClusters() {
	ctx := traceable_context.WithUUID(uuid.New())
	clusterList, err := database.GetAllClusters(ctx)
	if err != nil {
		log.Logger.Fatal("initializing clusters failed")
	}

	var tempClustList []domain.KCluster

	clusterLoop:
	for _, cluster := range clusterList {
		var brokerList []string
		var clustClient domain.KCluster
		clustClient.ClusterID = cluster.ID
		clustClient.ClusterName = cluster.ClusterName

		brokers, err := database.GetBrokersByClusterId(ctx, cluster.ID)
		if err != nil {
			log.Logger.ErrorContext(ctx, "fetching brokers failed for cluster", cluster.ClusterName)
			clustClient.Available = false
			tempClustList = append(tempClustList, clustClient)
			continue
		}

		for _, broker := range brokers {
			addr := broker.Host + ":" + strconv.Itoa(broker.Port)
			brokerList = append(brokerList, addr)
		}

		client, err := InitClient(ctx, brokerList)
		if err != nil {
			log.Logger.ErrorContext(ctx, "client could not be initialized for cluster", cluster.ClusterName, err)
			clustClient.Available = false
			tempClustList = append(tempClustList, clustClient)
			continue
		}

		saramaBrokers := client.Brokers()
		clustClient.Brokers = saramaBrokers

		saramaConsumer, err := InitClusterConfig(ctx, cluster.ClusterName, brokerList, "")
		if err != nil {
			log.Logger.ErrorContext(ctx,"cluster config could not be initialized for cluster", cluster.ClusterName, err)
			clustClient.Available = false
			tempClustList = append(tempClustList, clustClient)
			continue
		}

		topics, err := GetTopicList(ctx, saramaConsumer)
		if err != nil {
			log.Logger.ErrorContext(ctx, "topic list could not be fetched", cluster.ClusterName)
			clustClient.Available = false
			tempClustList = append(tempClustList, clustClient)
			continue
		}

		for _, topic := range topics {
			var clusterTopic domain.KTopic
			clusterTopic.Name = topic
			clusterTopic.Partitions, err = saramaConsumer.Partitions(topic)
			if err != nil {
				log.Logger.Error(fmt.Sprintf("partitions could not be fetched for %v topic in %v cluster", topic, cluster.ClusterName), err)
				clustClient.Available = false
				tempClustList = append(ClusterList, clustClient)
				continue clusterLoop
			}
			clustClient.Topics = append(clustClient.Topics, clusterTopic)
		}

		clustClient.Consumer = saramaConsumer
		clustClient.Client = client
		clustClient.Available = true
		tempClustList = append(tempClustList, clustClient)
	}

	ClusterList = tempClustList

	log.Logger.Trace("cluster initialization completed", fmt.Sprintf("No. of clusters : %v", len(ClusterList)))
}
