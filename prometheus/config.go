package prometheus

import (
	"context"
	"errors"
	"github.com/YasiruR/ktool-backend/domain"
	"github.com/YasiruR/ktool-backend/log"
	"github.com/YasiruR/ktool-backend/service"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"strconv"
)

type config struct{
	Global 			struct{
		ScrapeInterval		string			`yaml:"scrape_interval"`
		EvaluationInterval	string			`yaml:"evaluation_interval"`
	}	`yaml:"global"`
	ScrapeConfigs 	[]jobConfig				`yaml:"scrape_configs"`
	RemoteWrite 	[]map[string]string		`yaml:"remote_write"`
}

type jobConfig struct{
	JobName 		string					`yaml:"job_name"`
	ScrapeInterval	string					`yaml:"scrape_interval"`
	ScrapeTimeout	string					`yaml:"scrape_timeout"`
	StaticConfigs	[]map[string][]string	`yaml:"static_configs"`
}

func loadConfigurations() config {
	cfg := config{}
	//todo check if readfile is thread-safe
	yamlFile, err := ioutil.ReadFile("config/prometheus.yml")
	if err != nil {
		log.Logger.Fatal("reading prometheus.yml file failed")
	}
	err = yaml.Unmarshal(yamlFile, &cfg)
	if err != nil {
		log.Logger.Fatal("unmarshal error (prometheus.yml) : ", err)
	}

	log.Logger.Info("prometheus configurations initialized")
	return cfg
}

func AddNewJob(ctx context.Context, clusterName string, brokers []domain.Server) (err error) {
	cfg := loadConfigurations()

	var job jobConfig
	job.JobName = clusterName
	job.ScrapeInterval = service.Cfg.PromScrapeInterval
	job.ScrapeTimeout = service.Cfg.PromScrapeTimeout

	tmpMap := make(map[string][]string)
	for _, b := range brokers {
		tmpMap["targets"] = append(tmpMap["targets"], b.Host + ":" + strconv.Itoa(b.MetricsPort))
	}
	job.StaticConfigs = append(job.StaticConfigs, tmpMap)
	cfg.ScrapeConfigs = append(cfg.ScrapeConfigs, job)

	out, err := yaml.Marshal(&cfg)
	if err != nil {
		log.Logger.ErrorContext(ctx, err, "could not marshall new prom config to yaml", clusterName)
		return
	}

	//todo : check if writefile is thread safe. if not, use a lock here and in deletion. same for readfile
	err = ioutil.WriteFile("config/prometheus.yml", out, 0644)
	if err != nil {
		log.Logger.ErrorContext(ctx, err, "failed writing new prom configs to yaml file")
		return
	}

	log.Logger.TraceContext(ctx, "prometheus configs with new job updated successfully", clusterName)
	return
}

func DeleteJob(ctx context.Context, clusterName string) (err error) {
	cfg := loadConfigurations()

	newScrapeConfigs := cfg.ScrapeConfigs
	for index, job := range cfg.ScrapeConfigs {
		if job.JobName == clusterName {
			cfg.ScrapeConfigs[index] = cfg.ScrapeConfigs[len(newScrapeConfigs)-1]
			cfg.ScrapeConfigs[len(cfg.ScrapeConfigs)-1] = jobConfig{}
			cfg.ScrapeConfigs = cfg.ScrapeConfigs[:len(cfg.ScrapeConfigs)-1]

			out, err := yaml.Marshal(&cfg)
			if err != nil {
				log.Logger.ErrorContext(ctx, err, "could not marshall new prom config to yaml")
				return err
			}

			err = ioutil.WriteFile("config/prometheus.yml", out, 0644)
			if err != nil {
				log.Logger.ErrorContext(ctx, err, "failed writing new prom configs to yaml file")
				return err
			}

			log.Logger.TraceContext(ctx, "job removed from prom configs successfully", clusterName)
			return nil
		}
	}
	log.Logger.ErrorContext(ctx, "could not find the requested in prom configs", clusterName)
	return errors.New("failed to delete prometheus job")
}
