package domain

import (
	"context"
)

type Secret struct {
	ID         int
	Name       string
	OwnerId    int
	Provider   string
	Type       int
	Key        string
	CreatedOn  string
	CreatedBy  int
	ModifiedOn string
	ModifiedBy int
	Activated  bool
	Deleted    bool
	Encrypted  bool
	Tags       string
}

type DAOResult struct {
	SecretList []Secret
	Status     int
	Message    string
	Error      error
}

type SecretDAO interface {
	AddSecret(ctx context.Context, secretName string, userId string, service string, keyType int, key string, tags string) (result DAOResult)
	DeleteSecret(ctx context.Context, secretId string) (result DAOResult)
	GetAllSecretsByUser(ctx context.Context, userId string) (result DAOResult)
}
