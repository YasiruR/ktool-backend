package http

import (
	"encoding/json"
	"github.com/YasiruR/ktool-backend/database"
	"github.com/YasiruR/ktool-backend/domain"
	"github.com/YasiruR/ktool-backend/log"
	"github.com/google/uuid"
	traceable_context "github.com/pickme-go/traceable-context"
	"io/ioutil"
	"net/http"
)

func handleAddNewUser(res http.ResponseWriter, req *http.Request) {
	var addUserReq addUserReq

	ctx := traceable_context.WithUUID(uuid.New())

	content, err := ioutil.ReadAll(req.Body)
	if err != nil {
		log.Logger.ErrorContext(ctx, "error occurred while reading request", err)
		res.WriteHeader(http.StatusBadRequest)
		return
	}

	err = json.Unmarshal(content, &addUserReq)
	if err != nil {
		log.Logger.ErrorContext(ctx, "unmarshal error", err)
		res.WriteHeader(http.StatusBadRequest)
		return
	}

	token := generateToken()

	exists, err := database.AddNewUser(ctx, addUserReq.Username, addUserReq.Password, token, addUserReq.AccessLevel, addUserReq.FirstName, addUserReq.LastName, addUserReq.Email)
	if err != nil {
		if exists {
			log.Logger.ErrorContext(ctx, "add new user request failed", addUserReq.Username)
			res.WriteHeader(http.StatusPreconditionFailed)
			return
		}
		log.Logger.ErrorContext(ctx, "add new user request failed", addUserReq.Username)
		res.WriteHeader(http.StatusInternalServerError)
		return
	}

	var addUserRes userRes
	addUserRes.Token = token
	res.WriteHeader(http.StatusOK)
	err = json.NewEncoder(res).Encode(&addUserRes)
	if err != nil {
		log.Logger.ErrorContext(ctx, "encoding response into json failed", err)
		res.WriteHeader(http.StatusInternalServerError)
		return
	}

	res.WriteHeader(http.StatusOK)
	log.Logger.TraceContext(ctx, "add new user request was successful", addUserReq.Username)
}

func handleLogin(res http.ResponseWriter, req *http.Request) {
	var loginUserReq loginUserReq

	ctx := traceable_context.WithUUID(uuid.New())
	content, err := ioutil.ReadAll(req.Body)
	if err != nil {
		log.Logger.ErrorContext(ctx, "error occurred while reading request", err)
		res.WriteHeader(http.StatusBadRequest)
		return
	}

	err = json.Unmarshal(content, &loginUserReq)
	if err != nil {
		log.Logger.ErrorContext(ctx, "unmarshal error", err)
		res.WriteHeader(http.StatusBadRequest)
		return
	}

	id, ok, err := database.ValidateUserByPassword(ctx, loginUserReq.Username, loginUserReq.Password)
	if err != nil {
		if err.Error() == "incorrect credentials" {
			log.Logger.TraceContext(ctx, "no user encountered for the given credentials", loginUserReq.Username)
			res.WriteHeader(http.StatusUnauthorized)
			return
		}
		log.Logger.ErrorContext(ctx, "login request failed", err)
		res.WriteHeader(http.StatusInternalServerError)
		return
	}
	if ok {
		tokenRetry:
		token := generateToken()

		_, ok, err := database.GetUserByToken(ctx, token)
		if ok {
			log.Logger.ErrorContext(ctx, "generated token already exists in the database")
			goto tokenRetry
		}

		err = database.UpdateToken(ctx, loginUserReq.Username, token)
		if err != nil {
			log.Logger.ErrorContext(ctx, "login request failed")
			return
		}

		var userRes userRes
		userRes.Token = token
		res.WriteHeader(http.StatusOK)
		err = json.NewEncoder(res).Encode(&userRes)
		if err != nil {
			log.Logger.ErrorContext(ctx, "encoding response into json failed")
			res.WriteHeader(http.StatusInternalServerError)
			return
		}

		//check if user is already in connected list
		var exists bool
		var user domain.User
		for _, u := range domain.LoggedInUsers {
			if u.Username == loginUserReq.Username {
				exists = true
				user = u
				break
			}
		}

		//if exists, update only the token
		if !exists {
			user.Username = loginUserReq.Username
			user.Token = token
			user.Id = id
			domain.LoggedInUsers = append(domain.LoggedInUsers, user)
		} else {
			user.Token = token
		}

		log.Logger.TraceContext(ctx, "user logged in successfully", loginUserReq.Username)
	}
}

func generateToken() (token string) {
	sessionToken := uuid.New().String()
	return sessionToken
}