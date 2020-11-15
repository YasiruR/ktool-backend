package kubernetes

import (
	"context"
	"github.com/YasiruR/ktool-backend/database"
	"github.com/YasiruR/ktool-backend/domain"
	"github.com/YasiruR/ktool-backend/log"
	containerpb "google.golang.org/genproto/googleapis/container/v1"
	"math"
	"time"
)

// watcher has a list of jobs
var jobs []domain.AsyncCloudJob

func push(job domain.AsyncCloudJob) {
	jobs = append(jobs, job)
}

func pop(jobs []domain.AsyncCloudJob) (domain.AsyncCloudJob, []domain.AsyncCloudJob) {
	return jobs[0], jobs[1:]
}

func PushToJobList(job domain.AsyncCloudJob) {
	log.Logger.Trace("Async cloud job scheduled, ", job)
	push(job)
}

func ProcessAsyncJob(job *domain.AsyncCloudJob) {
	switch job.Provider {
	case "amazon":
		{
			if job.Status == domain.EKS_MASTER_CREATING {
				status := job.Information.(domain.EksClusterContext)
				result, err := CheckEksClusterCreationStatus(status.ClusterStatus.Name, status.ClusterRequest.Location, status.SecretID)
				if err == nil {
					if *result.Cluster.Status == "ACTIVE" {
						_, err = database.UpdateEksClusterCreationStatus(context.Background(), domain.EKS_MASTER_CREATED, status.ClusterRequest.Name)
						log.Logger.Trace("Control plane successfully created for cluster name: ", status.ClusterStatus.Name)
						nodeResult, err := CreateEksNodeGroup(status.SecretID, status)
						// todo: update the job
						if err == nil {
							_, err = database.UpdateEksClusterCreationStatus(context.Background(), domain.EKS_NODE_GROUP_CREATING, status.ClusterRequest.Name)
							job.Status = domain.EKS_NODE_GROUP_CREATING
							job.Information = nodeResult
							log.Logger.Trace("Node group creation job submitted for cluster name: ", status.ClusterStatus.Name)
						} else {
							log.Logger.Trace("Error occurred in node group creation for cluster name: ", status.ClusterStatus.Name)
						}
					} else if *result.Cluster.Status == "CREATE_FAILED" {
						//todo: process create failed
						_, err = database.UpdateEksClusterCreationStatus(context.Background(), domain.EKS_MASTER_FAILED, status.ClusterRequest.Name)
						log.Logger.Trace("Error occurred in control plane creation for cluster name: ", status.ClusterRequest.Name)
						PushToJobList(*job)
					} else {
						log.Logger.Trace("Control plane is still being created for cluster: ", status.ClusterStatus.Name)
					}
				} else {
					log.Logger.Trace("Control plane creation status check failed for cluster name: ", status.ClusterStatus.Name)
				}
				PushToJobList(*job) // schedule the job again
			} else if job.Status == domain.EKS_NODE_GROUP_CREATING {
				status := job.Information.(domain.EksNodeGroupContext)
				result, err := CheckEksNodeGroupCreationStatus(*status.Response.ClusterName, *status.Response.NodegroupName, status.Region, status.SecretId)
				if err == nil {
					if *result.Nodegroup.Status == "ACTIVE" {
						_, err = database.UpdateEksClusterCreationStatus(context.Background(), domain.COMPLETED, *status.Response.ClusterName)
						//todo: process create success
						job.Status = domain.COMPLETED
						log.Logger.Trace("Node group creation successful for cluster name: ", status.Response.ClusterName)
					} else if *result.Nodegroup.Status == "CREATE_FAILED" {
						//todo: process create failed
						_, err = database.UpdateEksClusterCreationStatus(context.Background(), domain.EKS_NODE_GROUP_FAILED, *status.Response.ClusterName)
						log.Logger.Trace("Error occurred in node group creation for cluster name: ", status.Response.ClusterName)
						PushToJobList(*job)
					} else {
						log.Logger.Trace("Node group is still being created for cluster name: ", status.Response.ClusterName)
						PushToJobList(*job)
					}
				} else {
					log.Logger.Trace("Node group creation status check failed for cluster name: ", status.Response.ClusterName)
					_, err = database.UpdateEksClusterCreationStatus(context.Background(), domain.EKS_NODE_GROUP_FAILED, *status.Response.ClusterName)
					//PushToJobList(*job)
				}
			}

		}
	case "google":
		{
			if job.Status == domain.GKE_CREATING {
				status := job.Information.(*containerpb.Operation)
				result, err := CheckGkeClusterCreationStatus(job.Reference, status.Name)
				if err == nil {
					if result.Status == "DONE" {
						log.Logger.Trace("Completed gke  creation for cluster name: ", status.Name)
						return
					} else if result.Status == "FAILED" {
						//todo: process create failed
						_, err = database.UpdateGkeClusterCreationStatus(context.Background(), domain.FAILED, status.Name)
						log.Logger.Trace("Error occurred in control plane creation for cluster name: ", status.Name)
						PushToJobList(*job)
					} else if result.Status == "RUNNING" {
						PushToJobList(*job)
						log.Logger.Trace("Cluster is still being created for cluster: ", status.Name)
					} else {
						log.Logger.Trace("Unhandled status received for cluster: ", status.Name)
					}
				} else {
					log.Logger.Trace("Cluster creation status check failed for cluster name: ", status.Name)
				}
			}
		}
	}
}

func ProcessAsyncCloudJobs() {
	var job domain.AsyncCloudJob
	var wait int64 = 5
	var counter int64 = 1
	var maxWait int64 = 60 // max wait is 1 min
	for {
		//loop:
		for len(jobs) != 0 {
			job, jobs = pop(jobs)
			go func(job domain.AsyncCloudJob) {
				ProcessAsyncJob(&job)
			}(job)
			counter = 1
			//goto loop
		}
		counter = counter * 2
		duration := math.Min(float64(maxWait), float64(counter*wait))
		log.Logger.Trace("No jobs to process. Sleeping", time.Duration(duration)*time.Second)
		time.Sleep(time.Duration(duration) * time.Second)
	}
}

func UpdateAllClusterStatus() {
	log.Logger.Info("Updating ktool managed kubernetes cluster status")
	clustersToCheck := database.GetAllRunningKubernetesClusters(context.Background())
	if clustersToCheck.Error != nil {
		log.Logger.Info("Error occurred while fetching cluster information. Check the database connection")
		return
	}
	eksClusters := make(map[string][]*domain.KubCluster, 0)
	gkeClusters := make(map[string][]*domain.KubCluster, 0)
	for _, cluster := range clustersToCheck.Clusters {
		if cluster.ServiceProvider == "google" {
			if gkeClusters[cluster.ClusterId] == nil {
				gkeClusters[cluster.ClusterId] = []*domain.KubCluster{&cluster}
			} else {
				gkeClusters[cluster.ClusterId] = append(gkeClusters[cluster.ClusterId], &cluster)
			}
		} else if cluster.ServiceProvider == "amazon" {
			if eksClusters[cluster.ClusterId] == nil {
				eksClusters[cluster.ClusterId] = []*domain.KubCluster{&cluster}
			} else {
				eksClusters[cluster.ClusterId] = append(eksClusters[cluster.ClusterId], &cluster)
			}
		}
	}
	go func() {
		processGkeClusters(gkeClusters)
	}()
	go func() {
		processEksClusters(eksClusters)
	}()
}

func processGkeClusters(clusters map[string][]*domain.KubCluster) {
	//for s, kubClusters := range clusters {
	//	runningClusters, _ := ListGkeClusters(s)
	//	for cluster := range runningClusters {
	//		log.Logger.Trace("cluster is running ", cluster)
	//	}
	//check against the list and update
	//}
}

func processEksClusters(clusters map[string][]*domain.KubCluster) {

}