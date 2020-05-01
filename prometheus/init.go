package prometheus

import (
	"context"
	"fmt"
	"github.com/YasiruR/ktool-backend/log"
	"os"
	"os/exec"
)

func Init() {
	//check if docker prometheus is installed already
	checkCmd := exec.Command("/bin/sh", "-c", "sudo docker images | grep prom/prometheus")
	promCheck, err := checkCmd.Output()
	if err != nil {
		log.Logger.Fatal(err,"looking for pre-installations of prometheus docker images failed")
	}
	log.Logger.Trace("checking for existing prometheus images", string(promCheck))

	//if does not exist, install it
	if string(promCheck) == "" {
		pullCmd := exec.Command("/bin/sh", "-c", "sudo docker pull prom/prometheus")
		_, err = pullCmd.Output()
		if err != nil {
			log.Logger.Fatal(err,"pulling docker prometheus failed")
		}
		log.Logger.Trace("pulled prometheus docker image successfully")
	}

	//getting file path for prometheus config
	pwd, err := os.Getwd()
	if err != nil {
		log.Logger.Fatal("failed in fetching the current working directory")
	}

	//start docker container using the config given
	runCmd := exec.Command("/bin/sh", "-c", "sudo docker run -p 9090:9090 --name prometheus -v " + pwd + "/config/prometheus.yml:/etc/prometheus/prometheus.yml prom/prometheus")
	runOutput, err := runCmd.CombinedOutput()
	if err != nil {
		log.Logger.Fatal(err,"could not start docker prometheus container (there might be a container already running (or stopped) as prometheus)", string(runOutput))
	}

	log.Logger.Trace("prometheus docker container is up and running")

	//todo should close docker container on termination
}

//iterate through brokers of cluster list
//job name should be cluster name and targets should be relevant brokers (db names and prom should be identical)
//query topics, leaders, replicas, mesgs, insync stats of brokers
//store all stats in broker table
//send a halting channel as this runs in a separate go routine. do these for all such process in project
func SyncBrokerData(ctx context.Context) {
	fmt.Println("broker metrics sync")
	err := setBrokerBytesIn(ctx)
	if err != nil {
		log.Logger.ErrorContext(ctx, "setting broker bytes in failed")
	}

	err = setBrokerBytesOut(ctx)
	if err != nil {
		log.Logger.ErrorContext(ctx, "setting broker bytes out failed")
	}

	log.Logger.TraceContext(ctx, "updated metrics")
}
