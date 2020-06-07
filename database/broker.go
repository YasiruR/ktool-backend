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

//returns metrics of a broker wrt a timestamp
func GetBrokerMetrics(ctx context.Context, host string, tsCount int) (brokerMetrics map[int64]domain.BrokerMetrics, err error) {
	brokerMetrics = make(map[int64]domain.BrokerMetrics)
	var metrics domain.BrokerMetrics
	rows, err := Db.Query("SELECT timestamp, partitions, leaders, act_controller, offline_part, under_replicated, bytes_in, bytes_out, mesg_rate, isr_exp_rate, isr_shrink_rate, send_time, queue_time, remote_time, local_time, total_time, net_proc_avg_idle_perc, max_lag, unclean_lead_elec, failed_fetch_rate, failed_prod_rate, total_messages, topics FROM " + brokerMetricsTable + ` WHERE host="` + host + `" ORDER BY ID DESC LIMIT ` + strconv.Itoa(tsCount) + `;`)
	if err != nil {
		log.Logger.ErrorContext(ctx, "get broker metrics query failed", err)
		return brokerMetrics, err
	}
	defer rows.Close()
	for rows.Next() {
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
			metrics.Topics = int(topics.Int64)
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

//func GetBrokerMetrics(ctx context.Context, host string) (byteRateIn, byteRateOut map[int64]int64, err error) {
//
//	byteRateIn = make(map[int64]int64)
//	byteRateOut = make(map[int64]int64)
//
//	rows, err := Db.Query("SELECT bytes_in, UNIX_TIMESTAMP(created_at) FROM " + brokerBytesInTable + ` WHERE host="` + host + `" ORDER BY ID DESC LIMIT ` + strconv.Itoa(metricsLimit) + `;`)
//	if err != nil {
//		log.Logger.ErrorContext(ctx, "get broker bytes in query failed", err)
//		return nil, nil, err
//	}
//	defer rows.Close()
//	for rows.Next() {
//		var byteRate, ts int64
//		err = rows.Scan(&byteRate, &ts)
//		if err != nil {
//			log.Logger.ErrorContext(ctx, "scanning rows in broker bytes in table failed", err)
//			return
//		}
//
//		byteRateIn[ts] = byteRate
//	}
//
//	err = rows.Err()
//	if err != nil {
//		log.Logger.ErrorContext(ctx, "error occurred when scanning rows", err)
//		return
//	}
//
//	rows, err = Db.Query("SELECT bytes_out, UNIX_TIMESTAMP(created_at) FROM " + brokerBytesOutTable + ` WHERE host="` + host + `" ORDER BY ID DESC LIMIT ` + strconv.Itoa(metricsLimit) + `;`)
//	if err != nil {
//		log.Logger.ErrorContext(ctx, "get broker bytes out query failed", err)
//		return nil, nil, err
//	}
//	defer rows.Close()
//	for rows.Next() {
//		var byteRate, ts int64
//		err = rows.Scan(&byteRate, &ts)
//		if err != nil {
//			log.Logger.ErrorContext(ctx, "scanning rows in broker bytes out table failed", err)
//			return
//		}
//
//		byteRateOut[ts] = byteRate
//	}
//
//	err = rows.Err()
//	if err != nil {
//		log.Logger.ErrorContext(ctx, "error occurred when scanning rows", err)
//		return
//	}
//
//	return
//}