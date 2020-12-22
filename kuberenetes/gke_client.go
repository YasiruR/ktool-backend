package kubernetes

import (
	container "cloud.google.com/go/container/apiv1"
	"context"
	"fmt"
	"github.com/YasiruR/ktool-backend/database"
	domain "github.com/YasiruR/ktool-backend/domain"
	iam "github.com/YasiruR/ktool-backend/iam"
	"github.com/YasiruR/ktool-backend/log"
	oauth2 "golang.org/x/oauth2/google"
	"golang.org/x/oauth2/jwt"
	resource "google.golang.org/api/cloudresourcemanager/v1"
	"strconv"

	//"golang.org/x/oauth2/google"
	"google.golang.org/api/option"
	containerpb "google.golang.org/genproto/googleapis/container/v1"
)

//TODO: this is the token source pool
var JWTConfigPool = make(map[string]jwt.Config)

func CheckGKEClusterStatus(secretId string, cluster_name string, project_name string, zone string) (isRunning bool) {
	ctx := context.Background()
	b, _, err := iam.GetGkeCredentialsForSecret(secretId)
	if err != nil {
		//return nil, err
	}
	c, err := container.NewClusterManagerClient(ctx, option.WithCredentialsJSON(b))
	if err != nil {
		//return nil, err
	}
	resp, err := c.GetCluster(ctx, &containerpb.GetClusterRequest{
		Name: fmt.Sprintf("projects/%s/locations/%s/clusters/%s", project_name, zone, cluster_name),
	})
	if err != nil {
		log.Logger.Info("Cluster with name, " + cluster_name + " not found in GCP. Updating status to stopped")
		return false
	}
	if resp.Status != containerpb.Cluster_ERROR && resp.Status != containerpb.Cluster_STATUS_UNSPECIFIED && resp.Status != containerpb.Cluster_STOPPING {
		//todo: process cluster not running
		return true
	}
	return false
	//log.Logger.Info(resp)
}

func CreateGkeCluster(clusterId string, secretId string, clusterOptions *domain.ClusterOptions) (*containerpb.Operation, error) {
	ctx := context.Background()
	b, cred, err := iam.GetGkeCredentialsForSecret(secretId)
	if err != nil {
		return nil, err
	}
	c, err := container.NewClusterManagerClient(ctx, option.WithCredentialsJSON(b))
	if err != nil {
		return nil, err
	}
	req, err := generateGKEClusterCreationRequest(&cred, clusterOptions)
	if err != nil {
		return nil, err
	}
	resp, err := c.CreateCluster(ctx, req)
	if err != nil { // todo: retry and if consistently failing, send delete cluster request
		log.Logger.ErrorContext(ctx, "Cluster creation failed with Google. Error, {}", err.Error())
		return nil, err
	}
	err = database.AddGkeLROperation(ctx, resp.Name, cred.ProjectId, clusterOptions.Location)
	if err != nil {
		log.Logger.ErrorContext(ctx, "Failed to add LRO to db.")
		return nil, err
	}
	err = database.AddGkeCluster(ctx, clusterId, clusterOptions.UserId, clusterOptions.Name, resp.Name, clusterOptions.Location, clusterOptions.SecretId)
	if err != nil {
		log.Logger.ErrorContext(ctx, "Failed to add cluster details to db.")
		return nil, err
	}
	return resp, nil
}

func CheckOperationStatus(secretId, operationName string) (response containerpb.Operation, err error) {
	ctx := context.Background()
	b, _, err := iam.GetGkeCredentialsForSecret(secretId)
	if err != nil {
		return containerpb.Operation{Status: containerpb.Operation_STATUS_UNSPECIFIED}, err
	}
	c, err := container.NewClusterManagerClient(ctx, option.WithCredentialsJSON(b))
	if err != nil {
		return containerpb.Operation{Status: containerpb.Operation_STATUS_UNSPECIFIED}, err
	}
	operation := database.GetGkeLROperation(ctx, operationName)
	if operation.Error != nil {
		log.Logger.Error("Error occurred while fetching operation id {}", operationName)
		return containerpb.Operation{Status: containerpb.Operation_STATUS_UNSPECIFIED}, err
	}
	return checkOperationStatus(c, operation.Zone, operation.Name, operation.ProjectId)
}

func checkOperationStatus(c *container.ClusterManagerClient, zone string, name string, projectId string) (response containerpb.Operation, err error) {
	ctx := context.Background()
	opReq := &containerpb.GetOperationRequest{
		//projects/ktool-280018/locations/us-central1-a/operations/operation-1592570796611-3735784d
		Name: fmt.Sprintf("projects/%s/locations/%s/operations/%s", projectId, zone, name),
	}
	resp, err := c.GetOperation(ctx, opReq)
	if err != nil {
		// TODO: Handle error.
		return containerpb.Operation{Status: containerpb.Operation_STATUS_UNSPECIFIED}, err
	}
	return *resp, nil
}

func ListGkeClusters(userId string) (*containerpb.ListClustersResponse, error) {
	ctx := context.Background()
	b, cred, err := iam.GetGkeCredentialsForUser(userId)
	if err != nil {
		return nil, err
	}
	conf, err := oauth2.JWTConfigFromJSON(b, resource.CloudPlatformScope)
	if err != nil {
		return nil, err
	}
	c, err := container.NewClusterManagerClient(ctx, option.WithTokenSource(conf.TokenSource(ctx)))
	if err != nil {
		return nil, err
	}
	req := &containerpb.ListClustersRequest{
		Parent: `projects/` + cred.ProjectId + `/locations/-`,
	}
	resp, err := c.ListClusters(ctx, req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func CheckGkeClusterCreationStatus(secretId string, operationName string) (status domain.GkeOperationStatusCheck, err error) {
	ctx := context.Background()
	b, _, err := iam.GetGkeCredentialsForSecret(secretId)
	if err != nil {
		return domain.GkeOperationStatusCheck{
			OperationName: "",
			Status:        "",
			Error:         err,
		}, err
	}
	c, err := container.NewClusterManagerClient(ctx, option.WithCredentialsJSON(b))
	if err != nil {
		return domain.GkeOperationStatusCheck{
			OperationName: "",
			Status:        "",
			Error:         err,
		}, err
	}
	operation := database.GetGkeLROperation(ctx, operationName)
	if operation.Error != nil {
		log.Logger.Error("Error occurred while fetching operation id {}", operationName)
	}
	if operation.Status != "DONE" {
		retriesLeft := 3
		resp, err := checkOperationStatus(c, operation.Zone, operation.Name, operation.ProjectId)
	updateFailed: // if update failed, we must retry
		operation.Status = resp.GetStatus().String()
		updateStatus, err := database.UpdateGkeLROperation(ctx, resp.Name, resp.GetStatus().String())
		updateStatus, err = database.UpdateGkeClusterCreationStatus(ctx, resp.GetStatus().String(), operationName)
		if !updateStatus {
			retriesLeft--
			if retriesLeft > 0 {
				goto updateFailed
			} else {
				log.Logger.Warn("Failed to update db on operation status change. Retry count exceeded.")
			}
		}
		if err != nil {
			return domain.GkeOperationStatusCheck{
				OperationName: resp.GetName(),
				Status:        operation.Status,
				Detail:        resp.GetDetail(),
				Error:         err,
			}, err
		}
	}

	return domain.GkeOperationStatusCheck{
		OperationName: operation.Name,
		Status:        operation.Status,
		Error:         err,
	}, nil
}

func DeleteGkeCluster(secretId string, clusterId int, clusterName string, zone string) (bool, error) {
	ctx := context.Background()
	b, cred, err := iam.GetGkeCredentialsForSecret(secretId)
	if err != nil {
		return false, err
	}
	c, err := container.NewClusterManagerClient(ctx, option.WithCredentialsJSON(b))
	if err != nil {
		return false, err
	}
	resp, err := c.DeleteCluster(ctx, generateDestroyClusterRequest(&cred, clusterName, clusterName, zone))
	if err != nil {
		log.Logger.Info("Cluster with name, " + clusterName + " not found in GCP. Updating status to stopped")
		_, err = database.UpdateClusterStatusById(ctx, 1, clusterId, domain.STOPPED)
		return false, err
	}
	_, err = database.UpdateClusterStatusById(ctx, 1, clusterId, domain.GKE_DELETING)
	_, err = database.UpdateGkeClusterMetaById(ctx, strconv.Itoa(clusterId), resp.GetName()) // updating the associated operation id
	if err != nil {
		log.Logger.Info("Failed to update cluster status in db. Cluster delete request sent.")
		return true, err
	}
	//if resp.Status != containerpb.Cluster_ERROR && resp.Status != containerpb.Cluster_STATUS_UNSPECIFIED && resp.Status != containerpb.Cluster_STOPPING {
	//	//todo: process cluster not running
	//	return true
	//}
	err = database.AddGkeLROperation(ctx, resp.GetName(), cred.ProjectId, zone)
	if err != nil {
		log.Logger.Info("Failed to add LRO to db. Cluster delete request sent.")
		return true, err
	}
	PushToJobList(domain.AsyncCloudJob{
		Provider:    "google",
		Status:      domain.GKE_DELETING,
		Reference:   secretId,
		Information: resp,
	})
	return true, nil
}

// helpers

func generateGKEClusterCreationRequest(credentials *domain.GkeSecret, clusterOptions *domain.ClusterOptions) (*containerpb.CreateClusterRequest, error) {
	nodePool1 := containerpb.NodePool{
		Name: clusterOptions.Name + "-pool-1",
		Config: &containerpb.NodeConfig{
			MachineType: clusterOptions.MachineType,
			//DiskSizeGb:             0,
			//OauthScopes:            nil,
			//ServiceAccount:         "",
			//Metadata:               nil,
			//ImageType:              "",
			//Labels:                 nil,
			//LocalSsdCount:          0,
			//Tags:                   nil,
			//Preemptible:            false,
		},
		InitialNodeCount: clusterOptions.InstanceCount,
	}
	return &containerpb.CreateClusterRequest{
		ProjectId: credentials.ProjectId,
		Cluster: &containerpb.Cluster{
			Name:        clusterOptions.Name,
			Description: clusterOptions.Description,
			Location:    clusterOptions.Location,
			NodePools: []*containerpb.NodePool{
				&nodePool1,
			},
			InitialClusterVersion: clusterOptions.KubVersion,
			//MasterAuth: containerpb.MasterAuth{
			//
			//}
		},
		Parent: `projects/` + credentials.ProjectId + `/locations/` + clusterOptions.Location,
	}, nil
}

func generateDestroyClusterRequest(credentials *domain.GkeSecret, name, clusterId, zone string) *containerpb.DeleteClusterRequest {
	return &containerpb.DeleteClusterRequest{
		ProjectId: credentials.ProjectId,
		Zone:      zone,
		ClusterId: clusterId,
		Name:      name,
	}
}
