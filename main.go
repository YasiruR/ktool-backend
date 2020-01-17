package main

import (
	"github.com/YasiruR/ktool-backend/http"
	"github.com/YasiruR/ktool-backend/log"
	"github.com/YasiruR/ktool-backend/service"
)

func main() {
	log.InitConfig()
	log.Init()
	service.Cfg.LoadConfigurations()
	http.InitRouter()
}
