package database

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/YasiruR/ktool-backend/log"
	"strconv"
)

func AddNewZookeeper(ctx context.Context, zookeeperHost string, zookeeperPort int, clusterName string) (err error) {
	clusterId, err := GetClusterIdByName(ctx, clusterName)
	if err != nil {
		log.Logger.ErrorContext(ctx, "getting cluster id to add new zookeeper failed")
		return err
	}

	query := "INSERT INTO " + zookeeperTable + " ( id, host, port, cluster_id ) " + ` VALUES ( null, "` + zookeeperHost + `", ` + fmt.Sprintf("%v", zookeeperPort) + ", " + fmt.Sprintf("%v", clusterId) +" )"

	insert, err := Db.Query(query)
	if err != nil {
		log.Logger.ErrorContext(ctx, fmt.Sprintf("insert to %s table failed", zookeeperTable), err)
		return err
	}

	defer insert.Close()
	log.Logger.TraceContext(ctx, "add new zookeeper db query was successful", zookeeperHost, zookeeperPort)

	return nil
}

func GetZookeeperByClusterId(ctx context.Context, clusterId int) (id int, zookeeperHost string, zookeeperPort int, err error) {
	query := "SELECT id, host, port FROM " + zookeeperTable + ` WHERE cluster_id="` + strconv.Itoa(clusterId) + `";`

	row := Db.QueryRow(query)

	switch err := row.Scan(&id, &zookeeperHost, &zookeeperPort); err {
	case sql.ErrNoRows:
		log.Logger.ErrorContext(ctx, "no zookeeper scanned for the cluster_id", clusterId)
		return 0, "", 0, errors.New("row scan failed")
	case nil:
		log.Logger.TraceContext(ctx, "fetched cluster by cluster id", clusterId)
		return id, zookeeperHost, zookeeperPort, nil
	default:
		log.Logger.ErrorContext(ctx, "unhandled error in row scan", clusterId)
		return 0, "", 0, errors.New("row scan failed")
	}
}