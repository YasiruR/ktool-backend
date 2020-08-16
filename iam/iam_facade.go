package iam

import (
	"context"
	"github.com/YasiruR/ktool-backend/domain"
	"github.com/YasiruR/ktool-backend/util"
	oauth2 "golang.org/x/oauth2/google"
)

func TestIamPermissions(cloudSecret *domain.CloudSecret) (isValid bool, err error) {
	switch cloudSecret.ServiceProvider {
	case "Google":
		credAsBytes, err := util.ConvertSecretToGKESecretBytes(*cloudSecret)
		cred, err := oauth2.CredentialsFromJSON(context.Background(), credAsBytes)
		//cred := oauth2.Credentials{
		//	ProjectID: cloudSecret.GkeProjectId,
		//	//JSON: cloudSecret.GkeProjectId.strin TODO: convert to json bytes
		//}
		if err != nil {

		}
		isValid, err = TestIamPermissionsGke(cred, credAsBytes)
	default:

	}
	if err != nil {
		return isValid, err
	}
	return isValid, nil
}
