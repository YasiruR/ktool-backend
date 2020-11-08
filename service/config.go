package service

import (
	"github.com/YasiruR/ktool-backend/log"
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

type Config struct {
	ServicePort 			string		`yaml:"service_port"`
	PingRetry 				int 		`yaml:"ping_retry"`
	PingTimeout				int			`yaml:"ping_timeout"`
	ClientInitTimeout		int			`yaml:"client_init_timeout"`
	ClusterRefreshInterval	int			`yaml:"cluster_refresh_interval"`
	MetricsUpdateInterval 	int			`yaml:"metrics_update_interval"`
	MetricsCleanInterval	int			`yaml:"metrics_table_clean_interval"`
	PromScrapeInterval		string		`yaml:"default_prom_scrape_interval"`
	PromScrapeTimeout		string		`yaml:"default_prom_scrape_timeout"`
	ConfigFilePath 			string		`yaml:"config_file_path"`
}

var Cfg = new(Config)

func (c *Config) LoadConfigurations() *Config {
	yamlFile, err := ioutil.ReadFile("config/service.yaml")
	if err != nil {
		log.Logger.Fatal("reading service.yaml file failed")
	}
	err = yaml.Unmarshal(yamlFile, &c)
	if err != nil {
		log.Logger.Fatal("unmarshal error (service.yaml) : ", err)
	}

	log.Logger.Info("service configurations initialized")
	return c
}
