package database

import (
	"context"
	"fmt"
	"github.com/YasiruR/ktool-backend/log"
)

func AddNewCluster(ctx context.Context, clusterName string, clusterVersion float64, zookeeperHost string, zookeeperPort int) (err error) {
	query := "INSERT INTO " + clusterTable + " ( id, cluster_name, cluster_version ) " + ` VALUES ( null, "` + clusterName + `", "` + fmt.Sprintf("%f", clusterVersion) + `" )`

	insert, err := Db.Query(query)
	if err != nil {
		log.Logger.ErrorContext(ctx, fmt.Sprintf("insert to %s table failed", clusterTable), err)
		return err
	}

	defer insert.Close()
	log.Logger.TraceContext(ctx, "add new cluster db query was successful", clusterName)

	return nil
}