package main

import (
	"github.com/YasiruR/ktool-backend/cloud"
	"github.com/YasiruR/ktool-backend/database"
	"github.com/YasiruR/ktool-backend/http"
	"github.com/YasiruR/ktool-backend/kafka"
	"github.com/YasiruR/ktool-backend/log"
	"github.com/YasiruR/ktool-backend/service"
)

func main() {
	log.Cfg.LoadConfigurations()
	log.Init()

	database.Cfg.LoadConfigurations()
	database.Init()

	cloud.Init()

	service.Cfg.LoadConfigurations()
	kafka.InitAllClusters()
	http.InitRouter()
}
