package database

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/YasiruR/ktool-backend/domain"
	"github.com/YasiruR/ktool-backend/log"
)

func AddSecret(ctx context.Context, request domain.CloudSecret) (result domain.Result) {
	UserId := request.UserId
	var err error
	switch request.ServiceProvider {
	case "Google":
		err = AddGkeSecret(ctx, UserId, request.Name, request.Tags, request.GkeType,
			request.GkeProjectId, request.GkePrivateKeyId, request.GkePrivateKey, request.GkeClientMail,
			request.GkeClientId, request.GkeAuthUri, request.GkeTokenUri, request.GkeAuthX509CertUrl,
			request.GkeClientX509CertUrl)
	case "Amazon":
		err = AddEksSecret(ctx, UserId, request.Name, request.Tags, request.EksAccessKeyId,
			request.EksSecretAccessKey)
	case "Microsoft":
		err = AddAksSecret(ctx, UserId, request.Name, request.Tags, request.AksClientId,
			request.AksClientSecret, request.AksTenantId, request.AksSubscriptionId)
	default:
		err = nil
	}

	if err != nil {
		result.Error = err
		result.Message = "Secret addition failed"
		result.Status = -1
		return result
	}

	return GetAllSecretsByUserExternal(ctx, UserId, "all")
}

func GetSecretInternal(ctx context.Context, OwnerId string, Provider string, SecretName string) (result domain.DAOResult) {

	switch Provider {
	case "Google":
		result = GetGkeSecret(ctx, OwnerId, SecretName)
	case "Amazon":
		result = GetEksSecret(ctx, OwnerId, SecretName)
	case "Microsoft":
		result = GetAksSecret(ctx, OwnerId, SecretName)
	default:

	}
	switch result.Error {
	case nil:
		log.Logger.InfoContext(ctx, "get secret query success")
	case sql.ErrNoRows:
		log.Logger.InfoContext(ctx, "no secrets found for userId %s", OwnerId)
		result.Status = -1
		result.Message = "no secrets found"
		return result
	default:
		log.Logger.InfoContext(ctx, "unhandled error occurred while fetching records for userId %s", OwnerId)
		result.Status = -1
		result.Message = "unhandled error occurred from db"
		return result
	}
	log.Logger.TraceContext(ctx, "get all secrets db query was successful")
	result.Status = 0
	result.Message = "Success"
	return result
}

func GetAksSecret(ctx context.Context, OwnerId string, SecretName string) (result domain.DAOResult) {

	query := "SELECT id, aksClientId, aksClientSecret, aksTenantId, aksSubscriptionId FROM " + cloudSecretTable +
		" WHERE OwnerId = " + OwnerId + " AND Name = '" + SecretName + "';"

	rows, err := Db.Query(query)

	if err != nil {
		result.Error = err
		return result
	}

	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(&result.Secret.ID, &result.Secret.AksClientId, &result.Secret.AksClientSecret, &result.Secret.AksTenantId,
			&result.Secret.AksSubscriptionId)
		if err != nil {
			log.Logger.ErrorContext(ctx, "scanning rows in secret table failed", err)
			result.Error = err
			return result
		}
	}
	return result
}

func GetGkeSecret(ctx context.Context, OwnerId string, SecretName string) (result domain.DAOResult) {

	query := "SELECT id, secretType, projectId, privateKeyId, privateKey, clientEmail, clientId, authUri, tokenUri," +
		" authCertUrl, clientCertUrl, secretKeyId, secretAccessKey FROM " + cloudSecretTable + " WHERE OwnerId = " +
		OwnerId + " AND Name = '" + SecretName + "';"

	rows, err := Db.Query(query)

	if err != nil {
		result.Error = err
		return result
	}

	defer rows.Close()

	for rows.Next() {

		err = rows.Scan(&result.Secret.ID, &result.Secret.GkeType, &result.Secret.GkeProjectId, &result.Secret.GkePrivateKeyId,
			&result.Secret.GkePrivateKey, &result.Secret.GkeClientMail, &result.Secret.GkeClientId, &result.Secret.GkeAuthUri,
			&result.Secret.GkeTokenUri, &result.Secret.GkeAuthX509CertUrl, &result.Secret.GkeClientX509CertUrl)
		if err != nil {
			log.Logger.ErrorContext(ctx, "scanning rows in secret table failed", err)
			result.Error = err
			return result
		}
	}

	return result
}

func GetEksSecret(ctx context.Context, OwnerId string, SecretName string) (result domain.DAOResult) {

	query := "SELECT id, accessKeyId, secretAccessKey FROM " + cloudSecretTable + " WHERE OwnerId = " +
		OwnerId + " AND Name = '" + SecretName + "';"

	rows, err := Db.Query(query)

	if err != nil {
		result.Error = err
		return result
	}

	defer rows.Close()

	for rows.Next() {

		err = rows.Scan(&result.Secret.ID, &result.Secret.EksAccessKeyId, &result.Secret.EksSecretAccessKey)
		if err != nil {
			log.Logger.ErrorContext(ctx, "scanning rows in secret table failed", err)
			result.Error = err
			return result
		}
	}
	return result
}

func AddEksSecret(ctx context.Context, UserId string, SecretName string, Tags string,
	EksAccessKeyId string, EksSecretAccessKey string) (err error) {
	query := "INSERT INTO kdb.cloud_secret (ownerId, name, provider, tags, createdBy, createdOn, modifiedBy," +
		" modifiedOn, activated, deleted, accessKeyId, secretAccessKey) VALUES(" +
		UserId + ",'" + SecretName + "','" + "amazon" + "','" + Tags + "'," +
		UserId + ", CURRENT_TIMESTAMP, '', '', 0, 0,'" + EksAccessKeyId + "','" + EksSecretAccessKey + "')"

	insert, err := Db.Query(query)

	if err != nil {
		log.Logger.ErrorContext(ctx, fmt.Sprintf("insert to %s table failed", cloudSecretTable), err)
		return err
	}

	defer insert.Close()
	log.Logger.TraceContext(ctx, "successfully added a new key ", SecretName)
	return nil
}

func AddAksSecret(ctx context.Context, UserId string, SecretName string, Tags string,
	AksClientId string, AksClientSecret string, AksTenantId string, AksSubscriptionId string) (err error) {

	//TODO: call a stored procedure
	query := "INSERT INTO kdb.cloud_secret (ownerId, name, provider, tags, createdBy, createdOn, modifiedBy," +
		" modifiedOn, activated, deleted, aksClientId, aksClientSecret, aksTenantId, aksSubscriptionId) VALUES(" +
		UserId + ",'" + SecretName + "','" + "microsoft" + "','" + Tags + "'," + UserId +
		", CURRENT_TIMESTAMP, '', '', 0, 0,'" + AksClientId + "','" + AksClientSecret + "','" + AksTenantId + "','" +
		AksSubscriptionId + "')"

	insert, err := Db.Query(query)

	if err != nil {
		log.Logger.ErrorContext(ctx, fmt.Sprintf("insert to %s table failed", cloudSecretTable), err)
		return err
	}

	defer insert.Close()
	log.Logger.TraceContext(ctx, "successfully added a new key ", SecretName)
	return nil
}

// valid for adding gke secrets
func AddGkeSecret(ctx context.Context, UserId string, SecretName string, Tags string, GkeType string, GkeProjectId string,
	GkePrivateKeyId string, GkePrivateKey string, GkeClientMail string, GkeClientId string, GkeAuthUri string,
	GkeTokenUri string, GkeAuthCertUrl string, GkeClientCertUrl string) (err error) {

	//TODO: call a stored procedure
	query := "INSERT INTO kdb.cloud_secret (ownerId, name, provider, tags, createdBy, createdOn, modifiedBy," +
		" modifiedOn, activated, deleted, secretType, projectId, privateKeyId, privateKey, clientEmail, clientId," +
		" authUri, tokenUri, authCertUrl, clientCertUrl) VALUES(" + UserId + ",'" + SecretName + "','" + "google" +
		"','" + Tags + "'," + UserId + ", CURRENT_TIMESTAMP, '', '', 0, 0,'" + GkeType + "','" + GkeProjectId + "','" +
		GkePrivateKeyId + "','" + GkePrivateKey + "','" + GkeClientMail + "','" + GkeClientId + "','" + GkeAuthUri +
		"','" + GkeTokenUri + "','" + GkeAuthCertUrl + "','" + GkeClientCertUrl + "')"

	insert, err := Db.Query(query)

	if err != nil {
		log.Logger.ErrorContext(ctx, fmt.Sprintf("insert to %s table failed", cloudSecretTable), err)
		return err
	}

	defer insert.Close()
	log.Logger.TraceContext(ctx, "successfully added a new key ", SecretName)
	return nil
}

func GetAllSecretsByUserExternal(ctx context.Context, OwnerId string, ServiceProvider string) (result domain.Result) {
	query := "SELECT id, ownerId, name, provider, tags, createdBy, createdOn, modifiedBy, modifiedOn, activated, deleted FROM " + cloudSecretTable + " WHERE OwnerId = " + OwnerId + ";"

	rows, err := Db.Query(query)

	switch err {
	case nil:
		log.Logger.InfoContext(ctx, "get secret query success")
	case sql.ErrNoRows:
		log.Logger.InfoContext(ctx, "no secrets found for userId %s", OwnerId)
		result.Error = err
		result.Status = -1
		result.Message = "no secrets found for userId " + OwnerId
		return result
	default:
		log.Logger.InfoContext(ctx, "unhandled error occurred while fetching records for userId %s", OwnerId)
		result.Error = err
		result.Status = -1
		result.Message = "unhandled error occurred from db"
		return result
	}

	defer rows.Close()
	secretList := make([]domain.Secret, 0)

	for rows.Next() {
		secret := domain.Secret{}

		err = rows.Scan(&secret.ID, &secret.OwnerId, &secret.Name, &secret.Provider, &secret.Tags, &secret.CreatedBy,
			&secret.CreatedOn, &secret.ModifiedBy, &secret.ModifiedOn, &secret.Activated, &secret.Deleted)
		if err != nil {
			log.Logger.ErrorContext(ctx, "scanning rows in secret table failed", err)
			result.Error = err
			result.Status = -1
			result.Message = "scanning rows in secret table failed"
			return result
		}
		secretList = append(secretList, secret)
	}

	log.Logger.TraceContext(ctx, "get all secrets db query was successful")
	result.SecretList = secretList
	result.Status = 0
	result.Message = "Success"
	return result
}

//TODO: implement delete query
