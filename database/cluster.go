package database

import (
	"context"
	"fmt"
	"github.com/YasiruR/ktool-backend/domain"
	"github.com/YasiruR/ktool-backend/log"
)

func AddNewCluster(ctx context.Context, clusterName string, kafkaVersion float64, zookeeperHost string, zookeeperPort int) (err error) {
	query := "INSERT INTO " + clusterTable + " ( id, cluster_name, kafka_version ) " + ` VALUES ( null, "` + clusterName + `", "` + fmt.Sprintf("%f", kafkaVersion) + `" )`

	insert, err := Db.Query(query)
	if err != nil {
		log.Logger.ErrorContext(ctx, fmt.Sprintf("insert to %s table failed", clusterTable), err)
		return err
	}

	defer insert.Close()
	log.Logger.TraceContext(ctx, "add new cluster db query was successful", clusterName)

	return nil
}

func GetAllClusters(ctx context.Context) (clusterList []domain.Cluster, err error) {
	query := "SELECT * FROM cluster;"

	rows, err := Db.Query(query)
	if err != nil {
		log.Logger.ErrorContext(ctx, "get all db query failed", err)
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		cluster := domain.Cluster{}

		err = rows.Scan(&cluster.ID, &cluster.ClusterName, &cluster.KafkaVersion, &cluster.ActiveControllers)
		if err != nil {
			log.Logger.ErrorContext(ctx, "scanning rows in cluster table failed", err)
			return nil, err
		}

		clusterList = append(clusterList, cluster)
	}

	err = rows.Err()
	if err != nil {
		log.Logger.ErrorContext(ctx, "error occurred when scanning rows", err)
		return nil, err
	}

	log.Logger.TraceContext(ctx, "get all clusters db query was successful")
	return clusterList, nil
}