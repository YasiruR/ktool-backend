package prometheus

import (
	"context"
	"github.com/YasiruR/ktool-backend/log"
	"os"
	"os/exec"
)

func Init() {
	//todo install docker if it does not exist
	//check if docker prometheus is installed already
	checkCmd := exec.Command("/bin/sh", "-c", "sudo docker images | grep prom/prometheus")
	promCheck, err := checkCmd.Output()
	if err != nil {
		log.Logger.Fatal(err,"looking for pre-installations of prometheus docker images failed")
	}

	//if does not exist, install it
	if string(promCheck) == "" {
		pullCmd := exec.Command("/bin/sh", "-c", "sudo docker pull prom/prometheus")
		_, err = pullCmd.Output()
		if err != nil {
			log.Logger.Fatal(err,"pulling docker prometheus failed")
		}
		log.Logger.Trace("pulled prometheus docker image successfully")
	} else {
		log.Logger.Trace("prometheus docker image already exists in the system")

		//if image exists already, check for running containers
		checkRunContCmd := exec.Command("/bin/sh", "-c", "sudo docker ps | grep prometheus")
		runContCheck, err := checkRunContCmd.Output()
		if err != nil {
			//todo think on this workaround
			if err.Error() != "exit status 1" && string(runContCheck) != "" {
				log.Logger.Fatal(err,"looking for running docker containers of prometheus failed", string(runContCheck))
			}
		}

		//if found, stop the running container
		if string(runContCheck) != "" {
			stopDocker := exec.Command("/bin/sh", "-c", "sudo docker stop prometheus")
			stopOutput, err := stopDocker.CombinedOutput()
			if err != nil {
				log.Logger.Fatal(err,"failed to stop docker prometheus container", string(stopOutput))
			}
			log.Logger.Trace("an instance of docker prometheus was found already running and it has been terminated")
		}

		//check for stopped containers
		checkStopContCmd := exec.Command("/bin/sh", "-c", "sudo docker ps -a | grep prometheus")
		stopContCheck, err := checkStopContCmd.CombinedOutput()
		if err != nil {
			//todo think on this workaround
			if err.Error() != "exit status 1" && string(stopContCheck) != "" {
				log.Logger.Fatal(err,"looking for stopped docker containers of prometheus failed", string(stopContCheck))
			}
		}

		if string(stopContCheck) != "" {
			rmDocker := exec.Command("/bin/sh", "-c", "sudo docker rm prometheus")
			rmOutput, err := rmDocker.CombinedOutput()
			if err != nil {
				log.Logger.Error(err,"failed to remove stopped docker prometheus container", string(rmOutput))
			} else {
				log.Logger.Trace("found an existing stopped prometheus docker container and it was removed")
			}
		}
	}

	//getting file path for prometheus config
	pwd, err := os.Getwd()
	if err != nil {
		log.Logger.Fatal("failed in fetching the current working directory")
	}

	//start docker container using the config given
	runCmd := exec.Command("/bin/sh", "-c", "sudo docker run -d -p 9090:9090 --name prometheus -v " + pwd + "/config/prometheus.yml:/etc/prometheus/prometheus.yml prom/prometheus")
	runOutput, err := runCmd.CombinedOutput()
	if err != nil {
		log.Logger.Fatal(err,"could not start docker prometheus container (there might be a container already running (or stopped) as prometheus)", string(runOutput))
	}

	log.Logger.Info("prometheus docker container is up and running", string(runOutput))
}

//iterate through brokers of cluster list
//job name should be cluster name and targets should be relevant brokers (db names and prom should be identical)
//query topics, leaders, replicas, mesgs, insync stats of brokers
//store all stats in broker table
//send a halting channel as this runs in a separate go routine. do these for all such process in project
func SyncBrokerData(ctx context.Context) {
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
