package kubernetes

import (
	container "cloud.google.com/go/container/apiv1"
	"context"
	//"encoding/json"
	//"fmt"
	//"github.com/YasiruR/ktool-backend/database"
	//"github.com/YasiruR/ktool-backend/domain"
	//"github.com/YasiruR/ktool-backend/log"
	iam "github.com/YasiruR/ktool-backend/iam"
	//"golang.org/x/oauth2/google"
	"google.golang.org/api/option"
	containerpb "google.golang.org/genproto/googleapis/container/v1"
)

// todo: maintain a map of credentials

//func main() {
//	ctx := context.Background()
//	res, err := ListGkeClusters("1")
//	if err != nil {
//		log.Logger.ErrorContext(ctx, "Could not retrieve cluster list")
//		return
//	}
//	log.Logger.InfoContext(ctx, "Successfully retrieved cluster list from GKE")
//	fmt.Println(res)
//}

func ListGkeClusters(userId string) (*containerpb.ListClustersResponse, error) {
	ctx := context.Background()
	b, cred, err := iam.GetGkeCredentialsForUser(userId)
	if err != nil {
		return nil, err
	}
	c, err := container.NewClusterManagerClient(ctx, option.WithCredentialsJSON(b))
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
