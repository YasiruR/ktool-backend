package iam

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/YasiruR/ktool-backend/database"
	"github.com/YasiruR/ktool-backend/domain"
	"github.com/YasiruR/ktool-backend/log"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/iam"
	"net/url"
	"strconv"
)

func GetEksCredentialsForUser(userId string) (*credentials.Credentials, error) {
	ctx := context.Background()
	secretDao := database.GetSecretInternal(ctx, userId, `Amazon`, `eks-ktool`)

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

func TestIamPermissionsEks(secret *domain.CloudSecret) (bool, error) {
	sess, _ := session.NewSession(&aws.Config{
		Credentials: credentials.NewStaticCredentials(secret.EksAccessKeyId, secret.EksSecretAccessKey, ""),
	})
	svc := iam.New(sess)
	input := &iam.ListGroupsInput{}
	result, err := svc.ListGroups(input)

	if err != nil {
		return false, err
	}

	for _, group := range result.Groups {
		policyInput := &iam.ListGroupPoliciesInput{
			GroupName: group.GroupName,
		}
		result, err := svc.ListGroupPolicies(policyInput)
		if err != nil {
			return false, err
		}
		for _, name := range result.PolicyNames {
			groupPolicyInput := &iam.GetGroupPolicyInput{
				GroupName:  group.GroupName,
				PolicyName: name,
			}
			result, err := svc.GetGroupPolicy(groupPolicyInput)
			if err != nil {
				return false, err
			}
			policyDoc, _ := url.QueryUnescape(*result.PolicyDocument)
			//todo: move to somewhere
			type Statement struct {
				Effect   string   `json:"Effect"`
				Action   []string `json:"Action"`
				Resource []string `json:"Resource"`
			}
			type Policy struct {
				Version   string      `json:"Version"`
				Statement []Statement `json:"Statement"`
			}
			policy := &Policy{}
			_ = json.Unmarshal([]byte(policyDoc), policy)
			okCF := false
			okEKS := false
			for _, statement := range policy.Statement {
				if statement.Effect == "Allow" && statement.Action[0] == "cloudformation:*" {
					okCF = true
				}
				if statement.Effect == "Allow" && statement.Action[0] == "eks:*" {
					okEKS = true
				}
			}
			return okCF && okEKS, nil
		}
	}
	return false, nil
}

func GetRoleArnForEks(secretId int) (arn string, err error) {
	cred, _ := GetEksCredentialsForSecretId(strconv.Itoa(secretId))
	//nodeGroupResp := domain.EksClusterStatus{}
	sess, _ := session.NewSession(&aws.Config{
		Credentials: cred,
	})

	svc := iam.New(sess)
	req := &iam.ListRolesInput{}
	res, err := svc.ListRoles(req)
	if err != nil {
		return "", err
	}
	for _, role := range res.Roles {
		if *role.RoleName == "eksClusterRole" {
			return *role.Arn, nil
		}
	}
	return "", awserr.New("404",
		"AWS role eksClusterRole not found. Please follow 'https://docs.aws.amazon.com/eks/latest/userguide/service_IAM_role.html'",
		errors.New("Role ARN not found for role eksClusterRole"))
}
