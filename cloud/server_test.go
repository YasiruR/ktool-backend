package cloud

import (
	"context"
	"github.com/YasiruR/ktool-backend/service"
	"testing"
)

func TestPingToServer(t *testing.T) {
	tests := [] struct{
		ipAddress	string
		out 		bool
	}{
		{"www.google.com", true},
		//{"35.247.188.238", false},
	}

	service.Cfg.PingRetry = 3
	ctx := context.Background()

	for _, test := range tests {
		ok, err := pingToServer(ctx, test.ipAddress)
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
