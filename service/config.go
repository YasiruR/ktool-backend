package service

import (
	"github.com/YasiruR/ktool-backend/log"
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

type Config struct {
	ServicePort 	string	`yaml:"service_port"`
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

	log.Logger.Trace("service configurations initialized")
	return c
}
