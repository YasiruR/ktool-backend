package database

import (
	"database/sql"
	"github.com/YasiruR/ktool-backend/log"
	"strconv"
)

var Database *sql.DB

func Init() {
	dataSource := Cfg.Username + ":" + Cfg.Password + "@tcp(" + Cfg.Host + ":" + strconv.Itoa(Cfg.Port) + "/" + Cfg.Database
	Database, err := sql.Open("mysql", dataSource)
	if err != nil {
		log.Logger.Fatal("failed in initializing mysql connection", err)
	}
}
