package database

import (
	"github.com/YasiruR/ktool-backend/log"
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

type Config struct {
	Database 		string		`yaml:"database"`
	Host 			string 		`yaml:"host"`
	Port 			int			`yaml:"port"`
	Username 		string 		`yaml:"username"`
	Password 		string 		`yaml:"password"`
}

var Cfg = new(Config)

func (c *Config) LoadConfigurations() *Config {
	yamlFile, err := ioutil.ReadFile("config/database.yaml")
	if err != nil {
		log.Logger.Fatal("reading database.yaml file failed")
	}
	err = yaml.Unmarshal(yamlFile, &c)
	if err != nil {
		log.Logger.Fatal("unmarshal error (database.yaml) : ", err)
	}

	log.Logger.Info("database configurations initialized")
	return c
}
