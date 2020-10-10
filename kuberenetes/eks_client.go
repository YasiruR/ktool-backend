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
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/eks"
	"github.com/hashicorp/go-uuid"
	"strconv"
)

//func main() {
//	region := "us-east-2"
//	cluster := "ktool-test-cluster"
//
//	config := aws.Config{
//		Credentials: credentials.NewStaticCredentials("AKIAQMZUT3KWPZ3BLHUO", "cqKaFp0AHf/KOoiHUJd01DPfxSkYcAE3h9+uMSot", ""),
//		//Credentials: credentials.NewStaticCredentials(id, secret, ""),
//		Region: &region,
//	}
//	session, _ := session.NewSession(&config)
//	svc := eks.New(session)
//
//	////list clusters
//	//input := &eks.ListClustersInput{}
//	//
//	//result, err := svc.ListClusters(input)
//
//	//describe cluster
//	input := &eks.DescribeClusterInput{Name: &cluster}
//	result, err := svc.DescribeCluster(input)
//	if err != nil {
//		if aerr, ok := err.(awserr.Error); ok {
//			switch aerr.Code() {
//			case eks.ErrCodeInvalidParameterException:
//				fmt.Println(eks.ErrCodeInvalidParameterException, aerr.Error())
//			case eks.ErrCodeClientException:
//				fmt.Println(eks.ErrCodeClientException, aerr.Error())
//			case eks.ErrCodeServerException:
//				fmt.Println(eks.ErrCodeServerException, aerr.Error())
//			case eks.ErrCodeServiceUnavailableException:
//				fmt.Println(eks.ErrCodeServiceUnavailableException, aerr.Error())
//			default:
//				fmt.Println(aerr.Error())
//			}
//		} else {
//			// Print the error, cast err to awserr.Error to get the Code and
//			// Message from an error.
//			fmt.Println(err.Error())
//		}
//		return
//	}
//
//	fmt.Println(result)
//}

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

func CheckEksClusterCreationStatus(clusterName string, region string, secretId int) (eks.DescribeClusterOutput, error) {
	cred, err := iam.GetEksCredentialsForSecretId(strconv.Itoa(secretId))
	if err != nil {
		return eks.DescribeClusterOutput{}, err
	}

	sess, _ := session.NewSession(&aws.Config{
		Credentials: cred,
		Region:      aws.String(region),
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

func CheckEksNodeGroupCreationStatus(clusterName string, nodeGroupName string, region string, secretId int) (eks.DescribeNodegroupOutput, error) {
	cred, err := iam.GetEksCredentialsForSecretId(strconv.Itoa(secretId))
	if err != nil {
		return eks.DescribeNodegroupOutput{}, err
	}

	sess, _ := session.NewSession(&aws.Config{
		Credentials: cred,
		Region:      aws.String(region),
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
			case eks.ErrCodeResourceNotFoundException:
				fmt.Println(eks.ErrCodeResourceNotFoundException, aerr.Error())
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

func CreateEksCluster(clusterId string, secretId int, createClusterRequest *domain.ClusterOptions) (domain.EksClusterContext, error) {
	cred, err := iam.GetEksCredentialsForSecretId(strconv.Itoa(secretId))
	//nodeGroupResp := domain.EksClusterStatus{}
	if err != nil {
		return domain.EksClusterContext{}, err
	}

	sess, _ := session.NewSession(&aws.Config{
		Credentials: cred,
		Region:      aws.String(createClusterRequest.Location),
	})

	// get ARN from here
	// https://docs.aws.amazon.com/eks/latest/userguide/service_IAM_role.html#create-service-role
	arn := "arn:aws:iam::899060911865:role/EKSManagerRole" //todo: fetch from api
	// fetch vpc config for user for region
	vpcConfig := getVPCConfigForUSerForRegion(secretId, createClusterRequest.Location)

	svc := eks.New(sess)
	ctrlResp, err := createEksControlPlane(svc, clusterId, createClusterRequest.Name, arn, createClusterRequest.KubVersion, vpcConfig)
	resp := domain.EksClusterContext{
		ClusterStatus:  ctrlResp,
		ClusterRequest: *createClusterRequest,
		SecretID:       secretId,
	}
	if err != nil {
		return resp, err
	}

	// submit job for the watcher
	//service.PushToJobList(service.AsyncCloudJob{
	//	Provider:    "amazon",
	//	Status:      service.EKS_MASTER_CREATING,
	//	Reference:   createClusterRequest.Name,
	//	Information: resp,
	//})

	//if *ctrlResp.CreateClusterOutput.Cluster.Status == "CREATING" {
	//	nodeGroupResp, err = createEksNodeGroup(svc, ctrlResp, createClusterRequest)
	//}

	// persist in db
	//err = database.AddGkeLROperation(context.Background(), createClusterRequest.Name, createClusterRequest.Name, createClusterRequest.Location)
	err = database.AddEksCluster(context.Background(), clusterId, createClusterRequest.UserId, createClusterRequest.Name,
		ctrlResp.RequestToken, ctrlResp.ClusterArn, ctrlResp.RoleArn, util.StringPointerListToEscapedCSV(ctrlResp.SubnetIds), ctrlResp.KubVersion, createClusterRequest.Location)
	//return nodeGroupResp, nil
	return resp, nil
}

func createEksControlPlane(svc *eks.EKS, id string, name string, arn string, kubVersion string, vpcConfig eks.VpcConfigRequest) (clusterCreationOutput domain.EksClusterStatus, err error) {
	input := &eks.CreateClusterInput{
		ClientRequestToken: aws.String(id),
		Name:               aws.String(name),
		ResourcesVpcConfig: &vpcConfig,
		RoleArn:            aws.String(arn),
		Version:            aws.String(kubVersion),
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

func CreateEksNodeGroup(secretId int, eksClusterContext domain.EksClusterContext) (nodeGroupResponse domain.EksNodeGroupContext, err error) {
	cred, err := iam.GetEksCredentialsForSecretId(strconv.Itoa(secretId))
	if err != nil {
		return domain.EksNodeGroupContext{}, err
	}

	sess, _ := session.NewSession(&aws.Config{
		Credentials: cred,
		Region:      aws.String(eksClusterContext.ClusterRequest.Location),
	})

	svc := eks.New(sess)
	ngResp, err := createEksNodeGroup(svc, eksClusterContext)
	if err != nil {
		return domain.EksNodeGroupContext{}, err
	}
	resp := domain.EksNodeGroupContext{
		SecretId: secretId,
		Response: *ngResp.Nodegroup,
		Region:   eksClusterContext.ClusterRequest.Location,
	}

	// persist in db
	//err = database.AddEksCluster(context.Background(), clusterId, createClusterRequest.UserId, createClusterRequest.Name, createClusterRequest.Name)
	//return nodeGroupResp, nil
	return resp, nil
}

func createEksNodeGroup(svc *eks.EKS, eksClusterContext domain.EksClusterContext) (nodeGroupResponse eks.CreateNodegroupOutput, err error) {
	groupName, _ := uuid.GenerateUUID()

	size := int64(eksClusterContext.ClusterRequest.InstanceCount)
	size2 := int64(eksClusterContext.ClusterRequest.InstanceCount) * 2

	input := &eks.CreateNodegroupInput{
		AmiType:            nil,
		ClientRequestToken: &eksClusterContext.ClusterStatus.RequestToken,
		ClusterName:        &eksClusterContext.ClusterStatus.Name,
		DiskSize:           nil,
		InstanceTypes:      eksClusterContext.ClusterRequest.MachineFamily,
		//InstanceTypes:  nil,
		Labels:         nil,
		NodeRole:       &eksClusterContext.ClusterStatus.RoleArn,
		NodegroupName:  &groupName,
		ReleaseVersion: nil,
		RemoteAccess:   nil,
		ScalingConfig: &eks.NodegroupScalingConfig{
			DesiredSize: &size,
			MaxSize:     &size2,
			MinSize:     &size,
		},
		Subnets: *eksClusterContext.ClusterStatus.SubnetIds,
		Tags:    nil,
		Version: &eksClusterContext.ClusterStatus.KubVersion,
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

func DeleteEksCluster(clusterName string, secretId string, region string) (out *eks.DeleteClusterOutput, err error) {
	//TODO; get the secret using id here
	cred, err := iam.GetEksCredentialsForSecretId(secretId)
	//if err != nil{
	//	return &eks.DeleteClusterOutput{}
	//}
	//id := "AKIA5CVBUZ342ISPRDVJ"
	//secret := "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"

	session, _ := session.NewSession(&aws.Config{
		Credentials: cred,
		Region:      aws.String(region),
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

// helper cloud services
func getVPCConfigForUSerForRegion(secretId int, region string) (vpcConfig eks.VpcConfigRequest) {
	cred, _ := iam.GetEksCredentialsForSecretId(strconv.Itoa(secretId))
	//nodeGroupResp := domain.EksClusterStatus{}
	sess, _ := session.NewSession(&aws.Config{
		Credentials: cred,
		Region:      aws.String(region),
	})

	svc := ec2.New(sess)

	result1, _ := svc.DescribeSubnets(&ec2.DescribeSubnetsInput{})

	resp := eks.VpcConfigRequest{}
	subnetsIds := make([]*string, 0)
	for _, subnet := range result1.Subnets {
		subnetsIds = append(subnetsIds, subnet.SubnetId)
	}
	resp.SetSubnetIds(subnetsIds)

	result2, _ := svc.DescribeSecurityGroups(&ec2.DescribeSecurityGroupsInput{})

	securityGroupIds := make([]*string, 0)
	for _, securityGroup := range result2.SecurityGroups {
		securityGroupIds = append(securityGroupIds, securityGroup.GroupId)
	}
	resp.SetSecurityGroupIds(securityGroupIds)
	//fmt.Printf(result.String())
	return resp
}

func generateEKSClusterCreationRequest(request *domain.ClusterOptions) *eks.CreateClusterInput {
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
