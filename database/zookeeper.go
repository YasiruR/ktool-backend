package database

import (
	"context"
	"fmt"
	"github.com/YasiruR/ktool-backend/log"
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