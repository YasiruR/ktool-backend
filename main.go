package main

import (
	"github.com/YasiruR/ktool-backend/database"
	"github.com/YasiruR/ktool-backend/http"
	"github.com/YasiruR/ktool-backend/kafka"
	"github.com/YasiruR/ktool-backend/log"
	"github.com/YasiruR/ktool-backend/prometheus"
	"github.com/YasiruR/ktool-backend/service"
	"github.com/google/uuid"
	traceable_context "github.com/pickme-go/traceable-context"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	log.Cfg.LoadConfigurations()
	log.Init()

	database.Cfg.LoadConfigurations()
	database.Init()

	service.Cfg.LoadConfigurations()

	kafka.InitAllClusters()
	prometheus.Init()

	//for termination of processes
	osChannel := make(chan os.Signal, 1)
	signal.Notify(osChannel, syscall.SIGINT, syscall.SIGKILL)

	//refresh cluster data
	ticker := time.NewTicker(time.Duration(service.Cfg.ClusterRefreshInterval) * time.Second)
	go func() {
		for {
			select {
				case <- ticker.C:
					kafka.InitAllClusters()
			}
		}
	}()

	//scrape prometheus metrics from jmx
	metricsTicker := time.NewTicker(time.Duration(service.Cfg.MetricsUpdateInterval) * time.Second)
	syncContext := traceable_context.WithUUID(uuid.New())
	go func() {
		for {
			select {
			case <- metricsTicker.C:
				prometheus.SyncBrokerMetrics(syncContext)
			}
		}
	}()

	//update topic metrics
	topicContext := traceable_context.WithUUID(uuid.New())
	go func() {
		for {
			select {
			case <- metricsTicker.C:
				prometheus.InitTopicMetrics(topicContext, osChannel)
			}
		}
	}()

	//update topic summary metrics
	topicSummaryContext := traceable_context.WithUUID(uuid.New())
	go func() {
		for {
			select {
			case <- metricsTicker.C:
				prometheus.InitSummaryMetrics(topicSummaryContext, osChannel)
			}
		}
	}()

	//update metrics ports of brokers
	metricsPortContext := traceable_context.WithUUID(uuid.New())
	go func() {
		for {
			select {
			case <- metricsTicker.C:
				prometheus.InitBrokerMetricsPorts(metricsPortContext)
			}
		}
	}()

	//run metrics clean job
	cleanTicker := time.NewTicker(time.Duration(service.Cfg.MetricsCleanInterval) * time.Second)
	cleanContext := traceable_context.WithUUID(uuid.New())
	go func() {
		for {
			select {
			case <- cleanTicker.C:
				database.CleanMetricsTable(cleanContext)
			}
		}
	}()

	http.InitRouter(osChannel)
}
