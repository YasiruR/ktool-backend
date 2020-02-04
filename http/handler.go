package http

import (
	"encoding/json"
	"github.com/YasiruR/ktool-backend/cloud"
	"github.com/YasiruR/ktool-backend/database"
	"github.com/YasiruR/ktool-backend/kafka"
	"github.com/YasiruR/ktool-backend/log"
	"github.com/google/uuid"
	//"github.com/gorilla/mux"
	traceable_context "github.com/pickme-go/traceable-context"
	"io/ioutil"
	"net/http"
)

//add existing kafka cluster
func handleAddCluster(res http.ResponseWriter, req *http.Request) {
	var addClusterReq addExistingCluster
	var reqFailed = false

	ctx := traceable_context.WithUUID(uuid.New())
	content, err := ioutil.ReadAll(req.Body)
	if err != nil {
		log.Logger.ErrorContext(ctx, "error occurred while reading request", err)
		reqFailed = true
		res.WriteHeader(http.StatusBadRequest)
	}

	err = json.Unmarshal(content, &addClusterReq)
	if err != nil {
		log.Logger.ErrorContext(ctx, "unmarshal error", err)
		reqFailed = true
		res.WriteHeader(http.StatusBadRequest)
	}

	//proceeds to db query
	//note : frontend validations should be added to request parameters
	err = database.AddNewCluster(ctx, addClusterReq.ClusterName, addClusterReq.KafkaVersion)
	if err != nil {
		log.Logger.ErrorContext(ctx, "add new cluster db transaction failed")
		reqFailed = true
		res.WriteHeader(http.StatusInternalServerError)
	} else {
		err = database.AddNewZookeeper(ctx, addClusterReq.ZookeeperHost, addClusterReq.ZookeeperPort, addClusterReq.ClusterName)
		if err != nil {
			log.Logger.ErrorContext(ctx, "add new zookeeper db transaction failed")
			//if adding new zookeeper failed reverts the adding cluster query as well
			err = database.DeleteCluster(ctx, addClusterReq.ClusterName)
			if err != nil {
				log.Logger.ErrorContext(ctx, "deleting newly added cluster failed")
				//cluster table modified but zookeeper table is not
				res.WriteHeader(http.StatusConflict)
			} else {
				log.Logger.TraceContext(ctx, "deleting newly added cluster was successful")
				res.WriteHeader(http.StatusInternalServerError)
			}
			reqFailed = true
		}
	}

	if reqFailed == false {
		log.Logger.TraceContext(ctx, "cluster stored in the database successfully", addClusterReq.ClusterName)
		res.WriteHeader(http.StatusOK)
	}
}

//handle ping to new server
func handlePingToZookeeper(res http.ResponseWriter, req *http.Request) {
	var testClusterReq server
	var reqFailed = false

	ctx := traceable_context.WithUUID(uuid.New())
	content, err := ioutil.ReadAll(req.Body)
	if err != nil {
		log.Logger.ErrorContext(ctx, "error occurred while reading request", err)
		reqFailed = true
		res.WriteHeader(http.StatusBadRequest)
	}

	err = json.Unmarshal(content, &testClusterReq)
	if err != nil {
		log.Logger.ErrorContext(ctx, "unmarshal error", err)
		reqFailed = true
		res.WriteHeader(http.StatusBadRequest)
	}

	//ssh ping to server
	//note : req address validations should be added in frontend
	ok, err := cloud.PingToServer(ctx, testClusterReq.Host)
	if err != nil {
		log.Logger.ErrorContext(ctx, "ping to server failed")
		reqFailed = true
		res.WriteHeader(http.StatusInternalServerError)
	} else {
		if ok {
			reqFailed = false
			log.Logger.TraceContext(ctx, "ping to server is successful")
		}
	}

	if reqFailed == false {
		res.WriteHeader(http.StatusOK)
	}
}

func handleTelnetToPort(res http.ResponseWriter, req *http.Request) {
	var testClusterReq server
	var reqFailed = false

	log.Logger.Trace("req : ", req)

	ctx := traceable_context.WithUUID(uuid.New())
	content, err := ioutil.ReadAll(req.Body)
	if err != nil {
		log.Logger.ErrorContext(ctx, "error occurred while reading request", err)
		reqFailed = true
		res.WriteHeader(http.StatusBadRequest)
	}

	err = json.Unmarshal(content, &testClusterReq)
	if err != nil {
		log.Logger.ErrorContext(ctx, "unmarshal error", err)
		reqFailed = true
		res.WriteHeader(http.StatusBadRequest)
	}

	ok, err := cloud.TelnetToPort(ctx, testClusterReq.Host, testClusterReq.Port)
	if err != nil {
		log.Logger.ErrorContext(ctx, "telnet to server failed")
		reqFailed = true
		res.WriteHeader(http.StatusInternalServerError)
	} else {
		if ok {
			reqFailed = false
			log.Logger.TraceContext(ctx, "telnet to server is successful")
		}
	}

	if reqFailed == false {
		res.WriteHeader(http.StatusOK)
	}
}

func handleGetAllClusters(res http.ResponseWriter, req *http.Request) {
	var reqFailed = false

	ctx := traceable_context.WithUUID(uuid.New())
	clusterList, err := database.GetAllClusters(ctx)
	if err != nil {
		log.Logger.ErrorContext(ctx, "get all clusters failed")
		reqFailed = true
		res.WriteHeader(http.StatusInternalServerError)
	}

	var clusterListRes []clusterInfo

	for _, cluster := range clusterList {
		clusterRes := clusterInfo{}
		clusterRes.Id = cluster.ID
		clusterRes.ClusterName = cluster.ClusterName
		clusterRes.KafkaVersion = cluster.KafkaVersion

		clusterListRes = append(clusterListRes, clusterRes)
	}

	err = json.NewEncoder(res).Encode(clusterListRes)
	if err != nil {
		log.Logger.ErrorContext(ctx, "encoding response into json failed", err)
		res.WriteHeader(http.StatusInternalServerError)
		reqFailed = true
	}

	if reqFailed == false {
		log.Logger.TraceContext(ctx, "get all clusters was successful")
		res.WriteHeader(http.StatusOK)
	}
}

func handleConnectToCluster(res http.ResponseWriter, req *http.Request) {
	var cluster connectToCluster
	var reqFailed = false

	ctx := traceable_context.WithUUID(uuid.New())
	content, err := ioutil.ReadAll(req.Body)
	if err != nil {
		log.Logger.ErrorContext(ctx, "error occurred while reading request", err)
		reqFailed = true
		res.WriteHeader(http.StatusBadRequest)
	}

	err = json.Unmarshal(content, &cluster)
	if err != nil {
		log.Logger.ErrorContext(ctx, "unmarshal error", err)
		reqFailed = true
		res.WriteHeader(http.StatusBadRequest)
	}

	err = kafka.InitClient(ctx, cluster.Brokers)
	if err != nil {
		log.Logger.ErrorContext(ctx, err)
		reqFailed = true
		res.WriteHeader(http.StatusServiceUnavailable)
	}

	err = kafka.InitClusterConfig(ctx, cluster.Brokers)
	if err != nil {
		log.Logger.ErrorContext(ctx, err)
		reqFailed = true
		res.WriteHeader(http.StatusServiceUnavailable)
	}

	if reqFailed == false {
		log.Logger.TraceContext(ctx, "connecting to cluster was successful", cluster.Brokers)
		res.WriteHeader(http.StatusOK)
	}
}

func handleGetAllBrokers(res http.ResponseWriter, req *http.Request) {
	reqFailed := false
	ctx := traceable_context.WithUUID(uuid.New())

	addrList, err := kafka.GetBrokerAddrList(ctx)
	if err != nil {
		log.Logger.ErrorContext(ctx, "fetching broker address list failed")
		reqFailed = true
		res.WriteHeader(http.StatusInternalServerError)
	}

	err = json.NewEncoder(res).Encode(addrList)
	if err != nil {
		log.Logger.ErrorContext(ctx, "encoding broker list response to json failed", err)
		reqFailed = true
		res.WriteHeader(http.StatusInternalServerError)
	}

	if reqFailed == false {
		res.WriteHeader(http.StatusOK)
	}
}

func handleGetTopicsForBroker(res http.ResponseWriter, req *http.Request) {
	//params := mux.Vars(req)
	reqFailed := false

	ctx := traceable_context.WithUUID(uuid.New())

	topics, err := kafka.GetTopicList(ctx)
	if err != nil {
		log.Logger.ErrorContext(ctx, "fetching topic data failed", err)
		reqFailed = true
		res.WriteHeader(http.StatusInternalServerError)
	}

	err = json.NewEncoder(res).Encode(topics)
	if err != nil {
		log.Logger.ErrorContext(ctx, "encoding topics response to json failed", err)
		reqFailed = true
		res.WriteHeader(http.StatusInternalServerError)
	}

	if reqFailed == false {
		res.WriteHeader(http.StatusOK)
	}
}