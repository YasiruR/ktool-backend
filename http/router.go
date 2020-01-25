package http

import (
	"fmt"
	"github.com/YasiruR/ktool-backend/cloud"
	"github.com/YasiruR/ktool-backend/database"
	"github.com/YasiruR/ktool-backend/service"
	"github.com/gorilla/mux"
	"github.com/pickme-go/log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

func InitRouter() {
	router := mux.NewRouter()
	router.HandleFunc("/add-cluster", handleAddCluster).Methods("POST")
	router.HandleFunc("/ping-server", handlePingToZookeeper).Methods("POST")

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

		//closing all server sessions
		for _, session := range cloud.SessionList {
			session.Close()
		}
		os.Exit(0)
	}()

	log.Fatal(http.ListenAndServe(service.Cfg.ServicePort, router))
}
