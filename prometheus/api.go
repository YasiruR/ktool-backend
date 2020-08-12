package prometheus

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/YasiruR/ktool-backend/domain"
	"github.com/YasiruR/ktool-backend/log"
	"io/ioutil"
	"net/http"
	"strconv"
)

func GetMetricsForRange(ctx context.Context, metricsName, host string, port, from, to int, cluster string) (metrics domain.PromResponse, err error) {

	//get step from 'from' and 'to'
	step, err := getStepFromTs(ctx, from, to)
	if err != nil {
		return metrics, err
	}

	var req string
	//pure queries
	if metricsName == partitionsQuery || metricsName == leadersQuery || metricsName == activeContQuery || metricsName == offlinePartQuery || metricsName == underReplQuery || metricsName == isrExpQuery || metricsName == isrShrinkQuery || metricsName == netProcIdleQuery || metricsName == maxLagQuery || metricsName == uncleanLeadElQuery {
		req = promUrl + "query_range?query=" + metricsName + "%7Binstance%3D%22" + host + "%3A" + strconv.Itoa(port) + "%22%7D&start=" + strconv.Itoa(from) + "&end=" + strconv.Itoa(to) + "&step=" + strconv.Itoa(step)
	//float and int updated queries sum by instance
	} else if metricsName == msgRateQuery || metricsName == resTimeQuery || metricsName == queueTimeQuery || metricsName == localTimeQuery || metricsName == remoteTimeQuery || metricsName == totalTimeQuery || metricsName == failedFetchQuery || metricsName == failedProdQuery || metricsName == bytesInQuery || metricsName == bytesOutQuery {
		req = promUrl + "query_range?query=rate(" + metricsName + "%7Binstance%3D%22" + host + "%3A" + strconv.Itoa(port) + "%22%7D%5B1m%5D)&start=" + strconv.Itoa(from) + "&end=" + strconv.Itoa(to) + "&step=" + strconv.Itoa(step)
	//total messages query
	} else if metricsName == totalMsgQuery {
		req = promUrl + "query_range?query=rate(" + metricsName + "%7Binstance%3D%22" + host + "%3A" + strconv.Itoa(port) + "%22%7D%5B1m%5D)&start=" + strconv.Itoa(from) + "&end=" + strconv.Itoa(to) + "&step=" + strconv.Itoa(step)
	//int updated queries count by instance
	} else if metricsName == totalTopicsQuery {
		req = promUrl + "query_range?query=count%20by%20(instance)%20(rate(" + metricsName + "%7Binstance%3D%22" + host + "%3A" + strconv.Itoa(port) + "%22%7D%5B1m%5D))&start=" + strconv.Itoa(from) + "&end=" + strconv.Itoa(to) + "&step=" + strconv.Itoa(step)
	} else if metricsName == totalBytesIn {
		req = promUrl + "query_range?query=sum%20by%20(job)%20(rate(" + bytesInQuery + "%7Bjob%3D%22" + cluster + "%22%7D%5B1m%5D))&start=" + strconv.Itoa(from) + "&end=" + strconv.Itoa(to) + "&step=" + strconv.Itoa(step)
	} else if metricsName == totalBytesOut {
		req = promUrl + "query_range?query=sum%20by%20(job)%20(rate(" + bytesOutQuery + "%7Bjob%3D%22" + cluster + "%22%7D%5B1m%5D))&start=" + strconv.Itoa(from) + "&end=" + strconv.Itoa(to) + "&step=" + strconv.Itoa(step)
	} else {
		log.Logger.ErrorContext(ctx, "undefined metrics name", metricsName)
		return metrics, errors.New("undefined metrics")
	}

	metrics, err = getResponseByEndpoint(ctx, req)
	if err != nil {
		log.Logger.DebugContext(ctx, req)
		return
	}
	return
}

func GetAllMetricsByTimestamp(ctx context.Context, host string, port, ts int, isFromTs bool) (brokerMetrics domain.BrokerMetrics, err error) {
	for key, _ := range queryReqList {
		var req string
		switch key {
		case partitions:
			//since this metrics will not be used in the ui
			if isFromTs {
				continue
			}
			req = promUrl + "query?query=kafka_server_replicamanager_partitioncount%7Binstance%3D%27" + host + "%3A" + strconv.Itoa(port) + "%27%7D&time=" + strconv.Itoa(ts)
			val, _, err := getValueByEndpoint(ctx, req, true)
			if err != nil {
				//log.Logger.DebugContext(ctx, "error fetching metrics", key)
				continue
			}
			brokerMetrics.NumReplicas = val
		case leaders:
			if isFromTs {
				continue
			}
			req = promUrl + "query?query=kafka_server_replicamanager_leadercount%7Binstance%3D%27" + host + "%3A" + strconv.Itoa(port) + "%27%7D&time=" + strconv.Itoa(ts)
			val, _, err := getValueByEndpoint(ctx, req, true)
			if err != nil {
				//log.Logger.DebugContext(ctx, "error fetching metrics", key)
				continue
			}
			brokerMetrics.NumLeaders = val
		case activeControllers:
			if isFromTs {
				continue
			}
			req = promUrl + "query?query=kafka_controller_kafkacontroller_activecontrollercount%7Binstance%3D%27" + host + "%3A" + strconv.Itoa(port) + "%27%7D&time=" + strconv.Itoa(ts)
			val, _, err := getValueByEndpoint(ctx, req, true)
			if err != nil {
				//log.Logger.DebugContext(ctx, "error fetching metrics", key)
				continue
			}
			brokerMetrics.NumActControllers = val
		case offlinePartitions:
			if isFromTs {
				continue
			}
			req = promUrl + "query?query=kafka_controller_kafkacontroller_offlinepartitionscount%7Binstance%3D%27" + host + "%3A" + strconv.Itoa(port) + "%27%7D&time=" + strconv.Itoa(ts)
			val, _, err := getValueByEndpoint(ctx, req, true)
			if err != nil {
				//log.Logger.DebugContext(ctx, "error fetching metrics", key)
				continue
			}
			brokerMetrics.OfflinePartitions = val
		case underReplicated:
			if isFromTs {
				continue
			}
			req = promUrl + "query?query=kafka_server_replicamanager_underreplicatedpartitions%7Binstance%3D%27" + host + "%3A" + strconv.Itoa(port) + "%27%7D&time=" + strconv.Itoa(ts)
			val, _, err := getValueByEndpoint(ctx, req, true)
			if err != nil {
				//log.Logger.DebugContext(ctx, "error fetching metrics", key)
				continue
			}
			brokerMetrics.UnderReplicated = val
		case messageRate:
			req = promUrl + "query?query=sum%20by%20(instance)%20(rate(kafka_server_brokertopicmetrics_messagesin_total%7Binstance%3D%27" + host + "%3A" + strconv.Itoa(port) + "%27%7D%5B1m%5D))&time=" + strconv.Itoa(ts)
			_, val, err := getValueByEndpoint(ctx, req, false)
			if err != nil {
				//log.Logger.DebugContext(ctx, "error fetching metrics", key)
				continue
			}
			brokerMetrics.MessageRate = val
		case isrExpansionRate:
			req = promUrl + "query?query=kafka_server_replicamanager_isrexpands_total%7Binstance%3D%27" + host + "%3A" + strconv.Itoa(port) + "%27%7D&time=" + strconv.Itoa(ts)
			_, val, err := getValueByEndpoint(ctx, req, false)
			if err != nil {
				//log.Logger.DebugContext(ctx, "error fetching metrics", key)
				continue
			}
			brokerMetrics.IsrExpansionRate = val
		case isrShrinkRate:
			req = promUrl + "query?query=kafka_server_replicamanager_isrshrinks_total%7Binstance%3D%27" + host + "%3A" + strconv.Itoa(port) + "%27%7D&time=" + strconv.Itoa(ts)
			_, val, err := getValueByEndpoint(ctx, req, false)
			if err != nil {
				//log.Logger.DebugContext(ctx, "error fetching metrics", key)
				continue
			}
			brokerMetrics.IsrShrinkRate = val
		case networkProcIdlePerc:
			req = promUrl + "query?query=kafka_network_socketserver_networkprocessoravgidlepercent%7Binstance%3D%27" + host + "%3A" + strconv.Itoa(port) + "%27%7D&time=" + strconv.Itoa(ts)
			_, val, err := getValueByEndpoint(ctx, req, false)
			if err != nil {
				//log.Logger.DebugContext(ctx, "error fetching metrics", key)
				continue
			}
			brokerMetrics.NetworkProcAvgIdlePercent = val
		case responseTime:
			req = promUrl + "query?query=sum%20by%20(instance)%20(rate(kafka_network_requestmetrics_responsesendtimems%7Binstance%3D%27" + host + "%3A" + strconv.Itoa(port) + "%27%7D%5B1m%5D))&time=" + strconv.Itoa(ts)
			_, val, err := getValueByEndpoint(ctx, req, false)
			if err != nil {
				//log.Logger.DebugContext(ctx, "error fetching metrics", key)
				continue
			}
			brokerMetrics.ResponseTime = val
		case queueTime:
			req = promUrl + "query?query=sum%20by%20(instance)%20(rate(kafka_network_requestmetrics_requestqueuetimems%7Binstance%3D%27" + host + "%3A" + strconv.Itoa(port) + "%27%7D%5B1m%5D))&time=" + strconv.Itoa(ts)
			_, val, err := getValueByEndpoint(ctx, req, false)
			if err != nil {
				//log.Logger.DebugContext(ctx, "error fetching metrics", key)
				continue
			}
			brokerMetrics.QueueTime = val
		case remoteTime:
			req = promUrl + "query?query=sum%20by%20(instance)%20(rate(kafka_network_requestmetrics_remotetimems%7Binstance%3D%27" + host + "%3A" + strconv.Itoa(port) + "%27%7D%5B1m%5D))&time=" + strconv.Itoa(ts)
			_, val, err := getValueByEndpoint(ctx, req, false)
			if err != nil {
				//log.Logger.DebugContext(ctx, "error fetching metrics", key)
				continue
			}
			brokerMetrics.RemoteTime = val
		case localTime:
			req = promUrl + "query?query=sum%20by%20(instance)%20(rate(kafka_network_requestmetrics_localtimems%7Binstance%3D%27" + host + "%3A" + strconv.Itoa(port) + "%27%7D%5B1m%5D))&time=" + strconv.Itoa(ts)
			_, val, err := getValueByEndpoint(ctx, req, false)
			if err != nil {
				//log.Logger.DebugContext(ctx, "error fetching metrics", key)
				continue
			}
			brokerMetrics.LocalTIme = val
		case maxLagBtwLeadAndRep:
			req = promUrl + "query?query=kafka_server_replicafetchermanager_minfetchrate%7Binstance%3D%27" + host + "%3A" + strconv.Itoa(port) + "%27%7D&time=" + strconv.Itoa(ts)
			_, val, err := getValueByEndpoint(ctx, req, false)
			if err != nil {
				//log.Logger.DebugContext(ctx, "error fetching metrics", key)
				continue
			}
			brokerMetrics.MaxLagBtwLeadAndRepl = val
		case uncleanLeadElec:
			req = promUrl + "query?query=kafka_controller_controllerstats_uncleanleaderelectionspersec%7Binstance%3D%27" + host + "%3A" + strconv.Itoa(port) + "%27%7D&time=" + strconv.Itoa(ts)
			_, val, err := getValueByEndpoint(ctx, req, false)
			if err != nil {
				//log.Logger.DebugContext(ctx, "error fetching metrics", key)
				continue
			}
			brokerMetrics.UncleanLeadElec = val
		case failedFetchRate:
			req = promUrl + "query?query=sum%20by%20(instance)%20(rate(kafka_server_brokertopicmetrics_failedfetchrequests_total%7Binstance%3D%27" + host + "%3A" + strconv.Itoa(port) + "%27%7D%5B1m%5D))&time=" + strconv.Itoa(ts)
			_, val, err := getValueByEndpoint(ctx, req, false)
			if err != nil {
				//log.Logger.DebugContext(ctx, "error fetching metrics", key)
				continue
			}
			brokerMetrics.FailedFetchReqRate = val
		case failedProdRate:
			req = promUrl + "query?query=sum%20by%20(instance)%20(rate(kafka_server_brokertopicmetrics_failedproducerequests_total%7Binstance%3D%27" + host + "%3A" + strconv.Itoa(port) + "%27%7D%5B1m%5D))&time=" + strconv.Itoa(ts)
			_, val, err := getValueByEndpoint(ctx, req, false)
			if err != nil {
				//log.Logger.DebugContext(ctx, "error fetching metrics", key)
				continue
			}
			brokerMetrics.FailedProdReqRate = val
		case bytesIn:
			req = promUrl + "query?query=sum%20by%20(instance)%20(rate(kafka_server_brokertopicmetrics_bytesin_total%7Binstance%3D%27" + host + "%3A" + strconv.Itoa(port) + "%27%7D%5B1m%5D))&time=" + strconv.Itoa(ts)
			_, val, err := getValueByEndpoint(ctx, req, false)
			if err != nil {
				//log.Logger.DebugContext(ctx, "error fetching metrics", key)
				continue
			}
			brokerMetrics.ByteInRate = int64(val)
		case bytesOut:
			req = promUrl + "query?query=sum%20by%20(instance)%20(rate(kafka_server_brokertopicmetrics_bytesout_total%7Binstance%3D%27" + host + "%3A" + strconv.Itoa(port) + "%27%7D%5B1m%5D))&time=" + strconv.Itoa(ts)
			_, val, err := getValueByEndpoint(ctx, req, false)
			if err != nil {
				//log.Logger.DebugContext(ctx, "error fetching metrics", key)
				continue
			}
			brokerMetrics.ByteOutRate = int64(val)
		case totalMessages:
			if isFromTs {
				continue
			}
			req = promUrl + "query?query=sum%20by%20(instance)%20(kafka_server_brokertopicmetrics_messagesin_total%7Binstance%3D%22" + host + "%3A" + strconv.Itoa(port) + "%22%7D)&time=" + strconv.Itoa(ts)
			val, _, err := getValueByEndpoint(ctx, req, true)
			if err != nil {
				//log.Logger.DebugContext(ctx, "error fetching metrics", key)
				continue
			}
			brokerMetrics.Messages = int64(val)
		case totalTopics:
			if isFromTs {
				continue
			}
			req = promUrl + "query?query=count%20by%20(instance)%20(rate(kafka_server_brokertopicmetrics_messagesin_total%7Binstance%3D%27" + host + "%3A" + strconv.Itoa(port) + "%27%7D%5B1m%5D))&time=" + strconv.Itoa(ts)
			val, _, err := getValueByEndpoint(ctx, req, true)
			if err != nil {
				//log.Logger.DebugContext(ctx, "error fetching metrics", key)
				continue
			}
			brokerMetrics.Topics = val
		}
	}
	return brokerMetrics, nil
}

func getResponseByEndpoint(ctx context.Context, req string) (metrics domain.PromResponse, err error) {
	res, err := http.Get(req)
	if err != nil {
		log.Logger.ErrorContext(ctx, err, "querying broker metrics failed")
		return metrics, err
	}

	content, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Logger.ErrorContext(ctx, err, "reading response failed")
		return
	}

	if res.StatusCode != http.StatusOK {
		log.Logger.ErrorContext(ctx, req, "error received for prometheus api call")
		return metrics, errors.New("api request failed")
	}

	err = json.Unmarshal(content, &metrics)
	if err != nil {
		log.Logger.ErrorContext(ctx, err, "unmarshalling response failed")
		return
	}

	return
}

func getValueByEndpoint(ctx context.Context, req string, typeInt bool) (intVal int, floatVal float64, err error) {
	res, err := http.Get(req)
	if err != nil {
		log.Logger.ErrorContext(ctx, err, "querying broker metrics failed")
		return intVal, floatVal, err
	}

	response, err := parseMetricsResponse(ctx, res)
	if err != nil {
		log.Logger.ErrorContext(ctx, "broker metrics call failed", req)
		return intVal, floatVal, err
	}

	if len(response.Data.Result) > 0 {
		if len(response.Data.Result[0].Value) < 2 {
			log.Logger.ErrorContext(ctx, "received insufficient values for query", response.Data.Result[0])
			return 0, 0.0, errors.New("received insufficient values for query")
		}

		if typeInt {
			val, err := strconv.ParseInt(response.Data.Result[0].Value[1].(string), 10, 64)
			if err != nil {
				log.Logger.ErrorContext(ctx, "converting val to int failed", response.Data.Result[0].Value[1])
				return intVal, floatVal, err
			}
			return int(val), floatVal, nil
		} else {
			floatVal, err := strconv.ParseFloat(response.Data.Result[0].Value[1].(string), 64)
			if err != nil {
				log.Logger.ErrorContext(ctx, "converting value to float failed", response.Data.Result[0].Value[1])
				return intVal, floatVal, err
			}
			return intVal, floatVal, nil
		}
	}
	return intVal, floatVal, errors.New("prom data result received insufficient data")
}

func getStepFromTs(ctx context.Context, from, to int) (step int, err error) {
	tsGap := to - from

	if tsGap < 0 {
		log.Logger.ErrorContext(ctx, fmt.Sprintf("invalid to (%v) and from (%v) ts", from, to))
		return 0, errors.New(fmt.Sprintf("invalid to (%v) and from (%v) ts", from, to))
	} else if tsGap < 3600 {
		return 15, nil
	} else if tsGap < 7200 {
		return 20, nil
	} else if tsGap < 10800 {
		return 30, nil
	} else if tsGap < 21600 {
		return 60, nil
	} else if tsGap < 43200 {
		return 120, nil
	} else {
		numOfDays := tsGap % 86400
		return 300*numOfDays, nil
	}
}
