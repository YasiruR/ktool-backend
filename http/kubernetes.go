package http

import (
	"encoding/json"
	"fmt"
	"github.com/YasiruR/ktool-backend/database"
	"github.com/YasiruR/ktool-backend/domain"
	kubernetes "github.com/YasiruR/ktool-backend/kuberenetes"
	"github.com/YasiruR/ktool-backend/log"
	"github.com/google/uuid"
	traceableContext "github.com/pickme-go/traceable-context"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

// facade methods
func handleGetAllKubClusters(res http.ResponseWriter, req *http.Request) {
	ctx := traceableContext.WithUUID(uuid.New())

	//user validation by token header
	token := req.Header.Get("Authorization")
	_, ok, err := database.ValidateUserByToken(ctx, strings.TrimSpace(strings.Split(token, "Bearer")[1]))
	if !ok {
		log.Logger.DebugContext(ctx, "invalid user", token)
		res.WriteHeader(http.StatusUnauthorized)
		return
	}
	if err != nil {
		log.Logger.ErrorContext(ctx, "error occurred in token validation", err)
		res.WriteHeader(http.StatusInternalServerError)
		return
	}

	UserId, _ := strconv.Atoi(req.FormValue("user_id"))

	// todo: replace with external call
	result := database.GetAllKubernetesClustersForUser(ctx, UserId)
	if result.Error != nil {
		log.Logger.ErrorContext(ctx, "Error occurred while retrieving cluster list")
		res.WriteHeader(http.StatusInternalServerError)
		return
	}
	log.Logger.InfoContext(ctx, "Successfully retrieved cluster list.")
	res.WriteHeader(http.StatusOK)
	err = json.NewEncoder(res).Encode(&result)
	if err != nil {
		res.WriteHeader(http.StatusOK)
		log.Logger.ErrorContext(ctx, "response json conversion failed")
	}
	log.Logger.TraceContext(ctx, "List kub clusters request successful")
}

func handleCreateKubCluster(res http.ResponseWriter, req *http.Request) {
	ctx := traceableContext.WithUUID(uuid.New())
	var createCluster domain.ClusterOptions

	//user validation by token header
	token := req.Header.Get("Authorization")
	_, ok, err := database.ValidateUserByToken(ctx, strings.TrimSpace(strings.Split(token, "Bearer")[1]))
	if !ok {
		log.Logger.DebugContext(ctx, "invalid user", token)
		res.WriteHeader(http.StatusUnauthorized)
		return
	}
	if err != nil {
		log.Logger.ErrorContext(ctx, "error occurred in token validation", err)
		res.WriteHeader(http.StatusInternalServerError)
		return
	}

	content, err := ioutil.ReadAll(req.Body)
	if err != nil {
		log.Logger.ErrorContext(ctx, "error occurred while reading request", err)
		res.WriteHeader(http.StatusBadRequest)
		return
	}

	err = json.Unmarshal(content, &createCluster)
	if err != nil {
		log.Logger.ErrorContext(ctx, "unmarshal error", err)
		res.WriteHeader(http.StatusBadRequest)
		return
	}
	//todo:refactor this
	if createCluster.Provider == "google" {
		handleCreateGkeKubClusters(res, createCluster)
	} else if createCluster.Provider == "amazon" {
		handleCreateEksKubClusters(res, createCluster)
	} else {
		handleCreateAksKubCluster(res, createCluster)
	}
}

func handleCheckClusterCreationStatus(res http.ResponseWriter, req *http.Request) {
	ctx := traceableContext.WithUUID(uuid.New())

	//user validation by token header
	token := req.Header.Get("Authorization")
	_, ok, err := database.ValidateUserByToken(ctx, strings.TrimSpace(strings.Split(token, "Bearer")[1]))
	if !ok {
		log.Logger.DebugContext(ctx, "invalid user", token)
		res.WriteHeader(http.StatusUnauthorized)
		return
	}
	if err != nil {
		log.Logger.ErrorContext(ctx, "error occurred in token validation", err)
		res.WriteHeader(http.StatusInternalServerError)
		return
	}

	provider := req.FormValue("service_provider")

	if provider == "google" {
		handleCheckGkeClusterCreationStatus(res, req)
	} else if provider == "amazon" {
		handleCheckEksClusterCreationStatus(res, req)
	} else {
		handleCheckAksClusterCreationStatus(res, req)
	}

}

func handleDeleteKubCluster(res http.ResponseWriter, req *http.Request) {
	ctx := traceableContext.WithUUID(uuid.New())

	//user validation by token header
	token := req.Header.Get("Authorization")
	_, ok, err := database.ValidateUserByToken(ctx, strings.TrimSpace(strings.Split(token, "Bearer")[1]))
	if !ok {
		log.Logger.DebugContext(ctx, "invalid user", token)
		res.WriteHeader(http.StatusUnauthorized)
		return
	}
	if err != nil {
		log.Logger.ErrorContext(ctx, "error occurred in token validation", err)
		res.WriteHeader(http.StatusInternalServerError)
		return
	}

	provider := req.FormValue("service_provider")

	if provider == "google" {
		handleDeleteGkeCluster(res, req)
	} else if provider == "amazon" {
		handleDeleteEksCluster(res, req)
	} else {
		handleDeleteAksCluster(res, req)
	}

}

func handleRemoveClusterEntry(res http.ResponseWriter, req *http.Request) {
	ctx := traceableContext.WithUUID(uuid.New())

	//user validation by token header
	token := req.Header.Get("Authorization")
	_, ok, err := database.ValidateUserByToken(ctx, strings.TrimSpace(strings.Split(token, "Bearer")[1]))
	if !ok {
		log.Logger.DebugContext(ctx, "invalid user", token)
		res.WriteHeader(http.StatusUnauthorized)
		return
	}
	if err != nil {
		log.Logger.ErrorContext(ctx, "error occurred in token validation", err)
		res.WriteHeader(http.StatusInternalServerError)
		return
	}

	id := req.FormValue("id")

	err = database.RemoveClusterEntry(ctx, id)
	if err != nil {
		log.Logger.ErrorContext(ctx, "error occurred in removing entry", err)
		res.WriteHeader(http.StatusInternalServerError)
		return
	}
	log.Logger.ErrorContext(ctx, "removed kub cluster entry, %s", id)
	res.WriteHeader(http.StatusOK)
	return
}

//GKE cluser commands
func handleGetAllGkeKubClusters(res http.ResponseWriter, req *http.Request) {
	ctx := traceableContext.WithUUID(uuid.New())

	//user validation by token header
	token := req.Header.Get("Authorization")
	_, ok, err := database.ValidateUserByToken(ctx, strings.TrimSpace(strings.Split(token, "Bearer")[1]))
	if !ok {
		log.Logger.DebugContext(ctx, "invalid user", token)
		res.WriteHeader(http.StatusUnauthorized)
		return
	}
	if err != nil {
		log.Logger.ErrorContext(ctx, "error occurred in token validation", err)
		res.WriteHeader(http.StatusInternalServerError)
		return
	}

	UserId := req.FormValue("user_id")

	// todo: replace with external call
	result, err := kubernetes.ListGkeClusters(UserId)
	if err != nil {
		log.Logger.ErrorContext(ctx, "Could not retrieve cluster list")
		res.WriteHeader(http.StatusInternalServerError)
		return
	}
	log.Logger.InfoContext(ctx, "Successfully retrieved cluster list from GKE")
	res.WriteHeader(http.StatusOK)
	err = json.NewEncoder(res).Encode(&result)
	if err != nil {
		res.WriteHeader(http.StatusOK)
		log.Logger.ErrorContext(ctx, "response json conversion failed")
	}
	log.Logger.TraceContext(ctx, "List kub clusters request successful")
}

func handleCheckGkeClusterCreationStatus(res http.ResponseWriter, req *http.Request) {
	ctx := traceableContext.WithUUID(uuid.New())

	//user validation by token header
	token := req.Header.Get("Authorization")
	_, ok, err := database.ValidateUserByToken(ctx, strings.TrimSpace(strings.Split(token, "Bearer")[1]))
	if !ok {
		log.Logger.DebugContext(ctx, "invalid user", token)
		res.WriteHeader(http.StatusUnauthorized)
		return
	}
	if err != nil {
		log.Logger.ErrorContext(ctx, "error occurred in token validation", err)
		res.WriteHeader(http.StatusInternalServerError)
		return
	}

	userId := req.FormValue("user_id")
	operationId := req.FormValue("op_id")

	fmt.Printf("Check cluster creation status request received %s\n", userId)
	//result, err := kubernetes.CheckGkeClusterCreationStatus(userId, operationId)
	result, err := database.CheckGkeClusterCreationStatus(operationId, userId)
	if err != nil {
		log.Logger.ErrorContext(ctx, "checking the operation status failed")
		res.WriteHeader(http.StatusInternalServerError)
		err = json.NewEncoder(res).Encode(&result)
		return
	}
	log.Logger.InfoContext(ctx, "Successfully retrieved op status from GKE")
	res.WriteHeader(http.StatusOK)
	err = json.NewEncoder(res).Encode(&result)
	if err != nil {
		res.WriteHeader(http.StatusOK)
		log.Logger.ErrorContext(ctx, "response json conversion failed")
	}
	log.Logger.TraceContext(ctx, "Check kub cluster creation status request successful")
}

func handleCreateGkeKubClusters(res http.ResponseWriter, createGkeCluster domain.ClusterOptions) {
	ctx := traceableContext.WithUUID(uuid.New())
	fmt.Println("Create Gke k8s cluster request received")
	clusterId := uuid.New().String()
	op, err := kubernetes.CreateGkeCluster(clusterId, strconv.Itoa(createGkeCluster.SecretId), &createGkeCluster)
	if err != nil {
		res.WriteHeader(http.StatusOK)
		result := domain.GkeClusterStatus{
			Name:      createGkeCluster.Name,
			OpId:      "",
			ClusterId: clusterId,
			Status:    "FAILED",
			Error:     err.Error(),
		}
		err = json.NewEncoder(res).Encode(&result)
		log.Logger.ErrorContext(ctx, "Cluster creation failed, check logs", createGkeCluster.Name)
		return
	}
	//submit request to watcher
	kubernetes.PushToJobList(domain.AsyncCloudJob{
		Provider:    "google",
		Status:      domain.GKE_CREATING,
		Reference:   strconv.Itoa(createGkeCluster.SecretId),
		Information: op,
	})
	log.Logger.InfoContext(ctx, "Cluster creation request sent to Google", createGkeCluster.Name)
	//result, err = database.UpdateGkeClusterCreationStatus(ctx, op.Name, 3)
	result := domain.GkeClusterStatus{
		Name:      createGkeCluster.Name,
		OpId:      op.Name,
		ClusterId: clusterId,
		Status:    op.GetStatus().String(),
		Error:     "",
	}
	res.WriteHeader(http.StatusOK)
	err = json.NewEncoder(res).Encode(&result)
	if err != nil {
		res.WriteHeader(http.StatusInternalServerError)
		log.Logger.ErrorContext(ctx, "response json conversion failed", createGkeCluster.Name)
		return
	}
	log.Logger.TraceContext(ctx, "add gke k8s cluster request successful", createGkeCluster.Name)
}

func handleDeleteGkeCluster(res http.ResponseWriter, req *http.Request) {
	ctx := traceableContext.WithUUID(uuid.New())

	secretId := req.FormValue("secret_id")
	clusterName := req.FormValue("cluster_name")
	//projectName := req.FormValue("project_name")
	zone := req.FormValue("zone")
	clusterId, err := strconv.Atoi(req.FormValue("cluster_id"))
	if err != nil {
		res.WriteHeader(http.StatusBadRequest)
		log.Logger.ErrorContext(ctx, "request param conversion failed", clusterName)
		return
	}

	fmt.Println("Delete GKE k8s cluster request received")
	ok, err := kubernetes.DeleteGkeCluster(secretId, clusterId, clusterName, zone)
	if err != nil || !ok {
		res.WriteHeader(http.StatusInternalServerError)
		result := domain.GkeClusterStatus{
			Name:      clusterName,
			ClusterId: strconv.Itoa(clusterId),
			Status:    "FAILED TO DELETE",
			Error:     err.Error(),
		}
		err = json.NewEncoder(res).Encode(&result)
		log.Logger.ErrorContext(ctx, "Cluster deletion failed, check logs", clusterName)
		return
	}
	log.Logger.InfoContext(ctx, "Cluster delete request sent to Google", clusterName)
	result := domain.GkeClusterStatus{
		Name:      clusterName,
		ClusterId: strconv.Itoa(clusterId),
		Status:    "DELETING",
	}
	res.WriteHeader(http.StatusOK)
	err = json.NewEncoder(res).Encode(&result)
	if err != nil {
		res.WriteHeader(http.StatusInternalServerError)
		log.Logger.ErrorContext(ctx, "response json conversion failed", clusterName)
		return
	}
	log.Logger.TraceContext(ctx, "delete gke k8s cluster request successful", clusterName)
}

func handleRecommendResource(res http.ResponseWriter, req *http.Request) {
	ctx := traceableContext.WithUUID(uuid.New())

	//user validation by token header
	token := req.Header.Get("Authorization")
	_, ok, err := database.ValidateUserByToken(ctx, strings.TrimSpace(strings.Split(token, "Bearer")[1]))
	if !ok {
		log.Logger.DebugContext(ctx, "invalid user", token)
		res.WriteHeader(http.StatusUnauthorized)
		return
	}
	if err != nil {
		log.Logger.ErrorContext(ctx, "error occurred in token validation", err)
		res.WriteHeader(http.StatusInternalServerError)
		return
	}

	arrays, _ := url.ParseQuery(req.URL.RawQuery)
	Continent := arrays["continent[]"]
	Type := arrays["type[]"]
	Network := arrays["network[]"]
	Provider := req.FormValue("service_provider")
	VCPU := req.FormValue("vcpu")
	RAM := req.FormValue("ram")
	MinNodes := req.FormValue("min_nodes")
	MaxNodes := req.FormValue("max_nodes")

	result := database.GetKubernetesResourcesRecommendation(ctx, Provider, Continent, VCPU, RAM, Network, Type, MinNodes, MaxNodes)

	if result.Status == 0 {
		res.WriteHeader(http.StatusOK)
	} else {
		res.WriteHeader(http.StatusBadRequest)
	}
	err = json.NewEncoder(res).Encode(&result)
	if err != nil {
		res.WriteHeader(http.StatusOK)
		log.Logger.ErrorContext(ctx, "response json conversion failed in get cluster recommendations")
	}
	log.Logger.TraceContext(ctx, "get cluster recommendations request successful")
}

func handleGetKubResource(res http.ResponseWriter, req *http.Request) {
	ctx := traceableContext.WithUUID(uuid.New())

	////user validation by token header
	token := req.Header.Get("Authorization")
	_, ok, err := database.ValidateUserByToken(ctx, strings.TrimSpace(strings.Split(token, "Bearer")[1]))
	if !ok {
		log.Logger.DebugContext(ctx, "invalid user", token)
		res.WriteHeader(http.StatusUnauthorized)
		return
	}
	if err != nil {
		log.Logger.ErrorContext(ctx, "error occurred in token validation", err)
		res.WriteHeader(http.StatusInternalServerError)
		return
	}

	Provider := req.FormValue("service_provider")
	result := database.GetKubernetesResources(ctx, Provider)

	res.WriteHeader(http.StatusOK)
	err = json.NewEncoder(res).Encode(&result)
	if err != nil {
		res.WriteHeader(http.StatusOK)
		log.Logger.ErrorContext(ctx, "response json conversion failed in get cluster recommendations")
	}
	log.Logger.TraceContext(ctx, "get cluster recommendations request successful")
}

func handleValidateClusterName(res http.ResponseWriter, req *http.Request) {
	ctx := traceableContext.WithUUID(uuid.New())

	//user validation by token header
	token := req.Header.Get("Authorization")
	_, ok, err := database.ValidateUserByToken(ctx, strings.TrimSpace(strings.Split(token, "Bearer")[1]))
	if !ok {
		log.Logger.DebugContext(ctx, "invalid user", token)
		res.WriteHeader(http.StatusUnauthorized)
		return
	}
	if err != nil {
		log.Logger.ErrorContext(ctx, "error occurred in token validation", err)
		res.WriteHeader(http.StatusInternalServerError)
		return
	}

	User := req.FormValue("user_id")
	Name := req.FormValue("name")
	ServiceProvider := req.FormValue("service_provider")

	result := database.ValidateClusterName(ctx, User, Name, ServiceProvider)

	res.WriteHeader(http.StatusOK)
	err = json.NewEncoder(res).Encode(&result)
	if err != nil {
		res.WriteHeader(http.StatusOK)
		log.Logger.ErrorContext(ctx, "response json conversion failed in get cluster recommendations")
	}
	log.Logger.TraceContext(ctx, "get cluster recommendations request successful")
}

//EKS cluster commands
func handleCreateEksKubClusters(res http.ResponseWriter, createEksCluster domain.ClusterOptions) {
	ctx := traceableContext.WithUUID(uuid.New())

	clusterId := uuid.New().String()
	result, err := kubernetes.CreateEksCluster(clusterId, createEksCluster.SecretId, &createEksCluster)
	if err != nil {
		res.WriteHeader(http.StatusOK)
		result := domain.GkeClusterStatus{
			Name:      createEksCluster.Name,
			OpId:      "nil",
			ClusterId: clusterId,
			Status:    "FAILED",
			Error:     err.Error(),
		}
		err = json.NewEncoder(res).Encode(&result)
		log.Logger.ErrorContext(ctx, "Cluster creation failed, check logs", createEksCluster.Name)
		return
	}
	log.Logger.InfoContext(ctx, "Cluster creation request sent to Amazon", createEksCluster.Name)
	// submit job for the watcher
	kubernetes.PushToJobList(domain.AsyncCloudJob{
		Provider:    "amazon",
		Status:      domain.EKS_MASTER_CREATING,
		Reference:   result.ClusterStatus.Name,
		Information: result,
	})
	_, err = database.UpdateEksClusterCreationStatus(ctx, domain.EKS_MASTER_CREATING, createEksCluster.Name)
	res.WriteHeader(http.StatusOK)
	err = json.NewEncoder(res).Encode(&result)
	if err != nil {
		res.WriteHeader(http.StatusInternalServerError)
		log.Logger.ErrorContext(ctx, "response json conversion failed", createEksCluster.Name)
		return
	}
	log.Logger.TraceContext(ctx, "add EKS k8s cluster request successful", createEksCluster.Name)
}

func handleCheckEksClusterCreationStatus(res http.ResponseWriter, req *http.Request) {
	ctx := traceableContext.WithUUID(uuid.New())

	//user validation by token header
	//token := req.Header.Get("Authorization")
	//_, ok, err := database.ValidateUserByToken(ctx, strings.TrimSpace(strings.Split(token, "Bearer")[1]))
	//if !ok {
	//	log.Logger.DebugContext(ctx, "invalid user", token)
	//	res.WriteHeader(http.StatusUnauthorized)
	//	return
	//}
	//if err != nil {
	//	log.Logger.ErrorContext(ctx, "error occurred in token validation", err)
	//	res.WriteHeader(http.StatusInternalServerError)
	//	return
	//}

	userId := req.FormValue("user_id")
	clusterName := req.FormValue("cluster_name")

	fmt.Printf("Check cluster creation status request received %s", userId)
	// todo: replace with external call
	result, err := database.CheckEksClusterCreationStatus(clusterName, userId)
	if err != nil {
		log.Logger.ErrorContext(ctx, "checking the operation status failed")
		res.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(res).Encode(&result)
		return
	}
	log.Logger.InfoContext(ctx, "Successfully retrieved op status from EKS")
	res.WriteHeader(http.StatusOK)
	err = json.NewEncoder(res).Encode(&result)
	if err != nil {
		res.WriteHeader(http.StatusOK)
		log.Logger.ErrorContext(ctx, "response json conversion failed")
	}
	log.Logger.TraceContext(ctx, "Check kub cluster creation status request successful")
}

func handleDeleteEksCluster(res http.ResponseWriter, req *http.Request) {
	ctx := traceableContext.WithUUID(uuid.New())

	secretId := req.FormValue("secret_id")
	clusterName := req.FormValue("cluster_name")
	region := req.FormValue("zone")
	nodeGroupName := req.FormValue("nodegroup_name")

	fmt.Println("Delete EKS k8s cluster request received")
	err := kubernetes.DeleteEksCluster(ctx, clusterName, nodeGroupName, secretId, region)
	if err != nil {
		res.WriteHeader(http.StatusOK)
		result := domain.GkeClusterStatus{
			Name:      clusterName,
			ClusterId: nodeGroupName,
			Status:    "FAILED TO DELETE",
			Error:     err.Error(),
		}
		err = json.NewEncoder(res).Encode(&result)
		log.Logger.ErrorContext(ctx, "Cluster deletion failed, check logs", clusterName)
		return
	}
	log.Logger.InfoContext(ctx, "Cluster deletion request sent to Amazon", clusterName)

	result := domain.GkeClusterStatus{
		Name:      clusterName,
		ClusterId: nodeGroupName,
		Status:    "NODE GROUP DELETING",
	}
	res.WriteHeader(http.StatusOK)
	err = json.NewEncoder(res).Encode(&result)
	if err != nil {
		res.WriteHeader(http.StatusInternalServerError)
		log.Logger.ErrorContext(ctx, "response json conversion failed", clusterName)
		return
	}
	log.Logger.TraceContext(ctx, "delete eks k8s cluster request successful", clusterName)
}

func handleCreateEksNodeGroup(res http.ResponseWriter, req *http.Request) {
	ctx := traceableContext.WithUUID(uuid.New())
	var createEksNodeGroup domain.EksClusterContext

	//user validation by token header
	token := req.Header.Get("Authorization")
	_, ok, err := database.ValidateUserByToken(ctx, strings.TrimSpace(strings.Split(token, "Bearer")[1]))
	if !ok {
		log.Logger.DebugContext(ctx, "invalid user", token)
		res.WriteHeader(http.StatusUnauthorized)
		return
	}
	if err != nil {
		log.Logger.ErrorContext(ctx, "error occurred in token validation", err)
		res.WriteHeader(http.StatusInternalServerError)
		return
	}

	content, err := ioutil.ReadAll(req.Body)
	if err != nil {
		log.Logger.ErrorContext(ctx, "error occurred while reading request", err)
		res.WriteHeader(http.StatusBadRequest)
		return
	}

	err = json.Unmarshal(content, &createEksNodeGroup)
	if err != nil {
		log.Logger.ErrorContext(ctx, "unmarshal error", err)
		res.WriteHeader(http.StatusBadRequest)
		return
	}
	clusterId := uuid.New().String()
	result, err := kubernetes.CreateEksNodeGroup(createEksNodeGroup.SecretID, createEksNodeGroup)
	if err != nil {
		res.WriteHeader(http.StatusOK)
		result := domain.GkeClusterStatus{
			Name:      createEksNodeGroup.ClusterStatus.Name,
			OpId:      "nil",
			ClusterId: clusterId,
			Status:    "FAILED",
			Error:     err.Error(),
		}
		err = json.NewEncoder(res).Encode(&result)
		log.Logger.ErrorContext(ctx, "Cluster creation failed, check logs", createEksNodeGroup.ClusterStatus.Name)
		return
	}
	log.Logger.InfoContext(ctx, "Cluster creation request sent to Amazon", createEksNodeGroup.ClusterStatus.Name)
	res.WriteHeader(http.StatusOK)
	err = json.NewEncoder(res).Encode(&result)
	if err != nil {
		res.WriteHeader(http.StatusInternalServerError)
		log.Logger.ErrorContext(ctx, "response json conversion failed", createEksNodeGroup.ClusterStatus.Name)
		return
	}
	log.Logger.TraceContext(ctx, "add eks k8s cluster request successful", createEksNodeGroup.ClusterStatus.Name)
}

func handleCheckEksNodeGroupCreationStatus(res http.ResponseWriter, req *http.Request) {
	ctx := traceableContext.WithUUID(uuid.New())

	//user validation by token header
	token := req.Header.Get("Authorization")
	_, ok, err := database.ValidateUserByToken(ctx, strings.TrimSpace(strings.Split(token, "Bearer")[1]))
	if !ok {
		log.Logger.DebugContext(ctx, "invalid user", token)
		res.WriteHeader(http.StatusUnauthorized)
		return
	}
	if err != nil {
		log.Logger.ErrorContext(ctx, "error occurred in token validation", err)
		res.WriteHeader(http.StatusInternalServerError)
		return
	}

	userId := req.FormValue("user_id")
	clusterName := req.FormValue("cluster_name")
	nodeGroupName := req.FormValue("node_group_name")
	region := req.FormValue("region")
	secretId, _ := strconv.Atoi(req.FormValue("secret_id"))

	fmt.Printf("Check cluster creation status request received %s", userId)
	// todo: replace with external call
	result, err := kubernetes.CheckEksNodeGroupCreationStatus(clusterName, nodeGroupName, region, secretId)
	if err != nil {
		log.Logger.ErrorContext(ctx, "checking the operation status failed")
		res.WriteHeader(http.StatusInternalServerError)
		err = json.NewEncoder(res).Encode(&result)
		return
	}
	log.Logger.InfoContext(ctx, "Successfully retrieved op status from EKS")
	res.WriteHeader(http.StatusOK)
	err = json.NewEncoder(res).Encode(&result)
	if err != nil {
		res.WriteHeader(http.StatusOK)
		log.Logger.ErrorContext(ctx, "response json conversion failed")
	}
	log.Logger.TraceContext(ctx, "Check kub cluster creation status request successful")
}

func handleGetAllEksKubClusters(res http.ResponseWriter, req *http.Request) {
	ctx := traceableContext.WithUUID(uuid.New())

	//user validation by token header
	token := req.Header.Get("Authorization")
	_, ok, err := database.ValidateUserByToken(ctx, strings.TrimSpace(strings.Split(token, "Bearer")[1]))
	if !ok {
		log.Logger.DebugContext(ctx, "invalid user", token)
		res.WriteHeader(http.StatusUnauthorized)
		return
	}
	if err != nil {
		log.Logger.ErrorContext(ctx, "error occurred in token validation", err)
		res.WriteHeader(http.StatusInternalServerError)
		return
	}

	UserId := req.FormValue("user_id")
	Region := req.FormValue("region")

	// todo: replace with external call
	result, err := kubernetes.ListEksClusters(UserId, Region)
	if err != nil {
		log.Logger.ErrorContext(ctx, "Could not retrieve cluster list")
		res.WriteHeader(http.StatusInternalServerError)
		return
	}
	log.Logger.InfoContext(ctx, "Successfully retrieved cluster list from EKS")
	res.WriteHeader(http.StatusOK)
	err = json.NewEncoder(res).Encode(&result)
	if err != nil {
		res.WriteHeader(http.StatusOK)
		log.Logger.ErrorContext(ctx, "response json conversion failed")
	}
	log.Logger.TraceContext(ctx, "List kub clusters request successful")
}

//AKS cluster commands
func handleCreateAksKubCluster(res http.ResponseWriter, createAksCluster domain.ClusterOptions) {
	ctx := traceableContext.WithUUID(uuid.New())
	//var createEksCluster domain.ClusterOptions

	//user validation by token header
	//token := req.Header.Get("Authorization")
	//_, ok, err := database.ValidateUserByToken(ctx, strings.TrimSpace(strings.Split(token, "Bearer")[1]))
	//if !ok {
	//	log.Logger.DebugContext(ctx, "invalid user", token)
	//	res.WriteHeader(http.StatusUnauthorized)
	//	return
	//}
	//if err != nil {
	//	log.Logger.ErrorContext(ctx, "error occurred in token validation", err)
	//	res.WriteHeader(http.StatusInternalServerError)
	//	return
	//}

	//content, err := ioutil.ReadAll(req.Body)
	//if err != nil {
	//	log.Logger.ErrorContext(ctx, "error occurred while reading request", err)
	//	res.WriteHeader(http.StatusBadRequest)
	//	return
	//}
	//
	//err = json.Unmarshal(content, &createEksCluster)
	//if err != nil {
	//	log.Logger.ErrorContext(ctx, "unmarshal error", err)
	//	res.WriteHeader(http.StatusBadRequest)
	//	return
	//}
	//fmt.Println("Create EKS k8s cluster request received")
	clusterId := uuid.New().String()
	result, err := kubernetes.CreateAKSCluster(&createAksCluster)
	if err != nil {
		res.WriteHeader(http.StatusOK)
		result := domain.GkeClusterStatus{
			Name:      createAksCluster.Name,
			OpId:      "nil",
			ClusterId: clusterId,
			Status:    "FAILED",
			Error:     err.Error(),
		}
		err = json.NewEncoder(res).Encode(&result)
		log.Logger.ErrorContext(ctx, "Cluster creation failed, check logs", createAksCluster.Name)
		return
	}
	log.Logger.InfoContext(ctx, "Cluster creation request sent to Microsoft", createAksCluster.Name)
	//_, err = database.AddAKsCluster(ctx, clusterId, createAksCluster.UserId, createAksCluster.Name, createAksCluster.ResourceGroupName, createAksCluster.Location)
	//if err != nil {
	//	res.WriteHeader(http.StatusInternalServerError)
	//	log.Logger.ErrorContext(ctx, "Could not add cluster creation request to db", createGkeCluster.Name)
	//	return
	//}
	// submit job for the watcher
	//kubernetes.PushToJobList(domain.AsyncCloudJob{
	//	Provider:    "microsoft",
	//	Status:      domain.AKS_CREATING,
	//	Reference:   result.Status,
	//	Information: result,
	//})
	//_, err = database.UpdateAksClusterCreationStatus(ctx, domain.AKS_CREATING, createAksCluster.Name)
	res.WriteHeader(http.StatusOK)
	err = json.NewEncoder(res).Encode(&result)
	if err != nil {
		res.WriteHeader(http.StatusInternalServerError)
		log.Logger.ErrorContext(ctx, "response json conversion failed", createAksCluster.Name)
		return
	}
	log.Logger.TraceContext(ctx, "add AKS k8s cluster request successful", createAksCluster.Name)
}

func handleCheckAksClusterCreationStatus(res http.ResponseWriter, req *http.Request) {
	ctx := traceableContext.WithUUID(uuid.New())

	userId := req.FormValue("user_id")
	clusterName := req.FormValue("cluster_name")
	resourceGroupName := req.FormValue("resource_group")

	fmt.Printf("Check cluster creation status request received %s", userId)
	// todo: replace with external call
	result, err := database.CheckAksClusterCreationStatus(clusterName, resourceGroupName, userId)
	if err != nil {
		log.Logger.ErrorContext(ctx, "checking the operation status failed")
		res.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(res).Encode(&result)
		return
	}
	log.Logger.InfoContext(ctx, "Successfully retrieved op status from EKS")
	res.WriteHeader(http.StatusOK)
	err = json.NewEncoder(res).Encode(&result)
	if err != nil {
		res.WriteHeader(http.StatusOK)
		log.Logger.ErrorContext(ctx, "response json conversion failed")
	}
	log.Logger.TraceContext(ctx, "Check kub cluster creation status request successful")
}

func handleDeleteAksCluster(res http.ResponseWriter, req *http.Request) {
	ctx := traceableContext.WithUUID(uuid.New())

	secretId := req.FormValue("secret_id")
	clusterName := req.FormValue("cluster_name")
	resourceGroup := req.FormValue("resource_group")

	err := kubernetes.DeleteAksCluster(clusterName, resourceGroup, secretId)
	log.Logger.InfoContext(ctx, "Cluster deletion request sent to Microsoft", clusterName)
	if err != nil {
		res.WriteHeader(http.StatusOK)
		result := domain.GkeClusterStatus{
			Name:      clusterName,
			OpId:      "",
			ClusterId: clusterName,
			Status:    "FAILED TO DELETE",
			Error:     err.Error(),
		}
		err = json.NewEncoder(res).Encode(&result)
		log.Logger.ErrorContext(ctx, "Cluster deletion failed, check logs", clusterName)
		return
	}
	result := domain.GkeClusterStatus{
		Name:      clusterName,
		ClusterId: clusterName,
		Status:    "DELETING",
		Error:     "",
	}
	res.WriteHeader(http.StatusOK)
	err = json.NewEncoder(res).Encode(&result)
	if err != nil {
		res.WriteHeader(http.StatusInternalServerError)
		log.Logger.ErrorContext(ctx, "response json conversion failed", clusterName)
		return
	}
	log.Logger.TraceContext(ctx, "delete eks k8s cluster request successful", clusterName)
}

func handleCheckExistenceResourceGroup(res http.ResponseWriter, req *http.Request) {
	ctx := traceableContext.WithUUID(uuid.New())

	secretId := req.FormValue("secret_id")
	region := req.FormValue("region")
	resourceGroup := req.FormValue("resource_group")

	result, err := kubernetes.CreateResourceGroupIfNotExist(ctx, resourceGroup, region, secretId)
	log.Logger.InfoContext(ctx, "Resource group check request sent to Microsoft", resourceGroup)
	if err != nil {
		res.WriteHeader(http.StatusOK)
		result := domain.GkeClusterStatus{
			Status: "FAILED TO CHECK",
			Error:  err.Error(),
		}
		err = json.NewEncoder(res).Encode(&result)
		log.Logger.ErrorContext(ctx, "Resource group check failed, check logs", resourceGroup)
		return
	}
	res.WriteHeader(http.StatusOK)
	err = json.NewEncoder(res).Encode(&result)
	if err != nil {
		res.WriteHeader(http.StatusInternalServerError)
		log.Logger.ErrorContext(ctx, "response json conversion failed", resourceGroup)
		return
	}
	log.Logger.TraceContext(ctx, "check resource group request successful", resourceGroup)
}
