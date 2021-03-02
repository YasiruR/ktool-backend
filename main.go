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

	service.Cfg.LoadConfigurations()

	kafka.InitAllClusters()
	//prometheus.Init()

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

	////scrape prometheus metrics from jmx
	//metricsTicker := time.NewTicker(time.Duration(service.Cfg.MetricsUpdateInterval) * time.Second)
	//syncContext := traceable_context.WithUUID(uuid.New())
	//go func() {
	//	for {
	//		select {
	//		case <- metricsTicker.C:
	//			prometheus.SyncBrokerMetrics(syncContext)
	//		}
	//	}
	//}()
	//
	////to update metrics ports of brokers
	//metricsPortContext := traceable_context.WithUUID(uuid.New())
	//go func() {
	//	for {
	//		select {
	//		case <- metricsTicker.C:
	//			prometheus.InitBrokerMetricsPorts(metricsPortContext)
	//		}
	//	}
	//}()
	//
	////run metrics clean job
	//cleanTicker := time.NewTicker(time.Duration(service.Cfg.MetricsCleanInterval) * time.Second)
	//cleanContext := traceable_context.WithUUID(uuid.New())
	//go func() {
	//	for {
	//		select {
	//		case <- cleanTicker.C:
	//			database.CleanMetricsTable(cleanContext)
	//		}
	//	}
	//}()

	//init web router
	http.InitRouter()
}
