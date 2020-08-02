package http

import (
	"encoding/json"
	"fmt"
	"github.com/YasiruR/ktool-backend/database"
	"github.com/YasiruR/ktool-backend/domain"
	"github.com/YasiruR/ktool-backend/iam"
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
	//if addSecretRequest.Validate {
	//	log.Logger.TraceContext(ctx, "secret is being validated ", addSecretRequest.UserId)
	//	valid, err := iam.TestIamPermissions(&addSecretRequest)
	//	if (err != nil) || !valid {
	//		log.Logger.ErrorContext(ctx, "error occurred while validating secret", err)
	//		res.WriteHeader(http.StatusBadRequest)
	//		return
	//	}
	//}
	result := database.AddSecret(ctx, &addSecretRequest)
	res.WriteHeader(http.StatusOK)
	err = json.NewEncoder(res).Encode(&result)
	if err != nil {
		res.WriteHeader(http.StatusInternalServerError)
		log.Logger.ErrorContext(ctx, "response json conversion failed", addSecretRequest.UserId)
	}
	log.Logger.TraceContext(ctx, "add secret request successful", addSecretRequest.UserId)
}

func handleGetAllSecrets(res http.ResponseWriter, req *http.Request) {
	ctx := traceableContext.WithUUID(uuid.New())
	//var searchSecretRequest SearchSecretsRequest

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

	OwnerId := req.FormValue("owner_id")
	ServiceProvider := req.FormValue("service_provider")

	if err != nil {
		log.Logger.ErrorContext(ctx, "unmarshal error", err)
		res.WriteHeader(http.StatusBadRequest)
		return
	}
	result := database.GetAllSecretsByUserExternal(ctx, OwnerId, ServiceProvider)

	res.WriteHeader(http.StatusOK)
	err = json.NewEncoder(res).Encode(&result)
	if err != nil {
		res.WriteHeader(http.StatusOK)
		log.Logger.ErrorContext(ctx, "response json conversion failed", OwnerId)
	}
	log.Logger.TraceContext(ctx, "search secret request successful", OwnerId)
}

func handleGetSecret(res http.ResponseWriter, req *http.Request) {
	ctx := traceableContext.WithUUID(uuid.New())
	//var searchSecretRequest SearchSecretsRequest

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

	//Name := req.FormValue("name")
	//OwnerId := req.FormValue("owner_id")
	Provider := req.FormValue("service_provider")
	SecretId := req.FormValue("secret_id")

	if err != nil {
		log.Logger.ErrorContext(ctx, "unmarshal error", err)
		res.WriteHeader(http.StatusBadRequest)
		return
	}

	//result := database.GetSecretInternal(ctx, Name, OwnerId, Provider)
	result := database.GetSecretExternal(ctx, SecretId, Provider)

	res.WriteHeader(http.StatusOK)
	err = json.NewEncoder(res).Encode(&result)
	if err != nil {
		res.WriteHeader(http.StatusOK)
		log.Logger.ErrorContext(ctx, "response json conversion failed", SecretId)
	}
	log.Logger.TraceContext(ctx, "search secret request successful", SecretId)
}

func handleDeleteSecret(res http.ResponseWriter, req *http.Request) {
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

	secretId := req.FormValue("secret_id")
	if err != nil {
		log.Logger.ErrorContext(ctx, "unmarshal error", err)
		res.WriteHeader(http.StatusBadRequest)
		return
	}

	log.Logger.TraceContext(ctx, "Delete secret request received")
	result, err := database.DeleteSecret(ctx, secretId)

	if (err != nil) || (result != true) {
		res.WriteHeader(http.StatusInternalServerError)
		log.Logger.ErrorContext(ctx, "delete secret failed. secretid: ", secretId)
	} else {
		res.WriteHeader(http.StatusOK)
		log.Logger.TraceContext(ctx, "delete secret request successful. secretid: ", secretId)
	}
}

func handleUpdateSecret(res http.ResponseWriter, req *http.Request) {
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

	var updateSecretRequest domain.CloudSecret
	err = json.Unmarshal(content, &updateSecretRequest)
	if err != nil {
		log.Logger.ErrorContext(ctx, "unmarshal error", err)
		res.WriteHeader(http.StatusBadRequest)
		return
	}

	log.Logger.TraceContext(ctx, "Delete secret request received")
	// todo: replace with external call
	result := database.UpdateSecret(ctx, &updateSecretRequest)

	res.WriteHeader(http.StatusOK)
	err = json.NewEncoder(res).Encode(&result)
	if err != nil {
		res.WriteHeader(http.StatusInternalServerError)
		log.Logger.ErrorContext(ctx, "response json conversion failed", updateSecretRequest.UserId)
	}
	log.Logger.TraceContext(ctx, "update secret request successful", updateSecretRequest.UserId)
}

func handleValidateSecret(res http.ResponseWriter, req *http.Request) {
	ctx := traceableContext.WithUUID(uuid.New())
	var validateSecretRequest domain.CloudSecret

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

	err = json.Unmarshal(content, &validateSecretRequest)
	if err != nil {
		log.Logger.ErrorContext(ctx, "unmarshal error", err)
		res.WriteHeader(http.StatusBadRequest)
		return
	}

	fmt.Println("Validate secret request received")
	log.Logger.TraceContext(ctx, "secret is being validated ", validateSecretRequest.UserId)
	valid, err := iam.TestIamPermissions(&validateSecretRequest)
	if (err != nil) || !valid {
		log.Logger.ErrorContext(ctx, "error occurred while validating secret", err)
		res.WriteHeader(http.StatusInternalServerError)
		return
	}
	res.WriteHeader(http.StatusOK)
	//err = json.NewEncoder(res).Encode(&result)
	if err != nil {
		res.WriteHeader(http.StatusInternalServerError)
		log.Logger.ErrorContext(ctx, "response json conversion failed", validateSecretRequest.UserId)
	}
	log.Logger.TraceContext(ctx, "validate secret request successful", validateSecretRequest.UserId)
}
