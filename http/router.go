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
	"os/signal"
	"syscall"
)

func InitRouter() {
	router := mux.NewRouter()

	router.HandleFunc("/cluster/ping", handlePingToServer).Methods("POST")
	router.HandleFunc("/cluster/telnet", handleTelnetToPort).Methods("POST")
	router.HandleFunc("/cluster/add", handleAddCluster).Methods("POST")
	router.HandleFunc("/clusters", handleGetAllClusters).Methods("GET")
	router.HandleFunc("/cluster/connect", handleConnectToCluster).Methods("POST")
	router.HandleFunc("/topics", handleGetTopicsForCluster).Methods("GET")
	router.HandleFunc("/brokers", handleGetBrokersForCluster).Methods("GET")

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
		for _, clustClient := range kafka.ClusterList {
			if clustClient.Client != nil {
				clustClient.Client.Close()
				clustClient.Consumer.Close()
			}
		}

		//closing all server sessions
		for _, session := range cloud.SessionList {
			session.Close()
		}
		os.Exit(0)
	}()

	handler := cors.Default().Handler(router)

	log.Fatal(http.ListenAndServe(":" + service.Cfg.ServicePort, handler))
}
