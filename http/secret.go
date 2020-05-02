package http

import (
	"encoding/json"
	"fmt"
	"github.com/YasiruR/ktool-backend/database"
	"github.com/YasiruR/ktool-backend/domain"
	"github.com/YasiruR/ktool-backend/log"
	"github.com/google/uuid"
	traceableContext "github.com/pickme-go/traceable-context"
	"io/ioutil"
	"net/http"
	"strings"
)

func handleAddSecret(res http.ResponseWriter, req *http.Request) {
	ctx := traceableContext.WithUUID(uuid.New())
	var addSecretRequest domain.CloudSecret

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

	err = json.Unmarshal(content, &addSecretRequest)
	if err != nil {
		log.Logger.ErrorContext(ctx, "unmarshal error", err)
		res.WriteHeader(http.StatusBadRequest)
		return
	}

	fmt.Println("Add secret request received")
	result := database.AddSecret(ctx, addSecretRequest)
	res.WriteHeader(http.StatusOK)
	err = json.NewEncoder(res).Encode(&result)
	if err != nil {
		res.WriteHeader(http.StatusOK)
		log.Logger.ErrorContext(ctx, "response json conversion failed", addSecretRequest.UserId)
	}
	log.Logger.TraceContext(ctx, "add secret request successful", addSecretRequest.UserId)
}

func handleGetAllSecrets(res http.ResponseWriter, req *http.Request) {
	ctx := traceableContext.WithUUID(uuid.New())
	var searchSecretRequest SearchSecretsRequest

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

	err = json.Unmarshal(content, &searchSecretRequest)
	if err != nil {
		log.Logger.ErrorContext(ctx, "unmarshal error", err)
		res.WriteHeader(http.StatusBadRequest)
		return
	}

	fmt.Println("Search secret request received")
	fmt.Println(&searchSecretRequest)
	// todo: replace with external call
	result := database.GetAllSecretsByUserExternal(ctx, searchSecretRequest.OwnerId, searchSecretRequest.ServiceProvider)

	res.WriteHeader(http.StatusOK)
	err = json.NewEncoder(res).Encode(&result)
	if err != nil {
		res.WriteHeader(http.StatusOK)
		log.Logger.ErrorContext(ctx, "response json conversion failed", searchSecretRequest.OwnerId)
	}
	log.Logger.TraceContext(ctx, "search secret request successful", searchSecretRequest.OwnerId)
}

func handleDeleteSecret(res http.ResponseWriter, req *http.Request) {
}
