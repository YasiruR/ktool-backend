package database

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/YasiruR/ktool-backend/domain"
	"github.com/YasiruR/ktool-backend/log"
)

// valid for adding gke secrets
func AddSecret(ctx context.Context, UserId string, SecretName string, ServiceProvider string, Tags string, GkeType string, GkeProjectId string, GkePrivateKeyId string, GkePrivateKey string, GkeClientMail string, GkeClientId string, GkeAuthUri string, GkeTokenUri string, GkeAuthCertUrl string, GkeClientCertUrl string) (result domain.Result) {
	//TODO: call a stored procedure

	query := ""
	switch ServiceProvider {
	case "Google":
		query = "INSERT INTO kdb.cloud_secret (ownerId, name, provider, tags, createdBy, createdOn, modifiedBy," +
			" modifiedOn, activated, deleted, secretType, projectId, privateKeyId, privateKey, clientEmail, clientId," +
			" authUri, tokenUri, authCertUrl, clientCertUrl) VALUES(" +
			UserId + ",'" + SecretName + "','" + ServiceProvider + "','" + Tags + "'," +
			UserId + ", CURRENT_TIMESTAMP, '', '', 0, 0,'" + GkeType + "','" + GkeProjectId +
			"','" + GkePrivateKeyId + "','" + GkePrivateKey + "','" + GkeClientMail + "','" +
			GkeClientId + "','" + GkeAuthUri + "','" + GkeTokenUri + "','" +
			GkeAuthCertUrl + "','" + GkeClientCertUrl + "')"

	default:
		query = ""
	}

	insert, err := Db.Query(query)

	if err != nil {
		log.Logger.ErrorContext(ctx, fmt.Sprintf("insert to %s table failed", cloudSecretTable), err)
		return domain.Result{
			SecretList: make([]domain.Secret, 0),
			Status:     1,
			Message:    "insert failed",
			Error:      nil,
		}
	}

	defer insert.Close()
	log.Logger.TraceContext(ctx, "successfully added a new key", SecretName)

	return domain.Result{
		SecretList: []domain.Secret{
			domain.Secret{
				ID:       1,
				Name:     SecretName,
				OwnerId:  UserId,
				Provider: ServiceProvider,
				Type:     GkeType,
				//CreatedOn : ,
				CreatedBy: UserId,
				//ModifiedOn : ,
				//ModifiedBy : ,
				Activated: false,
				Deleted:   false,
				Tags:      Tags,
			},
		},
		Status:  0,
		Message: "insert success",
		Error:   nil,
	}
}

func GetAllSecretsByUserInternal(ctx context.Context, OwnerId string, ServiceProvider string) (result domain.DAOResult) {
	query := "SELECT id, secretType, projectId, privateKeyId, privateKey, clientEmail, clientId, authUri, tokenUri, authCertUrl, clientCertUrl FROM " + cloudSecretTable + " WHERE OwnerId = " + OwnerId + ";"

	rows, err := Db.Query(query)

	switch err {
	case nil:
		log.Logger.InfoContext(ctx, "get secret query success")
	case sql.ErrNoRows:
		log.Logger.InfoContext(ctx, "no secrets found for userId %s", OwnerId)
		return domain.DAOResult{
			SecretList: make([]domain.CloudSecret, 0),
			Status:     -1,
			Message:    "no secrets found for userId",
			Error:      err,
		}
	default:
		log.Logger.InfoContext(ctx, "unhandled error occured while fetching records for userId %s", OwnerId)
		return domain.DAOResult{
			SecretList: make([]domain.CloudSecret, 0),
			Status:     -1,
			Message:    "unhandled error occured",
			Error:      err,
		}
	}

	defer rows.Close()
	secretList := make([]domain.CloudSecret, 0)

	for rows.Next() {
		secret := domain.CloudSecret{}

		err = rows.Scan(&secret.ID, &secret.Type, &secret.ProjectId, &secret.PrivateKeyId, &secret.PrivateKey, &secret.ClientMail,
			&secret.ClientId, &secret.AuthUri, &secret.TokenUri, &secret.AuthX509CertUrl, &secret.ClientX509CertUrl)
		if err != nil {
			log.Logger.ErrorContext(ctx, "scanning rows in secret table failed", err)
			return domain.DAOResult{
				SecretList: secretList,
				Status:     -1,
				Message:    "row scanning failed",
				Error:      err,
			}
		}
		secretList = append(secretList, secret)
	}

	log.Logger.TraceContext(ctx, "get all secrets db query was successful")
	return domain.DAOResult{
		SecretList: secretList,
		Status:     1,
		Message:    "get secrets request successful",
	}
}

func GetAllSecretsByUserExternal(ctx context.Context, OwnerId string, ServiceProvider string) (result domain.DAOResult) {
	query := "SELECT id, ownerId, name, provider, tags, createdBy, createdOn, modifiedBy, modifiedOn, activated, deleted FROM " + cloudSecretTable + " WHERE OwnerId = " + OwnerId + ";"

	rows, err := Db.Query(query)

	switch err {
	case nil:
		log.Logger.InfoContext(ctx, "get secret query success")
	case sql.ErrNoRows:
		log.Logger.InfoContext(ctx, "no secrets found for userId %s", OwnerId)
		return domain.DAOResult{
			SecretList: make([]domain.CloudSecret, 0),
			Status:     -1,
			Message:    "no secrets found for userId",
			Error:      err,
		}
	default:
		log.Logger.InfoContext(ctx, "unhandled error occured while fetching records for userId %s", OwnerId)
		return domain.DAOResult{
			SecretList: make([]domain.CloudSecret, 0),
			Status:     -1,
			Message:    "unhandled error occured",
			Error:      err,
		}
	}

	defer rows.Close()
	secretList := make([]domain.CloudSecret, 0)

	for rows.Next() {
		secret := domain.CloudSecret{}

		err = rows.Scan(&secret.ID, &secret.OwnerId, &secret.Name, &secret.Provider, &secret.Tags, &secret.CreatedBy,
			&secret.CreatedOn, &secret.ModifiedBy, &secret.ModifiedOn, &secret.Activated, &secret.Deleted)
		if err != nil {
			log.Logger.ErrorContext(ctx, "scanning rows in secret table failed", err)
			return domain.DAOResult{
				SecretList: secretList,
				Status:     -1,
				Message:    "row scanning failed",
				Error:      err,
			}
		}
		secretList = append(secretList, secret)
	}

	log.Logger.TraceContext(ctx, "get all secrets db query was successful")
	return domain.DAOResult{
		SecretList: secretList,
		Status:     1,
		Message:    "get secrets request successful",
	}
}

//TODO: implement delete query
