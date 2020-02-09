package kafka

import (
	"context"
	"errors"
	"fmt"
	"github.com/Shopify/sarama"
	"github.com/YasiruR/ktool-backend/log"
)

func InitClusterConfig(ctx context.Context, brokers []string) (consumer sarama.Consumer, err error) {
	config := sarama.NewConfig()
	config.Consumer.Return.Errors = true

	consumer, err = sarama.NewConsumer(brokers, config)
	if err != nil {
		log.Logger.ErrorContext(ctx, err, brokers)
		return nil, errors.New("kafka cluster config initialization failed")
	}

	//todo: close cluster connection on disconnect

	return consumer, nil
}

func GetTopicList(ctx context.Context, cluster sarama.Consumer) (topics []string, err error) {
	topics, err = cluster.Topics()
	if err != nil {
		log.Logger.ErrorContext(ctx, err, cluster)
		return nil, err
	}

	log.Logger.TraceContext(ctx, "all topics are fetched", fmt.Sprintf("no of topics : %v", len(topics)), fmt.Sprintf("cluster : %v", cluster))
	return topics, nil
}

func InitClient(ctx context.Context, brokers []string) (client sarama.Client, err error) {
	client, err = sarama.NewClient(brokers, nil)
	if err != nil {
		log.Logger.ErrorContext(ctx, fmt.Sprintf("creating new client failed for brokers : %v", brokers), err)
		return nil, err
	}

	log.Logger.TraceContext(ctx, "client initialized successfully", brokers)

	return client, nil
}

func GetBrokerAddrList(ctx context.Context, client sarama.Client) (addrList []string, err error) {
	brokers := client.Brokers()
	for _, broker := range brokers {
		addrList = append(addrList, broker.Addr())
	}

	log.Logger.TraceContext(ctx, "broker address list fetched")
	return addrList, nil
}