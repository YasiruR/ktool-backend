package database

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/YasiruR/ktool-backend/domain"
	"github.com/YasiruR/ktool-backend/log"
	"github.com/go-sql-driver/mysql"
	"strconv"
	"time"
)

func AddNewBrokers(ctx context.Context, hosts []string, ports []int, clusterName string) (err error) {
	clusterId, err := GetClusterIdByName(ctx, clusterName)
	if err != nil {
		log.Logger.ErrorContext(ctx, "getting cluster id to add new broker failed")
		return err
	}

	currentTime := time.Now()
	var rows string

	//creating multiple rows for the query
	for index, host := range hosts {
		rows += `(null, "` + host + `", ` + strconv.Itoa(ports[index]) +  ", " + strconv.Itoa(clusterId) + `, "` + currentTime.Format("2006-01-02 15:04:05") + `")`
		if index == len(hosts) - 1 {
			rows += ";"
		} else {
			rows += ","
		}
	}

	query := "INSERT INTO " + brokerTable + " ( id, host, port, cluster_id, created_at ) " + ` VALUES ` + rows

	insert, err := Db.Query(query)
	if err != nil {
		//check if the broker already exists in db
		if mysqlError, ok := err.(*mysql.MySQLError); ok {
			if mysqlError.Number == 1062 {
				log.Logger.ErrorContext(ctx, "at least one broker already exists in db")
				return errors.New("duplicate entry")
			}
		}
		log.Logger.ErrorContext(ctx, fmt.Sprintf("insert to %s table failed", brokerTable), err)
		return err
	}

	defer insert.Close()
	log.Logger.TraceContext(ctx, "add new brokers db query was successful", hosts, ports)

	return nil
}

func GetAllBrokers(ctx context.Context) (brokerList []domain.Broker, err error) {
	query := "SELECT * FROM " + brokerTable + ";"

	rows, err := Db.Query(query)
	if err != nil {
		log.Logger.ErrorContext(ctx, "get all brokers db query failed", err)
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		broker := domain.Broker{}

		err = rows.Scan(&broker.ID, &broker.Host, &broker.Port, &broker.CreatedAt, &broker.ClusterID)
		if err != nil {
			log.Logger.ErrorContext(ctx, "scanning rows in broker table failed", err)
			return nil, err
		}

		brokerList = append(brokerList, broker)
	}

	err = rows.Err()
	if err != nil {
		log.Logger.ErrorContext(ctx, "error occurred when scanning rows", err)
		return nil, err
	}

	log.Logger.TraceContext(ctx, "get all brokers db query was successful")
	return brokerList, nil
}

func GetBrokersByClusterId(ctx context.Context, clusterId int) (brokers []domain.Broker, err error) {
	query := "SELECT id, host, port, cluster_id FROM " + brokerTable + ` WHERE cluster_id="` + strconv.Itoa(clusterId) + `";`

	rows, err := Db.Query(query)
	if err != nil {
		log.Logger.ErrorContext(ctx, "get brokers for cluster db query failed", err)
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		broker := domain.Broker{}

		err = rows.Scan(&broker.ID, &broker.Host, &broker.Port, &broker.ClusterID)
		if err != nil {
			log.Logger.ErrorContext(ctx, "scanning rows in broker table failed", err)
			return nil, err
		}

		brokers = append(brokers, broker)
	}

	err = rows.Err()
	if err != nil {
		log.Logger.ErrorContext(ctx, "error occurred when scanning rows", err)
		return nil, err
	}

	//log.Logger.TraceContext(ctx, "get all brokers for cluster db query was successful", clusterId)
	return brokers, nil
}

func DeleteBrokersOfCluster(ctx context.Context, clusterID int) (err error) {
	query := "DELETE FROM " + brokerTable + " WHERE cluster_id=" + strconv.Itoa(clusterID) + ";"

	res, err := Db.Exec(query)
	if err != nil {
		log.Logger.ErrorContext(ctx, "deleting brokers for cluster failed", fmt.Sprintf("cluster id : %v", clusterID))
		return err
	}

	count, err := res.RowsAffected()
	if err != nil {
		log.Logger.WarnContext(ctx, "getting deleted broker count failed")
		return nil
	}

	log.Logger.TraceContext(ctx, fmt.Sprintf("%v brokers are delted for cluster id : %v", count, clusterID))
	return nil
}

func GetPreviousTimestamp(ctx context.Context, host string, ts int64) (prvTs int64, err error) {
	row := Db.QueryRow("SELECT timestamp FROM " + brokerMetricsTable + ` WHERE host="` + host +`" AND timestamp<=` + strconv.Itoa(int(ts)) + " ORDER BY ID DESC LIMIT 1;")
	err = row.Scan(&prvTs)
	switch err {
	case sql.ErrNoRows:
		log.Logger.ErrorContext(ctx, err, "query returned no rows", host, ts)
		return prvTs, err
	case nil:
		return prvTs, err
	default:
		log.Logger.ErrorContext(ctx, err, "query returned an error", host, ts)
		return prvTs, err
	}
}

func GetNextTimestamp(ctx context.Context, host string, ts int64) (nextTs int64, err error) {
	row := Db.QueryRow("SELECT timestamp FROM " + brokerMetricsTable + ` WHERE host="` + host +`" AND timestamp>=` + strconv.Itoa(int(ts)) + " ORDER BY ID LIMIT 1;")
	err = row.Scan(&nextTs)
	switch err {
	case sql.ErrNoRows:
		log.Logger.ErrorContext(ctx, err, "query returned no rows", host, ts)
		return nextTs, err
	case nil:
		return nextTs, err
	default:
		log.Logger.ErrorContext(ctx, err, "query returned an error", host, ts)
		return nextTs, err
	}
}

func GetBrokerMetricsByTimestampList(ctx context.Context, host string, tsList []int64) (brokerMetrics map[int64]domain.BrokerMetrics, err error) {
	brokerMetrics = make(map[int64]domain.BrokerMetrics)

	for _, ts := range tsList {
		var metrics domain.BrokerMetrics
		row := Db.QueryRow("SELECT timestamp, partitions, leaders, act_controller, offline_part, under_replicated, bytes_in, bytes_out, mesg_rate, isr_exp_rate, isr_shrink_rate, send_time, queue_time, remote_time, local_time, total_time, net_proc_avg_idle_perc, max_lag, unclean_lead_elec, failed_fetch_rate, failed_prod_rate, total_messages, topics FROM " + brokerMetricsTable + ` WHERE host="` + host + `" AND timestamp=` + strconv.Itoa(int(ts)) + `;`)

		var ts, bytesIn, bytesOut int64
		var partitions, leaders, actControllers, offlinePart, underReplicated, messages, topics sql.NullInt64
		var isrExp, isrShrink, sendTime, queueTime, localTime, remoteTime, totalTime, netIdle, maxLag, uncleanLeadElec, failedFetch, failedProd, mesgRate sql.NullFloat64
		err := row.Scan(&ts, &partitions, &leaders, &actControllers, &offlinePart, &underReplicated, &bytesIn, &bytesOut, &mesgRate, &isrExp, &isrShrink, &sendTime, &queueTime, &remoteTime, &localTime, &totalTime, &netIdle, &maxLag, &uncleanLeadElec, &failedFetch, &failedProd, &messages, &topics)
		switch err {
		case sql.ErrNoRows:
			log.Logger.ErrorContext(ctx, err, "query returned no rows", host, ts)
			//note : we can continue here if frontend validates ts when drawing the graph
			return nil, err
		case nil:
			if partitions.Valid {
				metrics.NumReplicas = int(partitions.Int64)
			}
			if leaders.Valid {
				metrics.NumLeaders = int(leaders.Int64)
			}
			if actControllers.Valid {
				metrics.NumActControllers = int(actControllers.Int64)
			}
			if offlinePart.Valid {
				metrics.OfflinePartitions = int(offlinePart.Int64)
			}
			if underReplicated.Valid {
				metrics.UnderReplicated = int(underReplicated.Int64)
			}
			metrics.ByteInRate = bytesIn
			metrics.ByteOutRate =- bytesOut
			if mesgRate.Valid {
				metrics.MessageRate = mesgRate.Float64
			}
			if isrExp.Valid {
				metrics.IsrExpansionRate = isrExp.Float64
			}
			if isrShrink.Valid {
				metrics.IsrShrinkRate = isrShrink.Float64
			}
			if sendTime.Valid {
				metrics.ResponseTime = sendTime.Float64
			}
			if queueTime.Valid {
				metrics.QueueTime = queueTime.Float64
			}
			if remoteTime.Valid {
				metrics.RemoteTime = remoteTime.Float64
			}
			if localTime.Valid {
				metrics.LocalTIme = localTime.Float64
			}
			if totalTime.Valid {
				metrics.TotalReqTime = totalTime.Float64
			}
			if netIdle.Valid {
				metrics.NetworkProcAvgIdlePercent = netIdle.Float64
			}
			if maxLag.Valid {
				metrics.MaxLagBtwLeadAndRepl = maxLag.Float64
			}
			if uncleanLeadElec.Valid {
				metrics.UncleanLeadElec = uncleanLeadElec.Float64
			}
			if failedFetch.Valid {
				metrics.FailedFetchReqRate = failedFetch.Float64
			}
			if failedProd.Valid {
				metrics.FailedProdReqRate = failedProd.Float64
			}
			if messages.Valid {
				metrics.Messages = messages.Int64
			}
			if topics.Valid {
				metrics.Topics = int(topics.Int64) - 1		//since prometheus adds one more stat per instance (aggregated count)
			}
			brokerMetrics[ts] = metrics

		default:
			log.Logger.ErrorContext(ctx, err, "query returned an error", host, ts)
			return nil, err
		}
	}
	log.Logger.DebugContext(ctx, brokerMetrics, tsList)
	return
}

func GetBrokerMetricsAverageValues(ctx context.Context, host string, startingTs, endingTs int64) (metrics domain.BrokerMetrics, err error) {
	//row := Db.QueryRow("SELECT AVG(bytes_in), AVG(bytes_out), AVG(mesg_rate), AVG(isr_exp_rate), AVG(isr_shrink_rate), AVG(send_time), AVG(queue_time), AVG(remote_time), AVG(local_time), AVG(total_time), AVG(net_proc_avg_idle_perc), AVG(max_lag), AVG(unclean_lead_elec), AVG(failed_fetch_rate), AVG(failed_prod_rate) FROM " + brokerMetricsTable + ` WHERE host="` + host + `" AND timestamp>` + strconv.Itoa(int(startingTs)) + ` AND timestamp<` + strconv.Itoa(int(endingTs)) + `;`)
	//row := Db.QueryRow("SELECT AVG(bytes_in), AVG(bytes_out), AVG(mesg_rate), AVG(isr_exp_rate), AVG(isr_shrink_rate), AVG(send_time), AVG(queue_time), AVG(remote_time), AVG(local_time), AVG(total_time), AVG(net_proc_avg_idle_perc), AVG(max_lag), AVG(unclean_lead_elec), AVG(failed_fetch_rate), AVG(failed_prod_rate) FROM " + brokerMetricsTable + ` WHERE host=? AND timestamp>? AND timestamp<?;`, host, startingTs, endingTs)

	stmt, err := Db.Prepare(`SELECT AVG(bytes_in), AVG(bytes_out), AVG(mesg_rate), AVG(isr_exp_rate), AVG(isr_shrink_rate), AVG(send_time), AVG(queue_time), AVG(remote_time), AVG(local_time), AVG(total_time), AVG(net_proc_avg_idle_perc), AVG(max_lag), AVG(unclean_lead_elec), AVG(failed_fetch_rate), AVG(failed_prod_rate) FROM ` + brokerMetricsTable + ` WHERE host=? AND timestamp>? AND timestamp<?;`)
	if err != nil {
		log.Logger.ErrorContext(ctx, err, "preparing db statement failed", host)
		return metrics, err
	}
	defer stmt.Close()

	//var bytesIn, bytesOut []uint8
	var bytesIn, bytesOut, isrExp, isrShrink, sendTime, queueTime, localTime, remoteTime, totalTime, netIdle, maxLag, uncleanLeadElec, failedFetch, failedProd, mesgRate sql.NullFloat64

	err = stmt.QueryRow(host, startingTs, endingTs).Scan(&bytesIn, &bytesOut, &mesgRate, &isrExp, &isrShrink, &sendTime, &queueTime, &remoteTime, &localTime, &totalTime, &netIdle, &maxLag, &uncleanLeadElec, &failedFetch, &failedProd)
	//err = row.Scan(&bytesIn, &bytesOut, &mesgRate, &isrExp, &isrShrink, &sendTime, &queueTime, &remoteTime, &localTime, &totalTime, &netIdle, &maxLag, &uncleanLeadElec, &failedFetch, &failedProd)
	switch err {
	case sql.ErrNoRows:
		log.Logger.ErrorContext(ctx, err, "query returned no rows", host)
		//note : we can continue here if frontend validates ts when drawing the graph
		return metrics, err
	case nil:
		if bytesOut.Valid {
			metrics.ByteOutRate = int64(bytesOut.Float64)
		}
		if bytesIn.Valid {
			metrics.ByteInRate = int64(bytesIn.Float64)
		}

		if mesgRate.Valid {
			metrics.MessageRate = mesgRate.Float64
		}
		if isrExp.Valid {
			metrics.IsrExpansionRate = isrExp.Float64
		}
		if isrShrink.Valid {
			metrics.IsrShrinkRate = isrShrink.Float64
		}
		if sendTime.Valid {
			metrics.ResponseTime = sendTime.Float64
		}
		if queueTime.Valid {
			metrics.QueueTime = queueTime.Float64
		}
		if remoteTime.Valid {
			metrics.RemoteTime = remoteTime.Float64
		}
		if localTime.Valid {
			metrics.LocalTIme = localTime.Float64
		}
		if totalTime.Valid {
			metrics.TotalReqTime = totalTime.Float64
		}
		if netIdle.Valid {
			metrics.NetworkProcAvgIdlePercent = netIdle.Float64
		}
		if maxLag.Valid {
			metrics.MaxLagBtwLeadAndRepl = maxLag.Float64
		}
		if uncleanLeadElec.Valid {
			metrics.UncleanLeadElec = uncleanLeadElec.Float64
		}
		if failedFetch.Valid {
			metrics.FailedFetchReqRate = failedFetch.Float64
		}
		if failedProd.Valid {
			metrics.FailedProdReqRate = failedProd.Float64
		}
	default:
		log.Logger.ErrorContext(ctx, err, "query returned an error", host)
		return metrics, err
	}
	return
}

//returns metrics of a broker wrt a timestamp
func GetBrokerMetrics(ctx context.Context, host string, tsCount int) (brokerMetrics map[int64]domain.BrokerMetrics, err error) {
	brokerMetrics = make(map[int64]domain.BrokerMetrics)
	rows, err := Db.Query("SELECT timestamp, partitions, leaders, act_controller, offline_part, under_replicated, bytes_in, bytes_out, mesg_rate, isr_exp_rate, isr_shrink_rate, send_time, queue_time, remote_time, local_time, total_time, net_proc_avg_idle_perc, max_lag, unclean_lead_elec, failed_fetch_rate, failed_prod_rate, total_messages, topics FROM " + brokerMetricsTable + ` WHERE host="` + host + `" ORDER BY ID DESC LIMIT ` + strconv.Itoa(tsCount) + `;`)
	if err != nil {
		log.Logger.ErrorContext(ctx, "get broker metrics query failed", err)
		return brokerMetrics, err
	}
	defer rows.Close()
	for rows.Next() {
		var metrics domain.BrokerMetrics
		var ts, bytesIn, bytesOut int64
		var partitions, leaders, actControllers, offlinePart, underReplicated, messages, topics sql.NullInt64
		var isrExp, isrShrink, sendTime, queueTime, localTime, remoteTime, totalTime, netIdle, maxLag, uncleanLeadElec, failedFetch, failedProd, mesgRate sql.NullFloat64
		err = rows.Scan(&ts, &partitions, &leaders, &actControllers, &offlinePart, &underReplicated, &bytesIn, &bytesOut, &mesgRate, &isrExp, &isrShrink, &sendTime, &queueTime, &remoteTime, &localTime, &totalTime, &netIdle, &maxLag, &uncleanLeadElec, &failedFetch, &failedProd, &messages, &topics)
		if err != nil {
			log.Logger.ErrorContext(ctx, "scanning rows in broker metrics table failed", err)
			return
		}

		if partitions.Valid {
			metrics.NumReplicas = int(partitions.Int64)
		}
		if leaders.Valid {
			metrics.NumLeaders = int(leaders.Int64)
		}
		if actControllers.Valid {
			metrics.NumActControllers = int(actControllers.Int64)
		}
		if offlinePart.Valid {
			metrics.OfflinePartitions = int(offlinePart.Int64)
		}
		if underReplicated.Valid {
			metrics.UnderReplicated = int(underReplicated.Int64)
		}
		metrics.ByteInRate = bytesIn
		metrics.ByteOutRate =- bytesOut
		if mesgRate.Valid {
			metrics.MessageRate = mesgRate.Float64
		}
		if isrExp.Valid {
			metrics.IsrExpansionRate = isrExp.Float64
		}
		if isrShrink.Valid {
			metrics.IsrShrinkRate = isrShrink.Float64
		}
		if sendTime.Valid {
			metrics.ResponseTime = sendTime.Float64
		}
		if queueTime.Valid {
			metrics.QueueTime = queueTime.Float64
		}
		if remoteTime.Valid {
			metrics.RemoteTime = remoteTime.Float64
		}
		if localTime.Valid {
			metrics.LocalTIme = localTime.Float64
		}
		if totalTime.Valid {
			metrics.TotalReqTime = totalTime.Float64
		}
		if netIdle.Valid {
			metrics.NetworkProcAvgIdlePercent = netIdle.Float64
		}
		if maxLag.Valid {
			metrics.MaxLagBtwLeadAndRepl = maxLag.Float64
		}
		if uncleanLeadElec.Valid {
			metrics.UncleanLeadElec = uncleanLeadElec.Float64
		}
		if failedFetch.Valid {
			metrics.FailedFetchReqRate = failedFetch.Float64
		}
		if failedProd.Valid {
			metrics.FailedProdReqRate = failedProd.Float64
		}
		if messages.Valid {
			metrics.Messages = messages.Int64
		}
		if topics.Valid {
			metrics.Topics = int(topics.Int64) - 1		//since prometheus adds one more stat per instance (aggregated count)
		}

		brokerMetrics[ts] = metrics
	}

	err = rows.Err()
	if err != nil {
		log.Logger.ErrorContext(ctx, "error occurred when scanning rows", err)
		return
	}
	return
}