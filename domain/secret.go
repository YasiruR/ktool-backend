package domain

import (
	"context"
	http "github.com/YasiruR/ktool-backend/http"
)

type Secret struct {
	ID         int
	Name       string
	OwnerId    int
	Provider   string
	Type       int
	CreatedOn  string
	CreatedBy  int
	ModifiedOn string
	ModifiedBy int
	Activated  bool
	Deleted    bool
	Encrypted  bool
	Tags       string
}

type GkeSecret struct {
	id                int
	Type              string
	ProjectId         string
	SecretId          string
	ProjectKeyId      string
	PrivateKey        string
	ClientMail        string
	ClientId          string
	ClientX509CertUrl string
}

type DAOResult struct {
	SecretList []Secret
	Status     int
	Message    string
	Error      error
}

type SecretDAO interface {
	AddSecret(ctx context.Context, addSecretRequest *http.AddSecretRequest) (result DAOResult)
	DeleteSecret(ctx context.Context, secretId string) (result DAOResult)
	GetAllSecretsByUser(ctx context.Context, userId string) (result DAOResult)
}
