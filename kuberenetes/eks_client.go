package kubernetes

import (
	"context"
	"fmt"
	"github.com/YasiruR/ktool-backend/database"
	"github.com/YasiruR/ktool-backend/log"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/eks"
)

func main() {
	region := "us-east-2"
	cluster := "ktool-test-cluster"
	config := aws.Config{
		Credentials: credentials.NewStaticCredentials("AKIAQMZUT3KWPZ3BLHUO", "cqKaFp0AHf/KOoiHUJd01DPfxSkYcAE3h9+uMSot", ""),
		Region:      &region,
	}
	svc := eks.New(session.New(&config))

	////list clusters
	//input := &eks.ListClustersInput{}
	//
	//result, err := svc.ListClusters(input)

	//describe cluster
	input := &eks.DescribeClusterInput{Name: &cluster}
	result, err := svc.DescribeCluster(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case eks.ErrCodeInvalidParameterException:
				fmt.Println(eks.ErrCodeInvalidParameterException, aerr.Error())
			case eks.ErrCodeClientException:
				fmt.Println(eks.ErrCodeClientException, aerr.Error())
			case eks.ErrCodeServerException:
				fmt.Println(eks.ErrCodeServerException, aerr.Error())
			case eks.ErrCodeServiceUnavailableException:
				fmt.Println(eks.ErrCodeServiceUnavailableException, aerr.Error())
			default:
				fmt.Println(aerr.Error())
			}
		} else {
			// Print the error, cast err to awserr.Error to get the Code and
			// Message from an error.
			fmt.Println(err.Error())
		}
		return
	}

	fmt.Println(result)
}

func ListEksClusers(userID string) eks.ListClustersOutput {
	region := "us-east-2" //TODO: how to get this?
	cred, err := GetEksCredentialsForUser(userID)
	if err != nil {
		log.Logger.ErrorContext(context.Background(), "Error occurred while fetching eks secret for client %s", userID)
		return eks.ListClustersOutput{}
	}
	config := aws.Config{
		Credentials: cred,
		Region:      &region,
	}
	svc := eks.New(session.New(&config))

	//list clusters
	input := &eks.ListClustersInput{}
	result, err := svc.ListClusters(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case eks.ErrCodeInvalidParameterException:
				fmt.Println(eks.ErrCodeInvalidParameterException, aerr.Error())
			case eks.ErrCodeClientException:
				fmt.Println(eks.ErrCodeClientException, aerr.Error())
			case eks.ErrCodeServerException:
				fmt.Println(eks.ErrCodeServerException, aerr.Error())
			case eks.ErrCodeServiceUnavailableException:
				fmt.Println(eks.ErrCodeServiceUnavailableException, aerr.Error())
			default:
				fmt.Println(aerr.Error())
			}
		} else {
			// Print the error, cast err to awserr.Error to get the Code and
			// Message from an error.
			fmt.Println(err.Error())
		}
		return eks.ListClustersOutput{}
	}
	return *result
}

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
