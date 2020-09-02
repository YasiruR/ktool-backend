package kubernetes

import (
	"context"
	"fmt"
	"github.com/YasiruR/ktool-backend/database"
	"github.com/YasiruR/ktool-backend/domain"
	"github.com/YasiruR/ktool-backend/iam"
	"github.com/YasiruR/ktool-backend/log"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/client"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/eks"
	"strconv"
)

func main() {
	region := "us-east-2"
	cluster := "ktool-test-cluster"

	config := aws.Config{
		Credentials: credentials.NewStaticCredentials("AKIAQMZUT3KWPZ3BLHUO", "cqKaFp0AHf/KOoiHUJd01DPfxSkYcAE3h9+uMSot", ""),
		//Credentials: credentials.NewStaticCredentials(id, secret, ""),
		Region: &region,
	}
	session, _ := session.NewSession(&config)
	svc := eks.New(session)

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

func CheckEksClusterCreationStatus(clusterName string, secretId int) (eks.DescribeClusterOutput, error) {
	//region := "us-east-2"
	//cluster := "ktool-test-cluster"
	//id := "AKIAY4OR54E7L5QR3QRF"
	//secret := "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"

	cred, err := iam.GetEksCredentialsForSecretId(strconv.Itoa(secretId))
	if err != nil {
		return eks.DescribeClusterOutput{}, err
	}
	//id := "AKIAY4OR54E7L5QR3QRF"
	//secret := "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"

	sess, _ := session.NewSession(&aws.Config{
		Credentials: cred,
		Region:      aws.String("us-east-2"),
	})
	svc := eks.New(sess)

	input := &eks.DescribeClusterInput{Name: &clusterName}
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
		return eks.DescribeClusterOutput{}, err
	}
	return *result, nil
}

func CreateEksCluster(clusterId string, secretId int, createClusterRequest *domain.GkeClusterOptions) (eks.CreateClusterOutput, error) {
	cred, err := iam.GetEksCredentialsForSecretId(strconv.Itoa(secretId))
	if err != nil {
		return eks.CreateClusterOutput{}, err
	}
	//id := "AKIAY4OR54E7L5QR3QRF"
	//secret := "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"

	sess, _ := session.NewSession(&aws.Config{
		Credentials: cred,
		Region:      aws.String("us-east-2"),
	})
	svc := eks.New(sess)
	input := &eks.CreateClusterInput{
		ClientRequestToken: aws.String(clusterId),
		Name:               aws.String("dev"),
		//Name:               aws.String(createClusterRequest.Name),
		ResourcesVpcConfig: &eks.VpcConfigRequest{
			SecurityGroupIds: []*string{
				aws.String("sg-6979fe18"),
			},
			SubnetIds: []*string{
				aws.String("subnet-6782e71e"),
				aws.String("subnet-e7e761ac"),
			},
		},
		//RoleArn: aws.String("arn:aws:iam::012345678910:role/eks-service-role-AWSServiceRoleForAmazonEKS-J7ONKE3BQ4PI"),
		RoleArn: aws.String("arn:aws:iam::610862489918:role/eks-admin"),
		Version: aws.String(createClusterRequest.KubVersion),
	}

	result, err := svc.CreateCluster(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case eks.ErrCodeResourceInUseException:
				fmt.Println(eks.ErrCodeResourceInUseException, aerr.Error())
			case eks.ErrCodeResourceLimitExceededException:
				fmt.Println(eks.ErrCodeResourceLimitExceededException, aerr.Error())
			case eks.ErrCodeInvalidParameterException:
				fmt.Println(eks.ErrCodeInvalidParameterException, aerr.Error())
			case eks.ErrCodeClientException:
				fmt.Println(eks.ErrCodeClientException, aerr.Error())
			case eks.ErrCodeServerException:
				fmt.Println(eks.ErrCodeServerException, aerr.Error())
			case eks.ErrCodeServiceUnavailableException:
				fmt.Println(eks.ErrCodeServiceUnavailableException, aerr.Error())
			case eks.ErrCodeUnsupportedAvailabilityZoneException:
				fmt.Println(eks.ErrCodeUnsupportedAvailabilityZoneException, aerr.Error())
			default:
				fmt.Println(aerr.Error())
			}
		} else {
			// Print the error, cast err to awserr.Error to get the Code and
			// Message from an error.
			fmt.Println(err.Error())
		}
		return eks.CreateClusterOutput{}, err
	}

	// persist in db
	err = database.AddEksCluster(context.Background(), clusterId, createClusterRequest.UserId, createClusterRequest.Name, *result.Cluster.Arn)
	return *result, nil
}

func DeleteEksCluster(clusterName string, session *client.ConfigProvider) *eks.DeleteClusterOutput {
	//TODO; get the secret using id here
	//cred, err := GetEksCredentialsForSecretId(strconv.Itoa(secretId))
	//if err != nil{
	//	return &eks.DeleteClusterOutput{}
	//}
	//id := "AKIA5CVBUZ342ISPRDVJ"
	//secret := "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"

	//session, _ := session.NewSession(&aws.Config{
	//	Credentials: cred,
	//	Region:      aws.String("us-east-2"),
	//})
	svc := eks.New(*session)
	input := &eks.DeleteClusterInput{
		Name: &clusterName,
	}

	result, err := svc.DeleteCluster(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case eks.ErrCodeResourceInUseException:
				fmt.Println(eks.ErrCodeResourceInUseException, aerr.Error())
			case eks.ErrCodeResourceNotFoundException:
				fmt.Println(eks.ErrCodeResourceNotFoundException, aerr.Error())
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
		return nil
	}

	return result
}

func generateEKSClusterCreationRequest(request *domain.GkeClusterOptions) *eks.CreateClusterInput {
	return &eks.CreateClusterInput{
		ClientRequestToken: nil,
		EncryptionConfig:   nil,
		Logging:            nil,
		Name:               &request.Name,
		ResourcesVpcConfig: nil,
		RoleArn:            nil,
		Tags:               nil,
		Version:            nil,
	}

}

func generateEKSClusterDeletionRequest(clusterId string) *eks.DeleteClusterInput {
	return &eks.DeleteClusterInput{
		Name: aws.String(clusterId),
	}
}

func ListEksClusers(userID string) eks.ListClustersOutput {
	region := "us-east-2" //TODO: global var
	cred, err := iam.GetEksCredentialsForUser(userID)
	if err != nil {
		log.Logger.ErrorContext(context.Background(), "Error occurred while fetching eks secret for client %s", userID)
		return eks.ListClustersOutput{}
	}
	sess, _ := session.NewSession(&aws.Config{
		Credentials: cred,
		Region:      &region,
	})
	svc := eks.New(sess)

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
