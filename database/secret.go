package database

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/YasiruR/ktool-backend/domain"
	"github.com/YasiruR/ktool-backend/log"
)

func GetAllSecretsByUserExternal(ctx context.Context, OwnerId string, ServiceProvider string) (result domain.Result) {
	query := "SELECT s.id, s.ownerId, s.name, s.provider, s.tags, s.createdBy, s.createdOn, u.username as modifiedBy, " +
		"s.modifiedOn, s.activated, s.deleted FROM `" + cloudSecretTable + "` s, `" + userTable + "` u WHERE s.OwnerId = " +
		OwnerId + " AND s.modifiedBy = u.id AND s.deleted = 0;"

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

func GetSecretExternal(ctx context.Context, SecretId string, Provider string) (result domain.DAOResult) {
	query := ""
	switch Provider {
	case "google":
		query = "SELECT id, name, provider, gkeSecretType, gkeProjectId, gkePrivateKeyId, gkePrivateKey, gkeClientEmail, " +
			"gkeClientId, gkeAuthUri, gkeTokenUri, gkeAuthCertUrl, gkeClientCertUrl FROM `" + cloudSecretTable + "` WHERE Id = " +
			SecretId + " AND deleted = 0;"
	case "amazon":
		query = "SELECT id, name, provider, eksAccessKeyId, eksSecretAccessKey FROM `" + cloudSecretTable + "` WHERE Id = " +
			SecretId + " AND deleted = 0;"
	case "microsoft":
		query = "SELECT id, name, provider, aksClientId, aksClientSecret, aksTenantId, aksSubscriptionId FROM `" +
			cloudSecretTable + "` WHERE Id = " + SecretId + " AND deleted = 0;"
	default:
		return
	}
	//query := "SELECT id, name, provider, gkeSecretType, gkeProjectId, gkePrivateKeyId, gkePrivateKey, gkeClientEmail, " +
	//	"gkeClientId, gkeAuthUri, gkeTokenUri, gkeAuthCertUrl, gkeClientCertUrl, aksClientId, aksClientSecret, aksTenantId, " +
	//	"aksSubscriptionId, eksAccessKeyId, eksSecretAccessKey FROM `" + cloudSecretTable + "` WHERE Id = " +
	//	SecretId + " AND deleted = 0;"

	rows, err := Db.Query(query)

	switch err {
	case nil:
		log.Logger.InfoContext(ctx, "get secret query success")
	case sql.ErrNoRows:
		log.Logger.InfoContext(ctx, "no secrets found for secretId %s", SecretId)
		result.Error = err
		result.Status = -1
		result.Message = "no secrets found for secretId " + SecretId
		return result
	default:
		log.Logger.InfoContext(ctx, "unhandled error occurred while fetching records for secretId %s", SecretId)
		result.Error = err
		result.Status = -1
		result.Message = "unhandled error occurred from db"
		return result
	}

	defer rows.Close()

	for rows.Next() {
		switch Provider {
		case "google":
			err = rows.Scan(&result.Secret.ID, &result.Secret.Name, &result.Secret.ServiceProvider, &result.Secret.GkeType, &result.Secret.GkeProjectId,
				&result.Secret.GkePrivateKeyId, &result.Secret.GkePrivateKey, &result.Secret.GkeClientMail, &result.Secret.GkeClientId,
				&result.Secret.GkeAuthUri, &result.Secret.GkeTokenUri, &result.Secret.GkeAuthX509CertUrl,
				&result.Secret.GkeClientX509CertUrl)
		case "amazon":
			err = rows.Scan(&result.Secret.ID, &result.Secret.Name, &result.Secret.ServiceProvider, &result.Secret.EksAccessKeyId, &result.Secret.EksSecretAccessKey)
		case "microsoft":
			err = rows.Scan(&result.Secret.ID, &result.Secret.Name, &result.Secret.ServiceProvider, &result.Secret.AksClientId, &result.Secret.AksClientSecret,
				&result.Secret.AksTenantId, &result.Secret.AksSubscriptionId)
		default:
			return
		}
		//err = rows.Scan(&secret.ID, &secret.Name, &secret.ServiceProvider, &result.Secret.GkeType, &result.Secret.GkeProjectId,
		//	&result.Secret.GkePrivateKeyId, &result.Secret.GkePrivateKey, &result.Secret.GkeClientMail, &result.Secret.GkeClientId,
		//	&result.Secret.GkeAuthUri, &result.Secret.GkeTokenUri, &result.Secret.GkeAuthX509CertUrl,
		//	&result.Secret.GkeClientX509CertUrl, &result.Secret.AksClientId, &result.Secret.AksClientSecret, &result.Secret.AksTenantId,
		//	&result.Secret.AksSubscriptionId, &result.Secret.EksAccessKeyId, &result.Secret.EksSecretAccessKey)
		if err != nil {
			log.Logger.ErrorContext(ctx, "scanning rows in secret table failed", err)
			result.Error = err
			result.Status = -1
			result.Message = "scanning rows in secret table failed"
			return result
		}
	}

	log.Logger.TraceContext(ctx, "get secret db query was successful")
	//result.Secret = secret
	result.Status = 0
	result.Message = "Success"
	return result
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

func GetSecretByIdInternal(ctx context.Context, secretId string, provider string) (result domain.DAOResult) {
	result = GetSecretById(ctx, secretId, provider)
	switch result.Error {
	case nil:
		log.Logger.InfoContext(ctx, "get secret query success")
	case sql.ErrNoRows:
		log.Logger.InfoContext(ctx, "no secrets found for secretId %s", secretId)
		result.Status = -1
		result.Message = "no secrets found"
		return result
	default:
		log.Logger.InfoContext(ctx, "unhandled error occurred while fetching records for secretId %s", secretId)
		result.Status = -1
		result.Message = "unhandled error occurred from db"
		return result
	}
	log.Logger.TraceContext(ctx, "get all secrets db query was successful")
	result.Status = 0
	result.Message = "Success"
	return result
}

func GetSecretById(ctx context.Context, secretId string, provider string) (result domain.DAOResult) {
	query := ""
	switch provider {
	case "google":
		query = "SELECT id, gkeSecretType, gkeProjectId, gkePrivateKeyId, gkePrivateKey, gkeClientEmail, " +
			"gkeClientId, gkeAuthUri, gkeTokenUri, gkeAuthCertUrl, gkeClientCertUrl FROM " + cloudSecretTable +
			" WHERE id = " + secretId + ";"
	case "amazon":
		query = "SELECT id, eksAccessKeyId, eksSecretAccessKey FROM " + cloudSecretTable +
			" WHERE id = " + secretId + ";"
	default:
		query = ""
	}

	rows, err := Db.Query(query)

	if err != nil {
		result.Error = err
		return result
	}

	defer rows.Close()

	for rows.Next() {
		if provider == "google" {
			err = rows.Scan(&result.Secret.ID, &result.Secret.GkeType, &result.Secret.GkeProjectId, &result.Secret.GkePrivateKeyId,
				&result.Secret.GkePrivateKey, &result.Secret.GkeClientMail, &result.Secret.GkeClientId, &result.Secret.GkeAuthUri,
				&result.Secret.GkeTokenUri, &result.Secret.GkeAuthX509CertUrl, &result.Secret.GkeClientX509CertUrl)
		} else if provider == "amazon" {
			err = rows.Scan(&result.Secret.ID, &result.Secret.EksAccessKeyId, &result.Secret.EksSecretAccessKey)
		}
		if err != nil {
			log.Logger.ErrorContext(ctx, "scanning rows in secret table failed", err)
			result.Error = err
			return result
		}
	}
	return result
}

func GetAksSecret(ctx context.Context, OwnerId string, SecretName string) (result domain.DAOResult) {

	query := "SELECT id, aksClientId, aksClientSecret, aksTenantId, aksSubscriptionId FROM " + cloudSecretTable +
		" WHERE OwnerId = " + OwnerId + " AND Provider = 'microsoft' AND Name = '" + SecretName + "';"

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

	query := "SELECT id, gkeSecretType, gkeProjectId, gkePrivateKeyId, gkePrivateKey, gkeClientEmail, " +
		"gkeClientId, gkeAuthUri, gkeTokenUri, gkeAuthCertUrl, gkeClientCertUrl FROM " + cloudSecretTable +
		" WHERE OwnerId = " + OwnerId + " AND Provider = 'google' AND Name = '" + SecretName + "';"

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

	query := "SELECT id, eksAccessKeyId, eksSecretAccessKey FROM " + cloudSecretTable + " WHERE OwnerId = " +
		OwnerId + " AND Provider = 'amazon' AND Name = '" + SecretName + "';"

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

func UpdateSecret(ctx context.Context, request *domain.CloudSecret) (result domain.Result) {
	OwnerId := request.OwnerId
	var err error
	switch request.ServiceProvider {
	case "Google":
		err = UpdateGkeSecret(ctx, OwnerId, request.Name, request.GkeType, request.GkeProjectId, request.GkePrivateKeyId,
			request.GkePrivateKey, request.GkeClientMail, request.GkeClientId, request.GkeAuthUri, request.GkeTokenUri,
			request.GkeAuthX509CertUrl, request.GkeClientX509CertUrl)
	case "Amazon":
		err = UpdateEksSecret(ctx, OwnerId, request.Name, request.EksAccessKeyId, request.EksSecretAccessKey)
	case "Microsoft":
		err = UpdateAksSecret(ctx, OwnerId, request.Name, request.AksClientId, request.AksClientSecret, request.AksTenantId,
			request.AksSubscriptionId)
	default:
		err = nil
	}

	if err != nil {
		result.Error = err
		result.Message = "Secret update failed"
		result.Status = -1
		return result
	}

	return GetAllSecretsByUserExternal(ctx, OwnerId, "all")
}

func UpdateEksSecret(ctx context.Context, UserId string, SecretName string, EksAccessKeyId string,
	EksSecretAccessKey string) (err error) {

	query := "UPDATE " + cloudSecretTable + " SET " +
		"eksAccessKeyId= '" + EksAccessKeyId + "', " +
		"eksSecretAccessKey='" + EksSecretAccessKey +
		"' WHERE OwnerId=" + UserId + " AND Name='" + SecretName + "' AND Provider='amazon';"

	insert, err := Db.Query(query)

	if err != nil {
		log.Logger.ErrorContext(ctx, fmt.Sprintf("update to %s table failed", cloudSecretTable), err)
		return err
	}

	defer insert.Close()
	log.Logger.TraceContext(ctx, "successfully updated key ", SecretName)
	return nil
}

func UpdateAksSecret(ctx context.Context, UserId string, SecretName string, AksClientId string, AksClientSecret string,
	AksTenantId string, AksSubscriptionId string) (err error) {
	query := "UPDATE " + cloudSecretTable + " SET " +
		"aksClientId= '" + AksClientId + "', " +
		"aksClientSecret='" + AksClientSecret + "', " +
		"aksTenantId='" + AksTenantId + "', " +
		"aksSubscriptionId='" + AksSubscriptionId +
		"' WHERE OwnerId=" + UserId + " AND Name='" + SecretName + "' AND Provider='microsoft';"

	insert, err := Db.Query(query)

	if err != nil {
		log.Logger.ErrorContext(ctx, fmt.Sprintf("update to %s table failed", cloudSecretTable), err)
		return err
	}

	defer insert.Close()
	log.Logger.TraceContext(ctx, "successfully updated key ", SecretName)
	return nil
}

// valid for adding gke secrets
func UpdateGkeSecret(ctx context.Context, UserId string, SecretName string, GkeType string, GkeProjectId string,
	GkePrivateKeyId string, GkePrivateKey string, GkeClientMail string, GkeClientId string, GkeAuthUri string,
	GkeTokenUri string, GkeAuthCertUrl string, GkeClientCertUrl string) (err error) {
	query := "UPDATE " + cloudSecretTable + " SET " +
		//"aksClientId= " + GkeType + ", " +
		"GkeProjectId='" + GkeProjectId + "', " +
		"GkePrivateKeyId='" + GkePrivateKeyId + "', " +
		"GkePrivateKey='" + GkePrivateKey + "', " +
		"GkeClientEmail='" + GkeClientMail + "', " +
		"GkeClientId='" + GkeClientId + "', " +
		"GkeAuthUri='" + GkeAuthUri + "', " +
		"GkeTokenUri='" + GkeTokenUri + "', " +
		"GkeAuthCertUrl='" + GkeAuthCertUrl + "', " +
		"GkeAuthCertUrl='" + GkeAuthCertUrl +
		"' WHERE OwnerId=" + UserId + " AND Name='" + SecretName + "' AND Provider='google';"

	insert, err := Db.Query(query)

	if err != nil {
		log.Logger.ErrorContext(ctx, fmt.Sprintf("update to %s table failed", cloudSecretTable), err)
		return err
	}

	defer insert.Close()
	log.Logger.TraceContext(ctx, "successfully updated key ", SecretName)
	return nil
}

func AddSecret(ctx context.Context, request *domain.CloudSecret) (result domain.Result) {
	OwnerId := request.OwnerId
	var err error
	switch request.ServiceProvider {
	case "Google":
		err = AddGkeSecret(ctx, OwnerId, request.Name, request.Tags, request.GkeType,
			request.GkeProjectId, request.GkePrivateKeyId, request.GkePrivateKey, request.GkeClientMail,
			request.GkeClientId, request.GkeAuthUri, request.GkeTokenUri, request.GkeAuthX509CertUrl,
			request.GkeClientX509CertUrl)
	case "Amazon":
		err = AddEksSecret(ctx, OwnerId, request.Name, request.Tags, request.EksAccessKeyId,
			request.EksSecretAccessKey)
	case "Microsoft":
		err = AddAksSecret(ctx, OwnerId, request.Name, request.Tags, request.AksClientId,
			request.AksClientSecret, request.AksTenantId, request.AksSubscriptionId)
	default:
		err = nil
	}

	if err != nil {
		result.ErrorMsg = err.Error()
		result.Error = err
		result.Message = "Secret addition failed"
		result.Status = -1
		return result
	}

	return GetAllSecretsByUserExternal(ctx, OwnerId, "all")
}

func AddEksSecret(ctx context.Context, UserId string, SecretName string, Tags string,
	EksAccessKeyId string, EksSecretAccessKey string) (err error) {
	//TODO: validate req params
	//TODO: call a stored procedure
	query := "INSERT INTO kdb." + cloudSecretTable + "(ownerId, name, provider, tags, createdBy, createdOn, modifiedBy," +
		" modifiedOn, activated, deleted, eksAccessKeyId, eksSecretAccessKey) VALUES(" +
		UserId + ",'" + SecretName + "','" + "amazon" + "','" + Tags + "'," +
		UserId + ", CURRENT_TIMESTAMP, " + UserId + ", CURRENT_TIMESTAMP, 0, 0,'" + EksAccessKeyId + "','" + EksSecretAccessKey + "')"

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
	//TODO: validate req params
	//TODO: call a stored procedure
	query := "INSERT INTO kdb." + cloudSecretTable + "(ownerId, name, provider, tags, createdBy, createdOn, modifiedBy," +
		" modifiedOn, activated, deleted, aksClientId, aksClientSecret, aksTenantId, aksSubscriptionId) VALUES(" +
		UserId + ",'" + SecretName + "','" + "microsoft" + "','" + Tags + "'," + UserId +
		", CURRENT_TIMESTAMP, " + UserId + ",CURRENT_TIMESTAMP, 0, 0,'" + AksClientId + "','" + AksClientSecret + "','" + AksTenantId + "','" +
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
	//TODO: validate req params
	//TODO: call a stored procedure
	query := "INSERT INTO kdb." + cloudSecretTable + "(ownerId, name, provider, tags, createdBy, createdOn, modifiedBy," +
		" modifiedOn, activated, deleted, gkeSecretType, gkeProjectId, gkePrivateKeyId, gkePrivateKey, gkeClientEmail, " +
		"gkeClientId, gkeAuthUri, gkeTokenUri, gkeAuthCertUrl, gkeClientCertUrl) VALUES(" + UserId + ",'" + SecretName +
		"','" + "google" + "','" + Tags + "'," + UserId + ", CURRENT_TIMESTAMP, " + UserId + ",CURRENT_TIMESTAMP, 0, 0,'" + GkeType + "','" +
		GkeProjectId + "','" + GkePrivateKeyId + "','" + GkePrivateKey + "','" + GkeClientMail + "','" + GkeClientId +
		"','" + GkeAuthUri + "','" + GkeTokenUri + "','" + GkeAuthCertUrl + "','" + GkeClientCertUrl + "')"

	insert, err := Db.Query(query)

	if err != nil {
		log.Logger.ErrorContext(ctx, fmt.Sprintf("insert to %s table failed", cloudSecretTable), err)
		return err
	}

	defer insert.Close()
	log.Logger.TraceContext(ctx, "successfully added a new key ", SecretName)
	return nil
}

func DeleteSecret(ctx context.Context, secretId string) (result bool, err error) {
	query := "UPDATE " + cloudSecretTable + " SET deleted=1 WHERE id=" + secretId + ";"

	_, err = Db.Query(query)

	switch err {
	case nil:
		log.Logger.InfoContext(ctx, "delete secret query success")
		return true, nil
	case sql.ErrNoRows:
		log.Logger.InfoContext(ctx, "delete secret failed. no secrets found")
		return false, err
	default:
		log.Logger.InfoContext(ctx, "unhandled error occurred while deleting secretId %s", secretId)
		return false, err
	}
}
