package database

import (
	"context"
	"fmt"
	"github.com/YasiruR/ktool-backend/domain"
	"github.com/YasiruR/ktool-backend/http"
	"github.com/YasiruR/ktool-backend/log"
	"strconv"
)

func (dao *domain.SecretDAO) AddSecret(ctx context.Context, request *http.AddSecretRequest) (result DAOResult) {
	//TODO: call a stored procedure
	query := "INSERT INTO " + secretTable + " ( Name, OwnerId, Provider, Type, CreatedBy, ModifiedBy, Tags) " +
		` VALUES ( "` + request.SecretName + `", "` + request.UserId + `", "` + request.ServiceProvider + `", "` + strconv.Itoa(request.KeyType) + `", "` + request.UserId + `", "` + request.UserId + `", "` + request.Tags + `" )`

	insert, err := Db.Query(query)

	switch request.ServiceProvider {
	case "Google":
		query = "INSERT INTO " + gkeSecretTable + " ( Type, ProjectId, SecretId, ProjectKeyId, PrivateKey, ClientMail, ClientId, ClientX509CertUrl) " +
			` VALUES ( "` + request.GkeType + `", "` + request.GkeProjectId + `", ` + LAST_INSERT_ID() + `, "` + strconv.Itoa(request.KeyType) + `", "` + request.UserId + `", "` + request.UserId + `", "` + request.Tags + `" )`
	default:
		query = ""
	}
	insert, err = Db.Query(query)

	if err != nil {
		log.Logger.ErrorContext(ctx, fmt.Sprintf("insert to %s table failed", clusterTable), err)
		return domain.DAOResult{
			make([]Secret, 0),
			1,
			"insert failed",
			err,
		}
	}

	defer insert.Close()
	log.Logger.TraceContext(ctx, "successfully added a new key", secretName)

	return domain.DAOResult{
		make([]Secret, 0),
		0,
		"insert success",
		nil,
	}
}

//TODO: implement delete query
//func DeleteSecret(ctx context.Context, secretId string)  (result DAOResult) {
//	query := "DELETE FROM " + clusterTable + `  WHERE cluster_name="` + clusterName + `";`
//
//	del, err := Db.Query(query)
//	if err != nil {
//		log.Logger.ErrorContext(ctx, fmt.Sprintf("deleting cluster %s failed", clusterName), err)
//		return err
//	}
//
//	defer del.Close()
//	log.Logger.TraceContext(ctx, "delete cluster db query was successful", clusterName)
//
//	return nil
//}

func (dao *domain.SecretDAO) GetAllSecretsByUser(ctx context.Context, userId string) (result DAOResult) {
	query := "SELECT * FROM " + clusterTable + " WHERE OwnerId = '" + userId + "';"

	rows, err := Db.Query(query)
	if err != nil {
		log.Logger.ErrorContext(ctx, "get all db query failed", err)
		return domain.DAOResult{
			make([]Secret, 0),
			1,
			"db query failed",
			err,
		}
	}

	defer rows.Close()
	secretList := make([]Secret, 1)

	for rows.Next() {
		secret := domain.Secret{}

		err = rows.Scan(&secret.ID, &secret.Name, &cluster.KafkaVersion, &cluster.ActiveControllers)
		if err != nil {
			log.Logger.ErrorContext(ctx, "scanning rows in secret table failed", err)
			return domain.DAOResult{
				secretList,
				1,
				"row scanning failed",
				err,
			}
		}

		secretList = append(secretList, secret)
	}

	log.Logger.TraceContext(ctx, "get all secrets db query was successful")
	return domain.DAOResult{
		secretList,
		1,
		"successful get secrets",
		nil,
	}
}
