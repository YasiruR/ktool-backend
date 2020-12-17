package kubernetes

import (
	"context"
	"fmt"
	"github.com/Azure/azure-sdk-for-go/profiles/latest/containerservice/mgmt/containerservice"
	"github.com/Azure/azure-sdk-for-go/services/resources/mgmt/2020-06-01/resources"
	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
	auth "github.com/Azure/go-autorest/autorest/azure/auth"
	"github.com/Azure/go-autorest/autorest/to"
	"github.com/YasiruR/ktool-backend/database"
	"github.com/YasiruR/ktool-backend/domain"
	"github.com/YasiruR/ktool-backend/iam"
	"github.com/YasiruR/ktool-backend/log"
	"github.com/YasiruR/ktool-backend/util"
	"strconv"
	"strings"
	"time"
)

func CheckAKSClusterStatus(clusterName string, resourceGroupName string, secretId string) bool {
	cred, err := iam.GetAksCredentialsForSecretId(secretId)
	aksClient := containerservice.NewManagedClustersClientWithBaseURI(azure.PublicCloud.ResourceManagerEndpoint, cred.SubscriptionId)
	a, err := auth.NewClientCredentialsConfig(cred.ClientId, cred.ClientSecret, cred.TenantId).Authorizer()
	if err != nil {
		log.Logger.Warn("AKs cluster status check failed")
		return true
	}
	aksClient.Authorizer = a
	aksClient.AddToUserAgent("ktool")
	aksClient.PollingDuration = time.Hour * 1
	conClust, err := aksClient.Get(context.Background(), resourceGroupName, clusterName)
	//conClust, err := aksClient.Get(context.Background(), "TestAKS", "TestCluster")
	if err != nil {
		return true
	}
	if conClust.ManagedClusterProperties.PowerState.Code == "Running" {
		return true
	}
	return false
}

func CreateAKSCluster(options *domain.ClusterOptions) (resp domain.AksClusterContext, err error) {
	//var sshKeyData string
	//if _, err = os.Stat(sshPublicKeyPath); err == nil {
	//	sshBytes, err := ioutil.ReadFile(sshPublicKeyPath)
	//	if err != nil {
	//		log.Fatalf("failed to read SSH key data: %v", err)
	//	}
	//	sshKeyData = string(sshBytes)
	//} else {
	//	sshKeyData = fakepubkey
	//}

	//if err != nil {
	//	return c, fmt.Errorf("cannot get AKS client: %v", err)
	//}
	ctx := context.Background()
	cred, err := iam.GetAksCredentialsForSecretId(strconv.Itoa(options.SecretId))
	aksClient := containerservice.NewManagedClustersClientWithBaseURI(azure.PublicCloud.ResourceManagerEndpoint, cred.SubscriptionId)
	a, err := auth.NewClientCredentialsConfig(cred.ClientId, cred.ClientSecret, cred.TenantId).Authorizer()
	if err != nil {
		return domain.AksClusterContext{
			//ClusterResponse:
			ClusterRequest: *options,
			SecretID:       options.SecretId,
		}, err
	}
	aksClient.Authorizer = a
	aksClient.AddToUserAgent("ktool")
	aksClient.PollingDuration = time.Hour * 1

	//TODO; get from request
	clientName, err := database.GetUserById(ctx, options.UserId)
	if err != nil {
		log.Logger.Warn("Could not fetch the username for the id, " + strconv.Itoa(options.UserId))
		log.Logger.Warn("Using default name ktool-admin")
		clientName = "ktool-admin"
	}

	pvtSSHKey, _ := util.GeneratePrivateKey(4096)
	publicSSHKey, _ := util.GeneratePublicKey(&pvtSSHKey.PublicKey)
	log.Logger.Info("Generated SSH pvt key for user " + clientName + ", key: ")
	//TODO: check whether resource group exists, if not create it

	//microsoft sync process -> async
	PushToJobList(domain.AsyncCloudJob{
		Provider:  "microsoft",
		Status:    domain.AKS_SUBMITTED,
		Reference: options.Name,
		Information: domain.AksAsyncJobParams{
			ClusterOptions: *options,
			CreateRequest: containerservice.ManagedCluster{
				Name:     &options.Name,
				Location: &options.Location,
				ManagedClusterProperties: &containerservice.ManagedClusterProperties{
					DNSPrefix: &options.Name,
					LinuxProfile: &containerservice.LinuxProfile{
						AdminUsername: to.StringPtr(clientName),
						SSH: &containerservice.SSHConfiguration{
							PublicKeys: &[]containerservice.SSHPublicKey{
								{
									KeyData: to.StringPtr(string(publicSSHKey)),
								},
							},
						},
					},
					AgentPoolProfiles: &[]containerservice.ManagedClusterAgentPoolProfile{ //todo: extend here if there needs to be more nodepools
						{
							Count:  to.Int32Ptr(options.InstanceCount),
							Name:   to.StringPtr(strings.ToLower(options.Name) + "p1"),
							VMSize: containerservice.StandardF2sV2,
							Mode:   containerservice.System,
						},
					},
					ServicePrincipalProfile: &containerservice.ManagedClusterServicePrincipalProfile{
						ClientID: to.StringPtr(cred.ClientId),
						Secret:   to.StringPtr(cred.ClientSecret),
					},
				},
			},
			Client: aksClient,
		},
	})
	err = database.AddAKsCluster(ctx, options.Name, options.UserId, options.Name, options.ResourceGroupName, options.Location, pvtSSHKey.D.String())
	if err != nil {
		return domain.AksClusterContext{
			ClusterResponse: domain.AksClusterStatus{
				Name:          options.Name,
				ResourceGroup: options.ResourceGroupName,
				UserName:      clientName,
				SSHPvtKey:     pvtSSHKey.D.String(),
				SSHPubKey:     string(publicSSHKey),
				Status:        "DATABASE INSERT FAILED", //todo: disable on prod
				Error:         err.Error(),
			},
			ClusterRequest: *options,
			SecretID:       options.SecretId,
		}, fmt.Errorf("cannot update the AKS cluster creation in database: %v", err)
	}

	return domain.AksClusterContext{
		ClusterResponse: domain.AksClusterStatus{
			Name:          options.Name,
			ResourceGroup: options.ResourceGroupName,
			UserName:      clientName,
			SSHPvtKey:     pvtSSHKey.D.String(),
			SSHPubKey:     string(publicSSHKey),
			Status:        domain.AKS_CREATING, //todo: disable on prod
		},
		ClusterRequest: *options,
		SecretID:       options.SecretId,
	}, nil
}

func DeleteAksCluster(clusterName, resourceGroupName, secretId string) (err error) {
	cred, err := iam.GetAksCredentialsForSecretId(secretId)
	aksClient := containerservice.NewManagedClustersClientWithBaseURI(azure.PublicCloud.ResourceManagerEndpoint, cred.SubscriptionId)
	a, err := auth.NewClientCredentialsConfig(cred.ClientId, cred.ClientSecret, cred.TenantId).Authorizer()
	if err != nil {
		return err
	}
	aksClient.Authorizer = a
	aksClient.AddToUserAgent("ktool")
	aksClient.PollingDuration = time.Hour * 1
	database.UpdateAksClusterCreationStatus(context.Background(), 1, "DELETING", clusterName, resourceGroupName)
	PushToJobList(domain.AsyncCloudJob{
		Provider:  "microsoft",
		Status:    domain.AKS_SUBMITTED_FOR_DELETION,
		Reference: clusterName,
		Information: domain.AksAsyncJobParams{
			ClusterOptions: domain.ClusterOptions{
				ResourceGroupName: resourceGroupName,
				Name:              clusterName,
			},
			Client: aksClient,
		},
	})
	return nil
}

func CreateResourceGroupIfNotExist(ctx context.Context, resourceGroupName, region, secretId string) (result domain.AksResourceGroup, err error) {
	cred, err := iam.GetAksCredentialsForSecretId(secretId)
	aksClient := resources.NewGroupsClientWithBaseURI(azure.PublicCloud.ResourceManagerEndpoint, cred.SubscriptionId)
	a, err := auth.NewClientCredentialsConfig(cred.ClientId, cred.ClientSecret, cred.TenantId).Authorizer()
	if err != nil {
		return result, err
	}
	aksClient.Authorizer = a
	aksClient.AddToUserAgent("ktool")
	aksClient.PollingDuration = time.Hour * 1
	res, err := aksClient.CheckExistence(ctx, resourceGroupName)
	//if err != nil {
	//	return result, err
	//}
	s := strings.Split(res.Status, " ")[0]
	if s == "404" {
		log.Logger.Info("Resource group %s not found. Attempting to create", resourceGroupName)
		group, err := aksClient.CreateOrUpdate(ctx, resourceGroupName, resources.Group{
			Response:   autorest.Response{},
			ID:         nil,
			Name:       &resourceGroupName,
			Type:       nil,
			Properties: nil,
			Location:   &region,
			ManagedBy:  nil,
			Tags:       nil,
		})
		if err != nil {
			return result, err
		}
		log.Logger.Info("Resource group %s created in region %s", resourceGroupName, region)
		log.Logger.Info("group resp, %s", group.Response.Status)
		return domain.AksResourceGroup{
			Groups: []string{resourceGroupName},
			Status: "CREATED",
			Error:  "",
		}, nil
	} else if s == "204" {
		return domain.AksResourceGroup{
			Groups: nil,
			Status: "EXISTING",
			Error:  "",
		}, nil
	} else {
		return domain.AksResourceGroup{
			Groups: nil,
			Status: "ERROR OCCURRED",
			Error:  err.Error(),
		}, err
	}
}

//internal helpers
func SyncDeleteAksCluster(ctx context.Context, aksClient containerservice.ManagedClustersClient, resourceGroupName, clusterName string) error {
	res, err := aksClient.Delete(context.Background(), resourceGroupName, clusterName)
	if err != nil {
		return err
	}
	err = res.WaitForCompletionRef(context.Background(), aksClient.Client)
	if err != nil {
		return err
	}
	return nil
}

func SyncCreateAksCluster(ctx context.Context, aksClient containerservice.ManagedClustersClient, options domain.ClusterOptions,
	params containerservice.ManagedCluster) (containerservice.ManagedClustersCreateOrUpdateFuture, error) {

	future, err := aksClient.CreateOrUpdate(
		ctx,
		options.ResourceGroupName,
		options.Name,
		params,
	)
	if err != nil {
		log.Logger.Error(fmt.Errorf("aks cluster creation failed by microsoft; %s", err))
		return future, err
	}
	return future, nil
}
