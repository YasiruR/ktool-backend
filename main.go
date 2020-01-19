package main

import (
	"github.com/YasiruR/ktool-backend/database"
	"github.com/YasiruR/ktool-backend/http"
	"github.com/YasiruR/ktool-backend/log"
	"github.com/YasiruR/ktool-backend/service"
)

func main() {
	log.Cfg.LoadConfigurations()
	log.Init()

	service.Cfg.LoadConfigurations()
	http.InitRouter()

	database.Cfg.LoadConfigurations()
	database.Init()
}
