package database

import (
	"context"
	"database/sql"
	"github.com/YasiruR/ktool-backend/log"
	log2 "github.com/pickme-go/log"
	"testing"
)

func TestGetZookeeperByClusterId(t *testing.T) {
	tests := []struct{
		clusterId 		int
		out 			struct{
			id 		int
			host 	string
			port 	int
		}
	}{
		{8, struct{
			id 		int
			host 	string
			port 	int
		}{1, "152.301.160.109", 3000}},
	}

	log.Logger = log2.Constructor.Log(log2.WithColors(true), log2.WithLevel("TRACE"), log2.WithFilePath(true))
	ctx := context.Background()

	db, err := sql.Open("mysql", "yasi:123@tcp(localhost:3306)/kdb")
	if err != nil {
		log.Logger.Fatal("failed in initializing mysql connection", err)
	}

	Db = db

	for _, test := range tests {
		id, host, port,  _ := GetZookeeperByClusterId(ctx, test.clusterId)
		if id != test.out.id {
			t.Error("id mismatch")
			t.Fail()
		}
		if host != test.out.host {
			t.Error("host mismatch")
			t.Fail()
		}
		if port != test.out.port {
			t.Error("port mismatch")
			t.Fail()
		}
	}
}
