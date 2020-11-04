package http

import (
	"encoding/json"
	"github.com/YasiruR/ktool-backend/database"
	"github.com/YasiruR/ktool-backend/domain"
	"github.com/YasiruR/ktool-backend/log"
	"github.com/YasiruR/ktool-backend/prometheus"
	"github.com/google/uuid"
	traceable_context "github.com/pickme-go/traceable-context"
	"net/http"
	"strconv"
	"strings"
)

func handleGetTopicMetrics(res http.ResponseWriter, req *http.Request) {
	ctx := traceable_context.WithUUID(uuid.New())
	//user validation by token header
	tokenHeader := req.Header.Get("Authorization")
	if len(strings.Split(tokenHeader, "Bearer")) < 2 {
		log.Logger.ErrorContext(ctx, "token format is invalid", tokenHeader)
		res.WriteHeader(http.StatusBadRequest)
		return
	}

	token := strings.TrimSpace(strings.Split(tokenHeader, "Bearer")[1])
	userID, ok, err := database.ValidateUserByToken(ctx, token)
	if !ok {
		log.Logger.DebugContext(ctx, "invalid user", token)
		res.WriteHeader(http.StatusUnauthorized)
		return
	}
	if err != nil {
		log.Logger.ErrorContext(ctx, "error occurred in token validation", err)
		res.WriteHeader(http.StatusInternalServerError)
		return
	}

	clusterID, err := strconv.Atoi(req.FormValue("cluster_id"))
	if err != nil {
		log.Logger.ErrorContext(ctx, "conversion of cluster id from string into int failed", err, req.FormValue("cluster_id"))
		res.WriteHeader(http.StatusBadRequest)
		return
	}

	_, ok = domain.LoggedInUserMap[userID]
	if ok {
		var metricsRes topicMetricsRes
		promCluster, promClusterExists := prometheus.PromClusterTopicMap[clusterID]
		promSummary, promSummaryClustExists := prometheus.PromSummaryMap[clusterID]

		log.Logger.TraceContext(ctx, "prom summary", promSummary)

		saramaTopics, saramaClusterExists := domain.ClusterTopicMap[clusterID]
		if promClusterExists && saramaClusterExists && promSummaryClustExists {
			//adding summary metrics
			metricsRes.TotalMessages = promSummary.TotalMessages
			metricsRes.MessageRate = promSummary.MessageRate
			metricsRes.BytesInRate = promSummary.BytesInRate
			metricsRes.BytesOutRate = promSummary.BytesOutRate

			for _, topic := range saramaTopics {
				var topicRes metricsTopic
				promTopic, topicExists := promCluster[topic.Name]
				if topicExists {
					topicRes.Brokers = promTopic.Brokers //to display no of brokers and for filter
					topicRes.Messages = promTopic.Messages
					topicRes.BytesIn = promTopic.BytesIn
					topicRes.BytesOut = promTopic.BytesOut
					topicRes.BytesRejected = promTopic.BytesRejected
					topicRes.ReplBytesIn = promTopic.ReplBytesIn
					topicRes.ReplBytesOut = promTopic.ReplBytesOut
				} else {
					log.Logger.ErrorContext(ctx, "could not find in prometheus topic map", topic.Name, clusterID)
				}

				topicRes.Name = topic.Name
				topicRes.WritablePartitions = len(topic.WritablePartitions)

				for _, partition := range topic.Partitions {
					var topicPartitionRes topicPartition
					topicPartitionRes.ID = int(partition.ID)
					topicPartitionRes.FirstOffset = int(partition.FirstOffset)			//show on hover
					topicPartitionRes.LastOffset = int(partition.NextOffset) - 1
					topicRes.Replicas += len(partition.Replicas)
					topicRes.InSyncReplicas += len(partition.InSyncReplicas)
					topicRes.OfflineReplicas += len(partition.OfflineReplicas)
					if partition.UnderReplicated {
						topicRes.UnderReplicatedPart += 1
					}
					topicRes.Partitions = append(topicRes.Partitions, topicPartitionRes)
				}
				metricsRes.Topics = append(metricsRes.Topics, topicRes)
			}

			res.WriteHeader(http.StatusOK)
			err = json.NewEncoder(res).Encode(metricsRes)
			if err != nil {
				log.Logger.ErrorContext(ctx, err, "marshalling response for get topic metrics request failed", clusterID, metricsRes)
				res.WriteHeader(http.StatusInternalServerError)
				return
			}
			log.Logger.TraceContext(ctx, "topic metrics fetched successfully", clusterID, len(metricsRes.Topics))
		} else {
			log.Logger.ErrorContext(ctx, "requested cluster id is not present in either maps", clusterID, prometheus.PromClusterTopicMap, domain.ClusterTopicMap)
			res.WriteHeader(http.StatusBadRequest)
		}
	} else {
		log.Logger.ErrorContext(ctx, "could not find a user from the logged in user list from token", token)
		res.WriteHeader(http.StatusForbidden)
	}
}

