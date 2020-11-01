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
	var createGkeCluster domain.ClusterOptions

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

	err = json.Unmarshal(content, &createGkeCluster)
	if err != nil {
		log.Logger.ErrorContext(ctx, "unmarshal error", err)
		res.WriteHeader(http.StatusBadRequest)
		return
	}
	//todo:refactor this
	if createGkeCluster.Provider == "google" {
		handleCreateGkeKubClusters(res, createGkeCluster)
	} else {
		handleCreateEksKubClusters(res, createGkeCluster)
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
	} else {
		handleCheckEksClusterCreationStatus(res, req)
	}

}

//GKE cluser commands
func handleGetAllGkeKubClusters(res http.ResponseWriter, req *http.Request) {
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
	//var createGkeCluster domain.ClusterOptions

	////user validation by token header
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
	//
	//content, err := ioutil.ReadAll(req.Body)
	//if err != nil {
	//	log.Logger.ErrorContext(ctx, "error occurred while reading request", err)
	//	res.WriteHeader(http.StatusBadRequest)
	//	return
	//}
	//
	//err = json.Unmarshal(content, &createGkeCluster)
	//if err != nil {
	//	log.Logger.ErrorContext(ctx, "unmarshal error", err)
	//	res.WriteHeader(http.StatusBadRequest)
	//	return
	//}
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
	//_, err = database.AddGkeCluster(ctx, clusterId, createGkeCluster.UserId, createGkeCluster.Name, op.Name)
	//if err != nil {
	//	res.WriteHeader(http.StatusInternalServerError)
	//	log.Logger.ErrorContext(ctx, "Could not add cluster creation request to db", createGkeCluster.Name)
	//	return
	//}

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

func handleRecommendGkeResource(res http.ResponseWriter, req *http.Request) {
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

	arrays, _ := url.ParseQuery(req.URL.RawQuery)
	Continent := arrays["continent[]"]
	Type := arrays["type[]"]
	Network := arrays["network[]"]
	Provider := req.FormValue("service_provider")
	VCPU := req.FormValue("vcpu")
	RAM := req.FormValue("ram")
	MinNodes := req.FormValue("min_nodes")
	MaxNodes := req.FormValue("max_nodes")

	//if err != nil {
	//	log.Logger.ErrorContext(ctx, "unmarshal error", err)
	//	res.WriteHeader(http.StatusBadRequest)
	//	return
	//}

	//result := database.GetSecretInternal(ctx, Name, OwnerId, Provider)
	//Continent := []string{"North America"}
	//Network := []string{"extra"}
	//Type := []string{"General purpose"}
	result := database.GetKubernetesResourcesRecommendation(ctx, Provider, Continent, VCPU, RAM, Network, Type, MinNodes, MaxNodes)

	if result.Status == 0 {
		res.WriteHeader(http.StatusOK)
	} else {
		res.WriteHeader(http.StatusBadRequest)
	}
	err := json.NewEncoder(res).Encode(&result)
	if err != nil {
		res.WriteHeader(http.StatusOK)
		log.Logger.ErrorContext(ctx, "response json conversion failed in get cluster recommendations")
	}
	log.Logger.TraceContext(ctx, "get cluster recommendations request successful")
}

func handleGetGkeResource(res http.ResponseWriter, req *http.Request) {
	ctx := traceableContext.WithUUID(uuid.New())

	////user validation by token header
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

	Provider := req.FormValue("service_provider")

	//if err != nil {
	//	log.Logger.ErrorContext(ctx, "unmarshal error", err)
	//	res.WriteHeader(http.StatusBadRequest)
	//	return
	//}

	//result := database.GetSecretInternal(ctx, Name, OwnerId, Provider)
	result := database.GetKubernetesResources(ctx, Provider)

	res.WriteHeader(http.StatusOK)
	err := json.NewEncoder(res).Encode(&result)
	if err != nil {
		res.WriteHeader(http.StatusOK)
		log.Logger.ErrorContext(ctx, "response json conversion failed in get cluster recommendations")
	}
	log.Logger.TraceContext(ctx, "get cluster recommendations request successful")
}

func handleValidateClusterName(res http.ResponseWriter, req *http.Request) {
	ctx := traceableContext.WithUUID(uuid.New())

	////user validation by token header
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

	User := req.FormValue("user_id")
	Name := req.FormValue("name")
	ServiceProvider := req.FormValue("service_provider")

	//if err != nil {
	//	log.Logger.ErrorContext(ctx, "unmarshal error", err)
	//	res.WriteHeader(http.StatusBadRequest)
	//	return
	//}

	//result := database.GetSecretInternal(ctx, Name, OwnerId, Provider)
	result := database.ValidateClusterName(ctx, User, Name, ServiceProvider)

	res.WriteHeader(http.StatusOK)
	err := json.NewEncoder(res).Encode(&result)
	if err != nil {
		res.WriteHeader(http.StatusOK)
		log.Logger.ErrorContext(ctx, "response json conversion failed in get cluster recommendations")
	}
	log.Logger.TraceContext(ctx, "get cluster recommendations request successful")
}

//EKS cluster commands
func handleCreateEksKubClusters(res http.ResponseWriter, createEksCluster domain.ClusterOptions) {
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
	fmt.Println("Create EKS k8s cluster request received")
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
	//_, err = database.AddGkeCluster(ctx, clusterId, createGkeCluster.UserId, createGkeCluster.Name, op.Name)
	//if err != nil {
	//	res.WriteHeader(http.StatusInternalServerError)
	//	log.Logger.ErrorContext(ctx, "Could not add cluster creation request to db", createGkeCluster.Name)
	//	return
	//}
	//service.PushToJobList(service.AsyncCloudJob{
	//	Provider:    "amazon",
	//	Status:      service.RUNNING,
	//	Reference:   result.ClusterStatus.Name,
	//	Information: result,
	//})
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

//func handleCheckEksClusterCreationStatus(res http.ResponseWriter, req *http.Request) {
//	ctx := traceableContext.WithUUID(uuid.New())
//
//	//user validation by token header
//	//token := req.Header.Get("Authorization")
//	//_, ok, err := database.ValidateUserByToken(ctx, strings.TrimSpace(strings.Split(token, "Bearer")[1]))
//	//if !ok {
//	//	log.Logger.DebugContext(ctx, "invalid user", token)
//	//	res.WriteHeader(http.StatusUnauthorized)
//	//	return
//	//}
//	//if err != nil {
//	//	log.Logger.ErrorContext(ctx, "error occurred in token validation", err)
//	//	res.WriteHeader(http.StatusInternalServerError)
//	//	return
//	//}
//
//	userId := req.FormValue("user_id")
//	clusterName := req.FormValue("cluster_name")
//	region := req.FormValue("region")
//	secretId, _ := strconv.Atoi(req.FormValue("secret_id"))
//
//	fmt.Printf("Check cluster creation status request received %s", userId)
//	// todo: replace with external call
//	result, err := kubernetes.CheckEksClusterCreationStatus(clusterName, region, secretId)
//	if err != nil {
//		log.Logger.ErrorContext(ctx, "checking the operation status failed")
//		res.WriteHeader(http.StatusInternalServerError)
//		err = json.NewEncoder(res).Encode(&result)
//		return
//	}
//	log.Logger.InfoContext(ctx, "Successfully retrieved op status from EKS")
//	res.WriteHeader(http.StatusOK)
//	err = json.NewEncoder(res).Encode(&result)
//	if err != nil {
//		res.WriteHeader(http.StatusOK)
//		log.Logger.ErrorContext(ctx, "response json conversion failed")
//	}
//	log.Logger.TraceContext(ctx, "Check kub cluster creation status request successful")
//}

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

func handleDeleteEksKubClusters(res http.ResponseWriter, req *http.Request) {
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

	secretId := req.FormValue("secret_id")
	clusterName := req.FormValue("cluster_name")
	region := req.FormValue("region")

	fmt.Println("Delete EKS k8s cluster request received")
	op, err := kubernetes.DeleteEksCluster(clusterName, secretId, region)
	log.Logger.InfoContext(ctx, "Cluster deletion request sent to Amazon", clusterName)
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
	//_, err = database.AddGkeCluster(ctx, clusterId, createGkeCluster.UserId, createGkeCluster.Name, op.Name)
	//if err != nil {
	//	res.WriteHeader(http.StatusInternalServerError)
	//	log.Logger.ErrorContext(ctx, "Could not add cluster creation request to db", createGkeCluster.Name)
	//	return
	//}
	//result, err = database.UpdateGkeClusterCreationStatus(ctx, op.Name, 3)
	result := domain.GkeClusterStatus{
		Name:      clusterName,
		ClusterId: clusterName,
		Status:    *op.Cluster.Status,
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

func handleCreateEksNodeGroup(res http.ResponseWriter, req *http.Request) {
	ctx := traceableContext.WithUUID(uuid.New())
	var createEksNodeGroup domain.EksClusterContext

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
	fmt.Println("Create EKS k8s node group request received")
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
	//_, err = database.AddGkeCluster(ctx, clusterId, createGkeCluster.UserId, createGkeCluster.Name, op.Name)
	//if err != nil {
	//	res.WriteHeader(http.StatusInternalServerError)
	//	log.Logger.ErrorContext(ctx, "Could not add cluster creation request to db", createGkeCluster.Name)
	//	return
	//}
	log.Logger.InfoContext(ctx, "Cluster creation request sent to Amazon", createEksNodeGroup.ClusterStatus.Name)
	//result, err = database.UpdateGkeClusterCreationStatus(ctx, op.Name, 3)
	//result := domain.GkeClusterStatus{
	//	Name:      createEksCluster.Name,
	//	ClusterId: clusterId,
	//	Status:    result.,
	//	Error:     "",
	//}
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

//func handleGetVPCConfigForRegion(res http.ResponseWriter, req *http.Request) {
//	ctx := traceableContext.WithUUID(uuid.New())
//
//	//user validation by token header
//	//token := req.Header.Get("Authorization")
//	//_, ok, err := database.ValidateUserByToken(ctx, strings.TrimSpace(strings.Split(token, "Bearer")[1]))
//	//if !ok {
//	//	log.Logger.DebugContext(ctx, "invalid user", token)
//	//	res.WriteHeader(http.StatusUnauthorized)
//	//	return
//	//}
//	//if err != nil {
//	//	log.Logger.ErrorContext(ctx, "error occurred in token validation", err)
//	//	res.WriteHeader(http.StatusInternalServerError)
//	//	return
//	//}
//
//	//secretId := req.FormValue("secret_id")
//	region := req.FormValue("region")
//	secretId, _ := strconv.Atoi(req.FormValue("secret_id"))
//
//	//fmt.Printf("Check cluster creation status request received %s", userId)
//	// todo: replace with external call
//	result := kubernetes.GetVPCConfigForUSerForRegion(secretId, region)
//	//if err != nil {
//	//	log.Logger.ErrorContext(ctx, "checking the operation status failed")
//	//	res.WriteHeader(http.StatusInternalServerError)
//	//	err = json.NewEncoder(res).Encode(&result)
//	//	return
//	//}
//	log.Logger.InfoContext(ctx, "Successfully retrieved op status from EKS")
//	res.WriteHeader(http.StatusOK)
//	err := json.NewEncoder(res).Encode(&result)
//	if err != nil {
//		res.WriteHeader(http.StatusOK)
//		log.Logger.ErrorContext(ctx, "response json conversion failed")
//	}
//	log.Logger.TraceContext(ctx, "Check kub cluster creation status request successful")
//}
