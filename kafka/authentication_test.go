package kafka

import (
	"context"
	"github.com/YasiruR/ktool-backend/log"
	log2 "github.com/pickme-go/log"
	"testing"
)

func TestGenerateSSLKeys(t *testing.T) {
	tests := []struct{
		clusterName 	string
	}{
		{"test_cluster"},
	}

	log.Logger = log2.Constructor.Log(log2.WithColors(true), log2.WithLevel("TRACE"), log2.WithFilePath(true))
	ctx := context.Background()

	//for _, test := range tests {
	//	_, _, _ = generateRSAKeys(ctx, test.clusterName)
	//}
}
