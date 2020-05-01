package database

import (
	"context"
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

	log.Logger.TraceContext(ctx, "get all brokers for cluster db query was successful", clusterId)
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

func UpdateBrokerByteInRate(ctx context.Context, bytesIn float64, host string) (err error) {
	query := "UPDATE " + brokerTable + " SET bytes_in=" + strconv.FormatFloat(bytesIn, 'f', -1, 64) + ` WHERE host="` + host + `";`

	tx, err := Db.Begin()
	if err != nil {
		log.Logger.ErrorContext(ctx, err,"starting the transaction failed", host)
		return err
	}
	defer tx.Rollback()

	statement, err := tx.Prepare(query)
	if err != nil {
		log.Logger.ErrorContext(ctx, err,"preparing the query failed", host)
		return err
	}

	_, err = statement.Exec()
	if err != nil {
		log.Logger.ErrorContext(ctx, err,"executing the update token query failed", host)
		return err
	}

	err = tx.Commit()
	if err != nil {
		log.Logger.ErrorContext(ctx, "committing the transaction failed", host)
		return err
	}

	return nil
}

func UpdateBrokerByteOutRate(ctx context.Context, bytesOut float64, host string) (err error) {
	query := "UPDATE " + brokerTable + " SET bytes_out=" + strconv.FormatFloat(bytesOut, 'f', -1, 64) + ` WHERE host="` + host + `";`

	tx, err := Db.Begin()
	if err != nil {
		log.Logger.ErrorContext(ctx, err,"starting the transaction failed", host)
		return err
	}
	defer tx.Rollback()

	statement, err := tx.Prepare(query)
	if err != nil {
		log.Logger.ErrorContext(ctx, "preparing the query failed", host)
		return err
	}

	_, err = statement.Exec()
	if err != nil {
		log.Logger.ErrorContext(ctx, err,"executing the update token query failed", host)
		return err
	}

	err = tx.Commit()
	if err != nil {
		log.Logger.ErrorContext(ctx, err,"committing the transaction failed", host)
		return err
	}

	return nil
}