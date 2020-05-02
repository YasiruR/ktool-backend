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

//todo run a separate job to clean these broker metrics table
func UpdateBrokerByteInRate(ctx context.Context, bytesIn float64, host string, ts int) (err error) {
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

func UpdateBrokerOfflinePartCount(ctx context.Context, count int, host string, ts int) (err error) {
	query := "UPDATE " + brokerMetricsTable + " SET offline_part=" + strconv.Itoa(count) + " WHERE timestamp=" + strconv.Itoa(ts) + ` AND host="` + host + `";`
	_, err = Db.Exec(query)
	if err != nil {
		log.Logger.ErrorContext(ctx, fmt.Sprintf("offline partition update to %s table failed", brokerMetricsTable), err)
		return err
	}
	return nil
}

func UpdateBrokerUnderReplicatedCount(ctx context.Context, count int, host string, ts int) (err error) {
	query := "UPDATE " + brokerMetricsTable + " SET under_replicated=" + strconv.Itoa(count) + " WHERE timestamp=" + strconv.Itoa(ts) + ` AND host="` + host + `";`
	_, err = Db.Exec(query)
	if err != nil {
		log.Logger.ErrorContext(ctx, fmt.Sprintf("under replicated partition update to %s table failed", brokerMetricsTable), err)
		return err
	}
	return nil
}

func UpdateBrokerMessageRate(ctx context.Context, count int, host string, ts int) (err error) {
	query := "UPDATE " + brokerMetricsTable + " SET mesg_rate=" + strconv.Itoa(count) + " WHERE timestamp=" + strconv.Itoa(ts) + ` AND host="` + host + `";`
	_, err = Db.Exec(query)
	if err != nil {
		log.Logger.ErrorContext(ctx, fmt.Sprintf("message rate update to %s table failed", brokerMetricsTable), err)
		return err
	}
	return nil
}

func UpdateBrokerIsrExpRate(ctx context.Context, count float64, host string, ts int) (err error) {
	query := "UPDATE " + brokerMetricsTable + " SET isr_exp_rate=" + strconv.FormatFloat(count, 'f', -1, 64) + " WHERE timestamp=" + strconv.Itoa(ts) + ` AND host="` + host + `";`
	_, err = Db.Exec(query)
	if err != nil {
		log.Logger.ErrorContext(ctx, fmt.Sprintf("isr expansion rate update to %s table failed", brokerMetricsTable), err)
		return err
	}
	return nil
}

func UpdateBrokerIsrShrinkRate(ctx context.Context, count float64, host string, ts int) (err error) {
	query := "UPDATE " + brokerMetricsTable + " SET isr_shrink_rate=" + strconv.FormatFloat(count, 'f', -1, 64) + " WHERE timestamp=" + strconv.Itoa(ts) + ` AND host="` + host + `";`
	_, err = Db.Exec(query)
	if err != nil {
		log.Logger.ErrorContext(ctx, fmt.Sprintf("isr shrink rate update to %s table failed", brokerMetricsTable), err)
		return err
	}
	return nil
}

func UpdateBrokerNetworkIdlePercentage(ctx context.Context, count float64, host string, ts int) (err error) {
	query := "UPDATE " + brokerMetricsTable + " SET net_proc_avg_idle_perc=" + strconv.FormatFloat(count, 'f', -1, 64) + " WHERE timestamp=" + strconv.Itoa(ts) + ` AND host="` + host + `";`
	_, err = Db.Exec(query)
	if err != nil {
		log.Logger.ErrorContext(ctx, fmt.Sprintf("net_proc_avg_idle_perc update to %s table failed", brokerMetricsTable), err)
		return err
	}
	return nil
}

func UpdateBrokerResponseTime(ctx context.Context, count float64, host string, ts int) (err error) {
	query := "UPDATE " + brokerMetricsTable + " SET send_time=" + strconv.FormatFloat(count, 'f', -1, 64) + " WHERE timestamp=" + strconv.Itoa(ts) + ` AND host="` + host + `";`
	_, err = Db.Exec(query)
	if err != nil {
		log.Logger.ErrorContext(ctx, fmt.Sprintf("response (send) time update to %s table failed", brokerMetricsTable), err)
		return err
	}
	return nil
}

func UpdateBrokerQueueTime(ctx context.Context, value float64, host string, ts int) (err error) {
	query := "UPDATE " + brokerMetricsTable + " SET queue_time=" + strconv.FormatFloat(value, 'f', -1, 64) + " WHERE timestamp=" + strconv.Itoa(ts) + ` AND host="` + host + `";`
	_, err = Db.Exec(query)
	if err != nil {
		log.Logger.ErrorContext(ctx, fmt.Sprintf("queue time update to %s table failed", brokerMetricsTable), err)
		return err
	}
	return nil
}

func UpdateBrokerRemoteTime(ctx context.Context, value float64, host string, ts int) (err error) {
	query := "UPDATE " + brokerMetricsTable + " SET remote_time=" + strconv.FormatFloat(value, 'f', -1, 64) + " WHERE timestamp=" + strconv.Itoa(ts) + ` AND host="` + host + `";`
	_, err = Db.Exec(query)
	if err != nil {
		log.Logger.ErrorContext(ctx, fmt.Sprintf("remote time update to %s table failed", brokerMetricsTable), err)
		return err
	}
	return nil
}

func UpdateBrokerLocalTime(ctx context.Context, value float64, host string, ts int) (err error) {
	query := "UPDATE " + brokerMetricsTable + " SET local_time=" + strconv.FormatFloat(value, 'f', -1, 64) + " WHERE timestamp=" + strconv.Itoa(ts) + ` AND host="` + host + `";`
	_, err = Db.Exec(query)
	if err != nil {
		log.Logger.ErrorContext(ctx, fmt.Sprintf("local time update to %s table failed", brokerMetricsTable), err)
		return err
	}
	return nil
}

func UpdateBrokerTotalTime(ctx context.Context, value float64, host string, ts int) (err error) {
	query := "UPDATE " + brokerMetricsTable + " SET total_time=" + strconv.FormatFloat(value, 'f', -1, 64) + " WHERE timestamp=" + strconv.Itoa(ts) + ` AND host="` + host + `";`
	_, err = Db.Exec(query)
	if err != nil {
		log.Logger.ErrorContext(ctx, fmt.Sprintf("total time update to %s table failed", brokerMetricsTable), err)
		return err
	}
	return nil
}

func UpdateBrokerMaxLag(ctx context.Context, value float64, host string, ts int) (err error) {
	query := "UPDATE " + brokerMetricsTable + " SET max_lag=" + strconv.FormatFloat(value, 'f', -1, 64) + " WHERE timestamp=" + strconv.Itoa(ts) + ` AND host="` + host + `";`
	_, err = Db.Exec(query)
	if err != nil {
		log.Logger.ErrorContext(ctx, fmt.Sprintf("max lag update to %s table failed", brokerMetricsTable), err)
		return err
	}
	return nil
}

func UpdateBrokerUncleanLeaderElection(ctx context.Context, value float64, host string, ts int) (err error) {
	query := "UPDATE " + brokerMetricsTable + " SET unclean_lead_elec=" + strconv.FormatFloat(value, 'f', -1, 64) + " WHERE timestamp=" + strconv.Itoa(ts) + ` AND host="` + host + `";`
	_, err = Db.Exec(query)
	if err != nil {
		log.Logger.ErrorContext(ctx, fmt.Sprintf("unclean leader election update to %s table failed", brokerMetricsTable), err)
		return err
	}
	return nil
}

func UpdateBrokerFailedFetchRate(ctx context.Context, value float64, host string, ts int) (err error) {
	query := "UPDATE " + brokerMetricsTable + " SET failed_fetch_rate=" + strconv.FormatFloat(value, 'f', -1, 64) + " WHERE timestamp=" + strconv.Itoa(ts) + ` AND host="` + host + `";`
	_, err = Db.Exec(query)
	if err != nil {
		log.Logger.ErrorContext(ctx, fmt.Sprintf("failed fetch rate election update to %s table failed", brokerMetricsTable), err)
		return err
	}
	return nil
}

func UpdateBrokerFailedProdRate(ctx context.Context, value float64, host string, ts int) (err error) {
	query := "UPDATE " + brokerMetricsTable + " SET failed_prod_rate=" + strconv.FormatFloat(value, 'f', -1, 64) + " WHERE timestamp=" + strconv.Itoa(ts) + ` AND host="` + host + `";`
	_, err = Db.Exec(query)
	if err != nil {
		log.Logger.ErrorContext(ctx, fmt.Sprintf("failed produce rate election update to %s table failed", brokerMetricsTable), err)
		return err
	}
	return nil
}