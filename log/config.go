package log

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
)

type Config struct {
	Level 			string	`yaml:"level"`
	RemoteLog		bool 	`yaml:"remoteLog"`
	FilePathEnabled bool	`yaml:"file_path_enabled"`
	Colors 			bool 	`yaml:"colors"`
}

var Cfg = new(Config)

func (c *Config) LoadConfigurations() *Config {
	yamlFile, err := ioutil.ReadFile("config/logger.yaml")
	if err != nil {
		log.Fatal("reading log.yaml file failed : ", err)
	}
	err = yaml.Unmarshal(yamlFile, &c)
	if err != nil {
		log.Fatal("unmarshal error (log.yaml) : ", err)
	}

	return c
}