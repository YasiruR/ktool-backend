package iam

import (
	"github.com/YasiruR/ktool-backend/domain"
	oauth2 "golang.org/x/oauth2/google"
)

func TestIamPermissions(cloudSecret *domain.CloudSecret) (isValid bool, err error) {
	switch cloudSecret.ServiceProvider {
	case "Google":
		cred := oauth2.Credentials{
			ProjectID: cloudSecret.GkeProjectId,
			//JSON: cloudSecret.GkeProjectId.strin TODO: convert to json bytes
		}
		isValid, err = TestIamPermissionsGke(&cred)
	default:

	}
	if err != nil {
		return isValid, err
	}
	return isValid, nil
}
