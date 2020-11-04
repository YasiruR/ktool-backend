package http

import (
	"encoding/json"
	"github.com/YasiruR/ktool-backend/database"
	"github.com/YasiruR/ktool-backend/kafka"
	"github.com/YasiruR/ktool-backend/log"
	"github.com/google/uuid"
	traceable_context "github.com/pickme-go/traceable-context"
	"net/http"
	"strconv"
	"strings"
)

func handleGetBrokersForCluster(res http.ResponseWriter, req *http.Request) {
	ctx := traceable_context.WithUUID(uuid.New())

	//user validation by token header
	tokenHeader := req.Header.Get("Authorization")
	if len(strings.Split(tokenHeader, "Bearer")) < 2 {
		log.Logger.ErrorContext(ctx, "token format is invalid", tokenHeader)
		res.WriteHeader(http.StatusBadRequest)
		return
	}

	token := strings.TrimSpace(strings.Split(tokenHeader, "Bearer")[1])
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

	//params := mux.Vars(req)
	clusterID, err := strconv.Atoi(req.FormValue("cluster_id"))
	if err != nil {
		log.Logger.ErrorContext(ctx, "cluster id conversion from string to int failed", err)
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

	//user validation by token header
	tokenHeader := req.Header.Get("Authorization")
	if len(strings.Split(tokenHeader, "Bearer")) < 2 {
		log.Logger.ErrorContext(ctx, "token format is invalid", tokenHeader)
		res.WriteHeader(http.StatusBadRequest)
		return
	}
	token := strings.TrimSpace(strings.Split(tokenHeader, "Bearer")[1])
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

	//params := mux.Vars(req)
	clusterID, err := strconv.Atoi(req.FormValue("cluster_id"))
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
