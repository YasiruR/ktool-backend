package kafka

import (
	"context"
	"errors"
	"fmt"
	"github.com/Shopify/sarama"
	"github.com/YasiruR/ktool-backend/log"
	"sync"
)

type topicClient struct {
	Name 		string
	ClusterID 	int
	Mesgs 		[]*sarama.ConsumerMessage
	Errors 		[]*sarama.ConsumerError
}

var TopicClientList []*topicClient
var topicListMu	*sync.Mutex

func init() {
	topicListMu = &sync.Mutex{}
}

func InitTopicConsumer(ctx context.Context, clusterID int, topic string) (err error) {
	//mesgs := make(chan *sarama.ConsumerMessage)
	//errors := make(chan *sarama.ConsumerError)

	var t topicClient
	t.Name = topic
	t.ClusterID = clusterID

	msgMu := &sync.Mutex{}
	errMu := &sync.Mutex{}

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
		partitionConsumer, err := client.Consumer.ConsumePartition(topic, partition, sarama.OffsetOldest)
		if err != nil {
			log.Logger.ErrorContext(ctx, fmt.Sprintf("getting consumer partition for topic %v and partition %v failed", topic, partition), err)
			return err
		}

		go func(topic string, consumer sarama.PartitionConsumer) {
			for {
				select {
				case consumerError := <- consumer.Errors():
					//errors <- consumerError
					errMu.Lock()
					t.Errors = append(t.Errors, consumerError)
					errMu.Unlock()
					log.Logger.WarnContext(ctx, fmt.Sprintf("consumer error for topic %v", topic), consumerError)
				case msg := <- consumer.Messages():
					//mesgs <- mesg
					msgMu.Lock()
					t.Mesgs = append(t.Mesgs, msg)
					msgMu.Unlock()
					log.Logger.TraceContext(ctx, fmt.Sprintf("consumed message for topic %v", topic), msg)
				}
			}
		}(topic, partitionConsumer)
	}

	topicListMu.Lock()
	TopicClientList = append(TopicClientList, &t)
	topicListMu.Unlock()

	return nil
}

func ReadMessages(ctx context.Context, start, end int32, topic string, clusterID int) (messages []*sarama.ConsumerMessage, err error) {
	for _, c := range TopicClientList {
		if c.Name == topic {
			if c.ClusterID == clusterID {
				messages = c.Mesgs[start:end]
				log.Logger.TraceContext(ctx, fmt.Sprintf("fetched messages for topic %v from %v to %v", topic, start, end))
				return messages, nil
			}
		}
	}

	log.Logger.ErrorContext(ctx, "topic or cluster id may not exist", topic, clusterID)
	return nil, errors.New(fmt.Sprintf("failed to fetch data for %v topic and cluster id %v", topic, clusterID))
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