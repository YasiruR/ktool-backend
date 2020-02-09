package database

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/YasiruR/ktool-backend/log"
	"github.com/go-sql-driver/mysql"
	"strconv"
	"time"
)

func AddNewBrokers(ctx context.Context, hosts []string, ports []int, clusterName string) (err error) {
	clusterId, err := GetClusterIdByName(ctx, clusterName)
	if err != nil {
		log.Logger.ErrorContext(ctx, "getting cluster id to add new zookeeper failed")
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

func GetBrokerByClusterId(ctx context.Context, clusterId int) (id int, brokerHost string, brokerPort int, err error) {
	query := "SELECT id, host, port FROM " + brokerTable + ` WHERE cluster_id="` + strconv.Itoa(clusterId) + `";`

	row := Db.QueryRow(query)

	switch err := row.Scan(&id, &brokerHost, &brokerPort); err {
	case sql.ErrNoRows:
		log.Logger.ErrorContext(ctx, "no broker scanned for the cluster_id", clusterId)
		return 0, "", 0, errors.New("row scan failed")
	case nil:
		log.Logger.TraceContext(ctx, "fetched broker by cluster id", clusterId)
		return id, brokerHost, brokerPort, nil
	default:
		log.Logger.ErrorContext(ctx, "unhandled error in row scan", clusterId)
		return 0, "", 0, errors.New("row scan failed")
	}
}