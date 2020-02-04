package kafka

import (
	"context"
	"fmt"
	"github.com/Shopify/sarama"
	"github.com/YasiruR/ktool-backend/log"
	log2 "github.com/pickme-go/log"
	"testing"
)

func TestGetAllClusters(t *testing.T) {
	tests := []struct{
		brokerList   []string
		brokers      []*sarama.Broker
		controllerId int32
	}{
		//{[]string{"capp-kafka-kfk-001.dev-mytaxi.com:9092"}, nil, 0},
		{[]string{"34.87.54.215:9092"}, nil, 0},
	}

	log.Logger = log2.Constructor.Log(log2.WithColors(true), log2.WithLevel("TRACE"), log2.WithFilePath(true))
	ctx := context.Background()

	for _, test := range tests {
		res, err := get(ctx, test.brokerList)
		if err != nil {
			t.Errorf("failed")
			t.Fail()
		}
		fmt.Printf("Brokers : %v", res)
		//fmt.Println("id : ", id)
	}
}
