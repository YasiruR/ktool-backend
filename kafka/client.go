package kafka

import (
	"context"
	"errors"
	"github.com/YasiruR/ktool-backend/log"
)

func GetClient(ctx context.Context, clusterID int) (clustClient *KCluster, err error) {
	var found bool
	for _, clus := range ClusterList {
		if clus.ClusterID == clusterID {
			clustClient = clus
			found = true
			break
		}
	}

	if !found {
		log.Logger.ErrorContext(ctx, "client not found for cluster id", clusterID)
		return clustClient, errors.New("client not found")
	}

	return clustClient, nil
}