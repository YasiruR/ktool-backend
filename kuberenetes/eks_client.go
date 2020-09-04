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
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/eks"
	"github.com/hashicorp/go-uuid"
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

func CreateEksCluster(clusterId string, secretId int, createClusterRequest *domain.GkeClusterOptions) (domain.EksClusterStatus, error) {
	cred, err := iam.GetEksCredentialsForSecretId(strconv.Itoa(secretId))
	nodeGroupResp := domain.EksClusterStatus{}
	if err != nil {
		return domain.EksClusterStatus{}, err
	}
	//id := "AKIAY4OR54E7L5QR3QRF"
	//secret := "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"

	sess, _ := session.NewSession(&aws.Config{
		Credentials: cred,
		Region:      aws.String("us-east-2"),
	})
	arn := "arn:aws:iam::899060911865:user/ktool-admin"

	svc := eks.New(sess)
	ctrlResp, err := createEksControlPlane(svc, clusterId, createClusterRequest.Name, arn, "1.15")
	if err != nil {
		return ctrlResp, err
	}
	if *ctrlResp.CreateClusterOutput.Cluster.Status == "SUCCESS" {
		nodeGroupResp, err = createEksNodeGroup(svc, ctrlResp, createClusterRequest)
	}
	if err != nil {
		return ctrlResp, err
	}
	// persist in db
	err = database.AddEksCluster(context.Background(), clusterId, createClusterRequest.UserId, createClusterRequest.Name, createClusterRequest.Name)
	return nodeGroupResp, nil
}

func createEksControlPlane(svc *eks.EKS, id string, name string, arn string, kubVersion string) (clusterCreationOutput domain.EksClusterStatus, err error) {
	input := &eks.CreateClusterInput{
		ClientRequestToken: aws.String(id),
		//Name:               aws.String("dev"),
		Name: aws.String(name),
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
		//RoleArn: aws.String("arn:aws:iam::610862489918:role/eks-admin"),
		RoleArn: aws.String(arn),
		Version: aws.String(kubVersion),
	}

	result, err := svc.CreateCluster(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case eks.ErrCodeResourceInUseException:
				log.Logger.Error(eks.ErrCodeResourceInUseException + aerr.Error())
			case eks.ErrCodeResourceLimitExceededException:
				log.Logger.Error(eks.ErrCodeResourceLimitExceededException + aerr.Error())
			case eks.ErrCodeInvalidParameterException:
				log.Logger.Error(eks.ErrCodeInvalidParameterException + aerr.Error())
			case eks.ErrCodeClientException:
				log.Logger.Error(eks.ErrCodeClientException + aerr.Error())
			case eks.ErrCodeServerException:
				log.Logger.Error(eks.ErrCodeServerException + aerr.Error())
			case eks.ErrCodeServiceUnavailableException:
				log.Logger.Error(eks.ErrCodeServiceUnavailableException + aerr.Error())
			case eks.ErrCodeUnsupportedAvailabilityZoneException:
				log.Logger.Error(eks.ErrCodeUnsupportedAvailabilityZoneException + aerr.Error())
			default:
				log.Logger.Error(aerr.Error())
			}
		} else {
			// Print the error, cast err to awserr.Error to get the Code and
			// Message from an error.
			log.Logger.Error(err.Error())
		}
		return domain.EksClusterStatus{}, err
	}
	return domain.EksClusterStatus{
		CreateClusterOutput: *result,
	}, nil
}

func createEksNodeGroup(svc *eks.EKS, ctrlplaneResponse domain.EksClusterStatus, clusterInput *domain.GkeClusterOptions) (nodeGroupResponse domain.EksClusterStatus, err error) {
	groupName, _ := uuid.GenerateUUID()
	input := &eks.CreateNodegroupInput{
		AmiType:            nil,
		ClientRequestToken: ctrlplaneResponse.CreateClusterOutput.Cluster.ClientRequestToken,
		ClusterName:        ctrlplaneResponse.CreateClusterOutput.Cluster.Name,
		DiskSize:           nil,
		InstanceTypes:      []*string{&clusterInput.MachineFamily},
		Labels:             nil,
		NodeRole:           nil,
		NodegroupName:      &groupName,
		ReleaseVersion:     nil,
		RemoteAccess:       nil,
		ScalingConfig:      nil,
		Subnets:            nil,
		Tags:               nil,
		Version:            &clusterInput.KubVersion,
	}

	result, err := svc.CreateNodegroup(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case eks.ErrCodeResourceInUseException:
				log.Logger.Error(eks.ErrCodeResourceInUseException + aerr.Error())
			case eks.ErrCodeResourceLimitExceededException:
				log.Logger.Error(eks.ErrCodeResourceLimitExceededException + aerr.Error())
			case eks.ErrCodeInvalidParameterException:
				log.Logger.Error(eks.ErrCodeInvalidParameterException + aerr.Error())
			case eks.ErrCodeClientException:
				log.Logger.Error(eks.ErrCodeClientException + aerr.Error())
			case eks.ErrCodeServerException:
				log.Logger.Error(eks.ErrCodeServerException + aerr.Error())
			case eks.ErrCodeServiceUnavailableException:
				log.Logger.Error(eks.ErrCodeServiceUnavailableException + aerr.Error())
			case eks.ErrCodeUnsupportedAvailabilityZoneException:
				log.Logger.Error(eks.ErrCodeUnsupportedAvailabilityZoneException + aerr.Error())
			default:
				log.Logger.Error(aerr.Error())
			}
		} else {
			// Print the error, cast err to awserr.Error to get the Code and
			// Message from an error.
			log.Logger.Error(err.Error())
		}
		return ctrlplaneResponse, err
	}
	ctrlplaneResponse.CreateNodGroupOutput = *result
	return ctrlplaneResponse, nil
}

func DeleteEksCluster(clusterName string, secretId string) (out *eks.DeleteClusterOutput, err error) {
	//TODO; get the secret using id here
	cred, err := iam.GetEksCredentialsForSecretId(secretId)
	//if err != nil{
	//	return &eks.DeleteClusterOutput{}
	//}
	//id := "AKIA5CVBUZ342ISPRDVJ"
	//secret := "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"

	session, _ := session.NewSession(&aws.Config{
		Credentials: cred,
		Region:      aws.String("us-east-2"),
	})
	svc := eks.New(session)
	input := &eks.DeleteClusterInput{
		Name: &clusterName,
	}

	result, err := svc.DeleteCluster(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case eks.ErrCodeResourceInUseException:
				log.Logger.Error(eks.ErrCodeResourceInUseException, aerr.Error())
			case eks.ErrCodeResourceNotFoundException:
				log.Logger.Error(eks.ErrCodeResourceNotFoundException, aerr.Error())
			case eks.ErrCodeClientException:
				log.Logger.Error(eks.ErrCodeClientException, aerr.Error())
			case eks.ErrCodeServerException:
				log.Logger.Error(eks.ErrCodeServerException, aerr.Error())
			case eks.ErrCodeServiceUnavailableException:
				log.Logger.Error(eks.ErrCodeServiceUnavailableException, aerr.Error())
			default:
				log.Logger.Error(aerr.Error())
			}
		} else {
			// Print the error, cast err to awserr.Error to get the Code and
			// Message from an error.
			log.Logger.Error(err.Error())
		}
		return nil, err
	}

	return result, nil
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
