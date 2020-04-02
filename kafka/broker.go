package kafka

import (
	"context"
	"github.com/YasiruR/ktool-backend/log"
)

func CheckIfBrokerExists(ctx context.Context, broker string) (exists bool) {
	for _, cluster := range ClusterList {
		for _, b := range cluster.Brokers {
			if broker == b.Addr() {
				log.Logger.WarnContext(ctx, "broker already exists", broker)
				return true
			}
		}
	}

	return false
}
