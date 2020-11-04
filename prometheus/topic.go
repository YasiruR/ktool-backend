package prometheus

import (
	"context"
	"github.com/YasiruR/ktool-backend/kafka"
	"github.com/YasiruR/ktool-backend/log"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

const (
	typeMessage = "message"
	typeBytesIn = "bytesIn"
	typeBytesOut = "bytesOut"
	typeBytesRej = "bytesRej"
	typeReplBytesIn = "replBytesIn"
	typeReplBytesOut = "replBytesOut"
	typeMessageRate = "messageRate"
	typeTotalMessages = "totalMessages"
	typeBytesInRate = "bytesInRate"
	typeBytesOutRate = "bytesOutRate"
)

var (
	topicQueryMap = map[string]string{
		typeMessage: "query?query=sum%20by%20(topic%2C%20job)(kafka_server_brokertopicmetrics_messagesin_total)&time=",
		typeBytesIn: "query?query=sum%20by%20(topic%2C%20job)(kafka_server_brokertopicmetrics_bytesin_total)&time=",
		typeBytesOut: "query?query=sum%20by%20(topic%2C%20job)(kafka_server_brokertopicmetrics_bytesout_total)&time=",
		typeBytesRej: "query?query=sum%20by%20(topic%2C%20job)(kafka_server_brokertopicmetrics_bytesrejected_total)&time=",
		typeReplBytesIn: "query?query=sum%20by%20(topic%2C%20job)(kafka_server_brokertopicmetrics_replicationbytesin_total)&time=",
		typeReplBytesOut: "query?query=sum%20by%20(topic%2C%20job)(kafka_server_brokertopicmetrics_replicationbytesout_total)&time=",
	}

	topicSummaryQueryMap = map[string]string {
		typeTotalMessages: "query?query=sum%20(kafka_server_brokertopicmetrics_messagesin_total)%20by%20(job%2C%20topic)&time=",
		typeMessageRate: "query?query=sum%20(irate(kafka_server_brokertopicmetrics_messagesin_total%5B1m%5D))%20by%20(job%2C%20topic)&time=",
		typeBytesInRate: "query?query=sum%20(irate(kafka_server_brokertopicmetrics_bytesin_total%5B1m%5D))%20by%20(job%2C%20topic)&time=",
		typeBytesOutRate: "query?query=sum%20(irate(kafka_server_brokertopicmetrics_bytesout_total%5B1m%5D))%20by%20(job%2C%20topic)&time=",
	}

	PromClusterTopicMap map[int]map[string]topicMetrics
	PromSummaryMap map[int]summaryMetrics
	//{cluster: {topic1: 123, topic2: 654}, ..}
	messageMap, bytesInMap, bytesOutMap, bytesRejMap, replBytesInMap, replBytesOutMap map[int]map[string]int
	totalMessagesMap, messageRateMap, bytesInRateMap, bytesOutRateMap map[int]int
	topicBrokerMap map[int]map[string][]string
)

type topicMetrics struct {
	Brokers 		[]string
	Messages 		int
	BytesIn 		int
	BytesOut 		int
	BytesRejected	int
	ReplBytesIn		int
	ReplBytesOut	int
}

type summaryMetrics struct {
	TotalMessages 	int
	MessageRate		int
	BytesInRate		int
	BytesOutRate	int
}

func InitTopicMetrics(ctx context.Context, osChannel chan os.Signal) {
	wg := &sync.WaitGroup{}
	currentTime := time.Now().Unix()
	for t, query := range topicQueryMap {
		wg.Add(1)
		go func(t string, query string) {
			select {
			case <-osChannel:
				log.Logger.InfoContext(ctx, "terminating topic metrics go routines")
				return
			default:
				req := promUrl + query + strconv.Itoa(int(currentTime))
				tmpMap, err := getMetricsByTopic(ctx, req)
				if err == nil {
					switch t {
					case typeMessage:
						topicBrokerMap, _ = getTopicBrokerMap(ctx, promUrl + "query?query=sum%20(kafka_server_brokertopicmetrics_messagesin_total)%20by%20(topic%2C%20instance%2C%20job)&time=" + strconv.Itoa(int(currentTime)))
						messageMap = tmpMap
					case typeBytesIn:
						bytesInMap = tmpMap
					case typeBytesOut:
						bytesOutMap = tmpMap
					case typeBytesRej:
						bytesRejMap = tmpMap
					case typeReplBytesIn:
						replBytesInMap = tmpMap
					case typeReplBytesOut:
						replBytesOutMap = tmpMap
					}
				}
				wg.Done()
			}
		}(t, query)
	}

	wg.Wait()
	tmpPromMap := make(map[int]map[string]topicMetrics)
	for _, cluster := range kafka.ClusterList {
		tmpPromMap[cluster.ClusterID] = make(map[string]topicMetrics)
		for _, topic := range cluster.Topics {
			tmpTopicMetrics := topicMetrics{}
			tmpTopicMetrics.Brokers = topicBrokerMap[cluster.ClusterID][topic.Name]
			tmpTopicMetrics.Messages = messageMap[cluster.ClusterID][topic.Name]
			tmpTopicMetrics.BytesIn = bytesInMap[cluster.ClusterID][topic.Name]
			tmpTopicMetrics.BytesOut = bytesOutMap[cluster.ClusterID][topic.Name]
			tmpTopicMetrics.BytesRejected = bytesRejMap[cluster.ClusterID][topic.Name]
			tmpTopicMetrics.ReplBytesIn = replBytesInMap[cluster.ClusterID][topic.Name]
			tmpTopicMetrics.ReplBytesOut = replBytesOutMap[cluster.ClusterID][topic.Name]

			tmpPromMap[cluster.ClusterID][topic.Name] = tmpTopicMetrics
		}
	}
	PromClusterTopicMap = tmpPromMap
}

func InitSummaryMetrics(ctx context.Context, osChannel chan os.Signal) {
	wg := &sync.WaitGroup{}
	currentTime := time.Now().Unix()
	for t, query := range topicSummaryQueryMap {
		wg.Add(1)
		go func(t string, query string) {
			select {
			case <-osChannel:
				log.Logger.InfoContext(ctx, "terminating topic metrics go routines")
				return
			default:
				req := promUrl + query + strconv.Itoa(int(currentTime))
				tmpMap, _ := getSummaryMetrics(ctx, req)
				switch t {
				case typeTotalMessages:
					totalMessagesMap = tmpMap
				case typeMessageRate:
					messageRateMap = tmpMap
				case typeBytesInRate:
					bytesInRateMap = tmpMap
				case typeBytesOutRate:
					bytesOutRateMap = tmpMap
				}
				wg.Done()
			}
		}(t, query)
	}

	wg.Wait()
	tmpPromSummaryMap := make(map[int]summaryMetrics)
	for _, cluster := range kafka.ClusterList {
		tmpClusterMap := summaryMetrics{}
		tmpClusterMap.TotalMessages = totalMessagesMap[cluster.ClusterID]
		tmpClusterMap.BytesInRate = bytesInRateMap[cluster.ClusterID]
		tmpClusterMap.BytesOutRate = bytesOutRateMap[cluster.ClusterID]
		tmpClusterMap.MessageRate = messageRateMap[cluster.ClusterID]
		tmpPromSummaryMap[cluster.ClusterID] = tmpClusterMap
	}
	PromSummaryMap = tmpPromSummaryMap
}

func getTopicBrokerMap(ctx context.Context, req string) (tmpMap map[int]map[string][]string, err error) {
	tmpMap = make(map[int]map[string][]string)

	metrics, err := getResponseByEndpoint(ctx, req)
	if err != nil {
		log.Logger.ErrorContext(ctx, err, "failed getting broker metrics by topic", req)
		return tmpMap, err
	}

	outerLoop:
	for _, res := range metrics.Data.Result {
		if res.Metric.Topic == "" {
			continue
		}
		for _, cluster := range kafka.ClusterList {
			if cluster.ClusterName == res.Metric.Job {
				_, ok := tmpMap[cluster.ClusterID]
				if !ok {
					tmpTopicMap := make(map[string][]string)
					tmpTopicMap[res.Metric.Topic] = append(tmpTopicMap[res.Metric.Topic], res.Metric.Instance)
					tmpMap[cluster.ClusterID] = tmpTopicMap
					continue outerLoop
				}
				tmpMap[cluster.ClusterID][res.Metric.Topic] = append(tmpMap[cluster.ClusterID][res.Metric.Topic], res.Metric.Instance)
				continue outerLoop
			}
			log.Logger.ErrorContext(ctx, "could not find cluster in prom response", cluster.ClusterName)
		}
	}

	return tmpMap,nil
}

func getSummaryMetrics(ctx context.Context, req string) (tmpMap map[int]int, err error) {
	tmpMap = make(map[int]int)

	metrics, err := getResponseByEndpoint(ctx, req)
	if err != nil {
		log.Logger.ErrorContext(ctx, err, "failed getting summary metrics by topic", req)
		return tmpMap, err
	}

	outerLoop:
	for _, res := range metrics.Data.Result {
		if res.Metric.Topic == "" {
			strVal, ok := res.Value[1].(string)
			if ok {
				intVal, err := strconv.Atoi(strVal)
				if err != nil {
					floatVal, err := strconv.ParseFloat(strVal, 64)
					if err != nil {
						log.Logger.ErrorContext(ctx, err, "res is not either int or float", res)
						continue
					}
					//check for message metrics whether it contains null topic
					for _, cluster := range kafka.ClusterList {
						if cluster.ClusterName ==  res.Metric.Job {
							tmpMap[cluster.ClusterID] = int(floatVal)
							continue outerLoop
						}
					}
				} else {
					//check for message metrics whether it contains null topic
					for _, cluster := range kafka.ClusterList {
						if cluster.ClusterName ==  res.Metric.Job {
							tmpMap[cluster.ClusterID] = intVal
							continue outerLoop
						}
					}
				}
				continue
			}
			log.Logger.ErrorContext(ctx, "fetched topic value is not a string", res)
		}
	}

	//log.Logger.TraceContext(ctx, "fetching topic summary metrics done")
	return tmpMap,nil
}

func getMetricsByTopic(ctx context.Context, req string) (tmpMap map[int]map[string]int, err error) {
	tmpMap = make(map[int]map[string]int)

	metrics, err := getResponseByEndpoint(ctx, req)
	if err != nil {
		log.Logger.ErrorContext(ctx, err, "failed getting metrics by topic", req)
		return tmpMap, err
	}

	outerLoop:
	for _, res := range metrics.Data.Result {
		if strings.Contains(req, "kafka_server_brokertopicmetrics_messagesin_total") {
			if res.Metric.Topic == "" {
				continue
			}
		}
		strVal, ok := res.Value[1].(string)
		if ok {
			intVal, err := strconv.Atoi(strVal)
			if err != nil {
				floatVal, err := strconv.ParseFloat(strVal, 64)
				if err != nil {
					log.Logger.ErrorContext(ctx, err, "res is not either int or float", res)
					continue
				}
				for _, cluster := range kafka.ClusterList {
					if cluster.ClusterName ==  res.Metric.Job {
						_, ok := tmpMap[cluster.ClusterID]
						if !ok {
							tmpMap[cluster.ClusterID] = make(map[string]int)
						}
						tmpMap[cluster.ClusterID][res.Metric.Topic] = int(floatVal)
						continue outerLoop
					}
				}
			} else {
				//check for message metrics whether it contains null topic
				for _, cluster := range kafka.ClusterList {
					if cluster.ClusterName ==  res.Metric.Job {
						_, ok := tmpMap[cluster.ClusterID]
						if !ok {
							tmpMap[cluster.ClusterID] = make(map[string]int)
						}
						tmpMap[cluster.ClusterID][res.Metric.Topic] = intVal
						continue outerLoop
					}
				}
			}
			continue
		}
		log.Logger.ErrorContext(ctx, "fetched topic value is not a string", res)
	}

	//log.Logger.TraceContext(ctx, "fetching topic metrics done")
	return tmpMap,nil
}

//func setBrokerTopicMap(ctx context.Context, req string) (err error) {
//	tmpMap := make(map[string][]string)
//
//	metrics, err := getResponseByEndpoint(ctx, req)
//	if err != nil {
//		log.Logger.ErrorContext(ctx, err, "failed getting broker metrics by topic", req)
//		return err
//	}
//
//	for _, res := range metrics.Data.Result {
//		if res.Metric.Topic == "" {
//			continue
//		}
//		tmpMap[res.Metric.Instance] = append(tmpMap[res.Metric.Instance], res.Metric.Topic)
//	}
//
//	BrokerTopicMap = tmpMap
//	return nil
//}
