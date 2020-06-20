package database

import (
	"database/sql"
	"github.com/YasiruR/ktool-backend/log"
	_ "github.com/go-sql-driver/mysql"
	"strconv"
)

const (
	clusterTable     = "cluster"
	zookeeperTable   = "zookeeper"
	brokerTable      = "broker"
	userTable        = "user"
	secretTable      = "secret"
	gkeSecretTable   = "gke_secret"
	cloudSecretTable = "cloud_secret"
	operationsTable  = "operations"
	k8sTable         = "k8s_clusters"
)

var Db *sql.DB

func Init() {
	dataSource := Cfg.Username + ":" + Cfg.Password + "@tcp(" + Cfg.Host + ":" + strconv.Itoa(Cfg.Port) + ")/" + Cfg.Database

	db, err := sql.Open("mysql", dataSource)
	if err != nil {
		log.Logger.Fatal("failed in initializing mysql connection", err)
	}

	Db = db
}
