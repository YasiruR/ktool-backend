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

	//"io/ioutil"
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
	result := database.GetAllKubernetesClusters(ctx, UserId)
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
	// todo:impl this refactor
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
	// todo: replace with external call
	result, err := kubernetes.CheckGkeClusterCreationStatus(userId, operationId)
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

func handleCreateGkeKubClusters(res http.ResponseWriter, req *http.Request) {
	ctx := traceableContext.WithUUID(uuid.New())
	var createGkeCluster domain.GkeClusterOptions

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
	result := database.GetGkeResourcesRecommendation(ctx, Provider, Continent, VCPU, RAM, Network, Type, MinNodes, MaxNodes)

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
	result := database.GetGkeResources(ctx, Provider)

	res.WriteHeader(http.StatusOK)
	err := json.NewEncoder(res).Encode(&result)
	if err != nil {
		res.WriteHeader(http.StatusOK)
		log.Logger.ErrorContext(ctx, "response json conversion failed in get cluster recommendations")
	}
	log.Logger.TraceContext(ctx, "get cluster recommendations request successful")
}

//EKS cluster commands
func handleCreateEksKubClusters(res http.ResponseWriter, req *http.Request) {
	ctx := traceableContext.WithUUID(uuid.New())
	var createEksCluster domain.GkeClusterOptions

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

	err = json.Unmarshal(content, &createEksCluster)
	if err != nil {
		log.Logger.ErrorContext(ctx, "unmarshal error", err)
		res.WriteHeader(http.StatusBadRequest)
		return
	}
	fmt.Println("Create EKS k8s cluster request received")
	clusterId := uuid.New().String()
	op, err := kubernetes.CreateEksCluster(clusterId, createEksCluster.SecretId, &createEksCluster)
	if err != nil {
		res.WriteHeader(http.StatusOK)
		result := domain.GkeClusterStatus{
			Name:      createEksCluster.Name,
			OpId:      "",
			ClusterId: clusterId,
			Status:    "FAILED",
			Error:     err.Error(),
		}
		err = json.NewEncoder(res).Encode(&result)
		log.Logger.ErrorContext(ctx, "Cluster creation failed, check logs", createEksCluster.Name)
		return
	}
	//_, err = database.AddGkeCluster(ctx, clusterId, createGkeCluster.UserId, createGkeCluster.Name, op.Name)
	//if err != nil {
	//	res.WriteHeader(http.StatusInternalServerError)
	//	log.Logger.ErrorContext(ctx, "Could not add cluster creation request to db", createGkeCluster.Name)
	//	return
	//}
	log.Logger.InfoContext(ctx, "Cluster creation request sent to Google", createEksCluster.Name)
	//result, err = database.UpdateGkeClusterCreationStatus(ctx, op.Name, 3)
	result := domain.GkeClusterStatus{
		Name:      createEksCluster.Name,
		ClusterId: clusterId,
		Status:    *op.Cluster.Status,
		Error:     "",
	}
	res.WriteHeader(http.StatusOK)
	err = json.NewEncoder(res).Encode(&result)
	if err != nil {
		res.WriteHeader(http.StatusInternalServerError)
		log.Logger.ErrorContext(ctx, "response json conversion failed", createEksCluster.Name)
		return
	}
	log.Logger.TraceContext(ctx, "add gke k8s cluster request successful", createEksCluster.Name)
}

func handleCheckEksClusterCreationStatus(res http.ResponseWriter, req *http.Request) {
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
	secretId, _ := strconv.Atoi(req.FormValue("secret_id"))

	fmt.Printf("Check cluster creation status request received %s\n", userId)
	// todo: replace with external call
	result, err := kubernetes.CheckEksClusterCreationStatus(operationId, secretId)
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
