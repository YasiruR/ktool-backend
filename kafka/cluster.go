package kafka

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"errors"
	"fmt"
	"github.com/Shopify/sarama"
	"github.com/YasiruR/ktool-backend/domain"
	"github.com/YasiruR/ktool-backend/log"
	"github.com/YasiruR/ktool-backend/service"
	"io/ioutil"
	"time"
)

func InitSaramaConfig(ctx context.Context,  clusterName string, networkSecurity string) (config *sarama.Config, err error) {
	config = sarama.NewConfig()
	config.Consumer.Return.Errors = true
	config.Version = sarama.V2_0_0_0	//todo: add this as a parameter

	if networkSecurity == "tsl" {
		caFile, certFile, keyFile, err := generateRSAKeys(ctx, clusterName)
		if err != nil {
			return config, err
		}

		cert, err := tls.LoadX509KeyPair(certFile, keyFile)
		if err != nil {
			log.Logger.ErrorContext(ctx, "loading X509 key pair failed", err)
			return config, err
		}

		ca, err := ioutil.ReadFile(caFile)
		if err != nil {
			log.Logger.ErrorContext(ctx, "reading ca pem file failed", err)
			return config, err
		}

		pool := x509.NewCertPool()
		pool.AppendCertsFromPEM(ca)

		tlsConfig := &tls.Config{
			Certificates: []tls.Certificate{cert},
			RootCAs: pool,
		}

		config.Net.TLS.Config = tlsConfig
	}

	err = config.Validate()
	if err != nil {
		log.Logger.ErrorContext(ctx, err, "invalid configuration")
		return
	}

	return
}

func InitSaramaConsumer(ctx context.Context, brokers []string, config *sarama.Config) (consumer sarama.Consumer, err error) {
	consumer, err = sarama.NewConsumer(brokers, config)
	if err != nil {
		log.Logger.ErrorContext(ctx, err, brokers)
		return nil, errors.New("kafka cluster config initialization failed")
	}

	//todo: close cluster connection on disconnect

	log.Logger.TraceContext(ctx, "new consumer initialized with config", config)

	return consumer, nil
}

func GetTopicList(ctx context.Context, cluster sarama.Consumer) (topics []string, err error) {
	topics, err = cluster.Topics()
	if err != nil {
		log.Logger.ErrorContext(ctx, err, cluster)
		return nil, err
	}

	log.Logger.TraceContext(ctx, "all topics are fetched", fmt.Sprintf("no of topics : %v", len(topics)), fmt.Sprintf("cluster : %v", cluster))
	return topics, nil
}

func InitClient(ctx context.Context, brokers []string, config *sarama.Config) (client sarama.Client, err error) {
	done := make(chan string)

	go func() {
		client, err = sarama.NewClient(brokers, config)
		if err != nil {
			log.Logger.ErrorContext(ctx, fmt.Sprintf("creating new client failed for brokers : %v", brokers), err)
			done <- "err"
		} else {
			done <- "done"
		}
	}()

	select {
	case out := <- done:
		if out == "done" {
			log.Logger.TraceContext(ctx, "client initialized successfully", brokers)
			return client, nil
		} else if out == "err" {
			return nil, err
		}
	case <- time.After(time.Duration(int64(service.Cfg.ClientInitTimeout)) * time.Second):
		log.Logger.ErrorContext(ctx, "client init timeout for brokers", brokers)
		return nil, errors.New("client timeout")
	}

	return
}

func GetBrokerAddrList(ctx context.Context, client sarama.Client) (addrList []string, err error) {
	brokers := client.Brokers()
	for _, broker := range brokers {
		addrList = append(addrList, broker.Addr())
	}

	log.Logger.TraceContext(ctx, "broker address list fetched")
	return addrList, nil
}

func DeleteCluster(ctx context.Context, clusterID int, users []domain.User) (err error) {
	for index, cluster := range ClusterList {
		if clusterID == cluster.ClusterID {
			//remove from all clusters
			ClusterList[index] = ClusterList[len(ClusterList)-1] // Copy last element to index i.
			ClusterList[len(ClusterList)-1] = domain.KCluster{}   // Erase last element (write zero value).
			ClusterList = ClusterList[:len(ClusterList)-1]   // Truncate slice.

			//for i, sCluster := range SelectedClusterList {
			//	if clusterID == sCluster.ClusterID {
			//		//remove from selected clusters, if selected
			//		SelectedClusterList[i] = SelectedClusterList[len(SelectedClusterList)-1] // Copy last element to index i.
			//		SelectedClusterList[len(SelectedClusterList)-1] = KCluster{}   // Erase last element (write zero value).
			//		SelectedClusterList = SelectedClusterList[:len(SelectedClusterList)-1]   // Truncate slice.
			//	}
			//}

			for _, user := range users {
				for i, cluster := range user.ConnectedClusters {
					if clusterID == cluster.ClusterID {
						//remove from selected clusters, if selected
						user.ConnectedClusters[i] = user.ConnectedClusters[len(user.ConnectedClusters)-1] // Copy last element to index i.
						user.ConnectedClusters[len(user.ConnectedClusters)-1] = domain.KCluster{}   // Erase last element (write zero value).
						user.ConnectedClusters = user.ConnectedClusters[:len(user.ConnectedClusters)-1]   // Truncate slice.
					}
				}
			}

			log.Logger.TraceContext(ctx, "removed cluster from cluster lists successfully")
			return
		}
	}

	log.Logger.ErrorContext(ctx, "could not find a cluster with the matched id", clusterID)
	return errors.New("unable to find cluster")
}