package cloud

import (
	"context"
	"github.com/YasiruR/ktool-backend/log"
	"github.com/YasiruR/ktool-backend/service"
	log2 "github.com/pickme-go/log"
	"testing"
)

func TestPingToServer(t *testing.T) {
	tests := [] struct{
		ipAddress	string
		out 		bool
	}{
		{"www.google.com", true},
		{"35.247.188.238", false},
	}

	log.Logger = log2.Constructor.Log(log2.WithColors(true), log2.WithLevel("TRACE"), log2.WithFilePath(true))
	service.Cfg.PingRetry = 3
	service.Cfg.PingTimeout = 5
	ctx := context.Background()

	for _, test := range tests {
		ok, err := PingToServer(ctx, test.ipAddress)
		if ok != test.out {
			if err != nil {
				t.Errorf("ping failed : %v", err)
				t.Fail()
			}
			t.Errorf("received unexpected value %v for test %v", ok, test.ipAddress)
			t.Fail()
		}
	}
}

func TestTelnetToPort(t *testing.T) {
	tests := [] struct {
		ipAddress 	string
		port 		int
		out 		bool
	}{
		{"localhost", 5000, true},
		{"www.google.com", 80, true},
		{"35.186.213.42", 500, false},
		{"127.0.0.1", 5000, true},
	}

	log.Logger = log2.Constructor.Log(log2.WithColors(true), log2.WithLevel("TRACE"), log2.WithFilePath(true))
	ctx := context.Background()
	service.Cfg.PingTimeout = 5

	for index, test := range tests {
		ok, _ := TelnetToPort(ctx, test.ipAddress, test.port)
		if ok != test.out {
			t.Error("unexpected output for test case ", index+1)
			t.Fail()
		}
	}
}
