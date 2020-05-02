package prometheus

import (
	"context"
	"encoding/json"
	"github.com/YasiruR/ktool-backend/log"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
)

func setIntMetrics(ctx context.Context, ts int, req string, dbFunc func(ctx context.Context, val int, host string, ts int)(err error)) (err error) {
	res, err := http.Get(req)
	if err != nil {
		log.Logger.ErrorContext(ctx, err, "querying broker metrics failed")
		return
	}

	response, err := parseMetricsResponse(ctx, res)
	if err != nil {
		log.Logger.ErrorContext(ctx, "broker metrics call failed")
		return
	}

	for _, result := range response.Data.Result {
		if len(result.Value) < 2 {
			log.Logger.ErrorContext(ctx, "received insufficient values for query", result)
			continue
		}

		s := strings.Split(result.Metric.Instance, ":")
		if len(s) < 2 {
			log.Logger.ErrorContext(ctx, "invalid format received for instance", result.Metric.Instance)
			continue
		}
		host := s[0]

		val, err := strconv.ParseInt(result.Value[1].(string), 10, 64)
		if err != nil {
			log.Logger.ErrorContext(ctx, "converting val to int failed", result.Value[1])
			continue
		}
		err = dbFunc(ctx, int(val), host, ts)
		if err != nil {
			log.Logger.ErrorContext(ctx, "db query to update broker metrics failed", result.Metric.Instance)
		}
	}
	return nil
}

func setFloatMetrics(ctx context.Context, ts int, req string, dbFunc func(ctx context.Context, val float64, host string, ts int)(err error)) (err error) {
	res, err := http.Get(req)
	if err != nil {
		log.Logger.ErrorContext(ctx, err, "querying broker metrics failed")
		return
	}

	response, err := parseMetricsResponse(ctx, res)
	if err != nil {
		log.Logger.ErrorContext(ctx, "broker metrics call failed")
		return
	}

	for _, result := range response.Data.Result {
		if len(result.Value) < 2 {
			log.Logger.ErrorContext(ctx, "received insufficient values for query", result)
			continue
		}

		s := strings.Split(result.Metric.Instance, ":")
		if len(s) < 2 {
			log.Logger.ErrorContext(ctx, "invalid format received for instance", result.Metric.Instance)
			continue
		}
		host := s[0]

		val, err := strconv.ParseFloat(result.Value[1].(string), 64)
		if err != nil {
			log.Logger.ErrorContext(ctx, "converting value to float failed", result.Value[1])
			continue
		}
		err = dbFunc(ctx, val, host, ts)
		if err != nil {
			log.Logger.ErrorContext(ctx, "db query to update broker metrics failed", result.Metric.Instance)
		}
	}
	return nil
}

func parseMetricsResponse(ctx context.Context, res *http.Response) (response BrokerMetricsResponse, err error) {
	content, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Logger.ErrorContext(ctx, err, "reading response failed")
		return
	}

	if res.StatusCode != http.StatusOK {
		log.Logger.ErrorContext(ctx, response, "error received for prometheus api call")
		return
	}

	err = json.Unmarshal(content, &response)
	if err != nil {
		log.Logger.ErrorContext(ctx, err, "unmarshalling response failed")
		return
	}
	return
}

//func setBrokerBytesIn(ctx context.Context, ts int) (err error) {
//	//query bytes in to the broker
//	res, err := http.Get(promUrl + "query?query=sum%20by%20(instance)%20(rate(kafka_server_brokertopicmetrics_bytesin_total%5B1m%5D))&time=" + strconv.Itoa(ts))
//	if err != nil {
//		log.Logger.ErrorContext(ctx, err, "querying broker total bytes in failed")
//		return
//	}
//
//	response, err := parseMetricsResponse(ctx, res)
//	if err != nil {
//		log.Logger.ErrorContext(ctx, "broker bytes in metrics call failed")
//		return
//	}
//
//	for _, result := range response.Data.Result {
//		if len(result.Value) < 2 {
//			log.Logger.ErrorContext(ctx, "received insufficient values for query", result)
//			continue
//		}
//		byteRate, err := strconv.ParseFloat(result.Value[1].(string), 64)
//		if err != nil {
//			log.Logger.ErrorContext(ctx, "converting byte in value to float failed", result.Value[1])
//			continue
//		}
//
//		s := strings.Split(result.Metric.Instance, ":")
//		if len(s) < 2 {
//			log.Logger.ErrorContext(ctx, "invalid format received for instance", result.Metric.Instance)
//			continue
//		}
//		host := s[0]
//
//		err = database.UpdateBrokerByteInRate(ctx, byteRate, host, ts)
//		if err != nil {
//			log.Logger.ErrorContext(ctx, "db query to update broker bytes in failed", result.Metric.Instance)
//		}
//	}
//	return nil
//}
//
//func setBrokerBytesOut(ctx context.Context, ts int) (err error) {
//	//query bytes out from the broker
//	res, err := http.Get(promUrl + "query?query=sum%20by%20(instance)%20(rate(kafka_server_brokertopicmetrics_bytesout_total%5B1m%5D))&time=" + strconv.Itoa(ts))
//	if err != nil {
//		log.Logger.ErrorContext(ctx, err, "querying broker total bytes out failed")
//		return
//	}
//
//	response, err := parseMetricsResponse(ctx, res)
//	if err != nil {
//		log.Logger.ErrorContext(ctx, "broker bytes out metrics call failed")
//		return
//	}
//
//	for _, result := range response.Data.Result {
//		if len(result.Value) < 2 {
//			log.Logger.ErrorContext(ctx, "received insufficient values for query", result)
//			continue
//		}
//		byteRate, err := strconv.ParseFloat(result.Value[1].(string), 64)
//		if err != nil {
//			log.Logger.ErrorContext(ctx, "converting byte out value to float failed", result.Value[1])
//			continue
//		}
//
//		s := strings.Split(result.Metric.Instance, ":")
//		if len(s) < 2 {
//			log.Logger.ErrorContext(ctx, "invalid format received for instance", result.Metric.Instance)
//			continue
//		}
//		host := s[0]
//
//		err = database.UpdateBrokerByteOutRate(ctx, byteRate, host, ts)
//		if err != nil {
//			log.Logger.ErrorContext(ctx, "db query to update broker bytes out failed", result.Metric.Instance)
//		}
//	}
//	return nil
//}
//
//func setPartitionCount(ctx context.Context, ts int) (err error) {
//	res, err := http.Get(promUrl + "query?query=kafka_server_replicamanager_partitioncount&time=" + strconv.Itoa(ts))
//	if err != nil {
//		log.Logger.ErrorContext(ctx, err, "querying broker partition count failed")
//		return
//	}
//
//	response, err := parseMetricsResponse(ctx, res)
//	if err != nil {
//		log.Logger.ErrorContext(ctx, "broker partitions metrics call failed")
//		return
//	}
//
//	for _, result := range response.Data.Result {
//		if len(result.Value) < 2 {
//			log.Logger.ErrorContext(ctx, "received insufficient values for query", result)
//			continue
//		}
//		partitions, err := strconv.ParseInt(result.Value[1].(string), 10, 64)
//		if err != nil {
//			log.Logger.ErrorContext(ctx, "converting partition count to int failed", result.Value[1])
//			continue
//		}
//
//		s := strings.Split(result.Metric.Instance, ":")
//		if len(s) < 2 {
//			log.Logger.ErrorContext(ctx, "invalid format received for instance", result.Metric.Instance)
//			continue
//		}
//		host := s[0]
//
//		err = database.UpdateBrokerPartitionCount(ctx, int(partitions), host, ts)
//		if err != nil {
//			log.Logger.ErrorContext(ctx, "db query to update broker partition count failed", result.Metric.Instance)
//		}
//	}
//	return nil
//}
//
//func setLeaderCount(ctx context.Context, ts int) (err error) {
//	res, err := http.Get(promUrl + "query?query=kafka_server_replicamanager_leadercount&time=" + strconv.Itoa(ts))
//	if err != nil {
//		log.Logger.ErrorContext(ctx, err, "querying broker leader count failed")
//		return
//	}
//
//	response, err := parseMetricsResponse(ctx, res)
//	if err != nil {
//		log.Logger.ErrorContext(ctx, "broker leader metrics call failed")
//		return
//	}
//
//	for _, result := range response.Data.Result {
//		if len(result.Value) < 2 {
//			log.Logger.ErrorContext(ctx, "received insufficient values for query", result)
//			continue
//		}
//		leaders, err := strconv.ParseInt(result.Value[1].(string), 10, 64)
//		if err != nil {
//			log.Logger.ErrorContext(ctx, "converting leader count to int failed", result.Value[1])
//			continue
//		}
//
//		s := strings.Split(result.Metric.Instance, ":")
//		if len(s) < 2 {
//			log.Logger.ErrorContext(ctx, "invalid format received for instance", result.Metric.Instance)
//			continue
//		}
//		host := s[0]
//
//		err = database.UpdateBrokerLeaderCount(ctx, int(leaders), host, ts)
//		if err != nil {
//			log.Logger.ErrorContext(ctx, "db query to update broker leader count failed", result.Metric.Instance)
//		}
//	}
//	return nil
//}
//
//func setActiveControllerCount(ctx context.Context, ts int) (err error) {
//	res, err := http.Get(promUrl + "query?query=kafka_controller_kafkacontroller_activecontrollercount&time=" + strconv.Itoa(ts))
//	if err != nil {
//		log.Logger.ErrorContext(ctx, err, "querying broker active controller count failed")
//		return
//	}
//
//	response, err := parseMetricsResponse(ctx, res)
//	if err != nil {
//		log.Logger.ErrorContext(ctx, "broker active controller metrics call failed")
//		return
//	}
//
//	for _, result := range response.Data.Result {
//		if len(result.Value) < 2 {
//			log.Logger.ErrorContext(ctx, "received insufficient values for query", result)
//			continue
//		}
//		activeControllers, err := strconv.ParseInt(result.Value[1].(string), 10, 64)
//		if err != nil {
//			log.Logger.ErrorContext(ctx, "converting controller count to int failed", result.Value[1])
//			continue
//		}
//
//		s := strings.Split(result.Metric.Instance, ":")
//		if len(s) < 2 {
//			log.Logger.ErrorContext(ctx, "invalid format received for instance", result.Metric.Instance)
//			continue
//		}
//		host := s[0]
//
//		err = database.UpdateBrokerLeaderCount(ctx, int(activeControllers), host, ts)
//		if err != nil {
//			log.Logger.ErrorContext(ctx, "db query to update broker active controller count failed", result.Metric.Instance)
//		}
//	}
//	return nil
//}