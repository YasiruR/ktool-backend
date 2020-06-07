package prometheus

import (
	"context"
	"github.com/YasiruR/ktool-backend/database"
	"github.com/YasiruR/ktool-backend/kafka"
	"github.com/YasiruR/ktool-backend/log"
	"os"
	"os/exec"
	"strconv"
	"time"
)

const (
	promUrl 			= "http://localhost:9090/api/v1/"
	partitions 			= "partitions"
	leaders				= "leaders"
	activeControllers	= "active_controllers"
	offlinePartitions	= "offline_partitions"
	underReplicated 	= "under_replicated"
	messageRate 		= "message_rate"
	isrExpansionRate	= "isr_expansion_rate"
	isrShrinkRate		= "isr_shrink_rate"
	networkProcIdlePerc = "network_processor_avg_idle_percent"
	responseTime 		= "response_time"
	queueTime			= "queue_time"
	remoteTime 			= "remote_time"
	localTime			= "local_time"
	totalTime 			= "total_time"
	maxLagBtwLeadAndRep = "max_message_lag_between_leader_and_replica"
	uncleanLeadElec		= "unclean_leader_election"
	failedFetchRate		= "failed_fetch_req_rate"
	failedProdRate		= "failed_prod_req_rate"
	bytesIn				= "bytes_in"
	bytesOut			= "bytes_out"
	totalMessages 		= "total_messages"
	totalTopics 		= "total_topics"
)

var (
	queryList = map[string]string{
		partitions: "query?query=kafka_server_replicamanager_partitioncount&time=",
		leaders: "query?query=kafka_server_replicamanager_leadercount&time=",
		activeControllers: "query?query=kafka_controller_kafkacontroller_activecontrollercount&time=",
		offlinePartitions: "query?query=kafka_controller_kafkacontroller_offlinepartitionscount&time=",
		underReplicated: "query?query=kafka_server_replicamanager_underreplicatedpartitions&time=",
		messageRate: "query?query=sum%20by%20(instance)%20(rate(kafka_server_brokertopicmetrics_messagesin_total%5B1m%5D))&time=",
		isrExpansionRate: "query?query=kafka_server_replicamanager_isrexpands_total&time=",
		isrShrinkRate: "query?query=kafka_server_replicamanager_isrshrinks_total&time=",
		networkProcIdlePerc: "query?query=kafka_network_socketserver_networkprocessoravgidlepercent&time=",
		responseTime: "query?query=sum%20by%20(instance)%20(rate(kafka_network_requestmetrics_responsesendtimems%5B1m%5D))&time=",
		queueTime: "query?query=sum%20by%20(instance)%20(rate(kafka_network_requestmetrics_requestqueuetimems%5B1m%5D))&time=",
		remoteTime: "query?query=sum%20by%20(instance)%20(rate(kafka_network_requestmetrics_remotetimems%5B1m%5D))&time=",
		localTime: "query?query=sum%20by%20(instance)%20(rate(kafka_network_requestmetrics_localtimems%5B1m%5D))&time=",
		totalTime: "query?query=sum%20by%20(instance)%20(rate(kafka_network_requestmetrics_totaltimems%5B1m%5D))&time=",
		maxLagBtwLeadAndRep: "query?query=kafka_server_replicafetchermanager_minfetchrate&time=",
		uncleanLeadElec: "query?query=kafka_controller_controllerstats_uncleanleaderelectionspersec&time=",
		failedFetchRate: "query?query=sum%20by%20(instance)%20(rate(kafka_server_brokertopicmetrics_failedfetchrequests_total%5B1m%5D))&time=",
		failedProdRate: "query?query=sum%20by%20(instance)%20(rate(kafka_server_brokertopicmetrics_failedproducerequests_total%5B1m%5D))&time=",
		bytesIn: "query?query=sum%20by%20(instance)%20(rate(kafka_server_brokertopicmetrics_bytesin_total%5B1m%5D))&time=",
		bytesOut: "query?query=sum%20by%20(instance)%20(rate(kafka_server_brokertopicmetrics_bytesout_total%5B1m%5D))&time=",
		totalMessages: "query?query=sum%20by%20(instance)%20(kafka_server_brokertopicmetrics_messagesin_total)&time=",
		totalTopics: "query?query=count%20by%20(instance)%20(kafka_server_brokertopicmetrics_messagesin_total)&time=",
	}
)

func Init() {
	//docker will anyway be available since the tool will be deployed as a docker container
	//check if docker prometheus is installed already
	checkCmd := exec.Command("/bin/sh", "-c", "sudo docker images | grep prom/prometheus")
	promCheck, err := checkCmd.Output()
	if err != nil {
		if string(promCheck) != "" {
			log.Logger.Fatal(err,"looking for pre-installations of prometheus docker images failed")
		}
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

	log.Logger.Info("prometheus docker container is up and running")
}

//iterate through brokers of cluster list
//job name should be cluster name and targets should be relevant brokers (db names and prom should be identical)
//query topics, leaders, replicas, mesgs, insync stats of brokers
//store all stats in broker table
//send a halting channel as this runs in a separate go routine. do these for all such process in project
func SyncBrokerMetrics(ctx context.Context) {
	currentTime := time.Now()
	ts := int(currentTime.Unix())

	err := initDbRows(ctx, ts)
	if err != nil {
		log.Logger.ErrorContext(ctx, "broker metrics update failed")
		return
	}

	for key, query := range queryList {
		req := promUrl + query + strconv.Itoa(ts)
		switch key {
		case partitions:
			go func() {
				err := setIntMetrics(ctx, ts, req, database.UpdateBrokerPartitionCount)
				if err != nil {
					log.Logger.ErrorContext(ctx, "broker partition metrics failed")
				}
			}()
		case leaders:
			go func() {
				err := setIntMetrics(ctx, ts, req, database.UpdateBrokerLeaderCount)
				if err != nil {
					log.Logger.ErrorContext(ctx, "broker leader metrics failed")
				}
			}()
		case activeControllers:
			go func() {
				err := setIntMetrics(ctx, ts, req, database.UpdateBrokerActControllerCount)
				if err != nil {
					log.Logger.ErrorContext(ctx, "broker active controller metrics failed")
				}
			}()
		case  offlinePartitions:
			go func() {
				err := setIntMetrics(ctx, ts, req, database.UpdateBrokerOfflinePartCount)
				if err != nil {
					log.Logger.ErrorContext(ctx, "broker active controller metrics failed")
				}
			}()
		case underReplicated:
			go func() {
				err := setIntMetrics(ctx, ts, req, database.UpdateBrokerUnderReplicatedCount)
				if err != nil {
					log.Logger.ErrorContext(ctx, "broker active controller metrics failed")
				}
			}()
		case messageRate:
			go func() {
				err := setFloatMetrics(ctx, ts, req, database.UpdateBrokerMessageRate)
				if err != nil {
					log.Logger.ErrorContext(ctx, "broker active controller metrics failed")
				}
			}()
		case isrExpansionRate:
			go func() {
				err := setFloatMetrics(ctx, ts, req, database.UpdateBrokerIsrExpRate)
				if err != nil {
					log.Logger.ErrorContext(ctx, "broker bytes in metrics failed")
				}
			}()
		case isrShrinkRate:
			go func() {
				err := setFloatMetrics(ctx, ts, req, database.UpdateBrokerIsrShrinkRate)
				if err != nil {
					log.Logger.ErrorContext(ctx, "broker bytes in metrics failed")
				}
			}()
		case networkProcIdlePerc:
			go func() {
				err := setFloatMetrics(ctx, ts, req, database.UpdateBrokerNetworkIdlePercentage)
				if err != nil {
					log.Logger.ErrorContext(ctx, "broker bytes in metrics failed")
				}
			}()
		case responseTime:
			go func() {
				err := setFloatMetrics(ctx, ts, req, database.UpdateBrokerResponseTime)
				if err != nil {
					log.Logger.ErrorContext(ctx, "broker bytes in metrics failed")
				}
			}()
		case queueTime:
			go func() {
				err := setFloatMetrics(ctx, ts, req, database.UpdateBrokerQueueTime)
				if err != nil {
					log.Logger.ErrorContext(ctx, "broker bytes in metrics failed")
				}
			}()
		case remoteTime:
			go func() {
				err := setFloatMetrics(ctx, ts, req, database.UpdateBrokerRemoteTime)
				if err != nil {
					log.Logger.ErrorContext(ctx, "broker bytes in metrics failed")
				}
			}()
		case localTime:
			go func() {
				err := setFloatMetrics(ctx, ts, req, database.UpdateBrokerLocalTime)
				if err != nil {
					log.Logger.ErrorContext(ctx, "broker bytes in metrics failed")
				}
			}()
		case totalTime:
			go func() {
				err := setFloatMetrics(ctx, ts, req, database.UpdateBrokerTotalTime)
				if err != nil {
					log.Logger.ErrorContext(ctx, "broker bytes in metrics failed")
				}
			}()
		case maxLagBtwLeadAndRep:
			go func() {
				err := setFloatMetrics(ctx, ts, req, database.UpdateBrokerMaxLag)
				if err != nil {
					log.Logger.ErrorContext(ctx, "broker bytes in metrics failed")
				}
			}()
		case uncleanLeadElec:
			go func() {
				err := setFloatMetrics(ctx, ts, req, database.UpdateBrokerUncleanLeaderElection)
				if err != nil {
					log.Logger.ErrorContext(ctx, "broker bytes in metrics failed")
				}
			}()
		case failedFetchRate:
			go func() {
				err := setFloatMetrics(ctx, ts, req, database.UpdateBrokerFailedFetchRate)
				if err != nil {
					log.Logger.ErrorContext(ctx, "broker bytes in metrics failed")
				}
			}()
		case failedProdRate:
			go func() {
				err := setFloatMetrics(ctx, ts, req, database.UpdateBrokerFailedProdRate)
				if err != nil {
					log.Logger.ErrorContext(ctx, "broker bytes in metrics failed")
				}
			}()
		case bytesIn:
			go func() {
				err := setFloatMetrics(ctx, ts, req, database.UpdateBrokerByteInRate)
				if err != nil {
					log.Logger.ErrorContext(ctx, "broker bytes in metrics failed")
				}
			}()
		case bytesOut:
			go func() {
				err := setFloatMetrics(ctx, ts, req, database.UpdateBrokerByteOutRate)
				if err != nil {
					log.Logger.ErrorContext(ctx, "broker bytes out metrics failed")
				}
			}()
		case totalMessages:
			go func() {
				err := setIntMetrics(ctx, ts, req, database.UpdateBrokerTotalMessages)
				if err != nil {
					log.Logger.ErrorContext(ctx, "broker total message metrics failed")
				}
			}()
		case totalTopics:
			go func() {
				err := setIntMetrics(ctx, ts, req, database.UpdateBrokerTotalTopics)
				if err != nil {
					log.Logger.ErrorContext(ctx, "broker total topics metrics failed")
				}
			}()
		}
	}
	log.Logger.TraceContext(ctx, "updated metrics")
}

func initDbRows(ctx context.Context, ts int) (err error) {
	for _, cluster := range kafka.ClusterList {
		if cluster.Available == true {
			//get all brokers for the cluster
			brokers, err := database.GetBrokersByClusterId(ctx, cluster.ClusterID)
			if err != nil {
				log.Logger.ErrorContext(ctx, "getting brokers for cluster failed", cluster.ClusterID)
				return err
			}
			for _, broker := range brokers {
				err = database.AddMetricsRow(ctx, broker.Host, ts)
				if err != nil {
					log.Logger.ErrorContext(ctx, "broker metrics update failed")
					return err
				}
			}
		}
	}
	return
}