package kubernetes

import (
	"context"
	"fmt"
	"github.com/YasiruR/ktool-backend/database"
	"github.com/YasiruR/ktool-backend/domain"
	"github.com/YasiruR/ktool-backend/iam"
	"github.com/YasiruR/ktool-backend/log"
	"github.com/YasiruR/ktool-backend/util"
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
	cred, err := iam.GetEksCredentialsForSecretId(strconv.Itoa(secretId))
	if err != nil {
		return eks.DescribeClusterOutput{}, err
	}

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

func CheckEksNodeGroupCreationStatus(clusterName string, nodeGroupName string, secretId int) (eks.DescribeNodegroupOutput, error) {
	cred, err := iam.GetEksCredentialsForSecretId(strconv.Itoa(secretId))
	if err != nil {
		return eks.DescribeNodegroupOutput{}, err
	}

	sess, _ := session.NewSession(&aws.Config{
		Credentials: cred,
		Region:      aws.String("us-east-2"),
	})
	svc := eks.New(sess)

	input := &eks.DescribeNodegroupInput{
		ClusterName:   &clusterName,
		NodegroupName: &nodeGroupName,
	}
	result, err := svc.DescribeNodegroup(input)
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
		return eks.DescribeNodegroupOutput{}, err
	}
	return *result, nil
}

func CreateEksCluster(clusterId string, secretId int, createClusterRequest *domain.GkeClusterOptions) (domain.EksClusterCreationResponse, error) {
	cred, err := iam.GetEksCredentialsForSecretId(strconv.Itoa(secretId))
	//nodeGroupResp := domain.EksClusterStatus{}
	if err != nil {
		return domain.EksClusterCreationResponse{}, err
	}

	sess, _ := session.NewSession(&aws.Config{
		Credentials: cred,
		Region:      aws.String("us-east-2"),
	})

	// get ARN from here
	// https://docs.aws.amazon.com/eks/latest/userguide/service_IAM_role.html#create-service-role
	arn := "arn:aws:iam::899060911865:role/EKSManagerRole"

	svc := eks.New(sess)
	ctrlResp, err := createEksControlPlane(svc, clusterId, createClusterRequest.Name, arn, "1.17")
	resp := domain.EksClusterCreationResponse{
		ClusterStatus: ctrlResp,
		SecretID:      secretId,
	}
	if err != nil {
		return resp, err
	}

	// we are not sending the node group creation request just yet

	//if *ctrlResp.CreateClusterOutput.Cluster.Status == "CREATING" {
	//	nodeGroupResp, err = createEksNodeGroup(svc, ctrlResp, createClusterRequest)
	//}

	//if err != nil {
	//	return resp, err
	//}
	// persist in db
	err = database.AddEksCluster(context.Background(), clusterId, createClusterRequest.UserId, createClusterRequest.Name,
		ctrlResp.RequestToken, ctrlResp.ClusterArn, ctrlResp.RoleArn, util.StringPointerListToEscapedCSV(ctrlResp.SubnetIds), ctrlResp.KubVersion)
	//return nodeGroupResp, nil
	return resp, nil
}

func createEksControlPlane(svc *eks.EKS, id string, name string, arn string, kubVersion string) (clusterCreationOutput domain.EksClusterStatus, err error) {
	input := &eks.CreateClusterInput{
		ClientRequestToken: aws.String(id),
		Name:               aws.String(name),
		ResourcesVpcConfig: &eks.VpcConfigRequest{
			SecurityGroupIds: []*string{
				aws.String("sg-43f0613f"),
			},
			SubnetIds: []*string{
				aws.String("subnet-b46865ce"),
				aws.String("subnet-a6b2c7ea"),
				aws.String("subnet-934a9ef8"),
			},
		},
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
		Name:         *result.Cluster.Name,
		ClusterArn:   *result.Cluster.Arn,
		RequestToken: id,
		RoleArn:      *result.Cluster.RoleArn,
		SubnetIds:    &result.Cluster.ResourcesVpcConfig.SubnetIds,
		KubVersion:   *result.Cluster.Version,
		Status:       *result.Cluster.Status,
		Error:        "nil",
	}, nil
}

func CreateEksNodeGroup(secretId int, eksClusterStatus domain.EksClusterStatus) (nodeGroupResponse domain.EksNodeGroupCreationResponse, err error) {
	cred, err := iam.GetEksCredentialsForSecretId(strconv.Itoa(secretId))
	if err != nil {
		return domain.EksNodeGroupCreationResponse{}, err
	}

	sess, _ := session.NewSession(&aws.Config{
		Credentials: cred,
		Region:      aws.String("us-east-2"),
	})

	svc := eks.New(sess)
	ngResp, err := createEksNodeGroup(svc, eksClusterStatus)
	if err != nil {
		return domain.EksNodeGroupCreationResponse{}, err
	}
	resp := domain.EksNodeGroupCreationResponse{
		SecretId: secretId,
		Response: *ngResp.Nodegroup,
	}

	// persist in db
	//err = database.AddEksCluster(context.Background(), clusterId, createClusterRequest.UserId, createClusterRequest.Name, createClusterRequest.Name)
	//return nodeGroupResp, nil
	return resp, nil
}

func createEksNodeGroup(svc *eks.EKS, ctrlPlaneResponse domain.EksClusterStatus) (nodeGroupResponse eks.CreateNodegroupOutput, err error) {
	groupName, _ := uuid.GenerateUUID()
	input := &eks.CreateNodegroupInput{
		AmiType:            nil,
		ClientRequestToken: &ctrlPlaneResponse.RequestToken,
		ClusterName:        &ctrlPlaneResponse.Name,
		DiskSize:           nil,
		//InstanceTypes:      []*string{&clusterInput.MachineFamily},
		InstanceTypes:  nil,
		Labels:         nil,
		NodeRole:       &ctrlPlaneResponse.RoleArn,
		NodegroupName:  &groupName,
		ReleaseVersion: nil,
		RemoteAccess:   nil,
		ScalingConfig:  nil,
		Subnets:        *ctrlPlaneResponse.SubnetIds,
		Tags:           nil,
		Version:        &ctrlPlaneResponse.KubVersion,
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
		return eks.CreateNodegroupOutput{}, err
	}
	return *result, nil
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
