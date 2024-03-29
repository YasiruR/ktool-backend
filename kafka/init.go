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

//taken from sarama library for histogram sample
const (
	metricsReservoirSize 	= 	1028
	metricsAlphaFactor   	= 	0.015
)

func init() {
	domain.LoggedInUserMap = make(map[int]domain.User)
}

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

		config, err := InitSaramaConfig(ctx, cluster.ClusterName, "")
		if err != nil {
			log.Logger.ErrorContext(ctx, "initializing sarama config failed and may proceed with default config for consumer and client init", cluster.ClusterName)
		}

		client, err := InitClient(ctx, brokerList, config)
		if err != nil {
			log.Logger.ErrorContext(ctx, "client could not be initialized for cluster", cluster.ClusterName, err)
			clustClient.Available = false
			tempClustList = append(tempClustList, clustClient)
			continue
		}

		saramaBrokers := client.Brokers()
		clustClient.Brokers = saramaBrokers

		saramaConsumer, err := InitSaramaConsumer(ctx, brokerList, config)
		if err != nil {
			log.Logger.ErrorContext(ctx, err,"cluster config could not be initialized for cluster", cluster.ClusterName)
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

		var numOfLeaders, numOfReplicas, numOfOfflineRepl, numOfInSyncRepl, numOfOnlinePartitions int
		for _, topic := range topics {
			var clusterTopic domain.KTopic
			clusterTopic.Name = topic
			clusterTopic.Partitions, err = saramaConsumer.Partitions(topic)
			if err != nil {
				log.Logger.ErrorContext(ctx, err, fmt.Sprintf("partitions could not be fetched for %v topic in %v cluster", topic, cluster.ClusterName))
				clustClient.Available = false
				tempClustList = append(ClusterList, clustClient)
				continue clusterLoop
			}
			numOfLeaders += len(clusterTopic.Partitions)
			clustClient.Topics = append(clustClient.Topics, clusterTopic)

			onlineReplicas, err := client.WritablePartitions(topic)
			if err != nil {
				log.Logger.ErrorContext(ctx, err, fmt.Sprintf("online partitions could not be fetched for %v topic in %v cluster", topic, cluster.ClusterName))
			} else {
				numOfOnlinePartitions += len(onlineReplicas)
			}

			//to fetch information about replicas
			partitionLoop:
			for _, partitionID := range clusterTopic.Partitions {
				replicas, err := client.Replicas(clusterTopic.Name, partitionID)
				if err != nil {
					log.Logger.ErrorContext(ctx, err, fmt.Sprintf("replicas could not be fetched for %v topic and %v paritition in %v cluster", topic, partitionID, cluster.ClusterName))
					continue partitionLoop
				}
				numOfReplicas += len(replicas)

				inSyncReplicas, err := client.InSyncReplicas(clusterTopic.Name, partitionID)
				if err != nil {
					log.Logger.ErrorContext(ctx, err, fmt.Sprintf("insync replicas could not be fetched for %v topic and %v paritition in %v cluster", topic, partitionID, cluster.ClusterName))
					continue partitionLoop
				}
				numOfInSyncRepl += len(inSyncReplicas)

				offlineReplicas, err := client.OfflineReplicas(clusterTopic.Name, partitionID)
				if err != nil {
					log.Logger.ErrorContext(ctx, err, fmt.Sprintf("offline replicas could not be fetched for %v topic and %v paritition in %v cluster", topic, partitionID, cluster.ClusterName))
					continue partitionLoop
				}
				numOfOfflineRepl += len(offlineReplicas)
			}
		}

		//getting cluster controller id
		controller, err := client.Controller()
		if err != nil {
			log.Logger.ErrorContext(ctx, err, fmt.Sprintf("fetching controller id for the cluster %v failed", cluster.ClusterName))
		} else {
			clustClient.ClusterOverview.ActiveController = controller.Addr()
		}

		//updating all collected broker metrics to the cluster
		clustClient.ClusterOverview.TotalLeaders = numOfLeaders
		clustClient.ClusterOverview.TotalTopics = len(topics)
		clustClient.ClusterOverview.TotalReplicas = numOfReplicas
		clustClient.ClusterOverview.OfflineReplicas = numOfOfflineRepl
		clustClient.ClusterOverview.ActiveBrokers = len(saramaBrokers)
		clustClient.ClusterOverview.UnderReplicatedPartitions = numOfReplicas - numOfInSyncRepl			//redundant since now val is fetched by sum in broker metrics api
		clustClient.ClusterOverview.OfflinePartitions = numOfLeaders - numOfOnlinePartitions			//redundant since now val is fetched by sum in broker metrics api

		clustClient.Consumer = saramaConsumer
		clustClient.Client = client
		clustClient.Available = true
		tempClustList = append(tempClustList, clustClient)
	}

	ClusterList = tempClustList

	//log.Logger.Trace("cluster initialization completed", fmt.Sprintf("No. of clusters : %v", len(ClusterList)))
}
