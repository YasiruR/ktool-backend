package main

import (
	container "cloud.google.com/go/container/apiv1"
	"context"
	"encoding/json"
	"fmt"
	"github.com/YasiruR/ktool-backend/database"
	"github.com/YasiruR/ktool-backend/domain"
	"github.com/YasiruR/ktool-backend/log"
	"google.golang.org/api/option"
	containerpb "google.golang.org/genproto/googleapis/container/v1"
)

// todo: maintain a map of credentials

func main() {
	ctx := context.Background()
	res, err := ListClusters("1")
	if err != nil {
		log.Logger.ErrorContext(ctx, "Could not retrieve cluster list")
		return
	}
	log.Logger.InfoContext(ctx, "Successfully retrieved cluster list from GKE")
	fmt.Println(res)
}

func ListClusters(userId string) (*containerpb.ListClustersResponse, error) {
	ctx := context.Background()
	b, cred, err := GetGkeCredentialsForUser(userId)
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

func GetGkeCredentialsForUser(userId string) ([]byte, domain.GkeSecret, error) {
	ctx := context.Background()
	secretDao := database.GetAllSecretsByUserInternal(ctx, userId, `Google`)

	if err := secretDao.Error; err != nil {
		log.Logger.ErrorContext(ctx, "Error occurred while fetching eks secret for client %s", userId)
		return nil, domain.GkeSecret{}, err
	}
	firstSecret := secretDao.SecretList[1]
	cred := domain.GkeSecret{
		Type:              firstSecret.Type,
		ProjectId:         firstSecret.ProjectId,
		PrivateKeyId:      firstSecret.PrivateKeyId,
		PrivateKey:        firstSecret.PrivateKey,
		ClientMail:        firstSecret.ClientMail,
		ClientId:          firstSecret.ClientId,
		AuthUri:           firstSecret.AuthUri,
		TokenUri:          firstSecret.TokenUri,
		AuthX509CertUrl:   firstSecret.AuthX509CertUrl,
		ClientX509CertUrl: firstSecret.ClientX509CertUrl,
	}
	bytes, err := json.Marshal(&cred)
	if err != nil {
		log.Logger.ErrorContext(ctx, "Could not marshall gke credentials for user %s", userId)
		return nil, cred, err
	}
	return bytes, cred, nil
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
