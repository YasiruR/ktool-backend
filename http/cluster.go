package http

import (
	"encoding/json"
	"fmt"
	"github.com/YasiruR/ktool-backend/database"
	"github.com/YasiruR/ktool-backend/domain"
	"github.com/YasiruR/ktool-backend/kafka"
	"github.com/YasiruR/ktool-backend/log"
	"github.com/google/uuid"
	traceable_context "github.com/pickme-go/traceable-context"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
)

const checkClusterRetryCount = 3

//add existing kafka cluster
func handleAddCluster(res http.ResponseWriter, req *http.Request) {
	ctx := traceable_context.WithUUID(uuid.New())
	var addClusterReq addExistingCluster

	//user validation by token header
	token := req.Header.Get("Authorization")
	_, ok, err := database.ValidateUserByToken(ctx, token)
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

	err = json.Unmarshal(content, &addClusterReq)
	if err != nil {
		log.Logger.ErrorContext(ctx, "unmarshal error", err)
		res.WriteHeader(http.StatusBadRequest)
		return
	}

	retryCount := 0
checkIfClusterExists:
	_, err = database.GetClusterIdByName(ctx, strings.TrimSpace(addClusterReq.ClusterName))
	if err == nil {
		log.Logger.ErrorContext(ctx, "cluster name already exists", addClusterReq.ClusterName)
		var errRes errorMessage
		res.WriteHeader(http.StatusPreconditionFailed)
		errRes.Mesg = "Cluster name already exists. Please provide a different name."
		err := json.NewEncoder(res).Encode(errRes)
		if err != nil {
			log.Logger.ErrorContext(ctx, "encoding error response for add cluster req failed")
		}
		return
	} else if err.Error() == "no rows found" {
		//when cluster is eligible to be added

		var hosts []string
		var ports []int
		var brokerAddrList []string
		for _, broker := range addClusterReq.Brokers {
			hosts = append(hosts, broker.Host)
			ports = append(ports, broker.Port)
			brokerAddrList = append(brokerAddrList, broker.Host + ":" + strconv.Itoa(broker.Port))
		}

		for _, brokerAddr := range brokerAddrList {
			ok := kafka.CheckIfBrokerExists(ctx, brokerAddr)
			if ok {
				var errRes errorMessage
				res.WriteHeader(http.StatusPreconditionFailed)
				errRes.Mesg = fmt.Sprintf("You have already added broker (%v) into a cluster", brokerAddr)
				err := json.NewEncoder(res).Encode(errRes)
				if err != nil {
					log.Logger.ErrorContext(ctx, "encoding error response for add cluster req failed")
				}
				return
			}
		}

		//proceeds to db query
		//note : frontend validations should be added to request parameters
		err = database.AddNewCluster(ctx, addClusterReq.ClusterName, addClusterReq.KafkaVersion)
		if err != nil {
			log.Logger.ErrorContext(ctx, "add new cluster db transaction failed")
			res.WriteHeader(http.StatusInternalServerError)
			return
		}

		err = database.AddNewBrokers(ctx, hosts, ports, addClusterReq.ClusterName)
		if err != nil {
			log.Logger.ErrorContext(ctx, "add new brokers db transaction failed", err)
			res.WriteHeader(http.StatusInternalServerError)

			//if adding new brokers failed, reverts the adding cluster query as well
			err = database.DeleteCluster(ctx, addClusterReq.ClusterName)
			if err != nil {
				log.Logger.ErrorContext(ctx, "deleting newly added cluster failed")
			} else {
				log.Logger.TraceContext(ctx, "deleting newly added cluster was successful")
			}
			return
		}
		log.Logger.TraceContext(ctx, "cluster stored in the database successfully", addClusterReq.ClusterName)
	} else {
		retryCount += 1
		if retryCount <= checkClusterRetryCount {
			goto checkIfClusterExists
		}
		log.Logger.ErrorContext(ctx, fmt.Sprintf("checking if cluster exists failed %v times", retryCount-1))
		res.WriteHeader(http.StatusInternalServerError)
		return
	}

	kafka.InitAllClusters()
}

func handleTestConnectionToCluster(res http.ResponseWriter, req *http.Request) {
	ctx := traceable_context.WithUUID(uuid.New())
	var testClusterReq addExistingCluster

	//user validation by token header
	token := req.Header.Get("Authorization")
	_, ok, err := database.ValidateUserByToken(ctx, token)
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

	err = json.Unmarshal(content, &testClusterReq)
	if err != nil {
		log.Logger.ErrorContext(ctx, "unmarshal error", err)
		res.WriteHeader(http.StatusBadRequest)
		return
	}

	log.Logger.TraceContext(ctx, testClusterReq, "request")

	var listOfBrokLists [][]string

	for _, b := range testClusterReq.Brokers {
		var tmpBrokers []string
		tmp := b.Host + ":" + strconv.Itoa(b.Port)
		tmpBrokers = append(tmpBrokers, tmp)

		client, err := kafka.InitClient(ctx, tmpBrokers)
		if err != nil {
			log.Logger.ErrorContext(ctx, "test connection to cluster failed", testClusterReq.ClusterName, b)
			var errRes errorMessage
			res.WriteHeader(http.StatusPreconditionFailed)
			errRes.Mesg = fmt.Sprintf("Could not find the cluster for broker - %v", tmp)
			err := json.NewEncoder(res).Encode(errRes)
			if err != nil {
				log.Logger.ErrorContext(ctx, "encoding error response for test cluster req failed")
			}
			return
		}

		tmpBrokList, err := kafka.GetBrokerAddrList(ctx, client)
		if err != nil {
			log.Logger.ErrorContext(ctx, "test connection to cluster failed", testClusterReq.ClusterName)
			var errRes errorMessage
			res.WriteHeader(http.StatusPreconditionFailed)
			errRes.Mesg = fmt.Sprintf("Could not fetch rest of the brokers for broker - %v", tmp)
			err := json.NewEncoder(res).Encode(errRes)
			if err != nil {
				log.Logger.ErrorContext(ctx, "encoding error response for test cluster req failed")
			}
			return
		}

		listOfBrokLists = append(listOfBrokLists, tmpBrokList)
	}

	//check if all brokers provided are from the same cluster
	for index, _ := range listOfBrokLists {
		if len(listOfBrokLists) >= 2 && index != 0 {
			//check by length
			if len(listOfBrokLists[index-1]) != len(listOfBrokLists[index]) {
				log.Logger.ErrorContext(ctx, "test connection to cluster failed", testClusterReq.ClusterName)
				var errRes errorMessage
				res.WriteHeader(http.StatusPreconditionFailed)
				errRes.Mesg = "Provided brokers are not from the same cluster"
				err := json.NewEncoder(res).Encode(errRes)
				if err != nil {
					log.Logger.ErrorContext(ctx, "encoding error response for test cluster req failed")
				}
				return
			}

			//check by elements
			exists := make(map[string]bool)
			for _, value := range listOfBrokLists[index-1] {
				exists[value] = true
			}
			for _, value := range listOfBrokLists[index] {
				if !exists[value] {
					log.Logger.ErrorContext(ctx, "test connection to cluster failed", testClusterReq.ClusterName)
					var errRes errorMessage
					res.WriteHeader(http.StatusPreconditionFailed)
					errRes.Mesg = "Provided brokers are not from the same cluster"
					err := json.NewEncoder(res).Encode(errRes)
					if err != nil {
						log.Logger.ErrorContext(ctx, "encoding error response for test cluster req failed")
					}
					return
				}
			}
		}
	}

	log.Logger.TraceContext(ctx, "telnet to server is successful")
	res.WriteHeader(http.StatusOK)
}

func handleDeleteCluster(res http.ResponseWriter, req *http.Request) {
	ctx := traceable_context.WithUUID(uuid.New())
	clusterName := req.FormValue("cluster_name")

	//user validation by token header
	token := req.Header.Get("Authorization")
	_, ok, err := database.ValidateUserByToken(ctx, token)
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

	clusterID, err := database.GetClusterIdByName(ctx, clusterName)
	if err != nil {
		log.Logger.ErrorContext(ctx, fmt.Sprintf("deleting cluster failed - %v due to being unable to get cluster id by name", clusterName), err)
		res.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = database.DeleteBrokersOfCluster(ctx, clusterID)
	if err != nil {
		log.Logger.ErrorContext(ctx, fmt.Sprintf("deleting cluster failed - %v due to failure in deleting brokers", clusterName))
		res.WriteHeader(http.StatusInternalServerError)
		return
	}

	users, err := database.GetAllUsers(ctx)
	if err != nil {
		log.Logger.ErrorContext(ctx, fmt.Sprintf("deleting cluster failed - %v due to failure in fetching all users", clusterName))
	}

	err = kafka.DeleteCluster(ctx, clusterID, users)
	if err != nil {
		log.Logger.ErrorContext(ctx, fmt.Sprintf("deleting cluster from list failed - %v", clusterName), err)
		res.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = database.DeleteCluster(ctx, clusterName)
	if err != nil {
		log.Logger.ErrorContext(ctx, fmt.Sprintf("deleting cluster from db failed - %v", clusterName), err)
		res.WriteHeader(http.StatusInternalServerError)
		return
	}

	log.Logger.TraceContext(ctx, "cluster deleted successfully", clusterName)
	res.WriteHeader(http.StatusOK)
}

func handleGetAllClusters(res http.ResponseWriter, req *http.Request) {
	ctx := traceable_context.WithUUID(uuid.New())
	var clusterListRes clusterRes

	//user validation by token header
	token := req.Header.Get("Authorization")
	_, ok, err := database.ValidateUserByToken(ctx, token)
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

	for _, cluster := range kafka.ClusterList {
		clusterRes := clusterInfo{}
		clusterRes.Id = cluster.ClusterID
		clusterRes.ClusterName = cluster.ClusterName
		clusterRes.Available = cluster.Available

		for _, b := range cluster.Brokers {
			clusterRes.Brokers = append(clusterRes.Brokers, b.Addr())
		}
		for _, t := range cluster.Topics {
			var topicRes topic
			topicRes.Name = t.Name
			topicRes.Partitions = t.Partitions
			clusterRes.Topics = append(clusterRes.Topics, topicRes)
		}

		clusterListRes.Clusters = append(clusterListRes.Clusters, clusterRes)
	}

	res.WriteHeader(http.StatusOK)
	err = json.NewEncoder(res).Encode(clusterListRes)
	if err != nil {
		log.Logger.ErrorContext(ctx, "encoding response into json failed", err)
		res.WriteHeader(http.StatusInternalServerError)
		return
	}

	log.Logger.TraceContext(ctx, "get all clusters was successful", fmt.Sprintf("no. of clusters : %v", len(clusterListRes.Clusters)))
}

func handleConnectToCluster(res http.ResponseWriter, req *http.Request) {
	ctx := traceable_context.WithUUID(uuid.New())

	//user validation by token header
	token := req.Header.Get("Authorization")
	_, ok, err := database.ValidateUserByToken(ctx, token)
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

	clusterID, err := strconv.Atoi(req.FormValue("cluster_id"))
	if err != nil {
		log.Logger.ErrorContext(ctx, "conversion of cluster id from string into int failed", err, req.FormValue("cluster_id"))
		res.WriteHeader(http.StatusBadRequest)
		return
	}

	username, _, err := database.GetUserByToken(ctx, token)
	if err != nil {
		log.Logger.ErrorContext(ctx, "fetching user for connect to cluster failed", token)
		res.WriteHeader(http.StatusInternalServerError)
		return
	}

	var user domain.User
	var userFound bool
	for _, u := range domain.LoggedInUsers {
		if u.Username == username {
			user = u
			userFound = true
			break
		}
	}

	if !userFound {
		log.Logger.ErrorContext(ctx, "could not find a user from the logged in user list from token", token)
		res.WriteHeader(http.StatusForbidden)
		return
	}

	for _, cluster := range kafka.ClusterList {
		if cluster.ClusterID == clusterID {
			if cluster.Available == true {
				//kafka.SelectedClusterList = append(kafka.SelectedClusterList, cluster)
				user.ConnectedClusters = append(user.ConnectedClusters, cluster)
				log.Logger.TraceContext(ctx, "connected to cluster successfully", cluster.ClusterName)
				res.WriteHeader(http.StatusOK)
				return
			}
			break
		}
	}

	log.Logger.ErrorContext(ctx, "could not find the cluster id from the cluster id", clusterID)
	//send error message
	res.WriteHeader(http.StatusBadRequest)
}

func handleDisconnectCluster(res http.ResponseWriter, req *http.Request) {
	ctx := traceable_context.WithUUID(uuid.New())

	//user validation by token header
	token := req.Header.Get("Authorization")
	_, ok, err := database.ValidateUserByToken(ctx, token)
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

	clusterID, err := strconv.Atoi(req.FormValue("cluster_id"))
	if err != nil {
		log.Logger.ErrorContext(ctx, "conversion od cluster id from string into int failed", err, req.FormValue("cluster_id"))
		res.WriteHeader(http.StatusBadRequest)
		return
	}

	username, _, err := database.GetUserByToken(ctx, token)
	if err != nil {
		log.Logger.ErrorContext(ctx, "fetching user for connect to cluster failed", token)
		res.WriteHeader(http.StatusInternalServerError)
		return
	}

	var user domain.User
	var userFound bool
	for _, u := range domain.LoggedInUsers {
		if u.Username == username {
			user = u
			userFound = true
			break
		}
	}

	if !userFound {
		log.Logger.ErrorContext(ctx, "could not find a user from the logged in user list from token", token)
		res.WriteHeader(http.StatusForbidden)
		return
	}

	//for index, cluster := range kafka.SelectedClusterList {
	//	if cluster.ClusterID == clusterID {
	//		//remove the cluster
	//		kafka.SelectedClusterList[index] = kafka.SelectedClusterList[len(kafka.SelectedClusterList)-1] // Copy last element to index i.
	//		kafka.SelectedClusterList[len(kafka.SelectedClusterList)-1] = kafka.KCluster{}   // Erase last element (write zero value).
	//		kafka.SelectedClusterList = kafka.SelectedClusterList[:len(kafka.SelectedClusterList)-1]   // Truncate slice.
	//
	//		log.Logger.TraceContext(ctx, "disconnected cluster successfully", cluster.ClusterName)
	//		res.WriteHeader(http.StatusOK)
	//		return
	//	}
	//}

	for index, cluster := range user.ConnectedClusters {
		if cluster.ClusterID == clusterID {
			//remove the cluster
			user.ConnectedClusters[index] = user.ConnectedClusters[len(user.ConnectedClusters)-1] // Copy last element to index i.
			user.ConnectedClusters[len(user.ConnectedClusters)-1] = domain.KCluster{}   // Erase last element (write zero value).
			user.ConnectedClusters = user.ConnectedClusters[:len(user.ConnectedClusters)-1]   // Truncate slice.

			log.Logger.TraceContext(ctx, "disconnected cluster successfully", cluster.ClusterName)
			res.WriteHeader(http.StatusOK)
			return
		}
	}

	log.Logger.ErrorContext(ctx, "could not find the cluster in selected clusters", clusterID)
	//send error message
	res.WriteHeader(http.StatusBadRequest)
}