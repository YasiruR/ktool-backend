package database

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/YasiruR/ktool-backend/domain"
	"github.com/YasiruR/ktool-backend/log"
	"github.com/go-sql-driver/mysql"
	"golang.org/x/crypto/bcrypt"
	"strconv"
)

func AddNewUser(ctx context.Context, username, password, token string, accessLevel int, firstName, lastName, email string) (exists bool, err error) {

	encryptedPass, err := hashPassword(ctx, password)
	if err != nil {
		log.Logger.ErrorContext(ctx, "adding new user to the db failed")
		return false, err
	}

	query := "INSERT INTO " + userTable + ` (id, username, password, token, access_level, first_name, last_name, email) VALUES (null, "` + username + `", "` + encryptedPass + `", "` + token + `", ` + strconv.Itoa(accessLevel) + `, "` + firstName + `", "` + lastName + `", "` + email + `");`

	insert, err := Db.Query(query)
	if err != nil {
		if err.(*mysql.MySQLError).Number == 1062 {
			log.Logger.ErrorContext(ctx, fmt.Sprintf("%v error for user %v", err, username))
			return true, err
		}
		log.Logger.ErrorContext(ctx, fmt.Sprintf("insert to %s table failed", userTable), err)
		return false, err
	}

	defer insert.Close()
	log.Logger.TraceContext(ctx, "add new user query was successful", username)

	return false, nil
}

func UpdateToken(ctx context.Context, username, token string) (err error) {
	tx, err := Db.Begin()
	if err != nil {
		log.Logger.ErrorContext(ctx, "starting the transaction failed", err, username)
		return err
	}
	defer tx.Rollback()

	query, err := tx.Prepare("UPDATE " + userTable + ` SET token="` + token + `" WHERE username="` + username +`";`)
	if err != nil {
		log.Logger.ErrorContext(ctx, "preparing the query failed", username)
		return err
	}

	_, err = query.Exec()
	if err != nil {
		log.Logger.ErrorContext(ctx, "executing the update token query failed", err, username)
		return err
	}

	err = tx.Commit()
	if err != nil {
		log.Logger.ErrorContext(ctx, "committing the transaction failed", username)
		return err
	}

	log.Logger.TraceContext(ctx, "updating user token query was successful", username)
	return nil
}

func GetUserByToken(ctx context.Context, token string) (username string, ok bool, err error) {
	query := "SELECT username from " + userTable + ` WHERE token="` + token + `";`

	row := Db.QueryRow(query)

	switch err := row.Scan(&username); err {
	case sql.ErrNoRows:
		log.Logger.ErrorContext(ctx, "no rows scanned for the token", token)
		return "", false, errors.New("no rows found")
	case nil:
		log.Logger.TraceContext(ctx, "fetched user by token", token)
		return token, true, nil
	default:
		log.Logger.ErrorContext(ctx, "unhandled error in row scan", token)
		return "", false, errors.New("row scan failed")
	}
}

func GetUserTokenByName(ctx context.Context, username string) (token string, err error) {
	query := "SELECT token from " + userTable + ` WHERE username="` + username + `";`

	row := Db.QueryRow(query)

	switch err := row.Scan(&token); err {
	case sql.ErrNoRows:
		log.Logger.ErrorContext(ctx, "no rows scanned for the user", username)
		return "", errors.New("no rows found")
	case nil:
		log.Logger.TraceContext(ctx, "fetched user by username", username)
		return token, nil
	default:
		log.Logger.ErrorContext(ctx, "unhandled error in row scan", username)
		return "", errors.New("row scan failed")
	}
}

func ValidateUserByPassword(ctx context.Context, username, password string) (id int, ok bool, err error) {

	var encryptedPass string
	query := "SELECT id, password from " + userTable + ` WHERE username="` + username + `";`
	row := Db.QueryRow(query)

	switch err := row.Scan(&id, &encryptedPass); err {
	case sql.ErrNoRows:
		log.Logger.ErrorContext(ctx, "no rows scanned for the user", username)
		return id,false, errors.New("incorrect username")
	case nil:
		if checkPasswordForHash(ctx, password, encryptedPass) {
			log.Logger.TraceContext(ctx, "fetched user by username and password", username)
			return id,true, nil
		}
		log.Logger.ErrorContext(ctx, "user found with incorrect password")
		return id, false, errors.New("incorrect password")
	default:
		log.Logger.ErrorContext(ctx, "unhandled error in row scan", username, err)
		return id,false, errors.New("row scan failed")
	}
}

func ValidateUserByToken(ctx context.Context, token string) (id int, ok bool, err error) {
	query := "SELECT id from " + userTable + ` WHERE token="` + token + `";`
	row := Db.QueryRow(query)

	switch err := row.Scan(&id); err {
	case sql.ErrNoRows:
		log.Logger.ErrorContext(ctx, "no rows scanned for the token", token)
		return id,false, errors.New("incorrect credentials")
	case nil:
		log.Logger.TraceContext(ctx, "fetched user by token", token)
		return id,true, nil
	default:
		log.Logger.ErrorContext(ctx, "unhandled error in row scan", token, err)
		return id,false, errors.New("row scan failed")
	}
}

func GetAllUsers(ctx context.Context) (userList []domain.User, err error) {
	query := "SELECT (id, name, token, access_level) FROM " + userTable + ";"

	rows, err := Db.Query(query)
	if err != nil {
		log.Logger.ErrorContext(ctx, "get all users db query failed", err)
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		user := domain.User{}

		err = rows.Scan(&user.Id, &user.Username, &user.Token, &user.AccessLevel)
		if err != nil {
			log.Logger.ErrorContext(ctx, "scanning rows in user table failed", err)
			return nil, err
		}

		userList = append(userList, user)
	}

	err = rows.Err()
	if err != nil {
		log.Logger.ErrorContext(ctx, "error occurred when scanning rows", err)
		return nil, err
	}

	log.Logger.TraceContext(ctx, "get all users db query was successful")
	return userList, nil
}

func hashPassword(ctx context.Context, password string) (hash string, err error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		log.Logger.ErrorContext(ctx, "generating hash for the password failed", err)
		return "", err
	}
	return string(bytes), nil
}

func checkPasswordForHash(ctx context.Context, password, hash string) (ok bool) {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	if err != nil {
		log.Logger.ErrorContext(ctx, "comparing hash and password failed", err)
		return false
	}
	return true
}