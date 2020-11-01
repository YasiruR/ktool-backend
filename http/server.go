package http

import (
	"encoding/json"
	"github.com/YasiruR/ktool-backend/cloud"
	"github.com/YasiruR/ktool-backend/domain"
	"github.com/YasiruR/ktool-backend/log"
	"github.com/google/uuid"
	//"github.com/gorilla/mux"
	traceable_context "github.com/pickme-go/traceable-context"
	"io/ioutil"
	"net/http"
)

func handleTelnetToPort(res http.ResponseWriter, req *http.Request) {
	//var testClusterReq server
	//
	//ctx := traceable_context.WithUUID(uuid.New())
	//
	//content, err := ioutil.ReadAll(req.Body)
	//if err != nil {
	//	log.Logger.ErrorContext(ctx, "error occurred while reading request", err)
	//	res.WriteHeader(http.StatusBadRequest)
	//	return
	//}
	//
	//err = json.Unmarshal(content, &testClusterReq)
	//if err != nil {
	//	log.Logger.ErrorContext(ctx, "unmarshal error", err)
	//	res.WriteHeader(http.StatusBadRequest)
	//	return
	//}
	//
	//log.Logger.TraceContext(ctx, testClusterReq, "request")
	//
	//ok, err := cloud.TelnetToPort(ctx, testClusterReq.Host, testClusterReq.Port)
	//if err != nil {
	//	log.Logger.ErrorContext(ctx, "telnet to server failed")
	//	res.WriteHeader(http.StatusInternalServerError)
	//	return
	//}
	//
	//if ok {
	//	log.Logger.TraceContext(ctx, "telnet to server is successful")
	//	res.WriteHeader(http.StatusOK)
	//}

	//res.WriteHeader(http.StatusOK)
}

//handle ping to new server
func handlePingToServer(res http.ResponseWriter, req *http.Request) {
	var testClusterReq domain.Server

	ctx := traceable_context.WithUUID(uuid.New())
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

	//ssh ping to server
	//note : req address validations should be added in frontend
	ok, err := cloud.PingToServer(ctx, testClusterReq.Host)
	if err != nil {
		log.Logger.ErrorContext(ctx, "ping to server failed")
		res.WriteHeader(http.StatusInternalServerError)
		return
	}

	if ok {
		log.Logger.TraceContext(ctx, "ping to server is successful")
		res.WriteHeader(http.StatusOK)
	}

}
