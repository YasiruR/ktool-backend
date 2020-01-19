package http

import (
	"fmt"
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
	router.HandleFunc("/add_connection", handleAddCluster).Methods("POST")

	osChannel := make(chan os.Signal, 1)
	signal.Notify(osChannel, syscall.SIGINT, syscall.SIGKILL)

	go func() {
		sig := <-osChannel
		log.Debug(fmt.Sprintf("\nprogram exits due to %v signal", sig))
		err := database.Database.Close()
		if err != nil {
			log.Error("error occurred in closing mysql connection")
		}
		os.Exit(0)
	}()

	log.Fatal(http.ListenAndServe(service.Cfg.ServicePort, router))
}
