package kafka

import (
	"context"
	"errors"
	"fmt"
	"github.com/Shopify/sarama"
	"github.com/YasiruR/ktool-backend/log"
)

var (
	Cluster 		sarama.Consumer
	Client			sarama.Client
)

func InitClusterConfig(ctx context.Context, brokers []string) (err error) {
	config := sarama.NewConfig()
	config.Consumer.Return.Errors = true

	cluster, err := sarama.NewConsumer(brokers, config)
	if err != nil {
		log.Logger.ErrorContext(ctx, err, brokers)
		return errors.New("kafka cluster config initialization failed")
	}

	//todo: close cluster connection on disconnect

	Cluster = cluster
	return nil
}

func GetTopicList(ctx context.Context) (topics []string, err error) {
	topics, err = Cluster.Topics()
	if err != nil {
		log.Logger.ErrorContext(ctx, err, Cluster)
		return nil, err
	}

	log.Logger.TraceContext(ctx, "all topics are fetched", fmt.Sprintf("no of topics : %v", len(topics)), fmt.Sprintf("cluster : %v", Cluster))
	return topics, nil
}

func InitClient(ctx context.Context, brokers []string) (err error) {
	client, err := sarama.NewClient(brokers, nil)
	if err != nil {
		log.Logger.ErrorContext(ctx, fmt.Sprintf("creating new client failed for brokers : %v", brokers), err)
		return err
	}

	Client = client
	log.Logger.TraceContext(ctx, "client initialized successfully", brokers)

	return nil
}

func GetBrokerAddrList(ctx context.Context) (addrList []string, err error) {
	brokers := Client.Brokers()
	for _, broker := range brokers {
		addrList = append(addrList, broker.Addr())
	}

	log.Logger.TraceContext(ctx, "broker address list fetched")
	return addrList, nil
}