package database

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/YasiruR/ktool-backend/domain"
	"github.com/YasiruR/ktool-backend/log"
)

func GetAllKubernetesClusters(ctx context.Context, userId int) (clusterResponse domain.ClusterResponse) {
	query := fmt.Sprintf("SELECT s.id, s.cluster_id, s.name, s.service_provider, s.status, s.created_on, u.zone,"+
		" u.project_id FROM %s s, %s u WHERE s.user_id = %d AND s.op_id = u.name AND s.active = 1;", k8sTable, operationsTable, userId)

	rows, err := Db.Query(query)

	switch err {
	case nil:
		log.Logger.InfoContext(ctx, "get secret query success")
	case sql.ErrNoRows:
		log.Logger.InfoContext(ctx, "no secrets found for userId %s", userId)
		clusterResponse.Error = err
		clusterResponse.Status = -1
		clusterResponse.Message = fmt.Sprintf("no secrets found for userId %d", userId)
		return clusterResponse
	default:
		log.Logger.InfoContext(ctx, "unhandled error occurred while fetching records for userId %s", userId)
		clusterResponse.Error = err
		clusterResponse.Status = -1
		clusterResponse.Message = "unhandled error occurred from db"
		return clusterResponse
	}

	defer rows.Close()
	kubClusterList := make([]domain.KubCluster, 0)

	for rows.Next() {
		cluster := domain.KubCluster{}

		err = rows.Scan(&cluster.Id, &cluster.ClusterId, &cluster.ClusterName, &cluster.ServiceProvider, &cluster.Status,
			&cluster.CreatedOn, &cluster.Location, &cluster.ProjectName)
		if err != nil {
			log.Logger.ErrorContext(ctx, "scanning rows in cluster table failed", err)
			clusterResponse.Error = err
			clusterResponse.Status = -1
			clusterResponse.Message = "scanning rows in cluster table failed"
			return clusterResponse
		}
		kubClusterList = append(kubClusterList, cluster)
	}

	log.Logger.TraceContext(ctx, "get all clusters db query was successful")
	clusterResponse.Clusters = kubClusterList
	clusterResponse.Status = 0
	clusterResponse.Message = "Success"
	return clusterResponse
}

func GetGkeLROperation(ctx context.Context, name string) (result domain.GkeLROperation) {

	query := fmt.Sprintf("SELECT id, project_id, name, zone FROM %s  WHERE name = '%s'", operationsTable, name)

	rows, err := Db.Query(query)

	if err != nil {
		result.Error = err
		return result
	}

	defer rows.Close()

	for rows.Next() {

		err = rows.Scan(&result.Id, &result.ProjectId, &result.Name, &result.Zone)
		if err != nil {
			log.Logger.ErrorContext(ctx, "scanning rows in operations table failed", err)
			result.Error = err
			return result
		}
	}

	return result
}

func AddGkeLROperation(ctx context.Context, Name string, ProjectId string, Zone string) (err error) {
	//TODO: validate req params
	//TODO: call a stored procedure
	query := fmt.Sprintf("INSERT INTO kdb.operations (name, project_id, zone) VALUES('%s', '%s', '%s')", Name, ProjectId, Zone)

	insert, err := Db.Query(query)

	if err != nil {
		log.Logger.ErrorContext(ctx, fmt.Sprintf("insert to %s table failed", operationsTable), err)
		return err
	}

	defer insert.Close()
	log.Logger.TraceContext(ctx, "successfully added a new operation ", Name)
	return nil
}

func UpdateGkeLROperation(ctx context.Context, name string, status string) (opStatus bool, err error) {
	query := fmt.Sprintf("UPDATE kdb.%s SET status='%s' WHERE name='%s'", operationsTable, status, name)

	insert, err := Db.Query(query)

	if err != nil {
		log.Logger.ErrorContext(ctx, fmt.Sprintf("update %s table failed", operationsTable), err)
		return false, err
	}

	defer insert.Close()
	log.Logger.TraceContext(ctx, "successfully updated operation ", name)
	return true, nil
}

func AddGkeCluster(ctx context.Context, clusterId string, userId int, clusterName string, operationName string) (err error) {
	query := fmt.Sprintf("INSERT INTO kdb.%s (cluster_id, user_id, name, op_id, service_provider, status, active) "+
		"VALUES ('%s', %d, '%s', '%s', '%s', 'CREATING', 1);", k8sTable, clusterId, userId, clusterName, operationName, "google")
	insert, err := Db.Query(query)

	if err != nil {
		log.Logger.ErrorContext(ctx, fmt.Sprintf("insert to %s table failed", k8sTable), err)
		return err
	}

	defer insert.Close()
	log.Logger.TraceContext(ctx, "successfully added a new cluster.")
	//todo: get id and return
	return nil
}

func UpdateGkeClusterCreationStatus(ctx context.Context, status string, operationId string) (opStatus bool, err error) {
	query := fmt.Sprintf("UPDATE kdb.%s SET status='%s' WHERE op_id='%s'", k8sTable, status, operationId)

	insert, err := Db.Query(query)

	if err != nil {
		log.Logger.ErrorContext(ctx, fmt.Sprintf("update %s table failed", k8sTable), err)
		return false, err
	}

	defer insert.Close()
	log.Logger.TraceContext(ctx, "successfully updated cluster ", operationId)
	return true, nil
}
