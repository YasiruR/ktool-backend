package kubernetes

import (
	"context"
	"fmt"
	"github.com/Azure/azure-sdk-for-go/profiles/latest/containerservice/mgmt/containerservice"
	"github.com/Azure/go-autorest/autorest/azure"
	auth "github.com/Azure/go-autorest/autorest/azure/auth"
	"github.com/Azure/go-autorest/autorest/to"
	"github.com/YasiruR/ktool-backend/domain"
	"github.com/YasiruR/ktool-backend/iam"
	"github.com/YasiruR/ktool-backend/log"
	"github.com/YasiruR/ktool-backend/util"
	"strconv"
	"time"
)

func GetAKSClusterStatus(clusterName string, resourceGroupName string, secretId string) (status containerservice.ManagedCluster, err error) {
	cred, err := iam.GetAksCredentialsForSecretId(secretId)
	aksClient := containerservice.NewManagedClustersClientWithBaseURI(azure.PublicCloud.ResourceManagerEndpoint, cred.SubscriptionId)
	a, err := auth.NewClientCredentialsConfig(cred.ClientId, cred.ClientSecret, cred.TenantId).Authorizer()
	if err != nil {
		return containerservice.ManagedCluster{}, err
	}
	aksClient.Authorizer = a
	aksClient.AddToUserAgent("ktool")
	aksClient.PollingDuration = time.Hour * 1
	conClust, err := aksClient.Get(context.Background(), resourceGroupName, clusterName)
	//conClust, err := aksClient.Get(context.Background(), "TestAKS", "TestCluster")
	if err != nil {
		print(conClust.Status)
	} else {
		print(conClust.ID)
	}
	return conClust, nil
}

func CreateAKSCluster(clusterName, resourceGroupName string, secretId int, options *domain.ClusterOptions) (c containerservice.ManagedCluster, err error) {
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
	cred, err := iam.GetAksCredentialsForSecretId(strconv.Itoa(secretId))
	aksClient := containerservice.NewManagedClustersClientWithBaseURI(azure.PublicCloud.ResourceManagerEndpoint, cred.SubscriptionId)
	a, err := auth.NewClientCredentialsConfig(cred.ClientId, cred.ClientSecret, cred.TenantId).Authorizer()
	if err != nil {
		return containerservice.ManagedCluster{}, err
	}
	aksClient.Authorizer = a
	aksClient.AddToUserAgent("ktool")
	aksClient.PollingDuration = time.Hour * 1

	//TODO; get from request
	clientName := "ktool-" + "admin"

	pvtSSHKey, _ := util.GeneratePrivateKey(4096)
	publicSSHKey, _ := util.GeneratePublicKey(&pvtSSHKey.PublicKey)
	log.Logger.Info("Generated SSH pvt key for user " + clientName + ", key: ")
	//TODO: check whetehr resource group exists, if not create it
	future, err := aksClient.CreateOrUpdate(
		ctx,
		resourceGroupName,
		clusterName,
		containerservice.ManagedCluster{
			Name:     &clusterName,
			Location: &options.Location,
			ManagedClusterProperties: &containerservice.ManagedClusterProperties{
				DNSPrefix: &clusterName,
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
				AgentPoolProfiles: &[]containerservice.ManagedClusterAgentPoolProfile{
					{
						Count:  to.Int32Ptr(options.InstanceCount),
						Name:   to.StringPtr("agentpool1"),
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
	)
	if err != nil {
		return c, fmt.Errorf("cannot create AKS cluster: %v", err)
	}

	err = future.WaitForCompletionRef(ctx, aksClient.Client)
	if err != nil {
		return c, fmt.Errorf("cannot get the AKS cluster create or update future response: %v", err)
	}

	return future.Result(aksClient)
}
