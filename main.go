package main

import (
	"github.com/YasiruR/ktool-backend/database"
	"github.com/YasiruR/ktool-backend/http"
	"github.com/YasiruR/ktool-backend/kafka"
	kubernetes "github.com/YasiruR/ktool-backend/kuberenetes"
	"github.com/YasiruR/ktool-backend/log"
	"github.com/YasiruR/ktool-backend/service"
	"time"
)

func main() {
	log.Cfg.LoadConfigurations()
	log.Init()

	database.Cfg.LoadConfigurations()
	database.Init()

	//cloud.Init()

	service.Cfg.LoadConfigurations()

	kafka.InitAllClusters()

	//refresh cluster data
	ticker := time.NewTicker(time.Duration(service.Cfg.ClusterRefreshInterval) * time.Second)
	go func() {
		for {
			select {
			case <-ticker.C:
				kafka.InitAllClusters()
			}
		}
	}()

	//process asynchronous background processes
	go func() {
		kubernetes.ProcessAsyncCloudJobs()
	}()

	//cloud deployment watcher
	anotherTicker := time.NewTicker(time.Duration(service.Cfg.ClusterRefreshInterval) * time.Second)
	go func() {
		for {
			select {
			case <-anotherTicker.C:
				kubernetes.UpdateAllClusterStatus()
			}
		}
	}()

	//init web router
	http.InitRouter()
}
