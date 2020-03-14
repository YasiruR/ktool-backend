package domain

import "github.com/YasiruR/ktool-backend/kafka"

var LoggedInUsers []User

type User struct {
	Id 					int
	Username			string
	Token				string
	AccessLevel 		int
	ConnectedClusters	[]kafka.KCluster
}
