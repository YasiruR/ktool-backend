package database

import (
	"context"
	"database/sql"
	"github.com/YasiruR/ktool-backend/log"
	log2 "github.com/pickme-go/log"
	"testing"
)

func TestGetClusterIdByName(t *testing.T) {
	tests := []struct{
		clusterName		string
		out 			int
	}{
		{"cluster_2", 1},
	}

	db, err := sql.Open("mysql", "yasi:123@tcp(localhost:3306)/kdb")
	if err != nil {
		log.Logger.Fatal("failed in initializing mysql connection", err)
	}

	Db = db

	log.Logger = log2.Constructor.Log(log2.WithColors(true), log2.WithLevel("TRACE"), log2.WithFilePath(true))
	ctx := context.Background()

	for _, test := range tests {
		res, _ := GetClusterIdByName(ctx, test.clusterName)
		if res != test.out {
			t.Error(res, test.out)
			t.Fail()
		}
	}
}
