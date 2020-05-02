package database

import (
	"context"
	"fmt"
	"github.com/YasiruR/ktool-backend/log"
	"strconv"
)

func AddMetricsRow(ctx context.Context, host string, ts int) (err error) {
	query := "INSERT INTO " + brokerMetricsTable + "(id, timestamp, host) VALUES (null, " + strconv.Itoa(ts) + `, "` + host + `");`
	_, err = Db.Exec(query)
	if err != nil {
		log.Logger.ErrorContext(ctx, err, "creating a row in broker metrics table failed")
		return err
	}
	return nil
}

func UpdateBrokerByteInRate(ctx context.Context, bytesIn float64, host string, ts int) (err error) {
	//todo run a separate job to clean these broker metrics table
	query := "UPDATE " + brokerMetricsTable + " SET bytes_in=" + strconv.FormatFloat(bytesIn, 'f', -1, 64) + " WHERE timestamp=" + strconv.Itoa(ts) + ` AND host="` + host + `";`
	_, err = Db.Exec(query)
	if err != nil {
		log.Logger.ErrorContext(ctx, fmt.Sprintf("bytes_in update to %s table failed", brokerMetricsTable), err)
		return err
	}
	return nil
}

func UpdateBrokerByteOutRate(ctx context.Context, bytesOut float64, host string, ts int) (err error) {
	query := "UPDATE " + brokerMetricsTable + " SET bytes_out=" + strconv.FormatFloat(bytesOut, 'f', -1, 64) + " WHERE timestamp=" + strconv.Itoa(ts) + ` AND host="` + host + `";`
	_, err = Db.Exec(query)
	if err != nil {
		log.Logger.ErrorContext(ctx, fmt.Sprintf("bytes_out update to %s table failed", brokerMetricsTable), err)
		return err
	}
	return nil
}

func UpdateBrokerPartitionCount(ctx context.Context, count int, host string, ts int) (err error) {
	query := "UPDATE " + brokerMetricsTable + " SET partitions=" + strconv.Itoa(count) + " WHERE timestamp=" + strconv.Itoa(ts) + ` AND host="` + host + `";`
	_, err = Db.Exec(query)
	if err != nil {
		log.Logger.ErrorContext(ctx, fmt.Sprintf("partition update to %s table failed", brokerMetricsTable), err)
		return err
	}
	return nil
}

func UpdateBrokerLeaderCount(ctx context.Context, count int, host string, ts int) (err error) {
	query := "UPDATE " + brokerMetricsTable + " SET leaders=" + strconv.Itoa(count) + " WHERE timestamp=" + strconv.Itoa(ts) + ` AND host="` + host + `";`
	_, err = Db.Exec(query)
	if err != nil {
		log.Logger.ErrorContext(ctx, fmt.Sprintf("leader update to %s table failed", brokerMetricsTable), err)
		return err
	}
	return nil
}

func UpdateBrokerActControllerCount(ctx context.Context, count int, host string, ts int) (err error) {
	query := "UPDATE " + brokerMetricsTable + " SET act_controller=" + strconv.Itoa(count) + " WHERE timestamp=" + strconv.Itoa(ts) + ` AND host="` + host + `";`
	_, err = Db.Exec(query)
	if err != nil {
		log.Logger.ErrorContext(ctx, fmt.Sprintf("act controller update to %s table failed", brokerMetricsTable), err)
		return err
	}
	return nil
}

