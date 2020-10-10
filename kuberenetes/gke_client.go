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
	//"golang.org/x/oauth2/google"
	"google.golang.org/api/option"
	containerpb "google.golang.org/genproto/googleapis/container/v1"
)

//TODO: this is the token source pool
var JWTConfigPool = make(map[string]jwt.Config)

//func CreateGkeCluster(clusterId string, userId string, clusterOptions *domain.ClusterOptions) (*containerpb.Operation, error) {
//	ctx := context.Background()
//	b, cred, err := iam.GetGkeCredentialsForUser(userId)
//	if err != nil {
//		return nil, err
//	}
//	c, err := container.NewClusterManagerClient(ctx, option.WithCredentialsJSON(b))
//	if err != nil {
//		return nil, err
//	}
//	req, err := generateGKEClusterCreationRequest(&cred, clusterOptions)
//	if err != nil {
//		return nil, err
//	}
//	resp, err := c.CreateCluster(ctx, req)
//	if err != nil { // todo: retry and if consistently failing, send delete cluster request
//		log.Logger.ErrorContext(ctx, "Cluster creation failed with Google. Error, {}", err.Error())
//		return nil, err
//	}
//	err = database.AddGkeLROperation(ctx, resp.Name, cred.ProjectId, clusterOptions.Location)
//	if err != nil {
//		log.Logger.ErrorContext(ctx, "Failed to add LRO to db.")
//		return nil, err
//	}
//	err = database.AddGkeCluster(ctx, clusterId, clusterOptions.UserId, clusterOptions.Name, resp.Name)
//	if err != nil {
//		log.Logger.ErrorContext(ctx, "Failed to add cluster details to db.")
//		return nil, err
//	}
//	return resp, nil
//}

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
	err = database.AddGkeCluster(ctx, clusterId, clusterOptions.UserId, clusterOptions.Name, resp.Name, clusterOptions.Location)
	if err != nil {
		log.Logger.ErrorContext(ctx, "Failed to add cluster details to db.")
		return nil, err
	}
	return resp, nil
}

//func CheckOperationStatus(b []byte, opName string) {
//	ctx := context.Background()
//
//	//c, err := oauth2.DefaultClient(ctx, cloudresourcemanager.CloudPlatformScope)
//	//if err != nil {
//	//	log.Logger.Fatal(err)
//	//}
//
//	cloudresourcemanagerService, err := cloudresourcemanager.NewService(ctx, option.WithCredentialsJSON(b))
//	if err != nil {
//		log.Logger.Fatal(err)
//	}
//
//	// The name of the operation resource.
//	name := "operations/" + opName // TODO: Update placeholder value.
//
//	resp, err := cloudresourcemanagerService.Operations.Get(name).Context(ctx).Do()
//	if err != nil {
//		log.Logger.Fatal(err)
//	}
//
//	// TODO: Change code below to process the `resp` object:
//	fmt.Printf("%#v\n", resp)
//}

func CheckOperationStatus(c *container.ClusterManagerClient, zone string, name string, projectId string) (response containerpb.Operation, err error) {
	ctx := context.Background()
	opReq := &containerpb.GetOperationRequest{
		//projects/ktool-280018/locations/us-central1-a/operations/operation-1592570796611-3735784d
		Name: fmt.Sprintf("projects/%s/locations/%s/operations/%s", projectId, zone, name),
	}
	resp, err := c.GetOperation(ctx, opReq)
	if err != nil {
		// TODO: Handle error.
		return containerpb.Operation{}, err
	}
	return *resp, nil
}

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

func generateDestroyClusterRequest(credentials *domain.GkeSecret, clusterOptions *domain.ClusterOptions) (*containerpb.DeleteClusterRequest, error) {
	return &containerpb.DeleteClusterRequest{
		ProjectId: credentials.ProjectId,
		Zone:      clusterOptions.Zone,
		ClusterId: clusterOptions.ClusterId,
		Name:      clusterOptions.Name,
	}, nil
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

func UpdateGkePendingOperations(userId string) {
	//query = ""
}

func CheckGkeClusterCreationStatus(userId string, operationName string) (status domain.GkeOperationStatusCheck, err error) {
	ctx := context.Background()
	b, _, err := iam.GetGkeCredentialsForUser(userId)
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
		log.Logger.Error("Error occured while fetching operation id {}", operationName)
	}
	if operation.Status != "DONE" {
		retriesLeft := 3
		resp, err := CheckOperationStatus(c, operation.Zone, operation.Name, operation.ProjectId)
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

//func GetGkeCredentialsForUser(userId string, cred *google.Credentials) {
//	ctx := context.Background()
//	secrets := database.GetAllSecretsByUserInternal(ctx, userId, `Google`)
//	jsonBytes, err := json.Marshal(&secrets.SecretList[0])
//	if err != nil {
//		log.Logger.ErrorContext(ctx, "Could not create gke credentials for user %s", userId)
//		return
//	}
//	credentials := google.Credentials{
//		ProjectID: secrets.SecretList[0].ProjectId,
//
//	}
//	option.WithCredentialsJSON()
//}
