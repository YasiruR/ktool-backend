package http

import (
	"context"
	"encoding/json"
	"github.com/YasiruR/ktool-backend/cloud"
	"github.com/YasiruR/ktool-backend/database"
	"github.com/YasiruR/ktool-backend/log"
	"github.com/gorilla/mux"
	"io/ioutil"
	"net/http"
)

//add existing kafka cluster
func handleAddCluster(res http.ResponseWriter, req *http.Request) {
	var addClusterReq reqAddExistingCluster
	var reqFailed = false

	ctx := context.Background()
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
	err = database.AddNewCluster(ctx, addClusterReq.ClusterName, addClusterReq.ClusterVersion, addClusterReq.ZookeeperHost, addClusterReq.ZookeeperPort)
	if err != nil {
		log.Logger.ErrorContext(ctx, "add new cluster db transaction failed", err)
		reqFailed = true
		res.WriteHeader(http.StatusInternalServerError)
	}

	if reqFailed == false {
		log.Logger.TraceContext(ctx, "cluster stored in the database successfully", addClusterReq.ClusterName)
		res.WriteHeader(http.StatusOK)
	}
}

//handle ping to new server
func handlePingToZookeeper(res http.ResponseWriter, req *http.Request) {
	var testClusterReq reqTestNewCluster
	var reqFailed = false

	ctx := context.Background()
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

func handleConnectToCluster(res http.ResponseWriter, req *http.Request) {
	params := mux.Vars(req)
	_ = params["name"]

	//db query to fetch cluster data

}