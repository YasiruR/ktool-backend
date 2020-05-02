package http

import (
	"encoding/json"
	"fmt"
	"github.com/YasiruR/ktool-backend/database"
	kubernetes "github.com/YasiruR/ktool-backend/kuberenetes"
	"github.com/YasiruR/ktool-backend/log"
	"github.com/google/uuid"
	traceableContext "github.com/pickme-go/traceable-context"
	"io/ioutil"
	"net/http"
	"strings"
)

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

	content, err := ioutil.ReadAll(req.Body)
	if err != nil {
		log.Logger.ErrorContext(ctx, "error occurred while reading request", err)
		res.WriteHeader(http.StatusBadRequest)
		return
	}

	UserId := string(content)
	err = json.Unmarshal(content, &UserId)
	if err != nil {
		log.Logger.ErrorContext(ctx, "unmarshal error", err)
		res.WriteHeader(http.StatusBadRequest)
		return
	}

	fmt.Printf("List Kub clusters request received %s\n", UserId)
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
