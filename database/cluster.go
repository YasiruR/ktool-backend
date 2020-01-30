package database

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/YasiruR/ktool-backend/domain"
	"github.com/YasiruR/ktool-backend/log"
)

func AddNewCluster(ctx context.Context, clusterName string, kafkaVersion string) (err error) {
	query := "INSERT INTO " + clusterTable + " ( id, cluster_name, kafka_version ) " + ` VALUES ( null, "` + clusterName + `", "` + kafkaVersion + `" )`

	insert, err := Db.Query(query)
	if err != nil {
		log.Logger.ErrorContext(ctx, fmt.Sprintf("insert to %s table failed", clusterTable), err)
		return err
	}

	defer insert.Close()
	log.Logger.TraceContext(ctx, "add new cluster db query was successful", clusterName)

	return nil
}

func DeleteCluster(ctx context.Context, clusterName string) (err error) {
	query := "DELETE FROM " + clusterTable + `  WHERE cluster_name="` + clusterName + `";`

	del, err := Db.Query(query)
	if err != nil {
		log.Logger.ErrorContext(ctx, fmt.Sprintf("deleting cluster %s failed", clusterName), err)
		return err
	}

	defer del.Close()
	log.Logger.TraceContext(ctx, "delete cluster db query was successful", clusterName)

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

func GetClusterIdByName(ctx context.Context, clusterName string) (clusterId int, err error) {
	query := "SELECT id FROM " + clusterTable + ` WHERE cluster_name="` + clusterName + `";`

	row := Db.QueryRow(query)

	switch err := row.Scan(&clusterId); err {
	case sql.ErrNoRows:
		log.Logger.ErrorContext(ctx, "no rows scanned for the cluster", clusterName)
		return 0, errors.New("row scan failed")
	case nil:
		log.Logger.TraceContext(ctx, "fetched cluster by cluster id", clusterName)
		return clusterId, nil
	default:
		log.Logger.ErrorContext(ctx, "unhandled error in row scan", clusterName)
		return 0, errors.New("row scan failed")
	}
}