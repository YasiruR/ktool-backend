package iam

import (
	"context"
	"github.com/YasiruR/ktool-backend/database"
	"github.com/YasiruR/ktool-backend/log"
	"github.com/aws/aws-sdk-go/aws/credentials"
)

func GetEksCredentialsForUser(userId string) (*credentials.Credentials, error) {
	ctx := context.Background()
	secretDao := database.GetSecretInternal(ctx, userId, `Amazon`, `aws-ktool`)

	if err := secretDao.Error; err != nil {
		log.Logger.ErrorContext(ctx, "Error occurred while fetching eks secret for client %s", userId)
		return &credentials.Credentials{}, err
	}
	cred := credentials.NewStaticCredentials(secretDao.Secret.EksAccessKeyId, secretDao.Secret.EksSecretAccessKey, "")
	return cred, nil
}

func GetEksCredentialsForSecretId(secretId string) (*credentials.Credentials, error) {
	ctx := context.Background()
	secretDao := database.GetSecretByIdInternal(ctx, secretId, "amazon")

	if err := secretDao.Error; err != nil {
		log.Logger.ErrorContext(ctx, "Error occurred while fetching eks secret for client %s", secretId)
		return &credentials.Credentials{}, err
	}
	cred := credentials.NewStaticCredentials(secretDao.Secret.EksAccessKeyId, secretDao.Secret.EksSecretAccessKey, "")
	return cred, nil
}
