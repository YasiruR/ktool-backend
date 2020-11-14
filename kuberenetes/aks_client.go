package kubernetes

import (
	"context"
	"github.com/Azure/azure-sdk-for-go/profiles/latest/containerservice/mgmt/containerservice"
	"github.com/Azure/go-autorest/autorest/azure"
	auth "github.com/Azure/go-autorest/autorest/azure/auth"
	"github.com/YasiruR/ktool-backend/iam"
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
	conClust, err := aksClient.Get(context.Background(), "TestAKS", "TestCluster")
	if err != nil {
		print(conClust.Status)
	} else {
		print(conClust.ID)
	}
	return conClust, nil
}
