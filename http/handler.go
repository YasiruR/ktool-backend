package http

import (
	"context"
	"encoding/json"
	"github.com/YasiruR/ktool-backend/log"
	"github.com/gorilla/mux"
	"io/ioutil"
	"net/http"
)

//add existing kafka cluster
func handleAddCluster(res http.ResponseWriter, req *http.Request) {
	var addClusterReq reqAddExistingCluster

	ctx := context.Background()
	content, err := ioutil.ReadAll(req.Body)
	if err != nil {
		log.Logger.ErrorContext(ctx, "error occurred while reading request ", err)
	}

	err = json.Unmarshal(content, &addClusterReq)
	if err != nil {
		log.Logger.ErrorContext(ctx, "unmarshal error", err)
	}
}

func handleConnectToCluster(res http.ResponseWriter, req *http.Request) {
	params := mux.Vars(req)
	_ = params["name"]

	//db query to fetch cluster data

}