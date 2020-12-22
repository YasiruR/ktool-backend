package kubernetes

import (
	"context"
	"fmt"
	"github.com/YasiruR/ktool-backend/database"
	"github.com/YasiruR/ktool-backend/domain"
	"github.com/YasiruR/ktool-backend/log"
	"github.com/aws/aws-sdk-go/service/eks"
	containerpb "google.golang.org/genproto/googleapis/container/v1"
	"math"
	"strconv"
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
			} else if job.Status == domain.EKS_SUBMITTED_FOR_DELETION {
				params := job.Information.(domain.EksAsyncJobParams)
				result, err := deleteNodeGroup(&params.Client, job.Reference)
				if err == nil {
					if *result.Nodegroup.Status == "DELETING" {
						log.Logger.Trace("Node group delete request successful for cluster name: ", job.Reference)
						_, err = database.UpdateEksClusterCreationStatus(context.Background(), domain.EKS_NODE_GROUP_DELETING, job.Reference)
						job.Status = domain.EKS_NODE_GROUP_DELETING
					} else {
						//todo: what are the possible states? handle them here
						//log.Logger.Trace("Node group is still being deleted for cluster name: ", status.Response.ClusterName)
						//PushToJobList(*job)
					}
				} else {
					log.Logger.Trace("Node group delete failed for cluster name: ", job.Reference)
					_, err = database.UpdateEksClusterCreationStatus(context.Background(), domain.COMPLETED, job.Reference)
					return
				}
				PushToJobList(*job)
				//} else if job.Status == domain.EKS_NODE_GROUP_DELETING {

			} else if job.Status == domain.EKS_NODE_GROUP_DELETED || job.Status == domain.EKS_NODE_GROUP_DELETING {
				params := job.Information.(domain.EksAsyncJobParams)
				result, err := deleteControlPlane(&params.Client, job.Reference, params.NodeGroupName)
				if err == nil {
					if *result.Cluster.Status == "DELETING" {
						log.Logger.Trace("Node group delete successful for cluster name: ", job.Reference)
						_, err = database.UpdateEksClusterCreationStatus(context.Background(), domain.EKS_MASTER_DELETING, job.Reference)
						job.Status = domain.EKS_MASTER_DELETING
					} else {
						//todo: handle other scenarios here
						//log.Logger.Trace("Node group is still being deleted for cluster name: ", status.Response.ClusterName)
						PushToJobList(*job)
					}
				} else {
					err = err.(*eks.ResourceInUseException)
					log.Logger.Trace("Node group is still being deleted for cluster name: ", job.Reference)
					//_, err = database.UpdateEksClusterCreationStatus(context.Background(), domain.EKS_NODE_GROUP_DELETE_FAILED, job.Reference)
					PushToJobList(*job)
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
					PushToJobList(*job)
				}
			} else if job.Status == domain.GKE_DELETING {
				opInfo := job.Information.(*containerpb.Operation)
				result, err := CheckOperationStatus(job.Reference, opInfo.GetName())
				if err != nil || result.Status == containerpb.Operation_STATUS_UNSPECIFIED {
					log.Logger.Trace("Error occurred while checking cluster delete operation status: ", opInfo.GetName())
					PushToJobList(*job)
				} else {
					if result.GetStatus() == containerpb.Operation_DONE {
						log.Logger.Trace("Cluster delete operation success: ", opInfo.GetName())
						_, err = database.UpdateGkeClusterCreationStatus(context.Background(), domain.DELETED, opInfo.GetName())
						_, err = database.UpdateGkeLROperation(context.Background(), opInfo.GetName(), "DONE")
					} else {
						// we have to check repeatedly
						PushToJobList(*job)
					}
				}
			}
		}
	case "microsoft":
		{
			if job.Status == domain.AKS_SUBMITTED {
				ctx := context.Background()
				params := job.Information.(domain.AksAsyncJobParams)
				// create the resource group if not exists
				_, err := CreateResourceGroupIfNotExist(ctx, params.ClusterOptions.ResourceGroupName, params.ClusterOptions.Zone,
					strconv.Itoa(params.ClusterOptions.SecretId))
				if err != nil {
					log.Logger.Error(fmt.Errorf("aks resource group creation failed by microsoft; %s", err))
					_, _ = database.UpdateAksClusterCreationStatus(ctx, 1, domain.FAILED, params.ClusterOptions.Name, params.ClusterOptions.ResourceGroupName)
					return
				}
				future, err := SyncCreateAksCluster(ctx, params.Client, params.ClusterOptions, params.CreateRequest)
				if err != nil {
					log.Logger.Error(fmt.Errorf("aks cluster creation failed by microsoft; %s", err))
					_, _ = database.UpdateAksClusterCreationStatus(ctx, 1, domain.FAILED, params.ClusterOptions.Name, params.ClusterOptions.ResourceGroupName)
					return
				}
				err = future.WaitForCompletionRef(ctx, params.Client.Client)
				if err != nil {
					log.Logger.Error(fmt.Errorf("aks cluster creation failed by microsoft; %s", err))
					_, _ = database.UpdateAksClusterCreationStatus(ctx, 1, domain.FAILED, params.ClusterOptions.Name, params.ClusterOptions.ResourceGroupName)
				}
				futureResolve, err := future.Result(params.Client) //we resolve the future
				if err != nil {
					log.Logger.Error(fmt.Errorf("aks cluster creation future resolve failed; %s", err))
					//database.UpdateAksClusterCreationStatus(ctx, 1,  "CREATION FAILED", params.ClusterOptions.Name, params.ClusterOptions.ResourceGroupName)
					PushToJobList(*job)
					return
				}
				if *futureResolve.ManagedClusterProperties.ProvisioningState == "Succeeded" {
					log.Logger.Info("aks cluster creation successful")
					_, _ = database.UpdateAksClusterCreationStatus(ctx, 1, domain.COMPLETED, params.ClusterOptions.Name, params.ClusterOptions.ResourceGroupName)
				} else {
					log.Logger.Error(fmt.Errorf("aks cluster creation failed; %s", err))
					_, _ = database.UpdateAksClusterCreationStatus(ctx, 1, domain.FAILED, params.ClusterOptions.Name, params.ClusterOptions.ResourceGroupName)
				}
				return
			} else if job.Status == domain.AKS_SUBMITTED_FOR_DELETION {
				ctx := context.Background()
				params := job.Information.(domain.AksAsyncJobParams)
				err := SyncDeleteAksCluster(ctx, params.Client, params.ClusterOptions.ResourceGroupName, params.ClusterOptions.Name)
				if err != nil {
					//todo: something happened
					PushToJobList(*job)
					log.Logger.Trace("aks deletion failed for cluster name: ", params.ClusterOptions.Name)
				} else {
					log.Logger.Trace("Completed aks deletion for cluster name: ", params.ClusterOptions.Name)
					_, _ = database.UpdateAksClusterCreationStatus(ctx, 1, domain.DELETED, params.ClusterOptions.Name, params.ClusterOptions.ResourceGroupName)
					return
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
		if counter == 0 { //todo: why?
			counter = 10
		}
		duration := math.Min(float64(maxWait), float64(counter*wait))
		log.Logger.Trace("No jobs to process. Sleeping", time.Duration(duration)*time.Second)
		time.Sleep(time.Duration(duration) * time.Second)
	}
}

// This method checks cluster status by pinging cloud
func UpdateAllClusterStatus() {
	log.Logger.Info("Updating ktool managed kubernetes cluster status")
	clustersToCheck := database.GetAllRunningKubernetesClusters(context.Background())
	if clustersToCheck.Error != nil {
		log.Logger.Info("Error occurred while fetching cluster information. Check the database connection")
		return
	}
	for _, cluster := range clustersToCheck.Clusters {
		if cluster.ServiceProvider == "google" {
			go checkGKEStatus(cluster)
		} else if cluster.ServiceProvider == "amazon" {
			go checkEKSStatus(cluster)
		} else {
			go checkAKSStatus(cluster)
		}
	}
	log.Logger.Info("Checks are performed on all active clusters.")
}

func checkGKEStatus(cluster domain.KubCluster) {
	isRunning := CheckGKEClusterStatus(cluster.SecretId, cluster.ClusterName, cluster.ProjectName, cluster.Location)
	if isRunning && cluster.Status != domain.COMPLETED {
		_, _ = database.UpdateClusterStatusById(context.Background(), domain.IsRunning, cluster.Id, domain.COMPLETED)
	} else if !isRunning {
		_, _ = database.UpdateClusterStatusById(context.Background(), domain.IsRunning, cluster.Id, domain.STOPPED)
	}
}

func checkEKSStatus(cluster domain.KubCluster) {
	isRunning := CheckEKSClusterStatus(cluster.ClusterName, cluster.Location, cluster.SecretId)
	if isRunning && cluster.Status != domain.COMPLETED {
		_, _ = database.UpdateClusterStatusById(context.Background(), domain.IsRunning, cluster.Id, domain.COMPLETED)
	} else if !isRunning {
		_, _ = database.UpdateClusterStatusById(context.Background(), domain.IsRunning, cluster.Id, domain.STOPPED)
	}
}

func checkAKSStatus(cluster domain.KubCluster) {
	isRunning := CheckAKSClusterStatus(cluster.ClusterName, cluster.ResourceGroup, cluster.SecretId)
	if isRunning && cluster.Status != domain.COMPLETED {
		_, _ = database.UpdateClusterStatusById(context.Background(), domain.IsRunning, cluster.Id, domain.COMPLETED)
	} else if !isRunning {
		_, _ = database.UpdateClusterStatusById(context.Background(), domain.IsRunning, cluster.Id, domain.STOPPED)
	}
}
