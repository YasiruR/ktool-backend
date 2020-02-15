package kafka

import (
	"context"
	"errors"
	"fmt"
	"github.com/Shopify/sarama"
	"github.com/YasiruR/ktool-backend/log"
)

type topicClient struct {
	Name 		string
	ClusterID 	int
	Mesgs 		[]*sarama.ConsumerMessage
	Errors 		[]*sarama.ConsumerError
}

var TopicClientList []*topicClient

func InitTopicConsumer(ctx context.Context, clusterID int, topic string) (err error) {
	//mesgs := make(chan *sarama.ConsumerMessage)
	//errors := make(chan *sarama.ConsumerError)

	var t topicClient
	t.Name = topic
	t.ClusterID = clusterID

	client, err := GetClient(ctx, clusterID)
	if err != nil {
		log.Logger.ErrorContext(ctx, "fetching client to get topic data failed")
		return
	}

	partitions, err	:= client.Consumer.Partitions(topic)
	if err != nil {
		log.Logger.ErrorContext(ctx, fmt.Sprintf("getting partitions for topic %v failed", topic), err)
		return
	}

	for _, partition := range partitions {
		consumerPartition, err := client.Consumer.ConsumePartition(topic, partition, sarama.OffsetOldest)
		if err != nil {
			log.Logger.ErrorContext(ctx, fmt.Sprintf("getting consumer partition for topic %v and partition %v failed", topic, partition), err)
			return err
		}

		go func(topic string, consumer sarama.PartitionConsumer) {
			for {
				select {
				case consumerError := <- consumer.Errors():
					//errors <- consumerError
					t.Errors = append(t.Errors, consumerError)
					log.Logger.WarnContext(ctx, fmt.Sprintf("consumer error for topic %v", topic), consumerError)
				case mesg := <- consumer.Messages():
					//mesgs <- mesg
					t.Mesgs = append(t.Mesgs, mesg)
					log.Logger.TraceContext(ctx, fmt.Sprintf("consumed message for topic %v", topic), mesg)
				}
			}
		}(topic, consumerPartition)
	}

	TopicClientList = append(TopicClientList, &t)

	return nil
}

func GetTopicData(ctx context.Context, clusterID int, topic string, start, end int) (mesgs []*sarama.ConsumerMessage, err error) {
	for _, topicClient := range TopicClientList {
		if topicClient.Name == topic && topicClient.ClusterID == clusterID {
			mesgs = topicClient.Mesgs[start:end]
			goto returnMesgs
		}
	}

	log.Logger.ErrorContext(ctx, fmt.Sprintf("could not find the topic %v", topic))
	return nil, errors.New("could not find the topic")

returnMesgs:
	log.Logger.TraceContext(ctx, fmt.Sprintf("topic data fetched for topic %v in cluster %v", topic, clusterID))
	return mesgs, nil
}