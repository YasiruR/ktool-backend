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

	router.HandleFunc("/user/register", handleAddNewUser).Methods("POST")
	router.HandleFunc("/user/login", handleLogin).Methods("POST")
	router.HandleFunc("/user/logout", handleLogout).Methods("GET")

	router.HandleFunc("/cluster/ping", handlePingToServer).Methods("POST")
	router.HandleFunc("/cluster/telnet", handleTestConnectionToCluster).Methods("POST")
	router.HandleFunc("/cluster/add", handleAddCluster).Methods("POST")
	router.HandleFunc("/cluster", handleDeleteCluster).Methods("DELETE")

	router.HandleFunc("/clusters", handleGetAllClusters).Methods("GET")
	router.HandleFunc("/cluster/connect", handleConnectToCluster).Methods("GET")
	router.HandleFunc("/cluster/disconnect", handleDisconnectCluster).Methods("GET")

	router.HandleFunc("/topics", handleGetTopicsForCluster).Methods("GET")
	router.HandleFunc("/brokers", handleGetBrokersForCluster).Methods("GET")

	router.HandleFunc("/secret/create", handleAddSecret).Methods("POST")
	router.HandleFunc("/secret/get/all", handleGetAllSecrets).Methods("GET")
	router.HandleFunc("/secret/get", handleGetSecret).Methods("GET")
	router.HandleFunc("/secret/delete", handleDeleteSecret).Methods("DELETE")
	router.HandleFunc("/secret/update", handleUpdateSecret).Methods("PATCH")
	router.HandleFunc("/secret/validate", handleValidateSecret).Methods("POST")

	router.HandleFunc("/kubernetes", handleGetAllKubClusters).Methods("GET")
	router.HandleFunc("/kubernetes", handleCreateKubCluster).Methods("POST")
	router.HandleFunc("/kubernetes/validate", handleValidateClusterName).Methods("GET")
	router.HandleFunc("/kubernetes/status", handleCheckClusterCreationStatus).Methods("GET")
	router.HandleFunc("/kubernetes/resources", handleGetGkeResource).Methods("GET")
	router.HandleFunc("/kubernetes/recommend", handleRecommendGkeResource).Methods("GET")

	router.HandleFunc("/kubernetes/gke", handleGetAllGkeKubClusters).Methods("GET")
	router.HandleFunc("/kubernetes/gke/status", handleCheckGkeClusterCreationStatus).Methods("GET")
	//router.HandleFunc("/kubernetes/gke", handleCreateGkeKubClusters).Methods("POST")

	router.HandleFunc("/kubernetes/eks", handleGetAllEksKubClusters).Methods("GET")
	router.HandleFunc("/kubernetes/eks/status", handleCheckEksClusterCreationStatus).Methods("GET")
	//router.HandleFunc("/kubernetes/eks", handleCreateEksKubClusters).Methods("POST")
	router.HandleFunc("/kubernetes/eks", handleDeleteEksKubClusters).Methods("DELETE")
	router.HandleFunc("/kubernetes/eks/nodegroup", handleCreateEksNodeGroup).Methods("POST")
	router.HandleFunc("/kubernetes/eks/nodegroup/status", handleCheckEksNodeGroupCreationStatus).Methods("GET")
	//router.HandleFunc("/kubernetes/ec2/vpc", handleGetVPCConfigForRegion).Methods("GET")

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
		log.Trace("closing all the initialized cluster connections")

		//closing all server sessions
		for _, session := range cloud.SessionList {
			session.Close()
		}
		os.Exit(0)
	}()

	handler := cors.AllowAll().Handler(router)

	log.Fatal(http.ListenAndServe(":"+service.Cfg.ServicePort, handler))
}
