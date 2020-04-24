package http

import (
	"encoding/json"
	"fmt"
	"github.com/YasiruR/ktool-backend/database"
	"github.com/YasiruR/ktool-backend/kafka"
	"github.com/YasiruR/ktool-backend/log"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	traceable_context "github.com/pickme-go/traceable-context"
	"io/ioutil"
	"net"
	"net/http"
	"strconv"
	"strings"
)

func handleAddSecret(res http.ResponseWriter, req *http.Request) {
	ctx := traceable_context.WithUUID(uuid.New())
	var addSecretRequest AddSecretRequest

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

	retryCount := 0
checkIfClusterExists:
	_, err = database.GetAllSecretsByUser(ctx, strings.TrimSpace(addClusterReq.ClusterName))

	if err == nil {
		log.Logger.ErrorContext(ctx, "cluster name already exists", addClusterReq.ClusterName)
		var errRes errorMessage
		res.WriteHeader(http.StatusPreconditionFailed)
		errRes.Mesg = "Cluster name already exists. Please provide a different name."
		err := json.NewEncoder(res).Encode(errRes)
		if err != nil {
			log.Logger.ErrorContext(ctx, "encoding error response for add cluster req failed")
		}
		return
	} else if err.Error() == "no rows found" {
		//when cluster is eligible to be added

		var brokerAddrList []string
		for _, broker := range addClusterReq.Brokers {
			brokerAddrList = append(brokerAddrList, broker.Host+":"+strconv.Itoa(broker.Port))
		}

		for _, brokerAddr := range brokerAddrList {
			ok := kafka.CheckIfBrokerExists(ctx, brokerAddr)
			if ok {
				var errRes errorMessage
				res.WriteHeader(http.StatusPreconditionFailed)
				errRes.Mesg = fmt.Sprintf("You have already added broker (%v) into a cluster", brokerAddr)
				err := json.NewEncoder(res).Encode(errRes)
				if err != nil {
					log.Logger.ErrorContext(ctx, "encoding error response for add cluster req failed")
				}
				return
			}
		}

		//get all relevant brokers
		client, err := kafka.InitClient(ctx, brokerAddrList)
		if err != nil {
			log.Logger.ErrorContext(ctx, "add cluster request failed", addClusterReq.ClusterName, brokerAddrList)
			var errRes errorMessage
			res.WriteHeader(http.StatusPreconditionFailed)
			errRes.Mesg = fmt.Sprintf("Could not find the cluster for brokers - %v", brokerAddrList)
			err := json.NewEncoder(res).Encode(errRes)
			if err != nil {
				log.Logger.ErrorContext(ctx, "encoding error response for test cluster req failed")
			}
			return
		}

		tmpBrokList, err := kafka.GetBrokerAddrList(ctx, client)
		if err != nil {
			log.Logger.ErrorContext(ctx, "test connection to cluster failed", addClusterReq.ClusterName)
			var errRes errorMessage
			res.WriteHeader(http.StatusPreconditionFailed)
			errRes.Mesg = fmt.Sprintf("Could not fetch rest of the brokers for cluster - %v", addClusterReq.ClusterName)
			err := json.NewEncoder(res).Encode(errRes)
			if err != nil {
				log.Logger.ErrorContext(ctx, "encoding error response for test cluster req failed")
			}
			return
		}

		var hosts []string
		var ports []int
		for _, tmpBroker := range tmpBrokList {
			host, portStr, err := net.SplitHostPort(tmpBroker)
			if err != nil {
				log.Logger.ErrorContext(ctx, fmt.Sprintf("splitting host and port failed for %v", tmpBroker), err)
				var errRes errorMessage
				res.WriteHeader(http.StatusInternalServerError)
				errRes.Mesg = fmt.Sprintf("Splitting host and port failed for cluster - %v", addClusterReq.ClusterName)
				err := json.NewEncoder(res).Encode(errRes)
				if err != nil {
					log.Logger.ErrorContext(ctx, "encoding error response for test cluster req failed")
				}
				return
			}

			hosts = append(hosts, host)
			port, err := strconv.Atoi(portStr)
			if err != nil {
				log.Logger.ErrorContext(ctx, "port does not contain an int", tmpBroker)
				var errRes errorMessage
				res.WriteHeader(http.StatusBadRequest)
				errRes.Mesg = fmt.Sprintf("At least one of the ports is not an integer - %v", addClusterReq.ClusterName)
				err := json.NewEncoder(res).Encode(errRes)
				if err != nil {
					log.Logger.ErrorContext(ctx, "encoding error response for test cluster req failed")
				}
				return
			}
			ports = append(ports, port)
		}

		//proceeds to db query
		//note : frontend validations should be added to request parameters
		err = database.AddNewCluster(ctx, addClusterReq.ClusterName, addClusterReq.KafkaVersion)
		if err != nil {
			log.Logger.ErrorContext(ctx, "add new cluster db transaction failed")
			res.WriteHeader(http.StatusInternalServerError)
			return
		}

		err = database.AddNewBrokers(ctx, hosts, ports, addClusterReq.ClusterName)
		if err != nil {
			log.Logger.ErrorContext(ctx, "add new brokers db transaction failed", err)
			res.WriteHeader(http.StatusInternalServerError)

			//if adding new brokers failed, reverts the adding cluster query as well
			err = database.DeleteCluster(ctx, addClusterReq.ClusterName)
			if err != nil {
				log.Logger.ErrorContext(ctx, "deleting newly added cluster failed")
			} else {
				log.Logger.TraceContext(ctx, "deleting newly added cluster was successful")
			}
			return
		}
		log.Logger.TraceContext(ctx, "cluster stored in the database successfully", addClusterReq.ClusterName)
	} else {
		retryCount += 1
		if retryCount <= checkClusterRetryCount {
			goto checkIfClusterExists
		}
		log.Logger.ErrorContext(ctx, fmt.Sprintf("checking if cluster exists failed %v times", retryCount-1))
		res.WriteHeader(http.StatusInternalServerError)
		return
	}

	kafka.InitAllClusters()
}

func handleGetAllSecrets(res http.ResponseWriter, req *http.Request) {
}

func handleDeleteSecret(res http.ResponseWriter, req *http.Request) {
}
