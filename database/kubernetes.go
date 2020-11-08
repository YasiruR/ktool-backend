package database

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/YasiruR/ktool-backend/domain"
	"github.com/YasiruR/ktool-backend/log"
	"github.com/YasiruR/ktool-backend/util"
	"strconv"
)

//facade
func GetAllKubernetesClustersForUser(ctx context.Context, userId int) (clusterResponse domain.ClusterResponse) {
	query := fmt.Sprintf("SELECT s.id, s.cluster_id, s.name, s.service_provider, s.status, s.created_on, s.zone, s.op_id "+
		" FROM %s s WHERE s.user_id = %d AND s.active = 1;", k8sTable, userId)

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
		log.Logger.InfoContext(ctx, fmt.Sprintf("unhandled error occurred while fetching records for userId %d", userId))
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
			&cluster.CreatedOn, &cluster.Location, &cluster.Reference)
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

func GetKubernetesResourcesRecommendation(ctx context.Context, Provider string, Continent []string, VCPU string, RAM string, Network []string, Type []string, MinNodes string, MaxNodes string) (result domain.GkeRecommendations) {
	nodeCount, _ := strconv.Atoi(MinNodes)
	memory, _ := strconv.Atoi(RAM)
	processor, _ := strconv.Atoi(VCPU)
	regions := util.StringListToEscapedCSV(Continent)
	category := util.StringListToEscapedCSV(Type)
	network := util.StringListToEscapedCSV(Network)
	//regionOk 	:= false
	//categoryOk 	:= false
	//networkOk 	:= false

	provider := 1
	if Provider == "google" {
		provider = 0
	}

	baseQuery :=
		"SELECT " +
			"p.`type` AS type, " +
			"r.region AS region, " +
			//"p.memory AS memory, " +
			"p.memory * %d AS memory, " +
			//"p.cpu AS cpu, " +
			"p.cpu * %d AS processor, " +
			"r.unit_price * %d AS cost, " +
			//"r.unit_price AS unit_price, " +
			"p.network AS network, " +
			"%d AS node_count, " +
			"'10 min' AS startup_time " +
			"FROM " +
			"%s p, " +
			"%s r " +
			"WHERE " +
			"p.cpu >= %d / %d AND " +
			"p.memory >= %d / %d AND " +
			"r.product=p.id AND " +
			"p.provider = %d "

	if len(regions) > 0 {
		baseQuery += fmt.Sprintf("AND r.region IN (SELECT region_id FROM %s WHERE continent IN (%s)) ", locationsTable, regions)
	}
	if len(category) > 0 {
		baseQuery += fmt.Sprintf("AND p.category IN (%s) ", category)
	}
	if len(network) > 0 {
		baseQuery += fmt.Sprintf("AND p.network IN (%s) ", network)
	}
	query := fmt.Sprintf(baseQuery+"ORDER BY cost LIMIT 6;", nodeCount, nodeCount, nodeCount, nodeCount, productsTable, priceTable, processor, nodeCount, memory, nodeCount, provider)

	rows, err := Db.Query(query)

	switch err {
	case nil:
		log.Logger.InfoContext(ctx, "get recommendations query success")
	case sql.ErrNoRows:
		log.Logger.InfoContext(ctx, "no recommendations found for parameters %s", Provider)
		result.Detail = fmt.Sprint("No recommendations found for parameters")
		result.Status = -1
		return result
	default:
		log.Logger.InfoContext(ctx, "no recommendations found for parameters %s", Provider)
		result.Detail = fmt.Sprint("No recommendations found for parameters")
		result.Status = -1
		return result
	}

	defer rows.Close()
	nodes := make([]domain.Node, 0)

	for rows.Next() {
		node := domain.Node{}

		err = rows.Scan(&node.Type, &node.Region, &node.Memory, &node.Processor, &node.Cost, &node.Network, &node.NodeCount, &node.StartupTime)
		if err != nil {
			log.Logger.ErrorContext(ctx, "scanning rows in recommendation query failed", err)
			result.Detail = fmt.Sprint("internal error occurred")
			result.Status = -1
			return result
		}
		nodes = append(nodes, node)
	}
	if len(nodes) > 0 {
		log.Logger.TraceContext(ctx, "get recommendation db query was successful")
		result.Nodes = nodes
		result.Detail = "Success"
		result.Status = 0
	} else {
		log.Logger.TraceContext(ctx, "get recommendation db query returned 0 rows")
		result.Nodes = nodes
		result.Detail = "Failed"
		result.Status = -1
	}
	return result
}

func GetKubernetesResources(ctx context.Context, Provider string) (result domain.GkeResources) {
	//continentQuery := "SELECT DISTINCT(continent) as continents FROM locations;"
	//regiondNameQuery := "SELECT DISTINCT(region_name) as continents FROM locations;"
	//regionIdQuery := "SELECT DISTINCT(region_id) as continents FROM locations;"

	provider := 1
	if Provider == "google" {
		provider = 0
	}

	query := fmt.Sprintf("SELECT continent, region_name, region_id FROM %s WHERE provider = %d;", locationsTable, provider)

	rows, err := Db.Query(query)

	switch err {
	case nil:
		log.Logger.InfoContext(ctx, "get locations query success")
	case sql.ErrNoRows:
		log.Logger.InfoContext(ctx, "no metadata found for locations")
		result.Detail = "no locations found"
		result.Status = -1
		return result
	default:
		log.Logger.InfoContext(ctx, "unhandled error occurred while fetching locations")
		result.Detail = "no locations found"
		result.Status = -1
		return result
	}

	defer rows.Close()
	continents := make(map[string]bool)
	regionNames := make([]string, 0)
	regionIds := make([]string, 0)

	for rows.Next() {
		location := domain.ResourceLocation{}

		err = rows.Scan(&location.Continent, &location.RegionName, &location.RegionId)
		if err != nil {
			log.Logger.ErrorContext(ctx, "scanning rows in location failed", err)
			result.Detail = "scanning rows in location failed"
			result.Status = -1
			return result
		}
		continents[location.Continent] = true
		regionNames = append(regionNames, location.RegionName)
		regionIds = append(regionIds, location.RegionId)
	}

	log.Logger.TraceContext(ctx, "get all resource locations db query was successful")
	for k := range continents {
		result.Continents = append(result.Continents, k)
	}
	//result.Continents = continents
	result.RegionIds = regionIds
	result.RegionNames = regionNames
	result.Detail = "Success"
	result.Status = 0
	return result
}

func GetAllRunningKubernetesClusters(ctx context.Context) (result domain.ClusterResponse) {

	query := fmt.Sprintf("SELECT k.id, k.cluster_id, k.name, k.user_id, k.status, k.service_provider, k.zone, s.id, IFNULL(s.gkeProjectId, '') AS project_id FROM %s k, %s s WHERE k.active = 1 AND k.service_provider = s.provider AND k.user_id = s.ownerId AND s.deleted = 0;", k8sTable, secretTable)

	rows, err := Db.Query(query)

	switch err {
	case nil:
		log.Logger.InfoContext(ctx, "get all kub clusters query success")
	case sql.ErrNoRows:
		log.Logger.InfoContext(ctx, "no running kub clusters found")
		result.Message = "no clusters found"
		result.Status = -1
		return result
	default:
		log.Logger.InfoContext(ctx, "unhandled error occurred while fetching kub clusters")
		result.Message = "error occurred"
		result.Status = -1
		return result
	}

	defer rows.Close()
	clusters := make([]domain.KubCluster, 0)

	for rows.Next() {
		cluster := domain.KubCluster{}

		err = rows.Scan(&cluster.Id, &cluster.ClusterId, &cluster.ClusterName, &cluster.ClusterId, &cluster.Status, &cluster.ServiceProvider, &cluster.Location, &cluster.SecretId, &cluster.ProjectName)
		if err != nil {
			log.Logger.ErrorContext(ctx, "scanning rows in kub clusters failed", err)
			result.Message = "scanning rows in kub clusters failed"
			result.Status = -1
			return result
		}
		clusters = append(clusters, cluster)
	}

	log.Logger.TraceContext(ctx, "get all kub clusters db query was successful")
	result.Clusters = clusters
	result.Message = "Success"
	result.Status = 0
	return result
}

func ValidateClusterName(ctx context.Context, userId string, name string, provider string) (result domain.ValidationResponse) {
	query := fmt.Sprintf("SELECT CASE WHEN list.count > 0 THEN 0 ELSE 1 END AS status FROM (SELECT count(*) as count FROM %s WHERE user_id = '%s' AND name = '%s' AND active = 1) list;", k8sTable, userId, name)

	rows, err := Db.Query(query)

	switch err {
	case nil:
		log.Logger.InfoContext(ctx, "validate cluster name query success")
	case sql.ErrNoRows:
		log.Logger.InfoContext(ctx, "no clusters found for name "+name)
		result.Message = "no clusters found"
		result.Status = -1
		return result
	default:
		log.Logger.InfoContext(ctx, "unhandled error occurred while validating cluster name")
		result.Message = "error occurred"
		result.Status = -1
		return result
	}

	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(&result.Status)
		if err != nil {
			log.Logger.ErrorContext(ctx, "scanning rows in validation request failed", err)
			result.Message = "scanning rows in kub clusters failed"
			result.Status = -1
			return result
		}
	}

	log.Logger.TraceContext(ctx, "validate cluster name db query was successful")
	result.Message = "Success"
	//result.Status = 0
	return result
}

//gke
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
	query := fmt.Sprintf("INSERT INTO kdb.%s (name, project_id, zone) VALUES('%s', '%s', '%s')", operationsTable, Name, ProjectId, Zone)

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

func AddGkeCluster(ctx context.Context, clusterId string, userId int, clusterName string, operationName string, zone string) (err error) {
	query := fmt.Sprintf("INSERT INTO kdb.%s (cluster_id, user_id, name, op_id, zone, service_provider, status, active) "+
		"VALUES ('%s', %d, '%s', '%s', '%s', '%s', 'CREATING', 1);", k8sTable, clusterId, userId, clusterName, operationName, zone, "google")
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

func CheckGkeClusterCreationStatus(opID string, userID string) (status domain.GkeClusterStatus, error error) {
	//status.Name = clusterName
	query := fmt.Sprintf("SELECT status, cluster_id, name FROM %s  WHERE op_id = '%s' AND user_id = '%s'", k8sTable, opID, userID)
	rows, err := Db.Query(query)
	if err != nil {
		status.Error = err.Error()
		status.Status = "CHECK FAILED"
		return status, err
	}
	defer rows.Close()
	for rows.Next() {

		err = rows.Scan(&status.Status, &status.ClusterId, &status.Name)
		if err != nil {
			status.Error = err.Error()
			status.Status = "CHECK FAILED"
			return status, err
		}
	}
	return status, nil
}

func UpdateGkeClusterCreationStatus(ctx context.Context, status string, operationId string) (opStatus bool, err error) {
	statusDesc := "UNSPECIFIED"
	switch status {
	case "SUBMITTED":
		statusDesc = "INITIALIZING"
	case "RUNNING":
		statusDesc = "CREATING"
	case "DONE":
		statusDesc = "RUNNING"
	default:
	}
	query := fmt.Sprintf("UPDATE kdb.%s SET status='%s' WHERE op_id='%s'", k8sTable, statusDesc, operationId)

	insert, err := Db.Query(query)

	if err != nil {
		log.Logger.ErrorContext(ctx, fmt.Sprintf("update %s table failed", k8sTable), err)
		return false, err
	}

	defer insert.Close()
	log.Logger.TraceContext(ctx, "successfully updated cluster ", operationId)
	return true, nil
}

func UpdateGkeClusterStatusById(ctx context.Context, isActive int, clusterId int, status string) (opStatus bool, err error) {

	query := fmt.Sprintf("UPDATE kdb.%s SET active = %d, status = '%s' WHERE id= %d ", k8sTable, isActive, status, clusterId)
	insert, err := Db.Query(query)

	if err != nil {
		log.Logger.ErrorContext(ctx, fmt.Sprintf("update %s table failed", k8sTable), err)
		return false, err
	}

	defer insert.Close()
	log.Logger.TraceContext(ctx, "successfully updated cluster ", clusterId)
	return true, nil
}

//eks
func AddEksCluster(ctx context.Context, clusterId string, userId int, clusterName string, reqToken string, arn string, roleArn string, ids string, version string, zone string) (err error) {
	query := fmt.Sprintf("INSERT INTO kdb.%s (cluster_id, user_id, name, request_token, service_provider, status, active, "+
		"cluster_arn, role_arn, subnet_ids, kub_version, op_id, zone) "+ //todo: remove this op_id
		"VALUES ('%s', %d, '%s', '%s', '%s', 'CREATING', 1, '%s', '%s', '%s', '%s', '%s', '%s');", k8sTable, clusterId, userId, clusterName, reqToken, "amazon", arn, roleArn, ids, version, clusterName, zone)
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

func UpdateEksClusterCreationStatus(ctx context.Context, status string, clusterName string) (opStatus bool, err error) {
	statusDesc := status
	//statusDesc := "UNSPECIFIED"
	//switch status {
	//case "SUBMITTED":
	//	statusDesc = "INITIALIZING"
	//case "CREATING":
	//	statusDesc = "CREATING"
	//case "ACTIVE":
	//	statusDesc = "RUNNING"
	//default:
	//}
	query := fmt.Sprintf("UPDATE kdb.%s SET status='%s' WHERE name='%s'", k8sTable, statusDesc, clusterName)

	insert, err := Db.Query(query)

	if err != nil {
		log.Logger.ErrorContext(ctx, fmt.Sprintf("update %s table failed", k8sTable), err)
		return false, err
	}

	defer insert.Close()
	log.Logger.TraceContext(ctx, "successfully updated cluster ", clusterName)
	return true, nil
}

func CheckEksClusterCreationStatus(clusterName string, userID string) (status domain.GkeClusterStatus, error error) {
	status.Name = clusterName
	query := fmt.Sprintf("SELECT status,cluster_id FROM %s  WHERE name = '%s' AND user_id = '%s'", k8sTable, clusterName, userID)
	rows, err := Db.Query(query)
	if err != nil {
		status.Error = err.Error()
		status.Status = "CHECK FAILED"
		return status, err
	}
	defer rows.Close()
	for rows.Next() {

		err = rows.Scan(&status.Status, &status.ClusterId)
		if err != nil {
			status.Error = err.Error()
			status.Status = "CHECK FAILED"
			return status, err
		}
	}
	return status, nil
}
