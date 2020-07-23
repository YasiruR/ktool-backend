package prometheus

import (
	"context"
	"errors"
	"github.com/YasiruR/ktool-backend/domain"
	"github.com/YasiruR/ktool-backend/log"
	"net/http"
	"strconv"
)

func GetMetricsByTimestamp(ctx context.Context, host string, port, ts int, toTs bool) (brokerMetrics domain.BrokerMetrics, err error) {
	for key, _ := range queryList {
		var req string
		switch key {
		case partitions:
			req = promUrl + "query?query=kafka_server_replicamanager_partitioncount%7Binstance%3D%27" + host + "%3A" + strconv.Itoa(port) + "%27%7D&time=" + strconv.Itoa(ts)
			val, _, err := getValueByEndpoint(ctx, req, true)
			if err != nil {
				log.Logger.DebugContext(ctx, "error fetching metrics", key)
				continue
			}
			brokerMetrics.NumReplicas = val
		case leaders:
			req = promUrl + "query?query=kafka_server_replicamanager_leadercount%7Binstance%3D%27" + host + "%3A" + strconv.Itoa(port) + "%27%7D&time=" + strconv.Itoa(ts)
			val, _, err := getValueByEndpoint(ctx, req, true)
			if err != nil {
				log.Logger.DebugContext(ctx, "error fetching metrics", key)
				continue
			}
			brokerMetrics.NumLeaders = val
		case activeControllers:
			req = promUrl + "query?query=kafka_controller_kafkacontroller_activecontrollercount%7Binstance%3D%27" + host + "%3A" + strconv.Itoa(port) + "%27%7D&time=" + strconv.Itoa(ts)
			val, _, err := getValueByEndpoint(ctx, req, true)
			if err != nil {
				log.Logger.DebugContext(ctx, "error fetching metrics", key)
				continue
			}
			brokerMetrics.NumActControllers = val
		case offlinePartitions:
			req = promUrl + "query?query=kafka_controller_kafkacontroller_offlinepartitionscount%7Binstance%3D%27" + host + "%3A" + strconv.Itoa(port) + "%27%7D&time=" + strconv.Itoa(ts)
			val, _, err := getValueByEndpoint(ctx, req, true)
			if err != nil {
				log.Logger.DebugContext(ctx, "error fetching metrics", key)
				continue
			}
			brokerMetrics.OfflinePartitions = val
		case underReplicated:
			req = promUrl + "query?query=kafka_server_replicamanager_underreplicatedpartitions%7Binstance%3D%27" + host + "%3A" + strconv.Itoa(port) + "%27%7D&time=" + strconv.Itoa(ts)
			val, _, err := getValueByEndpoint(ctx, req, true)
			if err != nil {
				log.Logger.DebugContext(ctx, "error fetching metrics", key)
				continue
			}
			brokerMetrics.UnderReplicated = val
		case messageRate:
			req = promUrl + "query?query=sum%20by%20(instance)%20(rate(kafka_server_brokertopicmetrics_messagesin_total%7Binstance%3D%27" + host + "%3A" + strconv.Itoa(port) + "%27%7D%5B1m%5D))&time=" + strconv.Itoa(ts)
			_, val, err := getValueByEndpoint(ctx, req, false)
			if err != nil {
				log.Logger.DebugContext(ctx, "error fetching metrics", key)
				continue
			}
			brokerMetrics.MessageRate = val
		case isrExpansionRate:
			req = promUrl + "query?query=kafka_server_replicamanager_isrexpands_total%7Binstance%3D%27" + host + "%3A" + strconv.Itoa(port) + "%27%7D&time=" + strconv.Itoa(ts)
			_, val, err := getValueByEndpoint(ctx, req, false)
			if err != nil {
				log.Logger.DebugContext(ctx, "error fetching metrics", key)
				continue
			}
			brokerMetrics.IsrExpansionRate = val
		case isrShrinkRate:
			req = promUrl + "query?query=kafka_server_replicamanager_isrshrinks_total%7Binstance%3D%27" + host + "%3A" + strconv.Itoa(port) + "%27%7D&time=" + strconv.Itoa(ts)
			_, val, err := getValueByEndpoint(ctx, req, false)
			if err != nil {
				log.Logger.DebugContext(ctx, "error fetching metrics", key)
				continue
			}
			brokerMetrics.IsrShrinkRate = val
		case networkProcIdlePerc:
			req = promUrl + "query?query=kafka_network_socketserver_networkprocessoravgidlepercent%7Binstance%3D%27" + host + "%3A" + strconv.Itoa(port) + "%27%7D&time=" + strconv.Itoa(ts)
			_, val, err := getValueByEndpoint(ctx, req, false)
			if err != nil {
				log.Logger.DebugContext(ctx, "error fetching metrics", key)
				continue
			}
			brokerMetrics.NetworkProcAvgIdlePercent = val
		case responseTime:
			req = promUrl + "query?query=sum%20by%20(instance)%20(rate(kafka_network_requestmetrics_responsesendtimems%7Binstance%3D%27" + host + "%3A" + strconv.Itoa(port) + "%27%7D%5B1m%5D))&time=" + strconv.Itoa(ts)
			_, val, err := getValueByEndpoint(ctx, req, false)
			if err != nil {
				log.Logger.DebugContext(ctx, "error fetching metrics", key)
				continue
			}
			brokerMetrics.ResponseTime = val
		case queueTime:
			req = promUrl + "query?query=sum%20by%20(instance)%20(rate(kafka_network_requestmetrics_requestqueuetimems%7Binstance%3D%27" + host + "%3A" + strconv.Itoa(port) + "%27%7D%5B1m%5D))&time=" + strconv.Itoa(ts)
			_, val, err := getValueByEndpoint(ctx, req, false)
			if err != nil {
				log.Logger.DebugContext(ctx, "error fetching metrics", key)
				continue
			}
			brokerMetrics.QueueTime = val
		case remoteTime:
			req = promUrl + "query?query=sum%20by%20(instance)%20(rate(kafka_network_requestmetrics_remotetimems%7Binstance%3D%27" + host + "%3A" + strconv.Itoa(port) + "%27%7D%5B1m%5D))&time=" + strconv.Itoa(ts)
			_, val, err := getValueByEndpoint(ctx, req, false)
			if err != nil {
				log.Logger.DebugContext(ctx, "error fetching metrics", key)
				continue
			}
			brokerMetrics.RemoteTime = val
		case localTime:
			req = promUrl + "query?query=sum%20by%20(instance)%20(rate(kafka_network_requestmetrics_localtimems%7Binstance%3D%27" + host + "%3A" + strconv.Itoa(port) + "%27%7D%5B1m%5D))&time=" + strconv.Itoa(ts)
			_, val, err := getValueByEndpoint(ctx, req, false)
			if err != nil {
				log.Logger.DebugContext(ctx, "error fetching metrics", key)
				continue
			}
			brokerMetrics.LocalTIme = val
		case maxLagBtwLeadAndRep:
			req = promUrl + "query?query=kafka_server_replicafetchermanager_minfetchrate%7Binstance%3D%27" + host + "%3A" + strconv.Itoa(port) + "%27%7D&time=" + strconv.Itoa(ts)
			_, val, err := getValueByEndpoint(ctx, req, false)
			if err != nil {
				log.Logger.DebugContext(ctx, "error fetching metrics", key)
				continue
			}
			brokerMetrics.MaxLagBtwLeadAndRepl = val
		case uncleanLeadElec:
			req = promUrl + "query?query=kafka_controller_controllerstats_uncleanleaderelectionspersec%7Binstance%3D%27" + host + "%3A" + strconv.Itoa(port) + "%27%7D&time=" + strconv.Itoa(ts)
			_, val, err := getValueByEndpoint(ctx, req, false)
			if err != nil {
				log.Logger.DebugContext(ctx, "error fetching metrics", key)
				continue
			}
			brokerMetrics.UncleanLeadElec = val
		case failedFetchRate:
			req = promUrl + "query?query=sum%20by%20(instance)%20(rate(kafka_server_brokertopicmetrics_failedfetchrequests_total%7Binstance%3D%27" + host + "%3A" + strconv.Itoa(port) + "%27%7D%5B1m%5D))&time=" + strconv.Itoa(ts)
			_, val, err := getValueByEndpoint(ctx, req, false)
			if err != nil {
				log.Logger.DebugContext(ctx, "error fetching metrics", key)
				continue
			}
			brokerMetrics.FailedFetchReqRate = val
		case failedProdRate:
			req = promUrl + "query?query=sum%20by%20(instance)%20(rate(kafka_server_brokertopicmetrics_failedproducerequests_total%7Binstance%3D%27" + host + "%3A" + strconv.Itoa(port) + "%27%7D%5B1m%5D))&time=" + strconv.Itoa(ts)
			_, val, err := getValueByEndpoint(ctx, req, false)
			if err != nil {
				log.Logger.DebugContext(ctx, "error fetching metrics", key)
				continue
			}
			brokerMetrics.FailedProdReqRate = val
		case bytesIn:
			req = promUrl + "query?query=sum%20by%20(instance)%20(rate(kafka_server_brokertopicmetrics_bytesin_total%7Binstance%3D%27" + host + "%3A" + strconv.Itoa(port) + "%27%7D%5B1m%5D))&time=" + strconv.Itoa(ts)
			val, _, err := getValueByEndpoint(ctx, req, true)
			if err != nil {
				log.Logger.DebugContext(ctx, "error fetching metrics", key)
				continue
			}
			brokerMetrics.ByteInRate = int64(val)
		case bytesOut:
			req = promUrl + "query?query=sum%20by%20(instance)%20(rate(kafka_server_brokertopicmetrics_bytesout_total%7Binstance%3D%27" + host + "%3A" + strconv.Itoa(port) + "%27%7D%5B1m%5D))&time=" + strconv.Itoa(ts)
			val, _, err := getValueByEndpoint(ctx, req, true)
			if err != nil {
				log.Logger.DebugContext(ctx, "error fetching metrics", key)
				continue
			}
			brokerMetrics.ByteOutRate = int64(val)
		case totalMessages:
			req = promUrl + "query?query=sum%20by%20(instance)%20(rate(kafka_server_brokertopicmetrics_messagesin_total%7Binstance%3D%27" + host + "%3A" + strconv.Itoa(port) + "%27%7D%5B1m%5D))&time=" + strconv.Itoa(ts)
			val, _, err := getValueByEndpoint(ctx, req, true)
			if err != nil {
				log.Logger.DebugContext(ctx, "error fetching metrics", key)
				continue
			}
			brokerMetrics.Messages = int64(val)
		case totalTopics:
			req = promUrl + "query?query=count%20by%20(instance)%20(rate(kafka_server_brokertopicmetrics_messagesin_total%7Binstance%3D%27" + host + "%3A" + strconv.Itoa(port) + "%27%7D%5B1m%5D))&time=" + strconv.Itoa(ts)
			val, _, err := getValueByEndpoint(ctx, req, true)
			if err != nil {
				log.Logger.DebugContext(ctx, "error fetching metrics", key)
				continue
			}
			brokerMetrics.Topics = val
		}
	}
	return brokerMetrics, nil
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
			intVal, ok := response.Data.Result[0].Value[1].(int)
			if ok {
				return intVal, floatVal, nil
			}
		} else {
			floatVal, ok := response.Data.Result[0].Value[1].(float64)
			if ok {
				return intVal, floatVal, nil
			}
		}

		log.Logger.ErrorContext(ctx, "received val is not the type expected", typeInt)
		return intVal, floatVal, errors.New("received val is not type of int")
	}

	log.Logger.ErrorContext(ctx, "prom data result received insufficient data")
	return intVal, floatVal, errors.New("prom data result received insufficient data")
}
