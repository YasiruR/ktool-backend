package http

import (
	"fmt"
	"github.com/YasiruR/ktool-backend/cloud"
	"github.com/YasiruR/ktool-backend/database"
	"github.com/YasiruR/ktool-backend/kafka"
	"github.com/YasiruR/ktool-backend/service"
	"github.com/gorilla/mux"
	"github.com/pickme-go/log"
	"github.com/rs/cors"
	"net/http"
	"os"
	"os/exec"
	"os/signal"
	"syscall"
)

func InitRouter() {
	router := mux.NewRouter()

	router.HandleFunc("/user/register", handleAddNewUser).Methods(http.MethodPost)
	router.HandleFunc("/user/login", handleLogin).Methods(http.MethodPost)
	router.HandleFunc("/user/logout", handleLogout).Methods(http.MethodGet)

	router.HandleFunc("/cluster/ping", handlePingToServer).Methods(http.MethodPost)
	router.HandleFunc("/cluster/telnet", handleTestConnectionToCluster).Methods(http.MethodPost)
	router.HandleFunc("/cluster/add", handleAddCluster).Methods(http.MethodPost)
	router.HandleFunc("/cluster", handleDeleteCluster).Methods(http.MethodDelete)

	router.HandleFunc("/clusters", handleGetAllClusters).Methods(http.MethodGet)
	router.HandleFunc("/cluster/connect", handleConnectToCluster).Methods(http.MethodGet)
	router.HandleFunc("/cluster/disconnect", handleDisconnectCluster).Methods(http.MethodGet)

	router.HandleFunc("/topics", handleGetTopicsForCluster).Methods(http.MethodGet)
	router.HandleFunc("/brokers", handleGetBrokersForCluster).Methods(http.MethodGet)
	router.HandleFunc("/cluster/broker_overview", handleGetBrokerOverview).Methods(http.MethodGet)

	osChannel := make(chan os.Signal, 1)
	signal.Notify(osChannel, syscall.SIGINT, syscall.SIGKILL)

	//handle OS kill and interrupt signals to close all connections
	go func() {
		sig := <-osChannel
		log.Debug(fmt.Sprintf("\nprogram exits due to %v signal", sig))
		err := database.Db.Close()
		if err != nil {
			log.Error("error occurred in closing mysql connection")
		}

		//closing cluster connections
		log.Trace("closing all the initialized cluster connections")
		for _, clustClient := range kafka.ClusterList {
			if clustClient.Client != nil {
				clustClient.Client.Close()
				clustClient.Consumer.Close()
			}
		}

		log.Trace("closing all spawned docker containers")
		//stop and remove docker prometheus container
		stopDocker := exec.Command("/bin/sh", "-c", "sudo docker stop prometheus")
		stopOutput, err := stopDocker.CombinedOutput()
		if err != nil {
			log.Info(err,"failed to stop docker prometheus container", string(stopOutput))
		} else {
			log.Trace("prometheus docker container is terminated")
			rmDocker := exec.Command("/bin/sh", "-c", "sudo docker rm prometheus")
			rmOutput, err := rmDocker.CombinedOutput()
			if err != nil {
				log.Error(err,"failed to remove stopped docker prometheus container", string(rmOutput))
			} else {
				log.Info("prometheus docker container is removed")
			}
		}

		//closing all server sessions
		for _, session := range cloud.SessionList {
			session.Close()
		}
		os.Exit(0)
	}()

	handler := cors.AllowAll().Handler(router)

	log.Fatal(http.ListenAndServe(":" + service.Cfg.ServicePort, handler))
}
