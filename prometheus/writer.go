package prometheus

import (
	"context"
	"encoding/json"
	"errors"
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
		log.Logger.ErrorContext(ctx, "broker metrics call failed", req)
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
		return response, errors.New("api request failed")
	}

	err = json.Unmarshal(content, &response)
	if err != nil {
		log.Logger.ErrorContext(ctx, err, "unmarshalling response failed")
		return
	}
	return
}