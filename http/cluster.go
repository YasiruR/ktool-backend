package http

import (
	"encoding/json"
	"fmt"
	"github.com/YasiruR/ktool-backend/database"
	"github.com/YasiruR/ktool-backend/domain"
	"github.com/YasiruR/ktool-backend/kafka"
	"github.com/YasiruR/ktool-backend/log"
	"github.com/YasiruR/ktool-backend/prometheus"
	"github.com/YasiruR/ktool-backend/service"
	"github.com/google/uuid"
	traceable_context "github.com/pickme-go/traceable-context"
	"io/ioutil"
	"net"
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
	if len(strings.Split(token, "Bearer")) < 2 {
		log.Logger.ErrorContext(ctx, "token format is invalid", token)
		res.WriteHeader(http.StatusBadRequest)
		return
	}
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

		var brokerAddrList []string
		for _, broker := range addClusterReq.Brokers {
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

		config, err := kafka.InitSaramaConfig(ctx, addClusterReq.ClusterName, "")
		if err != nil {
			log.Logger.ErrorContext(ctx, "initializing sarama config failed and may proceed with default config for client init", addClusterReq.ClusterName)
		}

		//get all relevant brokers
		client, err := kafka.InitClient(ctx, brokerAddrList, config)
		if err != nil {
			log.Logger.ErrorContext(ctx, "add cluster request failed", addClusterReq.ClusterName, brokerAddrList)
			var errRes errorMessage
			res.WriteHeader(http.StatusPreconditionFailed)
			errRes.Mesg = fmt.Sprintf("Could not find the cluster for brokers - %v", brokerAddrList)
			err := json.NewEncoder(res).Encode(errRes)
			if err != nil {
				log.Logger.ErrorContext(ctx, "encoding error response for test cluster req failed")
			}
			return
		}

		tmpBrokList, err := kafka.GetBrokerAddrList(ctx, client)
		if err != nil {
			log.Logger.ErrorContext(ctx, "test connection to cluster failed", addClusterReq.ClusterName)
			var errRes errorMessage
			res.WriteHeader(http.StatusPreconditionFailed)
			errRes.Mesg = fmt.Sprintf("Could not fetch rest of the brokers for cluster - %v", addClusterReq.ClusterName)
			err := json.NewEncoder(res).Encode(errRes)
			if err != nil {
				log.Logger.ErrorContext(ctx, "encoding error response for test cluster req failed")
			}
			return
		}

		var hosts []string
		var ports []int
		for _, tmpBroker := range tmpBrokList {
			host, portStr, err := net.SplitHostPort(tmpBroker)
			if err != nil {
				log.Logger.ErrorContext(ctx, fmt.Sprintf("splitting host and port failed for %v", tmpBroker), err)
				var errRes errorMessage
				res.WriteHeader(http.StatusInternalServerError)
				errRes.Mesg = fmt.Sprintf("Splitting host and port failed for cluster - %v", addClusterReq.ClusterName)
				err := json.NewEncoder(res).Encode(errRes)
				if err != nil {
					log.Logger.ErrorContext(ctx, "encoding error response for test cluster req failed")
				}
				return
			}

			hosts = append(hosts, host)
			port, err := strconv.Atoi(portStr)
			if err != nil {
				log.Logger.ErrorContext(ctx, "port does not contain an int", tmpBroker)
				var errRes errorMessage
				res.WriteHeader(http.StatusBadRequest)
				errRes.Mesg = fmt.Sprintf("At least one of the ports is not an integer - %v", addClusterReq.ClusterName)
				err := json.NewEncoder(res).Encode(errRes)
				if err != nil {
					log.Logger.ErrorContext(ctx, "encoding error response for test cluster req failed")
				}
				return
			}
			ports = append(ports, port)
		}

		//add new cluster and brokers to prom config if metrics enabled
		if addClusterReq.JmxEnabled == true {
			err = prometheus.AddNewJob(ctx, addClusterReq.ClusterName, addClusterReq.Brokers)
			if err != nil {
				log.Logger.ErrorContext(ctx, "failed adding new cluster to prom config", addClusterReq.ClusterName)
				var errRes errorMessage
				res.WriteHeader(http.StatusFailedDependency)
				errRes.Mesg = "Adding new cluster could not be completed. Please delete the added cluster and try creating it again."
				err := json.NewEncoder(res).Encode(errRes)
				if err != nil {
					log.Logger.ErrorContext(ctx, "encoding error response for add cluster req failed")
				}
				return
			}
		}

		//proceeds to db query
		//note : frontend validations should be added to request
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
	if len(strings.Split(token, "Bearer")) < 2 {
		log.Logger.ErrorContext(ctx, "token format is invalid", token)
		res.WriteHeader(http.StatusBadRequest)
		return
	}
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

		client, err := kafka.InitClient(ctx, tmpBrokers, nil)
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
	if len(strings.Split(token, "Bearer")) < 2 {
		log.Logger.ErrorContext(ctx, "token format is invalid", token)
		res.WriteHeader(http.StatusBadRequest)
		return
	}
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

	err = prometheus.DeleteJob(ctx, clusterName)
	if err != nil {
		log.Logger.ErrorContext(ctx, "failed deleting prom configs for the cluster", clusterID)
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
	if len(strings.Split(token, "Bearer")) < 2 {
		log.Logger.ErrorContext(ctx, "token format is invalid", token)
		res.WriteHeader(http.StatusBadRequest)
		return
	}

	userID, ok, err := database.ValidateUserByToken(ctx, strings.TrimSpace(strings.Split(token, "Bearer")[1]))
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

	for _, connectedCluster := range domain.LoggedInUserMap[userID].ConnectedClusters {
		for index, cluster := range clusterListRes.Clusters {
			if connectedCluster.ClusterName == cluster.ClusterName {
				clusterListRes.Clusters[index].Connected = true
				break
			}
		}
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
	tokenHeader := req.Header.Get("Authorization")
	if len(strings.Split(tokenHeader, "Bearer")) < 2 {
		log.Logger.ErrorContext(ctx, "token format is invalid", tokenHeader)
		res.WriteHeader(http.StatusBadRequest)
		return
	}
	token := strings.TrimSpace(strings.Split(tokenHeader, "Bearer")[1])
	userID, ok, err := database.ValidateUserByToken(ctx, token)
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

	user, ok := domain.LoggedInUserMap[userID]
	if ok {
		for _, cluster := range kafka.ClusterList {
			if cluster.ClusterID == clusterID {
				if cluster.Available == true {
					user.ConnectedClusters = append(user.ConnectedClusters, cluster)
					//note : check for concurrency issues
					domain.LoggedInUserMap[userID] = user
					log.Logger.TraceContext(ctx, "connected to cluster successfully", cluster.ClusterName)
					res.WriteHeader(http.StatusOK)
					return
				}
				log.Logger.WarnContext(ctx, "cluster is not available", clusterID)
				break
			}
		}
		log.Logger.WarnContext(ctx, "cluster does not exist", clusterID)
	} else {
		log.Logger.ErrorContext(ctx, "could not find a user from the logged in user list from token", userID)
		res.WriteHeader(http.StatusForbidden)
		return
	}

	//send error message
	res.WriteHeader(http.StatusBadRequest)
}

func handleDisconnectCluster(res http.ResponseWriter, req *http.Request) {
	ctx := traceable_context.WithUUID(uuid.New())

	//user validation by token header
	tokenHeader := req.Header.Get("Authorization")
	if len(strings.Split(tokenHeader, "Bearer")) < 2 {
		log.Logger.ErrorContext(ctx, "token format is invalid", tokenHeader)
		res.WriteHeader(http.StatusBadRequest)
		return
	}
	token := strings.TrimSpace(strings.Split(tokenHeader, "Bearer")[1])
	userID, ok, err := database.ValidateUserByToken(ctx, token)
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

	user, ok := domain.LoggedInUserMap[userID]
	if ok {
		for index, cluster := range user.ConnectedClusters {
			if cluster.ClusterID == clusterID {
				//remove the cluster
				user.ConnectedClusters[index] = user.ConnectedClusters[len(user.ConnectedClusters)-1] // Copy last element to index i.
				user.ConnectedClusters[len(user.ConnectedClusters)-1] = domain.KCluster{}   // Erase last element (write zero value).
				user.ConnectedClusters = user.ConnectedClusters[:len(user.ConnectedClusters)-1]   // Truncate slice.

				//note: check for concurrency issues
				domain.LoggedInUserMap[userID] = user
				log.Logger.TraceContext(ctx, "disconnected cluster successfully", cluster.ClusterName)
				res.WriteHeader(http.StatusOK)
				return
			}
		}
	} else {
		log.Logger.ErrorContext(ctx, "could not find a user from the logged in user list from token", token)
		res.WriteHeader(http.StatusForbidden)
		return
	}

	log.Logger.ErrorContext(ctx, "could not find the cluster in selected clusters", clusterID)
	//send error message
	res.WriteHeader(http.StatusBadRequest)
}

func handleGetGraphMetrics(res http.ResponseWriter, req *http.Request) {
	ctx := traceable_context.WithUUID(uuid.New())
	//user validation by token header
	tokenHeader := req.Header.Get("Authorization")
	if len(strings.Split(tokenHeader, "Bearer")) < 2 {
		log.Logger.ErrorContext(ctx, "token format is invalid", tokenHeader)
		res.WriteHeader(http.StatusBadRequest)
		return
	}

	token := strings.TrimSpace(strings.Split(tokenHeader, "Bearer")[1])
	userID, ok, err := database.ValidateUserByToken(ctx, token)
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

	//starting timestamp
	fromTs, err := strconv.Atoi(req.FormValue("from"))
	if err != nil {
		log.Logger.ErrorContext(ctx, "conversion of from_ts from string into int failed", err, req.FormValue("cluster_id"))
		res.WriteHeader(http.StatusBadRequest)
		return
	}

	//ending timestamp
	toTs, err := strconv.Atoi(req.FormValue("to"))
	if err != nil {
		log.Logger.ErrorContext(ctx, "conversion of to_ts from string into int failed", err, req.FormValue("cluster_id"))
		res.WriteHeader(http.StatusBadRequest)
		return
	}

	query := req.FormValue("query")
	host := req.FormValue("instance")

	user, ok := domain.LoggedInUserMap[userID]

	if ok {
		for _, userCluster := range user.ConnectedClusters {
			if userCluster.ClusterID == clusterID {
				// query directly from prometheus and send the res right away
				metrics, err := prometheus.GetMetricsForRange(ctx, query, host, domain.ClusterBrokerMetricsPortMap[userCluster.ClusterName][host], fromTs, toTs, userCluster.ClusterName)
				if err != nil {
					log.Logger.ErrorContext(ctx, "graph metrics api call failed", query, host)
					return
				}
				res.WriteHeader(http.StatusOK)
				err = json.NewEncoder(res).Encode(metrics)
				if err != nil {
					log.Logger.ErrorContext(ctx, err, "marshalling response for get graph metrics request failed", clusterID, metrics)
					res.WriteHeader(http.StatusInternalServerError)
					return
				}
				log.Logger.TraceContext(ctx, "graph metrics for requested cluster are fetched", clusterID)
				return
			}
		}
		log.Logger.ErrorContext(ctx, "user has not connected to the requested cluster to get broker overall metrics", clusterID)
		res.WriteHeader(http.StatusConflict)
	} else {
		log.Logger.ErrorContext(ctx, "could not find a user from the logged in user list from token", token)
		res.WriteHeader(http.StatusForbidden)
		return
	}
}

func handleGetBrokerOverview(res http.ResponseWriter, req *http.Request) {
	ctx := traceable_context.WithUUID(uuid.New())
	//user validation by token header
	tokenHeader := req.Header.Get("Authorization")
	if len(strings.Split(tokenHeader, "Bearer")) < 2 {
		log.Logger.ErrorContext(ctx, "token format is invalid", tokenHeader)
		res.WriteHeader(http.StatusBadRequest)
		return
	}

	token := strings.TrimSpace(strings.Split(tokenHeader, "Bearer")[1])
	userID, ok, err := database.ValidateUserByToken(ctx, token)
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

	//number of data points required
	clusterID, err := strconv.Atoi(req.FormValue("cluster_id"))
	if err != nil {
		log.Logger.ErrorContext(ctx, "conversion of cluster id from string into int failed", err, req.FormValue("cluster_id"))
		res.WriteHeader(http.StatusBadRequest)
		return
	}

	//starting timestamp
	fromTs, err := strconv.Atoi(req.FormValue("from"))
	if err != nil {
		log.Logger.ErrorContext(ctx, "conversion of from_ts from string into int failed", err, req.FormValue("cluster_id"))
		res.WriteHeader(http.StatusBadRequest)
		return
	}

	//ending timestamp
	toTs, err := strconv.Atoi(req.FormValue("to"))
	if err != nil {
		log.Logger.ErrorContext(ctx, "conversion of to_ts from string into int failed", err, req.FormValue("cluster_id"))
		res.WriteHeader(http.StatusBadRequest)
		return
	}

	//getting the number of timestamps
	count, err := strconv.Atoi(req.FormValue("count"))
	if err != nil {
		log.Logger.ErrorContext(ctx, "conversion of count from string into int failed", err, req.FormValue("count"))
		res.WriteHeader(http.StatusBadRequest)
		return
	}

	//setting default value to 1
	if count < 1 {
		count = 1
	}

	var overviewRes domain.ClusterOverview

	cluster, err := database.GetClusterByClusterID(ctx, clusterID)
	if err != nil {
		//even if error occurs, this may proceed without kafka version
		log.Logger.ErrorContext(ctx, "getting kafka version for broker overview failed", clusterID)
	} else {
		overviewRes.KafkaVersion = cluster.KafkaVersion
	}

	totalBytesIn := make(map[int64]int64)
	totalBytesOut := make(map[int64]int64)

	user, ok := domain.LoggedInUserMap[userID]
	if ok {
		for _, userCluster := range user.ConnectedClusters {
			if userCluster.ClusterID == clusterID {

				//get all brokers for the cluster
				brokers, err := database.GetBrokersByClusterId(ctx, clusterID)
				if err != nil {
					log.Logger.ErrorContext(ctx, "getting brokers for the requested cluster failed")
					res.WriteHeader(http.StatusInternalServerError)
					return
				}

				//fetching cluster object from kafkaList since user.connectedClusters do not update summary stats
				var metricsCluster domain.KCluster
				for _, c := range kafka.ClusterList {
					if c.ClusterID == clusterID {
						metricsCluster = c
						break
					}
				}

				for _, broker := range brokers {
					startingTs, err := database.GetNextTimestamp(ctx, broker.Host, int64(fromTs))
					if err != nil {
						log.Logger.ErrorContext(ctx, "getting broker metrics (starting ts) failed", broker.Host)
						continue
					}
					endingTs, err := database.GetPreviousTimestamp(ctx, broker.Host, int64(toTs))
					if err != nil {
						log.Logger.ErrorContext(ctx, "getting broker metrics (ending ts) failed", broker.Host)
						continue
					}

					//fetching all timestamps and intervals
					index := 1
					var avgMetricsEnabled bool		//to ignore the ts before the last one (when taking average, ts in between ending and last ts before ending have to be ignored)
					var metricsTs int64
					var tsList []int64
					var brokerOverview domain.BrokerOverview
					brokerOverview.Metrics = make(map[int64]domain.BrokerMetrics)
					brokerOverview.Host, brokerOverview.Port = broker.Host, broker.Port
					tsGapCount := int(endingTs - startingTs) / (service.Cfg.MetricsUpdateInterval * count)

					//from and to timestamps
					toMetrics, err := prometheus.GetAllMetricsByTimestamp(ctx, broker.Host, domain.ClusterBrokerMetricsPortMap[cluster.ClusterName][broker.Host], toTs, false)
					if err != nil {
						log.Logger.ErrorContext(ctx, "getting metrics for to ts failed", toTs)
					} else {
						//check if broker is in sync
						inSync := false
						if toMetrics.OfflinePartitions == 0 && toMetrics.UnderReplicated == 0 {
							inSync = true
							//in case the broker is down, to show 'not in sync'
							if toMetrics.Topics == 0 && toMetrics.NumLeaders == 0 && toMetrics.NumReplicas == 0 && toMetrics.Messages == 0 {
								inSync = false
							}
						}
						toMetrics.InSync = inSync
					}

					totalBytesOut[int64(toTs)] += toMetrics.ByteOutRate
					totalBytesIn[int64(toTs)] += toMetrics.ByteInRate
					brokerOverview.Metrics[int64(toTs)] = toMetrics

					//when count is 1 to ts metrics is sufficient
					if count == 1 {
						metricsCluster.ClusterOverview.Brokers = append(metricsCluster.ClusterOverview.Brokers, brokerOverview)
						continue
					}

					fromMetrics, err := prometheus.GetAllMetricsByTimestamp(ctx, broker.Host, domain.ClusterBrokerMetricsPortMap[cluster.ClusterName][broker.Host], fromTs, true)
					if err != nil {
						log.Logger.ErrorContext(ctx, "getting metrics for from ts failed", fromTs)
					}

					totalBytesOut[int64(fromTs)] += fromMetrics.ByteOutRate
					totalBytesIn[int64(fromTs)] += fromMetrics.ByteInRate
					brokerOverview.Metrics[int64(fromTs)] = fromMetrics

					//when count is 2 to ts metrics is sufficient
					if count == 2 {
						metricsCluster.ClusterOverview.Brokers = append(metricsCluster.ClusterOverview.Brokers, brokerOverview)
						continue
					}

					for t:=endingTs; t>=startingTs; t-=int64(service.Cfg.MetricsUpdateInterval) {
						if len(brokerOverview.Metrics) == count || (count == 1 && index == 2) {
							break
						}

						//todo: run these in separate go routines
						if t == endingTs {
							brokerMetricsMap, err := database.GetBrokerMetricsByTimestampList(ctx, broker.Host, []int64{endingTs})
							if err != nil {
								log.Logger.ErrorContext(ctx,"getting broker metrics failed", broker.Host)
								continue
							}

							totalBytesOut[t] += brokerMetricsMap[t].ByteOutRate
							totalBytesIn[t] += brokerMetricsMap[t].ByteInRate

							val := brokerMetricsMap[t]
							brokerOverview.Metrics[t] = val
						} else {
							tsList = append(tsList, t)
						}

						if tsGapCount != 0 && (index % tsGapCount == 0) {

							subEndingTs := t + int64(service.Cfg.MetricsUpdateInterval*(tsGapCount/2))
							subStartingTs := t - int64(service.Cfg.MetricsUpdateInterval*(tsGapCount/2))
							brokerMetrics, err := database.GetBrokerMetricsAverageValues(ctx, broker.Host, subStartingTs, subEndingTs)
							if err != nil {
								log.Logger.ErrorContext(ctx,"getting broker metrics failed", broker.Host)
								continue
							}

							//when taking average, ts in between ending and last ts before ending have to be ignored
							if !avgMetricsEnabled {
								avgMetricsEnabled = true
								tsList = []int64{}
								metricsTs = t		//to store the average of metrics in next interval in this ts
								continue
							}

							totalBytesOut[t] += brokerMetrics.ByteOutRate
							totalBytesIn[t] += brokerMetrics.ByteInRate
							brokerOverview.Metrics[metricsTs] = brokerMetrics	//broker metrics should be stored in last ts not the starting one (since ts values are decreased from the last ts to first ts)

							tsList = []int64{}
							metricsTs = t
						}
						index++
					}

					metricsCluster.ClusterOverview.Brokers = append(metricsCluster.ClusterOverview.Brokers, brokerOverview)
				}

				metricsCluster.ClusterOverview.TotalByteInRate = totalBytesIn
				metricsCluster.ClusterOverview.TotalByteOutRate = totalBytesOut
				overviewRes = metricsCluster.ClusterOverview
				res.WriteHeader(http.StatusOK)
				err = json.NewEncoder(res).Encode(overviewRes)
				if err != nil {
					log.Logger.ErrorContext(ctx, err, "marshalling response for get broker overview request failed", clusterID)
					res.WriteHeader(http.StatusInternalServerError)
				}
				log.Logger.TraceContext(ctx, "broker metrics for requested cluster are fetched", clusterID)
				return
			}
		}
		log.Logger.ErrorContext(ctx, "user has not connected to the requested cluster to get broker overall metrics", clusterID)
		res.WriteHeader(http.StatusConflict)
	} else {
		log.Logger.ErrorContext(ctx, "could not find a user from the logged in user list from token", token)
		res.WriteHeader(http.StatusForbidden)
		return
	}
}