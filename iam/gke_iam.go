package iam

import (
	adminpb "cloud.google.com/go/iam"
	admin "cloud.google.com/go/iam/admin/apiv1"
	"context"
	"encoding/json"
	"github.com/YasiruR/ktool-backend/database"
	"github.com/YasiruR/ktool-backend/domain"
	"github.com/YasiruR/ktool-backend/log"
	"google.golang.org/api/option"
	iampb "google.golang.org/genproto/googleapis/iam/v1"
)

//func main() {
//	ctx := context.Background()
//	res, err := validateGKESecret("1")
//	if err != nil {
//		log.Logger.ErrorContext(ctx, "Could not retrieve cluster list")
//		return
//	}
//	log.Logger.InfoContext(ctx, "Successfully retrieved cluster list from GKE")
//	fmt.Println(res)
//}

func TestIamPermissions(userId string) {
	//ctx := context.Background()
	//b, cred, err := GetGkeCredentialsForUser(userId)
	////log.Logger.Info(cred)
	//if err != nil {
	//	return nil, err
	//}
	//c, err := admin.NewIamClient(ctx, option.WithCredentialsJSON(b))
	//if err != nil {
	//	// TODO: Handle error.
	//}
	//
	////req := &adminpb.ListRolesRequest{
	////	// TODO: Fill request struct fields.
	////	Parent: "projects/ktool-280018",
	////	View: 1,
	////}
	//req := &iampb.GetIamPolicyRequest{
	//	// TODO: Fill request struct fields.
	//	Resource: "projects/" + cred.ProjectId + "/serviceAccounts/" + cred.ClientId,
	//}
	//resp, err := c.GetIamPolicy(ctx, req)
	//if err != nil {
	//	// TODO: Handle error.
	//}
	//// TODO: Use resp.
	//return resp, nil
}

func GetServiceAccountIamPolicies(userId string) (*adminpb.Policy, error) {
	ctx := context.Background()
	b, cred, err := GetGkeCredentialsForUser(userId)
	//log.Logger.Info(cred)
	if err != nil {
		return nil, err
	}
	c, err := admin.NewIamClient(ctx, option.WithCredentialsJSON(b))
	if err != nil {
		// TODO: Handle error.
	}

	//req := &adminpb.ListRolesRequest{
	//	// TODO: Fill request struct fields.
	//	Parent: "projects/ktool-280018",
	//	View: 1,
	//}
	req := &iampb.GetIamPolicyRequest{
		// TODO: Fill request struct fields.
		Resource: "projects/" + cred.ProjectId + "/serviceAccounts/" + cred.ClientId,
	}
	resp, err := c.GetIamPolicy(ctx, req)
	if err != nil {
		// TODO: Handle error.
	}
	// TODO: Use resp.
	return resp, nil
}

func GetGkeCredentialsForUser(userId string) ([]byte, domain.GkeSecret, error) {
	ctx := context.Background()
	secretDao := database.GetSecretInternal(ctx, userId, `Google`, `ktool-gke`)

	if err := secretDao.Error; err != nil {
		log.Logger.ErrorContext(ctx, "Error occurred while fetching eks secret for client %s", userId)
		return nil, domain.GkeSecret{}, err
	}
	cred := domain.GkeSecret{
		Type:              secretDao.Secret.GkeType,
		ProjectId:         secretDao.Secret.GkeProjectId,
		PrivateKeyId:      secretDao.Secret.GkePrivateKeyId,
		PrivateKey:        secretDao.Secret.GkePrivateKey,
		ClientMail:        secretDao.Secret.GkeClientMail,
		ClientId:          secretDao.Secret.GkeClientId,
		AuthUri:           secretDao.Secret.GkeAuthUri,
		TokenUri:          secretDao.Secret.GkeTokenUri,
		AuthX509CertUrl:   secretDao.Secret.GkeAuthX509CertUrl,
		ClientX509CertUrl: secretDao.Secret.GkeClientX509CertUrl,
	}
	bytes, err := json.Marshal(&cred)
	if err != nil {
		log.Logger.ErrorContext(ctx, "Could not marshall gke credentials for user %s", userId)
		return nil, cred, err
	}
	return bytes, cred, nil
}
