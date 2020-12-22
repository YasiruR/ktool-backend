package iam

import (
	"context"
	"github.com/YasiruR/ktool-backend/database"
	"github.com/YasiruR/ktool-backend/domain"
)

func GetAksCredentialsForSecretId(secretId string) (cred *domain.AksSecret, err error) {
	resp := database.GetSecretById(context.Background(), secretId, "microsoft")
	if resp.Error != nil {
		return nil, resp.Error
	}
	return &domain.AksSecret{
		ClientId:       resp.Secret.AksClientId,
		ClientSecret:   resp.Secret.AksClientSecret,
		TenantId:       resp.Secret.AksTenantId,
		SubscriptionId: resp.Secret.AksSubscriptionId,
	}, nil
}
