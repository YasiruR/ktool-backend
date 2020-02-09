package http

import (
	"encoding/json"
	"github.com/YasiruR/ktool-backend/kafka"
	"github.com/YasiruR/ktool-backend/log"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	traceable_context "github.com/pickme-go/traceable-context"
	"net/http"
	"strconv"
)

func handleGetBrokersForCluster(res http.ResponseWriter, req *http.Request) {
	ctx := traceable_context.WithUUID(uuid.New())

	params := mux.Vars(req)
	clusterID, err := strconv.Atoi(params["cluster_id"])
	if err != nil {
		log.Logger.ErrorContext(ctx, "cluster id param is not an int")
		res.WriteHeader(http.StatusBadRequest)
		return
	}

	clustClient, err := kafka.GetClient(ctx, clusterID)
	if err != nil {
		var errRes errorMessage
		res.WriteHeader(http.StatusBadRequest)
		errRes.Mesg = "Cluster id does not exist"
		err := json.NewEncoder(res).Encode(errRes)
		if err != nil {
			log.Logger.ErrorContext(ctx, "encoding error response for add cluster req failed")
		}
		return
	}

	addrList, err := kafka.GetBrokerAddrList(ctx, clustClient.Client)
	if err != nil {
		log.Logger.ErrorContext(ctx, "fetching broker address list failed")
		res.WriteHeader(http.StatusInternalServerError)
		return
	}

	res.WriteHeader(http.StatusOK)
	err = json.NewEncoder(res).Encode(addrList)
	if err != nil {
		log.Logger.ErrorContext(ctx, "encoding broker list response to json failed", err)
		res.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func handleGetTopicsForCluster(res http.ResponseWriter, req *http.Request) {
	ctx := traceable_context.WithUUID(uuid.New())
	params := mux.Vars(req)
	clusterID, err := strconv.Atoi(params["cluster_id"])
	if err != nil {
		log.Logger.ErrorContext(ctx, "cluster id param is not an int")
		res.WriteHeader(http.StatusBadRequest)
		return
	}

	clustClient, err := kafka.GetClient(ctx, clusterID)
	if err != nil {
		var errRes errorMessage
		res.WriteHeader(http.StatusBadRequest)
		errRes.Mesg = "Cluster id does not exist"
		err := json.NewEncoder(res).Encode(errRes)
		if err != nil {
			log.Logger.ErrorContext(ctx, "encoding error response for add cluster req failed")
		}
		return
	}

	topics, err := kafka.GetTopicList(ctx, clustClient.Consumer)
	if err != nil {
		log.Logger.ErrorContext(ctx, "fetching topic data failed", err)
		res.WriteHeader(http.StatusInternalServerError)
		return
	}

	res.WriteHeader(http.StatusOK)
	err = json.NewEncoder(res).Encode(topics)
	if err != nil {
		log.Logger.ErrorContext(ctx, "encoding topics response to json failed", err)
		res.WriteHeader(http.StatusInternalServerError)
		return
	}
}
