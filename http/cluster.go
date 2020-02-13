package http

import (
	"encoding/json"
	"fmt"
	"github.com/YasiruR/ktool-backend/database"
	"github.com/YasiruR/ktool-backend/kafka"
	"github.com/YasiruR/ktool-backend/log"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	traceable_context "github.com/pickme-go/traceable-context"
	"io/ioutil"
	"net/http"
	"strings"
)

const checkClusterRetryCount = 3

//add existing kafka cluster
func handleAddCluster(res http.ResponseWriter, req *http.Request) {
	var addClusterReq addExistingCluster

	ctx := traceable_context.WithUUID(uuid.New())
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

	////proceeds to db query
	////note : frontend validations should be added to request parameters
	//err = database.AddNewCluster(ctx, addClusterReq.ClusterName, addClusterReq.KafkaVersion)
	//if err != nil {
	//	log.Logger.ErrorContext(ctx, "add new cluster db transaction failed")
	//	reqFailed = true
	//	res.WriteHeader(http.StatusInternalServerError)
	//} else {
	//	err = database.AddNewZookeeper(ctx, addClusterReq.ZookeeperHost, addClusterReq.ZookeeperPort, addClusterReq.ClusterName)
	//	if err != nil {
	//		log.Logger.ErrorContext(ctx, "add new zookeeper db transaction failed")
	//		//if adding new zookeeper failed reverts the adding cluster query as well
	//		err = database.DeleteCluster(ctx, addClusterReq.ClusterName)
	//		if err != nil {
	//			log.Logger.ErrorContext(ctx, "deleting newly added cluster failed")
	//			//cluster table modified but zookeeper table is not
	//			res.WriteHeader(http.StatusConflict)
	//		} else {
	//			log.Logger.TraceContext(ctx, "deleting newly added cluster was successful")
	//			res.WriteHeader(http.StatusInternalServerError)
	//		}
	//		reqFailed = true
	//	}
	//}

	//if reqFailed == false {
	//	retryCount := 0
	//checkIfClusterExists:
	//	_, err = database.GetClusterIdByName(ctx, strings.TrimSpace(addClusterReq.ClusterName))
	//	if err == nil {
	//		log.Logger.ErrorContext(ctx, "cluster name already exists", addClusterReq.ClusterName)
	//		var errRes errorMessage
	//		res.WriteHeader(http.StatusPreconditionFailed)
	//		errRes.Mesg = "Cluster name already exists. Please provide a different name."
	//		err := json.NewEncoder(res).Encode(errRes)
	//		if err != nil {
	//			log.Logger.ErrorContext(ctx, "encoding error response for add cluster req failed")
	//		}
	//	} else if err.Error() == "no rows found" {
	//		//when cluster is eligible to be added
	//
	//		//proceeds to db query
	//		//note : frontend validations should be added to request parameters
	//		err = database.AddNewCluster(ctx, addClusterReq.ClusterName, addClusterReq.KafkaVersion)
	//		if err != nil {
	//			reqFailed = true
	//			log.Logger.ErrorContext(ctx, "add new cluster db transaction failed")
	//
	//		} else {
	//			var hosts []string
	//			var ports []int
	//			for _, broker := range addClusterReq.Brokers {
	//				hosts = append(hosts, broker.Host)
	//				ports = append(ports, broker.Port)
	//			}
	//
	//			err = database.AddNewBrokers(ctx, hosts, ports, addClusterReq.ClusterName)
	//			if err != nil {
	//				log.Logger.ErrorContext(ctx, "add new brokers db transaction failed", err)
	//
	//				if err.Error() == "duplicate entry" {
	//					var errRes errorMessage
	//					res.WriteHeader(http.StatusPreconditionFailed)
	//					errRes.Mesg = "You might have already added this cluster."
	//					err := json.NewEncoder(res).Encode(errRes)
	//					if err != nil {
	//						log.Logger.ErrorContext(ctx, "encoding error response for add cluster req failed")
	//					}
	//				} else {
	//					res.WriteHeader(http.StatusInternalServerError)
	//				}
	//
	//				//if adding new brokers failed, reverts the adding cluster query as well
	//				err = database.DeleteCluster(ctx, addClusterReq.ClusterName)
	//				if err != nil {
	//					log.Logger.ErrorContext(ctx, "deleting newly added cluster failed")
	//					//cluster table modified but zookeeper table is not
	//					//res.WriteHeader(http.StatusConflict)
	//				} else {
	//					log.Logger.TraceContext(ctx, "deleting newly added cluster was successful")
	//				}
	//				reqFailed = true
	//			}
	//		}
	//
	//		if reqFailed == false {
	//			log.Logger.TraceContext(ctx, "cluster stored in the database successfully", addClusterReq.ClusterName)
	//		}
	//	} else {
	//		retryCount += 1
	//		if retryCount <= checkClusterRetryCount {
	//			goto checkIfClusterExists
	//		}
	//		log.Logger.ErrorContext(ctx, fmt.Sprintf("checking if cluster exists failed %v times", retryCount-1))
	//		res.WriteHeader(http.StatusInternalServerError)
	//	}
	//}

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

		//proceeds to db query
		//note : frontend validations should be added to request parameters
		err = database.AddNewCluster(ctx, addClusterReq.ClusterName, addClusterReq.KafkaVersion)
		if err != nil {
			log.Logger.ErrorContext(ctx, "add new cluster db transaction failed")
			res.WriteHeader(http.StatusInternalServerError)
			return
		}

		fmt.Println("Req : ", addClusterReq.Brokers)

		var hosts []string
		var ports []int
		for _, broker := range addClusterReq.Brokers {
			hosts = append(hosts, broker.Host)
			ports = append(ports, broker.Port)
		}

		err = database.AddNewBrokers(ctx, hosts, ports, addClusterReq.ClusterName)
		if err != nil {
			log.Logger.ErrorContext(ctx, "add new brokers db transaction failed", err)

			if err.Error() == "duplicate entry" {
				var errRes errorMessage
				res.WriteHeader(http.StatusPreconditionFailed)
				errRes.Mesg = "You might have already added this cluster."
				err := json.NewEncoder(res).Encode(errRes)
				if err != nil {
					log.Logger.ErrorContext(ctx, "encoding error response for add cluster req failed")
				}
			} else {
				res.WriteHeader(http.StatusInternalServerError)
			}

			//if adding new brokers failed, reverts the adding cluster query as well
			err = database.DeleteCluster(ctx, addClusterReq.ClusterName)
			if err != nil {
				log.Logger.ErrorContext(ctx, "deleting newly added cluster failed")
				//cluster table modified but zookeeper table is not
				//res.WriteHeader(http.StatusConflict)
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
	}
}

//func handleGetAllClusters(res http.ResponseWriter, req *http.Request) {
//	ctx := traceable_context.WithUUID(uuid.New())
//	clusterList, err := database.GetAllClusters(ctx)
//	if err != nil {
//		log.Logger.ErrorContext(ctx, "get all clusters failed")
//		res.WriteHeader(http.StatusInternalServerError)
//		return
//	}
//
//	brokerList, err := database.GetAllBrokers(ctx)
//	if err != nil {
//		log.Logger.ErrorContext(ctx, "get all brokers failed")
//		res.WriteHeader(http.StatusInternalServerError)
//		return
//	}
//
//	var clusterListRes clusterRes
//
//	for _, cluster := range clusterList {
//		clusterRes := clusterInfo{}
//		clusterRes.Id = cluster.ID
//		clusterRes.ClusterName = cluster.ClusterName
//		clusterRes.KafkaVersion = cluster.KafkaVersion
//
//		for _, b := range brokerList {
//			if b.ClusterID == cluster.ID {
//				clusterRes.Brokers = append(clusterRes.Brokers, broker{b.Host, b.Port})
//				break
//			}
//		}
//
//		clusterListRes.Clusters = append(clusterListRes.Clusters, clusterRes)
//	}
//
//	res.WriteHeader(http.StatusOK)
//	err = json.NewEncoder(res).Encode(clusterListRes)
//	if err != nil {
//		log.Logger.ErrorContext(ctx, "encoding response into json failed", err)
//		res.WriteHeader(http.StatusInternalServerError)
//		return
//	}
//
//	log.Logger.TraceContext(ctx, "get all clusters was successful")
//}

func handleDeleteCluster(res http.ResponseWriter, req *http.Request) {
	ctx := traceable_context.WithUUID(uuid.New())
	params := mux.Vars(req)
	clusterName := params["cluster_id"]

	err := database.DeleteCluster(ctx, clusterName)
	if err != nil {
		log.Logger.ErrorContext(ctx, fmt.Sprintf("deleting cluster failed - %v", clusterName), err)
		res.WriteHeader(http.StatusInternalServerError)
		return
	}

	log.Logger.TraceContext(ctx, "cluster deleted successfully", clusterName)
	res.WriteHeader(http.StatusOK)
}

func handleGetAllClusters(res http.ResponseWriter, req *http.Request) {
	ctx := traceable_context.WithUUID(uuid.New())
	var clusterListRes clusterRes

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
	err := json.NewEncoder(res).Encode(clusterListRes)
	if err != nil {
		log.Logger.ErrorContext(ctx, "encoding response into json failed", err)
		res.WriteHeader(http.StatusInternalServerError)
		return
	}

	log.Logger.TraceContext(ctx, "get all clusters was successful")
}

//func handleConnectToCluster(res http.ResponseWriter, req *http.Request) {
//	var clusterReq connectToCluster
//
//	ctx := traceable_context.WithUUID(uuid.New())
//	content, err := ioutil.ReadAll(req.Body)
//	if err != nil {
//		log.Logger.ErrorContext(ctx, "error occurred while reading request", err)
//		res.WriteHeader(http.StatusBadRequest)
//		return
//	}
//
//	err = json.Unmarshal(content, &clusterReq)
//	if err != nil {
//		log.Logger.ErrorContext(ctx, "unmarshal error", err)
//		res.WriteHeader(http.StatusBadRequest)
//		return
//	}
//
//	client, err := kafka.InitClient(ctx, clusterReq.Brokers)
//	if err != nil {
//		log.Logger.ErrorContext(ctx, err)
//		res.WriteHeader(http.StatusServiceUnavailable)
//		return
//	}
//
//	cluster, err := kafka.InitClusterConfig(ctx, clusterReq.ClusterName, clusterReq.Brokers, "")
//	if err != nil {
//		log.Logger.ErrorContext(ctx, err)
//		res.WriteHeader(http.StatusServiceUnavailable)
//		return
//	}
//
//	var clustClient kafka.KCluster
//	clustClient.ClusterID = clusterReq.ClusterID
//	clustClient.Consumer = cluster
//	clustClient.Client = client
//	kafka.ClusterList = append(kafka.ClusterList, clustClient)
//
//	log.Logger.TraceContext(ctx, "connecting to clusterReq was successful", clusterReq.Brokers)
//	//res.Header().Set("ID", strconv.Itoa(clusterReq.ClusterID))
//	res.WriteHeader(http.StatusOK)
//}