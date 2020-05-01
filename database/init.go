package database

import (
	"database/sql"
	"github.com/YasiruR/ktool-backend/log"
	_ "github.com/go-sql-driver/mysql"
	"strconv"
)

const (
	clusterTable  		= 	"cluster"
	zookeeperTable 		= 	"zookeeper"
	brokerTable			=	"broker"
	userTable			= 	"user"
	brokerBytesInTable	= 	"broker_bytes_in"
	brokerBytesOutTable	= 	"broker_bytes_out"
	metricsLimit 		= 	5
)

var Db *sql.DB

func Init() {
	dataSource := Cfg.Username + ":" + Cfg.Password + "@tcp(" + Cfg.Host + ":" + strconv.Itoa(Cfg.Port) + ")/" + Cfg.Database

	db, err := sql.Open("mysql", dataSource)
	if err != nil {
		log.Logger.Fatal("failed in initializing mysql connection", err)
	}

	Db = db
	log.Logger.Info("connection to mysql database is established")
}
